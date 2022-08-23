package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
)

type cheEventMetricsListener struct {
	metricsSender  metrics.TechSender
	metricsStorage []metrics.CheEventMetric
	flushInterval  time.Duration
}

func (l *cheEventMetricsListener) Listen(ctx context.Context, metricsChan <-chan metrics.CheEventMetric) {
	ticker := time.NewTicker(l.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go l.metricsSender.SendCheEventBatch(ctx, l.metricsStorage)
			l.metricsStorage = make([]metrics.CheEventMetric, 0, len(l.metricsStorage))
		case metric, ok := <-metricsChan:
			if !ok {
				return
			}

			l.metricsStorage = append(l.metricsStorage, metric)
		}
	}
}
