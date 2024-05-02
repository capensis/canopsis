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
	eventCount := 0
	idleEventCount := 0
	defer func() {
		metric.Interval = time.Since(metric.Timestamp)
		metric.Events = int64(eventCount)
		metric.IdleEvents = int64(idleEventCount)
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

	resolveOkEvents, err := w.AlarmService.ResolveClosed(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot resolve ok alarms")

		return
	}

	eventCount += len(resolveOkEvents)
	w.publishEvents(ctx, resolveOkEvents)
	unsnoozeEvents, err := w.AlarmService.ResolveSnoozes(ctx, alarmConfig)
	if err != nil {
		w.Logger.Err(err).Msg("cannot unsnooze alarms")

		return
	}

	eventCount += len(unsnoozeEvents)
	w.publishEvents(ctx, unsnoozeEvents)
	resolveCanceledEvents, err := w.AlarmService.ResolveCancels(ctx, alarmConfig)
	if err != nil {
		w.Logger.Err(err).Msg("cannot resolve canceled alarms")

		return
	}

	eventCount += len(resolveCanceledEvents)
	w.publishEvents(ctx, resolveCanceledEvents)
	statusUpdateEvents, err := w.AlarmService.UpdateFlappingAlarms(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update flapping alarms")

		return
	}

	eventCount += len(statusUpdateEvents)
	w.publishEvents(ctx, statusUpdateEvents)
	idleEvents, err := w.IdleAlarmService.Process(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot process idle rules")

		return
	}

	idleEventCount = len(idleEvents)
	eventCount += idleEventCount
	w.publishEvents(ctx, idleEvents)
}

func (w *periodicalWorker) publishEvents(ctx context.Context, events []types.Event) {
	for _, event := range events {
		err := w.publishToEngineFIFO(ctx, event)
		if err != nil {
			w.Logger.Err(err).Msg("cannot publish event")
		}
	}
}

func (w *periodicalWorker) publishToEngineFIFO(ctx context.Context, event types.Event) error {
	bevent, err := w.Encoder.Encode(event)
	if err != nil {
		return fmt.Errorf("cannot encode event: %w", err)
	}

	return w.ChannelPub.PublishWithContext(
		ctx,
		"",
		canopsis.FIFOQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  canopsis.JsonContentType,
			Body:         bevent,
			DeliveryMode: amqp.Persistent,
		},
	)
}
