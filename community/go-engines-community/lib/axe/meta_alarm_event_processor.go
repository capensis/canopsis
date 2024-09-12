package axe

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	libaxeevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
)

const workers = 10

func NewMetaAlarmEventProcessor(
	dbClient mongo.DbClient,
	adapter libalarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmStatusService alarmstatus.Service,
	alarmConfigProvider config.AlarmConfigProvider,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	metricsSender metrics.Sender,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) libalarm.MetaAlarmEventProcessor {
	return &metaAlarmEventProcessor{
		dbClient:                  dbClient,
		alarmCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		metaAlarmStatesCollection: dbClient.Collection(mongo.MetaAlarmStatesCollection),
		metaAlarmStatesService:    metaAlarmStatesService,
		adapter:                   adapter,
		ruleAdapter:               ruleAdapter,
		alarmStatusService:        alarmStatusService,
		alarmConfigProvider:       alarmConfigProvider,
		encoder:                   encoder,
		amqpPublisher:             amqpPublisher,
		metricsSender:             metricsSender,
		templateExecutor:          templateExecutor,
		logger:                    logger,
	}
}

type metaAlarmEventProcessor struct {
	dbClient                  mongo.DbClient
	alarmCollection           mongo.DbCollection
	metaAlarmStatesCollection mongo.DbCollection
	metaAlarmStatesService    correlation.MetaAlarmStateService
	adapter                   libalarm.Adapter
	ruleAdapter               correlation.RulesAdapter

	alarmStatusService  alarmstatus.Service
	alarmConfigProvider config.AlarmConfigProvider

	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher

	metricsSender metrics.Sender

	logger zerolog.Logger

	templateExecutor template.Executor
}

