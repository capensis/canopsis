package techmetrics

//go:generate mockgen -destination=../../../mocks/lib/techmetrics/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics Sender

import (
	"context"
	"errors"
	"math"
	"os"
	"runtime/metrics"
	"slices"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/procfs"
	"github.com/rs/zerolog"
)

const maxUrlLength = 500

type Sender interface {
	Run(ctx context.Context)

	SendSimpleEvent(metricName string, metric EventMetric)

	SendQueue(metricName string, timestamp time.Time, length int64)

	SendCheEntityInfo(timestamp time.Time, name string)
	SendCheEvent(metric CheEventMetric)

	SendAxePeriodical(metric AxePeriodicalMetric)
	SendAxeEvent(metric AxeEventMetric)

	SendPBehaviorPeriodical(metric PbehaviorPeriodicalMetric)

	SendApiRequest(metric ApiRequestMetric)

	SendCorrelationRetries(metric CorrelationRetriesMetric)
}

func NewSender(
	instanceId string,
	configProvider config.TechMetricsConfigProvider,
	interval time.Duration,
	poolRetryCount int,
	poolRetryTimeout time.Duration,
	logger zerolog.Logger,
) Sender {
	return &sender{
		instanceId:       instanceId,
		configProvider:   configProvider,
		interval:         interval,
		poolRetryCount:   poolRetryCount,
		poolRetryTimeout: poolRetryTimeout,
		logger:           logger,

		batches: make(map[string][][]any),
	}
}

type sender struct {
	instanceId     string
	configProvider config.TechMetricsConfigProvider
	interval       time.Duration
	logger         zerolog.Logger

	poolRetryCount   int
	poolRetryTimeout time.Duration

	batchesMx sync.Mutex
	batches   map[string][][]any

	poolMx sync.Mutex
	pool   postgres.Pool
}

func (s *sender) Run(ctx context.Context) {
	defer func() {
		s.closePgPool()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.send(ctx)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		interval := s.configProvider.Get().GoMetricsInterval
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.collectGoMetrics(ctx)

				newInterval := s.configProvider.Get().GoMetricsInterval
				if newInterval != interval {
					ticker.Stop()
					interval = newInterval
					ticker = time.NewTicker(interval)
				}
			}
		}
	}()

	wg.Wait()
}

func (s *sender) SendQueue(metricName string, timestamp time.Time, length int64) {
	s.addBatch(metricName, []any{
		timestamp.UTC(),
		length,
	})
}

func (s *sender) SendAxePeriodical(metric AxePeriodicalMetric) {
	s.addBatch(AxePeriodical, []any{
		metric.Timestamp.UTC(),
		metric.Interval.Microseconds(),
		metric.Events,
	})
}

func (s *sender) SendPBehaviorPeriodical(metric PbehaviorPeriodicalMetric) {
	s.addBatch(PBehaviorPeriodical, []any{
		metric.Timestamp.UTC(),
		metric.Interval.Microseconds(),
		metric.Events,
		metric.Entities,
		metric.Pbehaviors,
	})
}

func (s *sender) SendCheEntityInfo(timestamp time.Time, name string) {
	s.addBatch(CheInfos, []any{
		timestamp.UTC(),
		name,
	})
}

func (s *sender) SendApiRequest(metric ApiRequestMetric) {
	url := metric.Url
	if len(url) > maxUrlLength {
		url = url[:maxUrlLength]
	}
	s.addBatch(ApiRequests, []any{
		metric.Timestamp.UTC(),
		metric.Method,
		url,
		metric.Interval.Microseconds(),
	})
}

func (s *sender) SendSimpleEvent(metricName string, metric EventMetric) {
	s.addBatch(metricName, []any{
		metric.Timestamp.UTC(),
		metric.EventType,
		metric.Interval.Microseconds(),
	})
}

