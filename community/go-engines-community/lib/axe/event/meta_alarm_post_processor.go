package event

//go:generate mockgen -destination=../../../mocks/lib/axe/event/event.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event MetaAlarmPostProcessor

import (
	"context"
	"fmt"
	"math"
	"strings"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
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
	eventGenerator event.Generator,
	amqpPublisher libamqp.Publisher,
	metricsSender metrics.Sender,
	logger zerolog.Logger,
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
		eventGenerator:            eventGenerator,
		amqpPublisher:             amqpPublisher,
		metricsSender:             metricsSender,
		logger:                    logger,
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
	eventGenerator            event.Generator
	amqpPublisher             libamqp.Publisher
	metricsSender             metrics.Sender
	logger                    zerolog.Logger
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

	var childAlarms []types.AlarmWithEntity
	err := p.adapter.GetOpenedAlarmsWithEntityByIDs(ctx, eventRes.Alarm.Value.Children, &childAlarms)
	if err != nil {
		return fmt.Errorf("cannot fetch children alarms: %w", err)
	}

	for _, childAlarm := range childAlarms {
		childEvent, err := p.getChildEventByMetaAlarmEvent(*eventRes.Alarm, childAlarm, event, eventRes)
		if err != nil {
			p.logger.Err(err).Str("entity", childAlarm.Entity.ID).Msg("cannot process child")
			continue
		}

		if childEvent.EventType == "" {
			continue
		}

		err = p.sendToFifo(ctx, childEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *metaAlarmPostProcessor) processComponent(ctx context.Context, event rpc.AxeEvent, eventRes rpc.AxeResultEvent) error {
	if !event.Parameters.TicketResources ||
		eventRes.AlarmChangeType != types.AlarmChangeTypeDeclareTicketWebhook ||
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

	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(componentAlarm.Value.Ticket.Message)
	outputBuilder.WriteString(" ")
	outputBuilder.WriteString(types.OutputComponentPrefix)
	outputBuilder.WriteString(componentAlarm.EntityID)
	outputBuilder.WriteRune('.')
	output := outputBuilder.String()
	for _, resource := range resources {
		if resource.Entity.Type != types.EntityTypeResource {
			continue
		}

		resourceEvent, err := p.eventGenerator.Generate(resource.Entity)
		if err != nil {
			p.logger.Err(err).Str("entity", resource.Entity.ID).Msg("cannot generate event")
			continue
		}

		resourceEvent.EventType = types.EventTypeDeclareTicketWebhook
		resourceEvent.Timestamp = datetime.NewCpsTime()
		resourceEvent.Output = output
		resourceEvent.TicketInfo = componentAlarm.Value.Ticket.TicketInfo
		resourceEvent.Author = event.Parameters.Author
		resourceEvent.UserID = event.Parameters.User
		resourceEvent.Initiator = event.Parameters.Initiator
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

func (p *metaAlarmPostProcessor) resolveParents(ctx context.Context, childAlarm types.Alarm, timestamp datetime.CpsTime) error {
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
					parentAlarm = types.AlarmWithEntity{}
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

					update := p.resolveMetaAlarm(&parentAlarm.Alarm, timestamp)
					_, err = p.alarmCollection.UpdateOne(ctx, bson.M{"_id": parentAlarm.Alarm.ID}, update)
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
					var newLastEventDate datetime.CpsTime
					if childState > parentState {
						newState = childState
					} else if childState < parentState {
						state, lastEventDate, err := p.adapter.GetWorstAlarmStateAndMaxLastEventDate(ctx,
							parentAlarm.Alarm.Value.Children)
						if err != nil {
							return fmt.Errorf("cannot fetch children state: %w", err)
						}

						newState = types.CpsNumber(state)
						newLastEventDate = datetime.NewCpsTime(lastEventDate)
					} else {
						return nil
					}

					setUpdate, pushUpdate, err := updateMetaAlarmState(&parentAlarm.Alarm, parentAlarm.Entity, childAlarm.Value.LastUpdateDate,
						newState, parentAlarm.Alarm.Value.Output, p.alarmStatusService)
					if err != nil {
						return fmt.Errorf("cannot update parent: %w", err)
					}

					if setUpdate == nil {
						setUpdate = bson.M{}
					}

					parentAlarm.Alarm.Value.LastEventDate = newLastEventDate
					setUpdate["v.last_event_date"] = newLastEventDate
					update := bson.M{
						"$set": setUpdate,
					}
					if len(pushUpdate) > 0 {
						update["$push"] = pushUpdate
					}

					_, err = p.alarmCollection.UpdateOne(ctx, bson.M{"_id": parentAlarm.Alarm.ID}, update)
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

func (p *metaAlarmPostProcessor) resolveMetaAlarm(metaAlarm *types.Alarm, timestamp datetime.CpsTime) bson.M {
	metaAlarm.Value.Resolved = &timestamp
	metaAlarm.Value.Duration = int64(timestamp.Sub(metaAlarm.Value.CreationDate.Time).Seconds())
	metaAlarm.Value.CurrentStateDuration = int64(timestamp.Sub(metaAlarm.Value.State.Timestamp.Time).Seconds())
	incUpdate := bson.M{}
	if metaAlarm.Value.Snooze != nil {
		snoozeDuration := int64(timestamp.Sub(metaAlarm.Value.Snooze.Timestamp.Time).Seconds())
		metaAlarm.Value.SnoozeDuration += snoozeDuration
		incUpdate["v.snooze_duration"] = snoozeDuration
	}

	if !metaAlarm.Value.PbehaviorInfo.IsActive() {
		enterTimestamp := datetime.CpsTime{}
		for i := len(metaAlarm.Value.Steps) - 2; i >= 0; i-- {
			if metaAlarm.Value.Steps[i].Type == types.AlarmStepPbhEnter {
				enterTimestamp = metaAlarm.Value.Steps[i].Timestamp
				break
			}
		}

		if !enterTimestamp.IsZero() {
			pbhDuration := int64(timestamp.Sub(enterTimestamp.Time).Seconds())
			metaAlarm.Value.PbehaviorInactiveDuration += pbhDuration
			incUpdate["v.pbh_inactive_duration"] = pbhDuration
		}
	}

	if (metaAlarm.Value.Snooze != nil || !metaAlarm.Value.PbehaviorInfo.IsActive()) && metaAlarm.Value.InactiveStart != nil {
		inactiveDuration := int64(timestamp.Sub(metaAlarm.Value.InactiveStart.Time).Seconds())
		metaAlarm.Value.InactiveDuration += inactiveDuration
		incUpdate["v.inactive_duration"] = inactiveDuration
	}

	metaAlarm.Value.ActiveDuration = metaAlarm.Value.Duration - metaAlarm.Value.InactiveDuration
	newStep := types.NewAlarmStep(types.AlarmStepResolve, timestamp, canopsis.DefaultEventAuthor, "", "", "",
		types.InitiatorSystem, false)
	update := bson.M{
		"$set": bson.M{
			"v.resolved":               metaAlarm.Value.Resolved,
			"v.duration":               metaAlarm.Value.Duration,
			"v.current_state_duration": metaAlarm.Value.CurrentStateDuration,
			"v.active_duration":        metaAlarm.Value.ActiveDuration,
		},
		"$unset": bson.M{
			"not_acked_metric_type":      "",
			"not_acked_metric_send_time": "",
			"not_acked_since":            "",
		},
		"$push": bson.M{
			"v.steps": newStep,
		},
	}
	if len(incUpdate) > 0 {
		update["$inc"] = incUpdate
	}

	return update
}

func (p *metaAlarmPostProcessor) getChildEventByMetaAlarmEvent(
	metaAlarm types.Alarm,
	childAlarm types.AlarmWithEntity,
	event rpc.AxeEvent,
	eventRes rpc.AxeResultEvent,
) (types.Event, error) {
	childEvent, err := p.eventGenerator.Generate(childAlarm.Entity)
	if err != nil {
		return childEvent, err
	}

	childEvent.EventType = event.EventType
	childEvent.Timestamp = datetime.NewCpsTime()
	childEvent.Author = event.Parameters.Author
	childEvent.UserID = event.Parameters.User
	childEvent.Initiator = event.Parameters.Initiator
	childEvent.TicketInfo = event.Parameters.TicketInfo
	childEvent.TicketInfo.TicketMetaAlarmID = eventRes.Alarm.ID
	output := event.Parameters.Output
	isTicket := false
	switch eventRes.AlarmChangeType {
	case types.AlarmChangeTypeDeclareTicketWebhook,
		types.AlarmChangeTypeAutoDeclareTicketWebhook:
		childEvent.EventType = types.EventTypeDeclareTicketWebhook
		output = event.Parameters.TicketInfo.GetStepMessage()
		isTicket = true
	case types.AlarmChangeTypeAssocTicket:
		isTicket = true
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

	childEvent.Output = getMetaAlarmChildEventOutput(metaAlarm, output, event.Parameters.Initiator, isTicket)

	return childEvent, nil
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
