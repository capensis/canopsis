package axe

import (
	"context"
	"fmt"
	"math"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10

func NewMetaAlarmEventProcessor(
	dbClient mongo.DbClient,
	adapter alarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmStatusService alarmstatus.Service,
	alarmConfigProvider config.AlarmConfigProvider,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	metricsSender metrics.Sender,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	templateExecutor template.Executor,
	logger zerolog.Logger,
) alarm.MetaAlarmEventProcessor {
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
	adapter                   alarm.Adapter
	ruleAdapter               correlation.RulesAdapter

	alarmStatusService  alarmstatus.Service
	alarmConfigProvider config.AlarmConfigProvider

	encoder       encoding.Encoder
	amqpPublisher libamqp.Publisher

	metricsSender metrics.Sender

	logger zerolog.Logger

	templateExecutor template.Executor
}

func (p *metaAlarmEventProcessor) CreateMetaAlarm(ctx context.Context, event types.Event) (*types.Alarm, []types.Alarm, error) {
	rule, err := p.ruleAdapter.GetRule(ctx, event.MetaAlarmRuleID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.MetaAlarmRuleID, err)
	} else if rule.ID == "" {
		return nil, nil, fmt.Errorf("meta alarm rule id=%q not found", event.MetaAlarmRuleID)
	}

	metaAlarm := newAlarm(event, p.alarmConfigProvider.Get())
	metaAlarm.Value.Meta = event.MetaAlarmRuleID
	metaAlarm.Value.MetaValuePath = event.MetaAlarmValuePath

	if event.DisplayName != "" {
		metaAlarm.Value.DisplayName = event.DisplayName
	}

	stateID := rule.ID
	if event.MetaAlarmValuePath != "" {
		stateID = fmt.Sprintf("%s&&%s", rule.ID, event.MetaAlarmValuePath)
	}

	var childEntityIDs []string
	var archived bool

	if rule.IsManual() {
		childEntityIDs = event.MetaAlarmChildren
	} else {
		metaAlarmState, err := p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID)
		if err != nil {
			return nil, nil, err
		}

		if metaAlarmState.MetaAlarmName != event.Resource {
			// try to get archived state
			metaAlarmState, err = p.metaAlarmStatesService.GetMetaAlarmState(ctx, stateID+"-"+event.Resource)
			if err != nil {
				return nil, nil, err
			}

			if metaAlarmState.ID == "" {
				return nil, nil, fmt.Errorf("meta alarm state for rule id=%q and meta alarm name=%q not found", event.MetaAlarmRuleID, event.Resource)
			}

			archived = true
		}

		childEntityIDs = metaAlarmState.ChildrenEntityIDs
	}

	var lastChild types.AlarmWithEntity
	worstState := types.CpsNumber(types.AlarmStateMinor)

	updatedChildrenAlarms := make([]types.Alarm, 0)
	if len(childEntityIDs) > 0 {
		childAlarms, err := p.getAlarmsWithEntityByEntityIDs(ctx, childEntityIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot fetch children alarms: %w", err)
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
					return nil, nil, err
				}
				metaAlarm.AddChild(childAlarm.Alarm.EntityID)
				updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
			}
		}
	}

	err = UpdateAlarmState(&metaAlarm, *event.Entity, event.Timestamp, worstState, event.Output, p.alarmStatusService)
	if err != nil {
		return nil, nil, err
	}

	output := ""
	if rule.IsManual() {
		output = event.Output
	} else {
		output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
			Rule:     rule,
			Count:    int64(len(updatedChildrenAlarms)),
			Children: lastChild,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	metaAlarm.UpdateOutput(output)

	err = p.adapter.Insert(ctx, metaAlarm)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create alarm: %w", err)
	}

	err = p.adapter.PartialMassUpdateOpen(ctx, updatedChildrenAlarms)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot update children alarms: %w", err)
	}

	if !rule.IsManual() && !archived {
		ok, err := p.metaAlarmStatesService.SwitchStateToCreated(ctx, stateID)
		if err != nil || !ok {
			return nil, nil, err
		}
	}

	return &metaAlarm, updatedChildrenAlarms, nil
}