func (p *metaAlarmEventProcessor) CreateMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, []types.Alarm, error) {
	if event.Entity == nil {
		return nil, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = make([]types.Alarm, 0)
		metaAlarm = types.Alarm{}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		metaAlarm = p.newAlarm(event.Parameters, *event.Entity, p.alarmConfigProvider.Get())
		metaAlarm.Value.Meta = event.Parameters.MetaAlarmRuleID
		metaAlarm.Value.MetaValuePath = event.Parameters.MetaAlarmValuePath
		metaAlarm.Value.LastEventDate = types.CpsTime{} // should be empty

		if event.Parameters.DisplayName != "" {
			metaAlarm.Value.DisplayName = event.Parameters.DisplayName
		}

		stateID := rule.ID
		if event.Parameters.MetaAlarmValuePath != "" {
			stateID = fmt.Sprintf("%s&&%s", rule.ID, event.Parameters.MetaAlarmValuePath)
		}

		var childEntityIDs []string
		var archived bool

		if rule.IsManual() {
			childEntityIDs = event.Parameters.MetaAlarmChildren
		} else {
			metaAlarmState, err := p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
			if err != nil {
				return err
			}

			if metaAlarmState.MetaAlarmName != event.Entity.Name {
				// try to get archived state
				metaAlarmState, err = p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID+"-"+event.Entity.Name)
				if err != nil {
					return err
				}

				if metaAlarmState.ID == "" {
					return fmt.Errorf("meta alarm state for rule id=%q and meta alarm name=%q not found", event.Parameters.MetaAlarmRuleID, event.Entity.Name)
				}

				archived = true
			}

			childEntityIDs = metaAlarmState.ChildrenEntityIDs
		}

		var lastChild types.AlarmWithEntity
		worstState := types.CpsNumber(types.AlarmStateMinor)
		eventsCount := types.CpsNumber(0)

		if len(childEntityIDs) > 0 {
			childAlarms, err := p.getAlarmsWithEntityByEntityIds(ctx, childEntityIDs)
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}

			if len(childAlarms) > 0 {
				lastChild = childAlarms[len(childAlarms)-1]
			}

			newStep := types.NewMetaAlarmAttachStep(metaAlarm, rule.Name)
			for _, childAlarm := range childAlarms {
				if childAlarm.Alarm.Value.State != nil {
					childState := childAlarm.Alarm.Value.State.Value
					if childState > worstState {
						worstState = childState
					}
				}

				if !childAlarm.Alarm.HasParentByEID(metaAlarm.EntityID) {
					childAlarm.Alarm.AddParent(metaAlarm.EntityID)
					err = childAlarm.Alarm.PartialUpdateAddStepWithStep(newStep)
					if err != nil {
						return err
					}
					metaAlarm.AddChild(childAlarm.Alarm.EntityID)
					updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)

					eventsCount += childAlarm.Alarm.Value.EventsCount
					if metaAlarm.Value.LastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
						metaAlarm.Value.LastEventDate = childAlarm.Alarm.Value.LastEventDate
					}
				}
			}
		}

		output := ""
		if rule.IsManual() {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
				Rule:     rule,
				Count:    int64(len(updatedChildrenAlarms)),
				Children: lastChild,
			})
			if err != nil {
				return err
			}
		}

		metaAlarm.UpdateOutput(output)
		err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, output, p.alarmStatusService)
		if err != nil {
			return err
		}

		metaAlarm.Value.EventsCount = eventsCount

		err = p.adapter.Insert(ctx, metaAlarm)
		if err != nil {
			return fmt.Errorf("cannot create alarm: %w", err)
		}

		err = p.adapter.PartialMassUpdateOpen(ctx, updatedChildrenAlarms)
		if err != nil {
			return fmt.Errorf("cannot update children alarms: %w", err)
		}

		if !rule.IsManual() && !archived {
			ok, err := p.metaAlarmStatesService.SwitchStateToCreated(ctx, stateID)
			if err != nil || !ok {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return &metaAlarm, updatedChildrenAlarms, nil
}

func (p *metaAlarmEventProcessor) AttachChildrenToMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, []types.Alarm, []types.Event, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil, nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var err error

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = make([]types.Alarm, 0)
		var lastChild types.AlarmWithEntity

		err = p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID}).Decode(&metaAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		alarms, err := p.getAlarmsWithEntityByEntityIds(ctx, event.Parameters.MetaAlarmChildren)
		if err != nil {
			return err
		}

		lastEventDate := metaAlarm.Value.LastEventDate
		newStep := types.NewMetaAlarmAttachStep(metaAlarm, rule.Name)
		worstState := types.CpsNumber(types.AlarmStateOK)
		eventsCount := types.CpsNumber(0)

		for _, childAlarm := range alarms {
			if !childAlarm.Alarm.AddParent(metaAlarm.EntityID) {
				continue
			}

			metaAlarm.AddChild(childAlarm.Entity.ID)
			eventsCount += childAlarm.Alarm.Value.EventsCount
			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}

			err = childAlarm.Alarm.PartialUpdateAddStepWithStep(newStep)
			if err != nil {
				return err
			}

			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
			lastChild = childAlarm
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		if metaAlarm.Value.Meta == "" {
			metaAlarm.SetMeta(event.Parameters.MetaAlarmRuleID)
			metaAlarm.SetMetaValuePath(event.Parameters.MetaAlarmValuePath)
		}

		if metaAlarm.Value.LastEventDate.Unix() != lastEventDate.Unix() {
			metaAlarm.PartialUpdateLastEventDate(lastEventDate)
		}

		metaAlarm.IncrementEventsCount(eventsCount)
		childrenCount, err := p.adapter.GetCountOpenedAlarmsByIDs(ctx, metaAlarm.Value.Children)
		if err != nil {
			return err
		}

		output := ""
		if rule.Type == correlation.RuleTypeManualGroup {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
				Rule:     rule,
				Count:    childrenCount,
				Children: lastChild,
			})
			if err != nil {
				return err
			}
		}

		if output == "" {
			output = metaAlarm.Value.Output
		}

		metaAlarm.UpdateOutput(output)
		if worstState > metaAlarm.CurrentState() {
			err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, output, p.alarmStatusService)
			if err != nil {
				return err
			}
		}

		return p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{metaAlarm}, updatedChildrenAlarms...))
	})
	if err != nil {
		return nil, nil, nil, err
	}

	if metaAlarm.ID == "" {
		return nil, nil, nil, nil
	}

	childrenEvents, err := p.applyActionsOnChildren(metaAlarm, updatedChildrenAlarms)
	if err != nil {
		return nil, nil, nil, err
	}

	return &metaAlarm, updatedChildrenAlarms, childrenEvents, nil
}