func (s *sender) SendCheEvent(metric CheEventMetric) {
	s.addBatch(CheEvent, []any{
		metric.Timestamp.UTC(),
		metric.EventType,
		metric.Interval.Microseconds(),
		metric.EntityType,
		metric.IsNewEntity,
		metric.IsInfosUpdated,
		metric.IsServicesUpdated,
	})
}

func (s *sender) SendAxeEvent(metric AxeEventMetric) {
	s.addBatch(AxeEvent, []any{
		metric.Timestamp.UTC(),
		metric.EventType,
		metric.Interval.Microseconds(),
		metric.EntityType,
		metric.AlarmChangeType,
	})

}

func (s *sender) SendCorrelationRetries(metric CorrelationRetriesMetric) {
	s.addBatch(CorrelationRetries, []any{
		metric.Timestamp.UTC(),
		metric.Type,
		metric.Retries,
	})
}

func (s *sender) send(ctx context.Context) {
	if !s.configProvider.Get().Enabled {
		s.closePgPool()
		s.cleanBatches()
		return
	}

	batches := s.flushBatches()
	pool := s.getPgPool(ctx)
	if pool == nil {
		return
	}

	bulkSize := canopsis.DefaultBulkSize
	for metricName, rows := range batches {
		columns := s.getColumns(metricName)
		if len(columns) == 0 {
			s.logger.Error().Msgf("unknown columns for %q", metricName)
			continue
		}

		rowsCount := len(rows)
		bulkCount := int(math.Ceil(float64(rowsCount) / float64(bulkSize)))
		for i := 0; i < bulkCount; i++ {
			begin := i * bulkSize
			end := (i + 1) * bulkSize
			if end > rowsCount {
				end = rowsCount
			}
			_, err := pool.CopyFrom(ctx, pgx.Identifier{metricName}, columns, pgx.CopyFromRows(rows[begin:end]))
			if err != nil {
				s.logger.Err(err).Msg("cannot send tech metrics")
				return
			}
		}
	}
}

func (s *sender) closePgPool() {
	s.poolMx.Lock()
	defer s.poolMx.Unlock()
	if s.pool != nil {
		s.pool.Close()
		s.pool = nil
	}
}

func (s *sender) getPgPool(ctx context.Context) postgres.Pool {
	s.poolMx.Lock()
	defer s.poolMx.Unlock()

	if s.pool == nil {
		var err error
		s.pool, err = postgres.NewTechMetricsPool(ctx, s.poolRetryCount, s.poolRetryTimeout)
		if err != nil {
			s.logger.Err(err).Msg("cannot connect tech metrics Postgres")
		}
	}

	return s.pool
}

func (s *sender) flushBatches() map[string][][]any {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	batches := make(map[string][][]any, len(s.batches))
	for metricName, items := range s.batches {
		if len(items) == 0 {
			continue
		}
		batches[metricName] = make([][]any, 0, len(items))
	}
	res := s.batches
	s.batches = batches

	return res
}

func (s *sender) cleanBatches() {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	if len(s.batches) > 0 {
		s.batches = make(map[string][][]any, 0)
	}
}

func (s *sender) addBatch(metricName string, args []any) {
	if !s.configProvider.Get().Enabled {
		return
	}

	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	if _, ok := s.batches[metricName]; !ok {
		s.batches[metricName] = make([][]any, 0, 1)
	}

	s.batches[metricName] = append(s.batches[metricName], args)
}

func (s *sender) getColumns(metricName string) []string {
	switch metricName {
	case CheEvent:
		return []string{
			"time",
			"type",
			"interval",
			"entity_type",
			"is_new_entity",
			"is_infos_updated",
			"is_services_updated",
		}
	case AxeEvent:
		return []string{
			"time",
			"type",
			"interval",
			"entity_type",
			"alarm_change_type",
		}
	case CanopsisEvent,
		FIFOEvent,
		CorrelationEvent,
		ServiceEvent,
		DynamicInfosEvent,
		ActionEvent:
		return []string{
			"time",
			"type",
			"interval",
		}
	case AxePeriodical:
		return []string{
			"time",
			"interval",
			"events",
		}
	case PBehaviorPeriodical:
		return []string{
			"time",
			"interval",
			"events",
			"entities",
			"pbehaviors",
		}
	case FIFOQueue:
		return []string{
			"time",
			"length",
		}
	case CheInfos:
		return []string{
			"time",
			"name",
		}
	case ApiRequests:
		return []string{
			"time",
			"method",
			"url",
			"interval",
		}
	case CorrelationRetries:
		return []string{
			"time",
			"type",
			"retries",
		}
	}

	return nil
}