func (p *metaAlarmEventProcessor) AttachChildrenToMetaAlarm(ctx context.Context, event types.Event) (*types.Alarm, []types.Alarm, []types.Event, error) {
	if len(event.MetaAlarmChildren) == 0 || event.Alarm == nil {
		return event.Alarm, nil, nil, nil
	}

	rule, err := p.ruleAdapter.GetRule(ctx, event.MetaAlarmRuleID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.MetaAlarmRuleID, err)
	} else if rule.ID == "" {
		return nil, nil, nil, fmt.Errorf("meta alarm rule id=%q not found", event.MetaAlarmRuleID)
	}

	alarms, err := p.getAlarmsWithEntityByEntityIDs(ctx, event.MetaAlarmChildren)
	if err != nil {
		return nil, nil, nil, err
	}

	newStep := types.NewMetaAlarmAttachStep(*event.Alarm, rule.Name)
	worstState := types.CpsNumber(types.AlarmStateOK)

	var lastChild types.AlarmWithEntity
	updatedChildrenAlarms := make([]types.Alarm, 0)

	for _, childAlarm := range alarms {
		if !childAlarm.Alarm.AddParent(event.Alarm.EntityID) {
			continue
		}

		event.Alarm.AddChild(childAlarm.Entity.ID)
		err = childAlarm.Alarm.PartialUpdateAddStepWithStep(newStep)
		if err != nil {
			return nil, nil, nil, err
		}

		if childAlarm.Alarm.Value.State.Value > worstState {
			worstState = childAlarm.Alarm.Value.State.Value
		}

		updatedChildrenAlarms = append(updatedChildrenAlarms, childAlarm.Alarm)
		lastChild = childAlarm
	}

	if len(updatedChildrenAlarms) == 0 {
		return nil, nil, nil, nil
	}

	if event.Alarm.Value.Meta == "" {
		event.Alarm.SetMeta(event.MetaAlarmRuleID)
		event.Alarm.SetMetaValuePath(event.MetaAlarmValuePath)
	}

	if worstState > event.Alarm.CurrentState() {
		err = UpdateAlarmState(event.Alarm, *event.Entity, event.Timestamp, worstState, event.Alarm.Value.Output, p.alarmStatusService)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	if event.Alarm.Value.LastEventDate.Before(event.Timestamp) {
		event.Alarm.PartialUpdateLastEventDate(event.Timestamp)
	}

	childrenCount, err := p.adapter.GetCountOpenedAlarmsByIDs(ctx, event.Alarm.Value.Children)
	if err != nil {
		return nil, nil, nil, err
	}

	output := ""
	if rule.Type == correlation.RuleTypeManualGroup {
		output = event.Output
	} else {
		output, err = p.executeOutputTpl(correlation.EventExtraInfosMeta{
			Rule:     rule,
			Count:    childrenCount,
			Children: lastChild,
		})
		if err != nil {
			return nil, nil, nil, err
		}
	}

	event.Alarm.UpdateOutput(output)

	err = p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{*event.Alarm}, updatedChildrenAlarms...))
	if err != nil {
		return nil, nil, nil, err
	}

	childrenEvents, err := p.applyActionsOnChildren(*event.Alarm, updatedChildrenAlarms)
	if err != nil {
		return nil, nil, nil, err
	}

	return event.Alarm, updatedChildrenAlarms, childrenEvents, nil
}