func (p *metaAlarmEventProcessor) DetachChildrenFromMetaAlarm(ctx context.Context, event rpc.AxeEvent) (*types.Alarm, error) {
	if len(event.Parameters.MetaAlarmChildren) == 0 {
		return nil, nil
	}

	var updatedChildrenAlarms []types.Alarm
	var metaAlarm types.Alarm
	var err error

	err = p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildrenAlarms = make([]types.Alarm, 0)

		err = p.alarmCollection.FindOne(ctx, bson.M{"d": event.Entity.ID}).Decode(&metaAlarm)
		if err != nil {
			if errors.Is(err, mongodriver.ErrNoDocuments) {
				return nil
			}

			return err
		}

		rule, err := p.ruleAdapter.GetRule(ctx, event.Parameters.MetaAlarmRuleID)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.Parameters.MetaAlarmRuleID, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", event.Parameters.MetaAlarmRuleID)
		}

		alarms, err := p.getAlarmsWithEntityByParentIdAndEntityIds(ctx, metaAlarm.EntityID, event.Parameters.MetaAlarmChildren)
		if err != nil {
			return err
		}

		eventsCount := types.CpsNumber(0)

		for _, childAlarm := range alarms {
			if childAlarm.Alarm.RemoveParent(event.Entity.ID) {
				metaAlarm.RemoveChild(childAlarm.Entity.ID)
				eventsCount -= childAlarm.Alarm.Value.EventsCount
				updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
			}
		}

		if len(updatedChildrenAlarms) == 0 {
			return nil
		}

		metaAlarmChildren, err := p.getAlarmsWithEntityByParentIdAndEntityIds(ctx, metaAlarm.EntityID, metaAlarm.Value.Children)
		if err != nil {
			return err
		}

		var lastEventDate types.CpsTime // should be empty
		worstState := types.CpsNumber(types.AlarmStateOK)

		for _, childAlarm := range metaAlarmChildren {
			if childAlarm.Alarm.Value.State.Value > worstState {
				worstState = childAlarm.Alarm.Value.State.Value
			}

			if lastEventDate.Before(childAlarm.Alarm.Value.LastEventDate) {
				lastEventDate = childAlarm.Alarm.Value.LastEventDate
			}
		}

		metaAlarm.PartialUpdateLastEventDate(lastEventDate)
		metaAlarm.IncrementEventsCount(eventsCount)
		infos := correlation.EventExtraInfosMeta{
			Rule:  rule,
			Count: int64(len(metaAlarmChildren)),
		}
		if len(metaAlarmChildren) > 0 {
			infos.Children = metaAlarmChildren[len(metaAlarmChildren)-1]
		}

		output := ""
		if rule.IsManual() {
			output = event.Parameters.Output
		} else {
			output, err = p.executeOutputTpl(infos)
			if err != nil {
				return err
			}
		}

		if output == "" {
			output = metaAlarm.Value.Output
		}

		metaAlarm.UpdateOutput(output)
		err = updateMetaAlarmState(&metaAlarm, *event.Entity, event.Parameters.Timestamp, worstState, output, p.alarmStatusService)
		if err != nil {
			return err
		}

		return p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{metaAlarm}, updatedChildrenAlarms...))
	})
	if err != nil {
		return nil, err
	}

	return &metaAlarm, nil
}

func (p *metaAlarmEventProcessor) applyActionsOnChildren(metaAlarm types.Alarm, childrenAlarms []types.Alarm) ([]types.Event, error) {
	var events []types.Event

	steps := metaAlarm.GetAppliedActions()

	for _, childAlarm := range childrenAlarms {
		childEvent := types.Event{
			Connector:     childAlarm.Value.Connector,
			ConnectorName: childAlarm.Value.ConnectorName,
			Resource:      childAlarm.Value.Resource,
			Component:     childAlarm.Value.Component,
			Timestamp:     types.NewCpsTime(),
		}
		childEvent.SourceType = childEvent.DetectSourceType()

		for _, step := range steps {
			childEvent.Output = step.Message
			childEvent.Author = step.Author
			childEvent.UserID = step.UserID
			childEvent.Initiator = step.Initiator
			childEvent.Role = step.Role
			switch step.Type {
			case types.AlarmStepAck:
				childEvent.EventType = types.EventTypeAck
			case types.AlarmStepSnooze:
				childEvent.EventType = types.EventTypeSnooze
				childEvent.Duration = types.CpsNumber(int64(step.Value) - childEvent.Timestamp.Unix())
				if childEvent.Duration <= 0 {
					continue
				}
			case types.AlarmStepAssocTicket:
				childEvent.EventType = types.EventTypeAssocTicket
				childEvent.TicketInfo = step.TicketInfo
				childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
			case types.AlarmStepDeclareTicket:
				childEvent.EventType = types.EventTypeDeclareTicketWebhook
				childEvent.TicketInfo = step.TicketInfo
				childEvent.TicketInfo.TicketMetaAlarmID = metaAlarm.ID
			case types.AlarmStepComment:
				childEvent.EventType = types.EventTypeComment
			default:
				continue
			}

			events = append(events, childEvent)
		}
	}

	return events, nil
}

