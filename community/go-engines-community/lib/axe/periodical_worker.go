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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type periodicalWorker struct {
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

	// Resolve the alarms marked as done.
	doneResolved, err := w.AlarmService.ResolveDone(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot resolve done alarms")
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

	for _, alarm := range statusUpdated {
		eventUpdateStatus := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			EventType:     types.EventTypeUpdateStatus,
			Author:        canopsis.DefaultEventAuthor,
			Output:        "",
		}
		eventUpdateStatus.SourceType = eventUpdateStatus.DetectSourceType()
		err = w.publishToEngineFIFO(eventUpdateStatus)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	for _, alarm := range closed {
		eventResolveClosed := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			EventType:     types.EventTypeResolveClose,
		}
		eventResolveClosed.SourceType = eventResolveClosed.DetectSourceType()
		err = w.publishToEngineFIFO(eventResolveClosed)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	for _, alarm := range cancelResolved {
		eventResolveCancel := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			EventType:     types.EventTypeResolveCancel,
		}
		eventResolveCancel.SourceType = eventResolveCancel.DetectSourceType()
		err = w.publishToEngineFIFO(eventResolveCancel)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	for _, alarm := range doneResolved {
		eventResolveDone := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			EventType:     types.EventTypeResolveDone,
		}
		eventResolveDone.SourceType = eventResolveDone.DetectSourceType()
		err = w.publishToEngineFIFO(eventResolveDone)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	for _, alarm := range unsnoozedAlarms {
		eventUnsnooze := types.Event{
			Connector:     alarm.Value.Connector,
			ConnectorName: alarm.Value.ConnectorName,
			Component:     alarm.Value.Component,
			Resource:      alarm.Value.Resource,
			Timestamp:     types.CpsTime{Time: time.Now()},
			EventType:     types.EventTypeUnsnooze,
		}
		eventUnsnooze.SourceType = eventUnsnooze.DetectSourceType()
		err = w.publishToEngineFIFO(eventUnsnooze)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}

	events, err := w.IdleAlarmService.Process(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot process idle rules")
	}
	for _, event := range events {
		err = w.publishToEngineFIFO(event)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}
}

func (w *periodicalWorker) publishToEngineFIFO(event types.Event) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event : %w", err)
	}
	return errt.NewIOError(w.ChannelPub.Publish(
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json", // this type is mandatory to avoid bad conversions into Python.
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	))
}
