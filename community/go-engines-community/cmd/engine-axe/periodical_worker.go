package main

import (
	"context"
	"fmt"
	"runtime/trace"
	"time"

	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	libalarm "git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlealarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/errt"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

const PeriodicalLockKey = "axe-periodical-lock-key"

type periodicalWorker struct {
	PeriodicalInterval            time.Duration
	LockerClient                  redis.LockClient
	ChannelPub                    libamqp.Channel
	AlarmService                  libalarm.Service
	AlarmBaggotTime               time.Duration
	AlarmCancelAutosolveDelay     time.Duration
	Encoder                       encoding.Encoder
	IdleAlarmService              idlealarm.Service
	DisableActionSnoozeDelayOnPbh bool
	Logger                        zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work() error {
	_, err := w.LockerClient.Obtain(PeriodicalLockKey, w.GetInterval(), nil)
	if err == redislock.ErrNotObtained {
		w.Logger.Debug().Msg("Could not obtain lock! Skip periodical process")

		return nil
	} else if err != nil {
		w.Logger.Error().Err(err).Msg("Obtain redis lock: unexpected error")

		return nil
	}

	ctx, task := trace.NewTask(context.Background(), "axe.PeriodicalProcess")
	defer task.End()

	// Resolve the alarms whose state is info.
	w.Logger.Debug().Msg("Closing alarms")
	closed, err := w.AlarmService.ResolveAlarms(ctx, w.AlarmBaggotTime)
	if err != nil {
		w.Logger.Error().Err(err).Msg("cannot resolve ok alarms")
		return nil
	}

	// Process the snoozed alarms.
	// Note that this may unsnooze some alarms, but it will not resolve any.
	// This is the reason why the snoozedResolved alarms are not added to the
	// resolvedAlarms slice.
	w.Logger.Debug().Msg("Resolve snooze")
	unsnoozedAlarms, err := w.AlarmService.ResolveSnoozes(ctx, w.DisableActionSnoozeDelayOnPbh)
	if err != nil {
		w.Logger.Error().Err(err).Msg("cannot unsnooze alarms")
		return nil
	}

	// Resolve the alarms marked as canceled.
	w.Logger.Debug().Msg("Resolve cancel")
	cancelResolved, err := w.AlarmService.ResolveCancels(ctx, w.AlarmCancelAutosolveDelay)
	if err != nil {
		w.Logger.Error().Err(err).Msg("cannot resolve canceled alarms")
		return nil
	}

	// Resolve the alarms marked as done.
	w.Logger.Debug().Msg("Resolve done")
	doneResolved, err := w.AlarmService.ResolveDone(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("cannot resolve done alarms")
		return nil
	}

	// Process the flapping alarms.
	// Note that this may change the status of some alarms, but it will not
	// resolve any.
	// This is the reason why the statusUpdated alarms are not added to the
	// resolvedAlarms slice.
	w.Logger.Debug().Msg("Update flapping alarms")
	statusUpdated, err := w.AlarmService.UpdateFlappingAlarms(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("cannot update flapping alarms")
		return nil
	}

	w.Logger.Info().
		Int("closed", len(closed)).
		Int("unsnoozed", len(unsnoozedAlarms)).
		Int("cancel_resolved", len(cancelResolved)).
		Int("done_resolved", len(doneResolved)).
		Int("flapping_updated", len(statusUpdated)).
		Msg("updated alarms")

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

		if alarm.Value.Resource == "" {
			eventUpdateStatus.SourceType = types.SourceTypeComponent
		} else {
			eventUpdateStatus.SourceType = types.SourceTypeResource
		}

		err = w.publishToEngineFIFO(eventUpdateStatus)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish update_status event to FIFO")
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

		if alarm.Value.Resource == "" {
			eventResolveClosed.SourceType = types.SourceTypeComponent
		} else {
			eventResolveClosed.SourceType = types.SourceTypeResource
		}

		err = w.publishToEngineFIFO(eventResolveClosed)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish resolve_close event to FIFO")
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

		if alarm.Value.Resource == "" {
			eventResolveCancel.SourceType = types.SourceTypeComponent
		} else {
			eventResolveCancel.SourceType = types.SourceTypeResource
		}

		err = w.publishToEngineFIFO(eventResolveCancel)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish resolve_cancel event to FIFO")
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

		if alarm.Value.Resource == "" {
			eventResolveDone.SourceType = types.SourceTypeComponent
		} else {
			eventResolveDone.SourceType = types.SourceTypeResource
		}

		err = w.publishToEngineFIFO(eventResolveDone)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish resolve_done event to FIFO")
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

		if alarm.Value.Resource == "" {
			eventUnsnooze.SourceType = types.SourceTypeComponent
		} else {
			eventUnsnooze.SourceType = types.SourceTypeResource
		}

		err = w.publishToEngineFIFO(eventUnsnooze)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish unsnooze event to FIFO")
		}
	}

	events := w.IdleAlarmService.Process()
	for _, event := range events {
		err = w.publishToEngineFIFO(event)
		if err != nil {
			w.Logger.Error().Err(err).Msg("Failed publish idle event to FIFO")
		}
	}

	return nil
}

func (w *periodicalWorker) publishToEngineFIFO(event types.Event) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("publishEvent(): error while encoding event %+v", err)
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