func (p *metaAlarmEventProcessor) ProcessAxeRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if eventRes.Alarm == nil {
		return nil
	}

	alarm := eventRes.Alarm

	err := p.processComponentRpc(ctx, event, eventRes)
	if err != nil {
		return err
	}

	if alarm.IsMetaAlarm() {
		return p.processParentRpc(ctx, event, eventRes)
	}

	if alarm.IsMetaChild() {
		return p.processChildRpc(ctx, event, eventRes)
	}

	return nil
}

func (p *metaAlarmEventProcessor) processParentRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !applyOnChild(eventRes.AlarmChangeType) {
		return nil
	}

	childEvent := types.Event{
		EventType:  event.EventType,
		Timestamp:  types.NewCpsTime(),
		Output:     event.Parameters.Output,
		Author:     event.Parameters.Author,
		UserID:     event.Parameters.User,
		Initiator:  event.Parameters.Initiator,
		TicketInfo: event.Parameters.TicketInfo,
	}

	childEvent.TicketInfo.TicketMetaAlarmID = eventRes.Alarm.ID

	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		childEvent.EventType = types.EventTypeDeclareTicketWebhook
	}

	if event.Parameters.State != nil {
		childEvent.State = *event.Parameters.State
	}

	if event.Parameters.Duration != nil {
		seconds, err := event.Parameters.Duration.To("s")
		if err == nil {
			childEvent.Duration = types.CpsNumber(seconds.Value)
		}
	}

	err := p.sendChildrenEvents(ctx, eventRes.Alarm.Value.Children, childEvent)
	if err != nil {
		return err
	}

	return nil
}

func (p *metaAlarmEventProcessor) processComponentRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !event.Parameters.TicketResources ||
		(eventRes.AlarmChangeType != types.AlarmChangeTypeDeclareTicketWebhook && eventRes.AlarmChangeType != types.AlarmChangeTypeAutoDeclareTicketWebhook) ||
		event.Entity.Type != types.EntityTypeComponent {
		return nil
	}

	componentAlarm := eventRes.Alarm
	if componentAlarm == nil || componentAlarm.Value.Ticket == nil {
		return nil
	}

	resources, err := p.adapter.GetAlarmsWithoutTicketByComponent(ctx, eventRes.Alarm.Value.Component)
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, resource := range resources {
		if resource.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent := types.Event{
			EventType:     types.EventTypeDeclareTicketWebhook,
			Connector:     resource.Alarm.Value.Connector,
			ConnectorName: resource.Alarm.Value.ConnectorName,
			Resource:      resource.Alarm.Value.Resource,
			Component:     resource.Alarm.Value.Component,
			Timestamp:     types.NewCpsTime(),
			Output:        componentAlarm.Value.Ticket.Message,
			TicketInfo:    componentAlarm.Value.Ticket.TicketInfo,
			Author:        event.Parameters.Author,
			UserID:        event.Parameters.User,
			Initiator:     event.Parameters.Initiator,
		}
		resourceEvent.SourceType = resourceEvent.DetectSourceType()

		err = p.sendToFifo(ctx, resourceEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) processChildRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeChangeState:
		err := p.updateParentState(ctx, *eventRes.Alarm)
		if err != nil {
			return err
		}
	case types.AlarmChangeTypeResolve:
		err := p.incrementParentEventsCount(ctx, eventRes.Alarm.Value.Parents, -eventRes.Alarm.Value.EventsCount)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}

		err = p.resolveParents(ctx, *eventRes.Alarm, event.Parameters.Timestamp)
		if err != nil {
			return err
		}

		if eventRes.Alarm.Value.State.Value != types.AlarmStateOK {
			err := p.updateParentState(ctx, *eventRes.Alarm)
			if err != nil {
				return err
			}
		}
	}

	if event.EventType == types.EventTypeCheck {
		err := p.incrementParentEventsCount(ctx, eventRes.Alarm.Value.Parents, 1)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}

		err = p.adapter.UpdateLastEventDate(ctx, eventRes.Alarm.Value.Parents, event.Parameters.Timestamp)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) incrementParentEventsCount(ctx context.Context, parentIDs []string, count types.CpsNumber) error {
	_, err := p.alarmCollection.UpdateMany(
		ctx,
		bson.M{
			"d":          bson.M{"$in": parentIDs},
			"v.resolved": nil,
		},
		bson.M{"$inc": bson.M{"v.events_count": count}},
	)

	return err
}