func (s *sender) collectGoMetrics(ctx context.Context) {
	conf := s.configProvider.Get()
	if !conf.Enabled {
		s.closePgPool()
		return
	}

	if len(conf.GoMetrics) == 0 {
		return
	}

	pool := s.getPgPool(ctx)
	if pool == nil {
		return
	}

	samples := make([]metrics.Sample, 0, len(conf.GoMetrics))
	descs := metrics.All()
	for _, desc := range descs {
		if slices.Contains(conf.GoMetrics, desc.Name) {
			samples = append(samples, metrics.Sample{
				Name: desc.Name,
			})
		}
	}

	rows := make([][]any, len(samples))
	now := time.Now().UTC()
	if len(samples) > 0 {
		metrics.Read(samples)
		for i, sample := range samples {
			name, value := sample.Name, sample.Value
			switch value.Kind() {
			case metrics.KindUint64:
				rows[i] = []any{now, name, value.Uint64(), s.instanceId}
			case metrics.KindFloat64:
				rows[i] = []any{now, name, value.Float64(), s.instanceId}
			case metrics.KindFloat64Histogram:
				v, err := s.medianBucket(value.Float64Histogram())
				if err != nil {
					s.logger.Err(err).Msg("cannot update go metrics")
					return
				}

				rows[i] = []any{now, name, v, s.instanceId}
			default:
				s.logger.Error().Msgf("unexpected %q metric kind: %v\n", name, value.Kind())
				return
			}
		}
	}

	pid := os.Getpid()
	proc, err := procfs.NewProc(pid)
	if err != nil {
		s.logger.Err(err).Msg("cannot update get proc metrics")
		return
	}

	procStat, err := proc.Stat()
	if err != nil {
		s.logger.Err(err).Msg("cannot update get proc metrics")
		return
	}

	if slices.Contains(conf.GoMetrics, "proc_cpu_total") {
		rows = append(rows, []any{now, "proc_cpu_total", procStat.CPUTime(), s.instanceId})
	}

	if slices.Contains(conf.GoMetrics, "proc_virtual_memory") {
		rows = append(rows, []any{now, "proc_virtual_memory", procStat.VirtualMemory(), s.instanceId})
	}

	if slices.Contains(conf.GoMetrics, "proc_resident_memory") {
		rows = append(rows, []any{now, "proc_resident_memory", procStat.ResidentMemory(), s.instanceId})
	}

	if slices.Contains(conf.GoMetrics, "proc_file_descriptors") {
		fileDescriptors, err := proc.FileDescriptorsLen()
		if err != nil {
			s.logger.Err(err).Msg("cannot update get proc metrics")
			return
		}

		rows = append(rows, []any{now, "proc_file_descriptors", fileDescriptors, s.instanceId})
	}

	if len(rows) == 0 {
		return
	}

	_, err = pool.CopyFrom(ctx, pgx.Identifier{GoMetrics}, []string{"time", "metric", "value", "source"}, pgx.CopyFromRows(rows))
	if err != nil {
		s.logger.Err(err).Msg("cannot update go metrics")
	}
}

func (s *sender) medianBucket(h *metrics.Float64Histogram) (float64, error) {
	total := uint64(0)
	for _, count := range h.Counts {
		total += count
	}
	thresh := total / 2
	total = 0
	for i, count := range h.Counts {
		total += count
		if total >= thresh {
			return h.Buckets[i], nil
		}
	}

	return 0, errors.New("unexpected value")
}
