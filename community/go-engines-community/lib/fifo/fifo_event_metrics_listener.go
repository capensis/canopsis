package fifo

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
)

type fifoEventMetricsListener struct {
	metricsSender  metrics.TechSender
	metricsStorage []metrics.FifoEventMetric
	flushInterval  time.Duration
}

func (l *fifoEventMetricsListener) Listen(ctx context.Context, metricsChan <-chan metrics.FifoEventMetric) {
	ticker := time.NewTicker(l.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go l.metricsSender.SendFifoEventBatch(ctx, l.metricsStorage)
			l.metricsStorage = make([]metrics.FifoEventMetric, 0, len(l.metricsStorage))
		case metric, ok := <-metricsChan:
			if !ok {
				return
			}

			l.metricsStorage = append(l.metricsStorage, metric)
		}
	}
}