func (p *metaAlarmEventProcessor) DetachChildrenFromMetaAlarm(ctx context.Context, event types.Event) (*types.Alarm, error) {
	if len(event.MetaAlarmChildren) == 0 || event.Alarm == nil {
		return event.Alarm, nil
	}

	rule, err := p.ruleAdapter.GetRule(ctx, event.MetaAlarmRuleID)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", event.MetaAlarmRuleID, err)
	} else if rule.ID == "" {
		return nil, fmt.Errorf("meta alarm rule id=%q not found", event.MetaAlarmRuleID)
	}

	alarms, err := p.getAlarmsWithEntityByEntityIDs(ctx, event.MetaAlarmChildren)
	if err != nil {
		return nil, err
	}

	updatedChildAlarms := make([]types.Alarm, 0)
	for _, childAlarm := range alarms {
		if childAlarm.Alarm.RemoveParent(event.Entity.ID) {
			event.Alarm.RemoveChild(childAlarm.Entity.ID)
			updatedChildAlarms = append(updatedChildAlarms, childAlarm.Alarm)
		}
	}

	if len(updatedChildAlarms) == 0 {
		return nil, nil
	}

	metaAlarmChildren, err := p.getAlarmsWithEntityByEntityIDs(ctx, event.Alarm.Value.Children)
	if err != nil {
		return nil, err
	}

	worstState := types.CpsNumber(types.AlarmStateOK)

	for _, childAlarm := range metaAlarmChildren {
		if childAlarm.Alarm.Value.State.Value > worstState {
			worstState = childAlarm.Alarm.Value.State.Value
		}
	}

	err = UpdateAlarmState(event.Alarm, *event.Entity, event.Timestamp, worstState, event.Alarm.Value.Output, p.alarmStatusService)
	if err != nil {
		return nil, err
	}

	infos := correlation.EventExtraInfosMeta{
		Rule:  rule,
		Count: int64(len(metaAlarmChildren)),
	}
	if len(metaAlarmChildren) > 0 {
		infos.Children = metaAlarmChildren[len(metaAlarmChildren)-1]
	}

	output := ""
	if rule.IsManual() {
		output = event.Output
	} else {
		output, err = p.executeOutputTpl(infos)
		if err != nil {
			return nil, err
		}
	}

	event.Alarm.UpdateOutput(output)

	err = p.adapter.PartialMassUpdateOpen(ctx, append([]types.Alarm{*event.Alarm}, updatedChildAlarms...))
	if err != nil || len(updatedChildAlarms) == 0 {
		return nil, err
	}

	return event.Alarm, nil
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
			Initiator:     types.InitiatorSystem,
			Timestamp:     types.NewCpsTime(),
		}
		childEvent.SourceType = childEvent.DetectSourceType()

		for _, step := range steps {
			childEvent.Output = step.Message
			childEvent.Author = step.Author
			childEvent.UserID = step.UserID
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

func (p *metaAlarmEventProcessor) Process(ctx context.Context, event types.Event) error {
	if event.Alarm == nil || event.AlarmChange == nil {
		return nil
	}

	alarm := event.Alarm

	if alarm.IsMetaAlarm() {
		return p.processParent(ctx, event)
	}

	if alarm.IsMetaChildren() {
		return p.processChild(ctx, event)
	}

	return nil
}

func (p *metaAlarmEventProcessor) ProcessAxeRpc(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if event.Alarm == nil || eventRes.AlarmChangeType == "" {
		return nil
	}

	alarm := event.Alarm

	err := p.processComponentRpc(ctx, event, eventRes)
	if err != nil {
		return err
	}

	if alarm.IsMetaAlarm() {
		return p.processParentRpc(ctx, event, eventRes)
	}

	if alarm.IsMetaChildren() {
		return p.processChildRpc(ctx, eventRes)
	}

	return nil
}

func (p *metaAlarmEventProcessor) processParent(ctx context.Context, event types.Event) error {
	if !applyOnChild(event.AlarmChange.Type) {
		return nil
	}

	childEvent := types.Event{
		EventType:  event.EventType,
		Timestamp:  event.Timestamp,
		Output:     event.Output,
		LongOutput: event.LongOutput,
		Author:     event.Author,
		UserID:     event.UserID,
		Debug:      event.Debug,
		Role:       event.Role,
		Initiator:  types.InitiatorSystem,
		Status:     event.Status,
		State:      event.State,
		TicketInfo: event.TicketInfo,
		Duration:   event.Duration,
	}

	childEvent.TicketInfo.TicketMetaAlarmID = event.Alarm.ID

	err := p.sendChildrenEvents(ctx, event.Alarm.Value.Children, childEvent)
	if err != nil {
		return err
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
		Initiator:  types.InitiatorSystem,
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

	err := p.sendChildrenEvents(ctx, event.Alarm.Value.Children, childEvent)
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

	resources, err := p.adapter.GetAlarmsWithoutTicketByComponent(ctx, event.Alarm.Value.Component)
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
			Initiator:     types.InitiatorSystem,
		}
		resourceEvent.SourceType = resourceEvent.DetectSourceType()

		err = p.sendToFifo(ctx, resourceEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) processChild(ctx context.Context, event types.Event) error {
	switch event.AlarmChange.Type {
	case types.AlarmChangeTypeResolve:
		err := p.resolveParents(ctx, *event.Alarm, event.Timestamp)
		if err != nil {
			return err
		}

		if event.Alarm.Value.State.Value != types.AlarmStateOK {
			err := p.updateParentState(ctx, *event.Alarm)
			if err != nil {
				return err
			}
		}
	case types.AlarmChangeTypeStateIncrease, types.AlarmChangeTypeStateDecrease, types.AlarmChangeTypeChangeState:
		err := p.updateParentState(ctx, *event.Alarm)
		if err != nil {
			return err
		}
	}

	if event.EventType == types.EventTypeCheck && p.alarmConfigProvider.Get().EnableLastEventDate {
		err := p.adapter.UpdateLastEventDate(ctx, event.Alarm.Value.Parents, event.Timestamp)
		if err != nil {
			return fmt.Errorf("cannot update parent alarms: %w", err)
		}
	}

	return nil
}

func (p *metaAlarmEventProcessor) processChildRpc(ctx context.Context, eventRes rpc.AxeResultEvent) error {
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeChangeState:
		err := p.updateParentState(ctx, *eventRes.Alarm)
		if err != nil {
			return err
		}
	}

	return nil
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

					err = parentAlarm.Alarm.PartialUpdateResolve(types.NewCpsTime())
					if err != nil {
						return fmt.Errorf("cannot update alarm: %w", err)
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
						r, err := p.adapter.GetWorstAlarmState(ctx, parentAlarm.Alarm.Value.Children)
						if err != nil {
							return fmt.Errorf("cannot fetch children state: %w", err)
						}

						newState = types.CpsNumber(r)
					} else {
						return nil
					}

					err = UpdateAlarmState(&parentAlarm.Alarm, parentAlarm.Entity, childAlarm.Value.LastUpdateDate,
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

func (p *metaAlarmEventProcessor) getAlarmsWithEntityByEntityIDs(ctx context.Context, entityIDs []string) ([]types.AlarmWithEntity, error) {
	var alarms []types.AlarmWithEntity

	cursor, err := p.alarmCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{
				"d":          bson.M{"$in": entityIDs},
				"v.resolved": nil,
			},
		},
		{
			"$project": bson.M{
				"alarm": "$$ROOT",
				"_id":   0,
			},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.EntityMongoCollection,
				"localField":   "alarm.d",
				"foreignField": "_id",
				"as":           "entity",
			},
		},
		{
			"$unwind": "$entity",
		},
		{
			"$sort": bson.M{
				"alarm.v.last_update_date": 1,
			},
		},
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
