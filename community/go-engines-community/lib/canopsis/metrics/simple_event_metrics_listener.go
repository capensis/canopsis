package metrics

import (
	"context"
	"time"
)

type SimpleEventMetricsListener struct {
	metricsSender  TechSender
	metricsStorage []SimpleEventMetric
	flushInterval  time.Duration
}

func NewSimpleEventMetricsListener(sender TechSender, flushInterval time.Duration) *SimpleEventMetricsListener {
	return &SimpleEventMetricsListener{
		metricsSender: sender,
		flushInterval: flushInterval,
	}
}

func (l *SimpleEventMetricsListener) Listen(ctx context.Context, metricsChan <-chan SimpleEventMetric, metric string) {
	ticker := time.NewTicker(l.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go l.metricsSender.SendSimpleEventBatch(ctx, l.metricsStorage, metric)
			l.metricsStorage = make([]SimpleEventMetric, 0, len(l.metricsStorage))
		case metric, ok := <-metricsChan:
			if !ok {
				return
			}

			l.metricsStorage = append(l.metricsStorage, metric)
		}
	}
}
