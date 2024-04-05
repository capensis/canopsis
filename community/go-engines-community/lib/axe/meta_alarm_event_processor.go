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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
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
	fifoExchange, fifoQueue string,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
) alarm.MetaAlarmEventProcessor {
	return &metaAlarmEventProcessor{
		dbClient:            dbClient,
		adapter:             adapter,
		ruleAdapter:         ruleAdapter,
		alarmStatusService:  alarmStatusService,
		alarmConfigProvider: alarmConfigProvider,
		encoder:             encoder,
		amqpPublisher:       amqpPublisher,
		fifoExchange:        fifoExchange,
		fifoQueue:           fifoQueue,
		metricsSender:       metricsSender,
		logger:              logger,
	}
}

type metaAlarmEventProcessor struct {
	dbClient    mongo.DbClient
	adapter     alarm.Adapter
	ruleAdapter correlation.RulesAdapter

	alarmStatusService  alarmstatus.Service
	alarmConfigProvider config.AlarmConfigProvider

	encoder                 encoding.Encoder
	amqpPublisher           libamqp.Publisher
	fifoExchange, fifoQueue string

	metricsSender metrics.Sender

	logger zerolog.Logger
}

func (p *metaAlarmEventProcessor) CreateMetaAlarm(ctx context.Context, event types.Event) (*types.Alarm, error) {
	var updatedChildAlarms []types.Alarm
	var metaAlarm types.Alarm
	now := types.NewCpsTime()

	err := p.dbClient.WithTransaction(ctx, func(ctx context.Context) error {
		updatedChildAlarms = make([]types.Alarm, 0)
		metaAlarm = types.Alarm{}

		ruleIdentifier := event.MetaAlarmRuleID
		rule, err := p.ruleAdapter.GetRule(ctx, ruleIdentifier)
		if err != nil {
			return fmt.Errorf("cannot fetch meta alarm rule id=%q: %w", ruleIdentifier, err)
		} else if rule.ID == "" {
			return fmt.Errorf("meta alarm rule id=%q not found", ruleIdentifier)
		} else {
			ruleIdentifier = rule.Name
		}

		metaAlarm = newAlarm(event, p.alarmConfigProvider.Get(), now)
		metaAlarm.Value.Meta = event.MetaAlarmRuleID
		metaAlarm.Value.MetaValuePath = event.MetaAlarmValuePath

		if event.DisplayName != "" {
			metaAlarm.Value.DisplayName = event.DisplayName
		}

		var childAlarms []types.Alarm
		worstState := types.CpsNumber(types.AlarmStateMinor)

		if len(event.MetaAlarmChildren) > 0 {
			err := p.adapter.GetOpenedAlarmsByIDs(ctx, event.MetaAlarmChildren, &childAlarms)
			if err != nil {
				return fmt.Errorf("cannot fetch children alarms: %w", err)
			}

			newStep := types.NewMetaAlarmAttachStep(metaAlarm, ruleIdentifier)
			for i := range childAlarms {
				if childAlarms[i].Value.State != nil {
					childState := childAlarms[i].Value.State.Value
					if childState > worstState {
						worstState = childState
					}
				}

				if !childAlarms[i].HasParentByEID(metaAlarm.EntityID) {
					childAlarms[i].AddParent(metaAlarm.EntityID)
					err = childAlarms[i].PartialUpdateAddStepWithStep(newStep)
					if err != nil {
						return err
					}
					metaAlarm.AddChild(childAlarms[i].EntityID)
					updatedChildAlarms = append(updatedChildAlarms, childAlarms[i])
				}
			}
		}

		err = UpdateAlarmState(&metaAlarm, *event.Entity, event.Timestamp, worstState, event.Output, p.alarmStatusService)
		if err != nil {
			return err
		}

		err = p.adapter.Insert(ctx, metaAlarm)
		if err != nil {
			return fmt.Errorf("cannot create alarm: %w", err)
		}

		err = p.adapter.PartialMassUpdateOpen(ctx, childAlarms)
		if err != nil {
			return fmt.Errorf("cannot update children alarms: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, child := range updatedChildAlarms {
		p.metricsSender.SendCorrelation(event.Timestamp.Time, child)
	}

	return &metaAlarm, nil
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

func (p *metaAlarmEventProcessor) ProcessAckResources(ctx context.Context, event types.Event) error {
	if !event.AckResources || event.Alarm == nil || event.AlarmChange == nil ||
		event.AlarmChange.Type != types.AlarmChangeTypeAck || event.SourceType != types.SourceTypeComponent {
		return nil
	}

	alarms, err := p.adapter.GetUnacknowledgedAlarmsByComponent(ctx, event.Component)
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if alarm.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent := types.Event{
			EventType:     event.EventType,
			Connector:     alarm.Alarm.Value.Connector,
			ConnectorName: alarm.Alarm.Value.ConnectorName,
			Resource:      alarm.Alarm.Value.Resource,
			Component:     alarm.Alarm.Value.Component,
			Timestamp:     event.Timestamp,
			TicketInfo:    event.TicketInfo,
			Output:        event.Output,
			LongOutput:    event.LongOutput,
			Author:        event.Author,
			UserID:        event.UserID,
			Debug:         event.Debug,
			Role:          event.Role,
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

func (p *metaAlarmEventProcessor) ProcessTicketResources(ctx context.Context, event types.Event) error {
	if !event.TicketResources || event.Alarm == nil || event.AlarmChange == nil ||
		event.AlarmChange.Type != types.AlarmChangeTypeAssocTicket || event.SourceType != types.SourceTypeComponent {
		return nil
	}

	alarms, err := p.adapter.GetAlarmsWithoutTicketByComponent(ctx, event.Component)
	if err != nil {
		return fmt.Errorf("cannot fetch alarms: %w", err)
	}

	for _, alarm := range alarms {
		if alarm.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent := types.Event{
			EventType:     event.EventType,
			Connector:     alarm.Alarm.Value.Connector,
			ConnectorName: alarm.Alarm.Value.ConnectorName,
			Resource:      alarm.Alarm.Value.Resource,
			Component:     alarm.Alarm.Value.Component,
			Timestamp:     event.Timestamp,
			TicketInfo:    event.TicketInfo,
			Output:        event.Output,
			LongOutput:    event.LongOutput,
			Author:        event.Author,
			UserID:        event.UserID,
			Debug:         event.Debug,
			Role:          event.Role,
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

					err = parentAlarm.Alarm.PartialUpdateResolve(timestamp)
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
		p.fifoExchange,
		p.fifoQueue,
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
