package engine

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"github.com/rs/zerolog"
)

func NewRunInfoMetricsPeriodicalWorker(
	periodicalInterval time.Duration,
	manager RunInfoManager,
	info InstanceRunInfo,
	channel amqp.Channel,
	techMetricsSender techmetrics.Sender,
	techMetric string,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &runInfoMetricsPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		manager:            manager,
		info:               info,
		channel:            channel,
		techMetricsSender:  techMetricsSender,
		techMetric:         techMetric,
		logger:             logger,
	}
}

type runInfoMetricsPeriodicalWorker struct {
	periodicalInterval time.Duration
	manager            RunInfoManager
	info               InstanceRunInfo
	channel            amqp.Channel
	techMetricsSender  techmetrics.Sender
	techMetric         string
	logger             zerolog.Logger
}

func (w *runInfoMetricsPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *runInfoMetricsPeriodicalWorker) Work(ctx context.Context) {
	info := updateInstanceRunInfo(ctx, w.GetInterval(), w.manager, w.info, w.channel, w.logger)

	if info.QueueLength > 0 {
		w.techMetricsSender.SendQueue(w.techMetric, time.Now(), int64(info.QueueLength))
	}
}
