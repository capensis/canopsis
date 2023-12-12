package main

import (
	"context"
	"fmt"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	libentity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type periodicalWorker struct {
	TechMetricsSender      techmetrics.Sender
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
	metric := techmetrics.PbehaviorPeriodicalMetric{}
	metric.Timestamp = time.Now()
	eventsCount := 0
	entitiesCount := 0
	pbehaviorsCount := 0
	defer func() {
		metric.Interval = time.Since(metric.Timestamp)
		metric.Events = int64(eventsCount)
		metric.Entities = int64(entitiesCount)
		metric.Pbehaviors = int64(pbehaviorsCount)
		w.TechMetricsSender.SendPBehaviorPeriodical(metric)
	}()

	now := time.Now().In(w.TimezoneConfigProvider.Get().Location)
	newSpan := timespan.New(now, now.Add(w.FrameDuration))

	resolver, recomputedCount, err := w.PbhService.Compute(ctx, newSpan)
	if err != nil {
		w.Logger.Err(err).Msg("compute pbehavior's frames failed")
		return
	}

	if recomputedCount >= 0 {
		w.Logger.Info().
			Time("interval_from", newSpan.From()).
			Time("interval_to", newSpan.To()).
			Int("count", recomputedCount).
			Msg("pbehaviors are recomputed")
	}

	computedEntityIDs, err := resolver.GetComputedEntityIDs()
	entitiesCount = len(computedEntityIDs)
	if err != nil {
		w.Logger.Err(err).Msg("cannot get entities which have pbehavior")
		return
	}

	var processedEntityIds []string
	processedEntityIds, eventsCount = w.processAlarms(ctx, now, computedEntityIDs, resolver)
	eventsCount += w.processEntities(ctx, now, computedEntityIDs, processedEntityIds, resolver)

	pbehaviorsCount, err = resolver.GetPbehaviorsCount(ctx, now)
	if err != nil {
		w.Logger.Err(err).Msg("cannot get pbehaviors count")
	}
}

func (w *periodicalWorker) processAlarms(
	ctx context.Context, computedAt time.Time,
	computedEntityIDs []string,
	resolver pbehavior.ComputedEntityTypeResolver,
) ([]string, int) {
	eventsCount := 0
	cursor, err := w.AlarmAdapter.FindToCheckPbehaviorInfo(ctx, datetime.CpsTime{Time: computedAt}, computedEntityIDs)
	if err != nil {
		w.Logger.Err(err).Msg("get alarms from mongo failed")
		return nil, eventsCount
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
			return processedEntityIds, eventsCount
		}

		event := w.EventManager.GetEvent(resolveResult, alarm, now)
		if event.EventType != "" {
			eventsCount++
			ech <- PublishEventMsg{
				event:   event,
				id:      alarm.ID,
				pbhID:   resolveResult.ResolvedPbhID,
				pbhType: resolveResult.ResolvedType,
			}
		}
	}

	return processedEntityIds, eventsCount
}

func (w *periodicalWorker) processEntities(
	ctx context.Context,
	computedAt time.Time,
	computedEntityIDs,
	processedEntityIds []string,
	resolver pbehavior.ComputedEntityTypeResolver,
) int {
	eventsCount := 0
	cursor, err := w.EntityAdapter.FindToCheckPbehaviorInfo(ctx, computedEntityIDs, processedEntityIds)
	if err != nil {
		w.Logger.Err(err).Msg("get alarms from mongo failed")
		return eventsCount
	}

	defer cursor.Close(ctx)

	eventGenerator := libevent.NewGenerator("engine", "pbehavior")

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
			return eventsCount
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
			return eventsCount
		}

		if lastAlarm == nil {
			event, err = eventGenerator.Generate(entity)
			if err != nil {
				w.Logger.Err(err).Msg("cannot generate event")
				return eventsCount
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
		event.Timestamp = datetime.CpsTime{Time: now}
		event.PbehaviorInfo = pbehavior.NewPBehaviorInfo(event.Timestamp, resolveResult)

		eventsCount++
		ech <- PublishEventMsg{
			event:   event,
			id:      entity.ID,
			pbhID:   resolveResult.ResolvedPbhID,
			pbhType: resolveResult.ResolvedType,
		}
	}

	return eventsCount
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
				Str("resolve_pbehavior", ms.pbhID).
				Str("resolve_type", fmt.Sprintf("%+v", ms.pbhType)).
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
		return fmt.Errorf("publishTo(): error while encoding event %w", err)
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