func (p *metaAlarmEventProcessor) sendChildrenEvents(ctx context.Context, childrenIds []string, childEvent types.Event) error {
	var alarms []types.Alarm
	err := p.adapter.GetOpenedAlarmsByIDs(ctx, childrenIds, &alarms)
	if err != nil {
		return fmt.Errorf("cannot fetch children alarms: %w", err)
	}

	for _, alarm := range alarms {
		childEvent.Connector = alarm.Value.Connector
		childEvent.ConnectorName = alarm.Value.ConnectorName
		childEvent.Resource = alarm.Value.Resource
		childEvent.Component = alarm.Value.Component
		childEvent.SourceType = childEvent.DetectSourceType()

		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) resolveParents(ctx context.Context, childAlarm types.Alarm, timestamp types.CpsTime) error {
	ch := make(chan string)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		for _, p := range childAlarm.Value.Parents {
			select {
			case <-ctx.Done():
				return nil
			case ch <- p:
			}
		}

		return nil
	})

	w := int(math.Min(float64(workers), float64(len(childAlarm.Value.Parents))))
	for i := 0; i < w; i++ {
		g.Go(func() error {
			for parentId := range ch {
				var parentAlarm types.AlarmWithEntity
				err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
					alarms := make([]types.AlarmWithEntity, 0)
					err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, []string{parentId}, &alarms)
					if err != nil {
						return fmt.Errorf("cannot fetch parent: %w", err)
					}
					if len(alarms) == 0 {
						return nil
					}
					parentAlarm = alarms[0]

					rule, err := p.ruleAdapter.GetRule(ctx, parentAlarm.Alarm.Value.Meta)
					if err != nil {
						return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
					}
					if rule.ID == "" {
						return fmt.Errorf("meta alarm rule %s not found", parentAlarm.Alarm.Value.Meta)
					}
					if !rule.AutoResolve {
						return nil
					}

					resolvedCount, err := p.adapter.CountResolvedAlarm(ctx, parentAlarm.Alarm.Value.Children)
					if err != nil {
						return fmt.Errorf("cannot fetch alarms: %w", err)
					}

					if resolvedCount < len(parentAlarm.Alarm.Value.Children) {
						return nil
					}

					err = parentAlarm.Alarm.PartialUpdateResolve(timestamp)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					err = p.adapter.PartialUpdateOpen(ctx, &parentAlarm.Alarm)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					return libaxeevent.RemoveMetaAlarmState(ctx, parentAlarm.Alarm, rule, p.metaAlarmStatesService)
				})
				if err != nil {
					return err
				}

				if parentAlarm.Alarm.IsResolved() {
					err = p.adapter.CopyAlarmToResolvedCollection(ctx, parentAlarm.Alarm)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					p.metricsSender.SendResolve(parentAlarm.Alarm, parentAlarm.Entity, timestamp.Time)
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func (p *metaAlarmEventProcessor) updateParentState(ctx context.Context, childAlarm types.Alarm) error {
	ch := make(chan string)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer close(ch)
		for _, p := range childAlarm.Value.Parents {
			select {
			case <-ctx.Done():
				return nil
			case ch <- p:
			}
		}

		return nil
	})

	w := int(math.Min(float64(workers), float64(len(childAlarm.Value.Parents))))
	for i := 0; i < w; i++ {
		g.Go(func() error {
			for parentId := range ch {
				err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
					alarms := make([]types.AlarmWithEntity, 0)
					err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, []string{parentId}, &alarms)
					if err != nil {
						return fmt.Errorf("cannot fetch parent: %w", err)
					}
					if len(alarms) == 0 {
						return nil
					}
					parentAlarm := alarms[0]

					rule, err := p.ruleAdapter.GetRule(ctx, parentAlarm.Alarm.Value.Meta)
					if err != nil {
						return fmt.Errorf("cannot fetch meta alarm rule: %w", err)
					}
					if rule.ID == "" {
						return fmt.Errorf("meta alarm rule %s not found", parentAlarm.Alarm.Value.Meta)
					}

					parentState := parentAlarm.Alarm.Value.State.Value
					childState := childAlarm.Value.State.Value
					if childAlarm.IsResolved() {
						childState = types.AlarmStateOK
					}
					var newState types.CpsNumber

					if childState > parentState {
						newState = childState
					} else if childState < parentState {
						state, lastEventDate, err := p.adapter.GetWorstAlarmStateAndMaxLastEventDate(ctx, parentAlarm.Alarm.Value.Children)
						if err != nil {
							return fmt.Errorf("cannot fetch children state: %w", err)
						}

						newState = types.CpsNumber(state)
						parentAlarm.Alarm.PartialUpdateLastEventDate(types.NewCpsTime(lastEventDate))
					} else {
						return nil
					}

					err = updateMetaAlarmState(&parentAlarm.Alarm, parentAlarm.Entity, childAlarm.Value.LastUpdateDate,
						newState, parentAlarm.Alarm.Value.Output, p.alarmStatusService)
					if err != nil {
						return fmt.Errorf("cannot update parent: %w", err)
					}

					err = p.adapter.PartialUpdateOpen(ctx, &parentAlarm.Alarm)
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
					}

					return nil
				})
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	return g.Wait()
}

