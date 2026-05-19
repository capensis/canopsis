package engine

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics"
	"github.com/rs/zerolog"
)

func NewQueueMetricsPeriodicalWorker(
	periodicalInterval time.Duration,
	channel amqp.Channel,
	techMetricsSender techmetrics.Sender,
	consumeQueues []string,
	techMetric string,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &queueMetricsPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		channel:            channel,
		techMetricsSender:  techMetricsSender,
		consumeQueues:      consumeQueues,
		techMetric:         techMetric,
		logger:             logger,
	}
}

type queueMetricsPeriodicalWorker struct {
	periodicalInterval time.Duration
	channel            amqp.Channel
	techMetricsSender  techmetrics.Sender
	consumeQueues      []string
	techMetric         string
	logger             zerolog.Logger
}

func (w *queueMetricsPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *queueMetricsPeriodicalWorker) Work(ctx context.Context) {
	queueLength := 0

	for i := range w.consumeQueues {
		queue, err := w.channel.QueueDeclarePassive(ctx, w.consumeQueues[i], false, false, false, false, nil)
		if err != nil {
			w.logger.Err(err).Msg("cannot get consume queue length")
		}

		queueLength += queue.Messages
	}

	if queueLength > 0 {
		w.techMetricsSender.SendQueue(w.techMetric, time.Now(), int64(queueLength))
	}
}
