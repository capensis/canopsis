package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
)

type axeEventMetricsListener struct {
	metricsSender  metrics.TechSender
	metricsStorage []metrics.AxeEventMetric
	flushInterval  time.Duration
}

func (l *axeEventMetricsListener) Listen(ctx context.Context, metricsChan <-chan metrics.AxeEventMetric) {
	ticker := time.NewTicker(l.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go l.metricsSender.SendAxeEventBatch(ctx, l.metricsStorage)
			l.metricsStorage = make([]metrics.AxeEventMetric, 0, len(l.metricsStorage))
		case metric, ok := <-metricsChan:
			if !ok {
				return
			}

			l.metricsStorage = append(l.metricsStorage, metric)
		}
	}
}
