package main

import (
	"context"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type periodicalWorker struct {
	ChannelPub             libamqp.Channel
	PeriodicalInterval     time.Duration
	PbhService             pbehavior.Service
	AlarmAdapter           libalarm.Adapter
	EntityAdapter          libentity.Adapter
	EventManager           pbehavior.EventManager
	FrameDuration          time.Duration
	TimezoneConfigProvider config.TimezoneConfigProvider
	Encoder                encoding.Encoder
	Logger                 zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) {
	now := time.Now().In(w.TimezoneConfigProvider.Get().Location)
	newSpan := timespan.New(now, now.Add(w.FrameDuration))

	resolver, count, err := w.PbhService.Compute(ctx, newSpan)
	if err != nil {
		w.Logger.Err(err).Msg("compute pbehavior's frames failed")
		return
	}

	if count >= 0 {
		w.Logger.Info().
			Time("interval from", newSpan.From()).
			Time("interval to", newSpan.To()).
			Int("count", count).
			Msg("pbehaviors are recomputed")
	}

	computedEntityIDs, err := resolver.GetComputedEntityIDs()
	if err != nil {
		w.Logger.Err(err).Msg("cannot get entities which have pbehavior")
		return
	}

	processedEntityIds := w.processAlarms(ctx, now, computedEntityIDs, resolver)
	w.processEntities(ctx, now, computedEntityIDs, processedEntityIds, resolver)
}

func (w *periodicalWorker) processAlarms(
	ctx context.Context, computedAt time.Time,
	computedEntityIDs []string,
	resolver pbehavior.ComputedEntityTypeResolver,
) []string {
	cursor, err := w.AlarmAdapter.FindToCheckPbehaviorInfo(ctx, types.CpsTime{Time: computedAt}, computedEntityIDs)
	if err != nil {
		w.Logger.Err(err).Msg("get alarms from mongo failed")
		return nil
	}

	defer cursor.Close(ctx)

	ech := make(chan PublishEventMsg, 1)
	defer close(ech)
	go w.publishToFifoChan("alarm", ech)

	processedEntityIds := make([]string, 0)
	for cursor.Next(ctx) {
		var alarmWithEntity types.AlarmWithEntity

		err = cursor.Decode(&alarmWithEntity)
		if err != nil {
			w.Logger.Err(err).Msg("decode alarm with entity failed")
			continue
		}

		alarm := alarmWithEntity.Alarm
		entity := alarmWithEntity.Entity
		processedEntityIds = append(processedEntityIds, alarm.EntityID)

		if len(alarm.Value.Steps) > 0 {
			lastStep := alarm.Value.Steps[len(alarm.Value.Steps)-1]
			if lastStep.Timestamp.Unix() >= computedAt.Unix() {
				continue
			}
		}

		now := time.Now()
		resolveResult, err := resolver.Resolve(ctx, entity, now)
		if err != nil {
			w.Logger.Err(err).Str("entity_id", entity.ID).Msg("resolve an entity failed")
			return processedEntityIds
		}

		event := w.EventManager.GetEvent(resolveResult, alarm, now)
		if event.EventType != "" {
			ech <- PublishEventMsg{
				event:   event,
				id:      alarm.ID,
				pbhID:   resolveResult.ResolvedPbhID,
				pbhType: resolveResult.ResolvedType,
			}
		}
	}

	return processedEntityIds
}

func (w *periodicalWorker) processEntities(
	ctx context.Context,
	computedAt time.Time,
	computedEntityIDs,
	processedEntityIds []string,
	resolver pbehavior.ComputedEntityTypeResolver,
) {
	cursor, err := w.EntityAdapter.FindToCheckPbehaviorInfo(ctx, computedEntityIDs, processedEntityIds)
	if err != nil {
		w.Logger.Err(err).Msg("get alarms from mongo failed")
		return
	}

	defer cursor.Close(ctx)

	eventGenerator := libevent.NewGenerator(w.EntityAdapter)

	ech := make(chan PublishEventMsg, 1)
	defer close(ech)

	go w.publishToFifoChan("entity", ech)

	for cursor.Next(ctx) {
		var entity types.Entity

		err = cursor.Decode(&entity)
		if err != nil {
			w.Logger.Err(err).Msg("decode alarm with entity failed")
			continue
		}

		if entity.PbehaviorInfo.Timestamp != nil && entity.PbehaviorInfo.Timestamp.Unix() >= computedAt.Unix() {
			continue
		}

		now := time.Now()
		resolveResult, err := resolver.Resolve(ctx, entity, now)
		if err != nil {
			w.Logger.Err(err).Str("entity_id", entity.ID).Msg("resolve an entity failed")
			return
		}

		eventType, output := w.EventManager.GetEventType(resolveResult, entity.PbehaviorInfo)
		if eventType == "" {
			continue
		}

		event := types.Event{
			Initiator: types.InitiatorSystem,
		}
		lastAlarm, err := w.AlarmAdapter.GetLastAlarmByEntityID(ctx, entity.ID)
		if err != nil {
			w.Logger.Err(err).Msg("cannot fetch last alarm")
			return
		}

		if lastAlarm == nil {
			event, err = eventGenerator.Generate(ctx, entity)
			if err != nil {
				w.Logger.Err(err).Msg("cannot generate event")
				return
			}
		} else {
			event.Connector = lastAlarm.Value.Connector
			event.ConnectorName = lastAlarm.Value.ConnectorName
			event.Component = lastAlarm.Value.Component
			event.Resource = lastAlarm.Value.Resource
			event.SourceType = event.DetectSourceType()
		}

		event.EventType = eventType
		event.Output = output
		event.Timestamp = types.CpsTime{Time: now}
		event.PbehaviorInfo = pbehavior.NewPBehaviorInfo(event.Timestamp, resolveResult)

		ech <- PublishEventMsg{
			event:   event,
			id:      entity.ID,
			pbhID:   resolveResult.ResolvedPbhID,
			pbhType: resolveResult.ResolvedType,
		}
	}
}

type PublishEventMsg struct {
	event     types.Event
	pbhType   pbehavior.Type
	id, pbhID string
}

func (w *periodicalWorker) publishToFifoChan(idTitle string, msgs <-chan PublishEventMsg) {
	for ms := range msgs {
		err := w.publishToEngineFIFO(context.Background(), ms.event)
		if err != nil {
			w.Logger.Err(err).Str(idTitle, ms.id).Msgf("failed to send %s event", ms.event.EventType)
		} else {
			w.Logger.Debug().
				Str("resolve pbehavior", ms.pbhID).
				Str("resolve type", fmt.Sprintf("%+v", ms.pbhType)).
				Str(idTitle, ms.id).
				Msgf("send %s event", ms.event.EventType)
		}
	}
}

func (w *periodicalWorker) publishToEngineFIFO(ctx context.Context, event types.Event) error {
	return w.publishTo(ctx, event, canopsis.FIFOQueueName)
}

func (w *periodicalWorker) publishTo(ctx context.Context, event types.Event, queue string) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("publishTo(): error while encoding event %+v", err)
	}

	return errt.NewIOError(w.ChannelPub.PublishWithContext(
		ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	))
}
