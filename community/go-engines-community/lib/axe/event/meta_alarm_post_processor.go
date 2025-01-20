package event

//go:generate mockgen -destination=../../../mocks/lib/axe/event/event.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event MetaAlarmPostProcessor

import (
	"context"
	"fmt"
	"math"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/sync/errgroup"
)

const workers = 10

type MetaAlarmPostProcessor interface {
	// Process handles related meta alarm parents and children after alarm change.
	Process(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error
}

func NewMetaAlarmPostProcessor(
	dbClient mongo.DbClient,
	adapter libalarm.Adapter,
	ruleAdapter correlation.RulesAdapter,
	alarmStatusService alarmstatus.Service,
	metaAlarmStatesService correlation.MetaAlarmStateService,
	encoder encoding.Encoder,
	amqpPublisher libamqp.Publisher,
	metricsSender metrics.Sender,
) MetaAlarmPostProcessor {
	return &metaAlarmPostProcessor{
		dbClient:                  dbClient,
		alarmCollection:           dbClient.Collection(mongo.AlarmMongoCollection),
		metaAlarmStatesCollection: dbClient.Collection(mongo.MetaAlarmStatesCollection),
		entityCollection:          dbClient.Collection(mongo.EntityMongoCollection),
		adapter:                   adapter,
		ruleAdapter:               ruleAdapter,
		alarmStatusService:        alarmStatusService,
		metaAlarmStatesService:    metaAlarmStatesService,
		encoder:                   encoder,
		amqpPublisher:             amqpPublisher,
		metricsSender:             metricsSender,
	}
}

type metaAlarmPostProcessor struct {
	dbClient                  mongo.DbClient
	alarmCollection           mongo.DbCollection
	metaAlarmStatesCollection mongo.DbCollection
	entityCollection          mongo.DbCollection
	adapter                   libalarm.Adapter
	ruleAdapter               correlation.RulesAdapter
	alarmStatusService        alarmstatus.Service
	metaAlarmStatesService    correlation.MetaAlarmStateService
	encoder                   encoding.Encoder
	amqpPublisher             libamqp.Publisher
	metricsSender             metrics.Sender
}

func (p *metaAlarmPostProcessor) Process(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if eventRes.Alarm == nil {
		return nil
	}

	// todo: processComponent is not in metaalarm context, move it to somewhere else.
	err := p.processComponent(ctx, event, eventRes)
	if err != nil {
		return err
	}

	if eventRes.Alarm.IsMetaAlarm() {
		return p.processParent(ctx, event, eventRes)
	}

	if eventRes.Alarm.IsMetaChild() {
		return p.processChild(ctx, event, eventRes)
	}

	return nil
}

func (p *metaAlarmPostProcessor) processParent(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !p.applyOnChild(eventRes.AlarmChangeType) {
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
		seconds, err := event.Parameters.Duration.To(types.DurationUnitSecond)
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

func (p *metaAlarmPostProcessor) processComponent(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
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
			SourceType:    types.SourceTypeResource,
		}

		err = p.sendToFifo(ctx, resourceEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmPostProcessor) processChild(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
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

func (p *metaAlarmPostProcessor) incrementParentEventsCount(ctx context.Context, parentIDs []string, count types.CpsNumber) error {
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

func (p *metaAlarmPostProcessor) sendChildrenEvents(ctx context.Context, childrenIds []string, childEvent types.Event) error {
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

func (p *metaAlarmPostProcessor) resolveParents(ctx context.Context, childAlarm types.Alarm, timestamp types.CpsTime) error {
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

					return removeMetaAlarmState(ctx, parentAlarm.Alarm, rule, p.metaAlarmStatesService)
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

func (p *metaAlarmPostProcessor) updateParentState(ctx context.Context, childAlarm types.Alarm) error {
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

func (p *metaAlarmPostProcessor) sendToFifo(ctx context.Context, event types.Event) error {
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

func (p *metaAlarmPostProcessor) applyOnChild(changeType types.AlarmChangeType) bool {
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