func (p *metaAlarmEventProcessor) sendToFifo(ctx context.Context, event types.Event) error {
	body, err := p.encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	err = p.amqpPublisher.PublishWithContext(
		ctx,
		canopsis.FIFOExchangeName,
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot send child event: %w", err)
	}

	return nil
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByEntityIds(ctx context.Context, entityIDs []string) ([]types.AlarmWithEntity, error) {
	return p.getAlarmsWithEntityByMatch(ctx, bson.M{
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": nil,
	})
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByParentIdAndEntityIds(ctx context.Context, parentId string, entityIDs []string) ([]types.AlarmWithEntity, error) {
	return p.getAlarmsWithEntityByMatch(ctx, bson.M{
		"v.parents":  parentId,
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": nil,
	})
}

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByMatch(ctx context.Context, match bson.M) ([]types.AlarmWithEntity, error) {
	var alarms []types.AlarmWithEntity

	cursor, err := p.alarmCollection.Aggregate(ctx, []bson.M{
		{"$match": match},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         mongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$sort": bson.M{
			"alarm.v.last_update_date": 1,
		}},
	})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &alarms)
	if err != nil {
		return nil, err
	}

	return alarms, err
}

func (p *metaAlarmEventProcessor) executeOutputTpl(data correlation.EventExtraInfosMeta) (string, error) {
	rule := data.Rule
	if rule.OutputTemplate == "" {
		return "", nil
	}

	res, err := p.templateExecutor.Execute(rule.OutputTemplate, data)
	if err != nil {
		return "", fmt.Errorf("unable to execute output template for metaalarm rule %s: %+v", rule.ID, err)
	}

	return res, nil
}

func (p *metaAlarmEventProcessor) newAlarm(
	params rpc.AxeParameters,
	entity types.Entity,
	alarmConfig config.AlarmConfig,
) types.Alarm {
	now := types.NewCpsTime()
	alarm := types.Alarm{
		EntityID: entity.ID,
		ID:       utils.NewID(),
		Time:     now,
		Value: types.AlarmValue{
			CreationDate:      now,
			DisplayName:       types.GenDisplayName(alarmConfig.DisplayNameScheme),
			InitialOutput:     params.Output,
			Output:            params.Output,
			InitialLongOutput: params.LongOutput,
			LongOutput:        params.LongOutput,
			LongOutputHistory: []string{params.LongOutput},
			LastUpdateDate:    params.Timestamp,
			LastEventDate:     now,
			Parents:           []string{},
			Children:          []string{},
			UnlinkedParents:   []string{},
			Infos:             map[string]map[string]interface{}{},
			RuleVersion:       map[string]string{},
		},
	}

	switch entity.Type {
	case types.EntityTypeResource:
		alarm.Value.Resource = entity.Name
		alarm.Value.Component = entity.Component
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.Connector, "/")
	case types.EntityTypeComponent, types.EntityTypeService:
		alarm.Value.Component = entity.Name
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.Connector, "/")
	case types.EntityTypeConnector:
		alarm.Value.Connector, alarm.Value.ConnectorName, _ = strings.Cut(entity.ID, "/")
	}

	return alarm
}

