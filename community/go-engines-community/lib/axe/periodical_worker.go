package axe

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type periodicalWorker struct {
	TechMetricsSender   techmetrics.Sender
	PeriodicalInterval  time.Duration
	ChannelPub          libamqp.Channel
	AlarmService        libalarm.Service
	AlarmAdapter        libalarm.Adapter
	Encoder             encoding.Encoder
	IdleAlarmService    idlealarm.Service
	AlarmConfigProvider config.AlarmConfigProvider
	Logger              zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(parentCtx context.Context) {
	metric := techmetrics.AxePeriodicalMetric{}
	metric.Timestamp = time.Now()
	eventsCount := 0
	defer func() {
		metric.Interval = time.Since(metric.Timestamp)
		metric.Events = int64(eventsCount)
		w.TechMetricsSender.SendAxePeriodical(metric)
	}()

	ctx, task := trace.NewTask(parentCtx, "axe.PeriodicalProcess")
	defer task.End()

	alarmConfig := w.AlarmConfigProvider.Get()
	if alarmConfig.TimeToKeepResolvedAlarms > 0 {
		err := w.AlarmAdapter.DeleteResolvedAlarms(ctx, alarmConfig.TimeToKeepResolvedAlarms)
		if err != nil {
			w.Logger.Err(err).Msg("cannot delete resolved alarms")
			return
		}
	}

	// Resolve the alarms whose state is info.
	closed, err := w.AlarmService.ResolveClosed(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot resolve ok alarms")
		return
	}

	// Process the snoozed alarms.
	// Note that this may unsnooze some alarms, but it will not resolve any.
	// This is the reason why the snoozedResolved alarms are not added to the
	// resolvedAlarms slice.
	unsnoozedAlarms, err := w.AlarmService.ResolveSnoozes(ctx, alarmConfig)
	if err != nil {
		w.Logger.Err(err).Msg("cannot unsnooze alarms")
		return
	}

	// Resolve the alarms marked as canceled.
	cancelResolved, err := w.AlarmService.ResolveCancels(ctx, alarmConfig)
	if err != nil {
		w.Logger.Err(err).Msg("cannot resolve canceled alarms")
		return
	}

	// Process the flapping alarms.
	// Note that this may change the status of some alarms, but it will not
	// resolve any.
	// This is the reason why the statusUpdated alarms are not added to the
	// resolvedAlarms slice.
	statusUpdated, err := w.AlarmService.UpdateFlappingAlarms(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update flapping alarms")
		return
	}

	eventsCount += len(statusUpdated)
	for _, alarm := range statusUpdated {
		eventUpdateStatus := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     datetime.NewCpsTime(),
			EventType:     types.EventTypeUpdateStatus,
			Author:        canopsis.DefaultEventAuthor,
			Output:        "",
			Initiator:     types.InitiatorSystem,
		}
		eventUpdateStatus.SourceType = eventUpdateStatus.DetectSourceType()
		err = w.publishToEngineFIFO(ctx, eventUpdateStatus)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	eventsCount += len(closed)
	for _, alarm := range closed {
		eventResolveClosed := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     datetime.NewCpsTime(),
			EventType:     types.EventTypeResolveClose,
			Initiator:     types.InitiatorSystem,
		}
		eventResolveClosed.SourceType = eventResolveClosed.DetectSourceType()
		err = w.publishToEngineFIFO(ctx, eventResolveClosed)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	eventsCount += len(cancelResolved)
	for _, alarm := range cancelResolved {
		eventResolveCancel := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     datetime.NewCpsTime(),
			EventType:     types.EventTypeResolveCancel,
			Initiator:     types.InitiatorSystem,
		}
		eventResolveCancel.SourceType = eventResolveCancel.DetectSourceType()
		err = w.publishToEngineFIFO(ctx, eventResolveCancel)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	eventsCount += len(unsnoozedAlarms)
	for _, alarm := range unsnoozedAlarms {
		eventUnsnooze := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     datetime.NewCpsTime(),
			EventType:     types.EventTypeUnsnooze,
			Initiator:     types.InitiatorSystem,
		}
		eventUnsnooze.SourceType = eventUnsnooze.DetectSourceType()
		err = w.publishToEngineFIFO(ctx, eventUnsnooze)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	events, err := w.IdleAlarmService.Process(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot process idle rules")
	}
	eventsCount += len(events)
	for _, event := range events {
		err = w.publishToEngineFIFO(ctx, event)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}
}

func (w *periodicalWorker) publishToEngineFIFO(ctx context.Context, event types.Event) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event : %w", err)
	}
	return w.ChannelPub.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	)
}
