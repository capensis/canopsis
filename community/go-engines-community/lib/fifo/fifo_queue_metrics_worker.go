package fifo

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
)

type fifoQueueMetricsWorker struct {
	techMetricsSender  metrics.TechSender
	periodicalInterval time.Duration
	channel            amqp.Channel
	logger             zerolog.Logger
	consumeQueue       string
}

func (w *fifoQueueMetricsWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *fifoQueueMetricsWorker) Work(ctx context.Context) {
	queue, err := w.channel.QueueInspect(w.consumeQueue)
	if err != nil {
		w.logger.Err(err).Msg("cannot get consume queue length")
		return
	}

	w.techMetricsSender.SendFifoQueue(ctx, time.Now(), int64(queue.Messages))
}