func applyOnChild(changeType types.AlarmChangeType) bool {
	switch changeType {
	case types.AlarmChangeTypeAck,
		types.AlarmChangeTypeAckremove,
		types.AlarmChangeTypeAssocTicket,
		types.AlarmChangeTypeCancel,
		types.AlarmChangeTypeChangeState,
		types.AlarmChangeTypeComment,
		types.AlarmChangeTypeSnooze,
		types.AlarmChangeTypeUncancel,
		types.AlarmChangeTypeUpdateStatus,
		types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		return true
	}

	return false
}

func updateMetaAlarmState(alarm *types.Alarm, entity types.Entity, timestamp types.CpsTime, state types.CpsNumber, output string,
	service alarmstatus.Service) error {
	var currentState, currentStatus types.CpsNumber
	if alarm.Value.State != nil {
		currentState = alarm.Value.State.Value
		currentStatus = alarm.Value.Status.Value
	}

	author := strings.Replace(entity.Connector, "/", ".", 1)
	if state != currentState {
		// Event is an OK, so the alarm should be resolved anyway
		if alarm.IsStateLocked() && state != types.AlarmStateOK {
			return nil
		}

		// Create new Step to keep track of the alarm history
		newStep := types.NewAlarmStep(types.AlarmStepStateIncrease, timestamp, author, output, "", "", types.InitiatorSystem)
		newStep.Value = state

		if state < currentState {
			newStep.Type = types.AlarmStepStateDecrease
		}

		alarm.Value.State = &newStep
		err := alarm.Value.Steps.Add(newStep)
		if err != nil {
			return err
		}

		alarm.Value.TotalStateChanges++
		alarm.Value.LastUpdateDate = timestamp
	}

	newStatus := service.ComputeStatus(*alarm, entity)

	if newStatus == currentStatus {
		if state != currentState {
			alarm.Value.StateChangesSinceStatusUpdate++

			alarm.AddUpdate("$set", bson.M{
				"v.state":                             alarm.Value.State,
				"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
				"v.total_state_changes":               alarm.Value.TotalStateChanges,
				"v.last_update_date":                  alarm.Value.LastUpdateDate,
			})
			alarm.AddUpdate("$push", bson.M{"v.steps": alarm.Value.State})
		}

		return nil
	}

	// Create new Step to keep track of the alarm history
	newStepStatus := types.NewAlarmStep(types.AlarmStepStatusIncrease, timestamp, author, output, "", "", types.InitiatorSystem)
	newStepStatus.Value = newStatus

	if newStatus < currentStatus {
		newStepStatus.Type = types.AlarmStepStatusDecrease
	}

	alarm.Value.Status = &newStepStatus
	err := alarm.Value.Steps.Add(newStepStatus)
	if err != nil {
		return err
	}

	alarm.Value.StateChangesSinceStatusUpdate = 0
	alarm.Value.LastUpdateDate = timestamp

	set := bson.M{
		"v.status":                            alarm.Value.Status,
		"v.state_changes_since_status_update": alarm.Value.StateChangesSinceStatusUpdate,
		"v.last_update_date":                  alarm.Value.LastUpdateDate,
	}
	newSteps := bson.A{}
	if state != currentState {
		set["v.total_state_changes"] = alarm.Value.TotalStateChanges
		set["v.state"] = alarm.Value.State
		newSteps = append(newSteps, alarm.Value.State)
	}
	newSteps = append(newSteps, alarm.Value.Status)

	alarm.AddUpdate("$set", set)
	alarm.AddUpdate("$push", bson.M{"v.steps": bson.M{"$each": newSteps}})

	return nil
}
