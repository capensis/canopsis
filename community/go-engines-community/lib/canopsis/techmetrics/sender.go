package techmetrics

//go:generate mockgen -destination=../../../mocks/lib/techmetrics/sender.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/techmetrics Sender

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

const maxUrlLength = 255

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
}

func NewSender(
	configProvider config.TechMetricsConfigProvider,
	interval time.Duration,
	poolRetryCount int,
	poolRetryTimeout time.Duration,
	logger zerolog.Logger,
) Sender {
	return &sender{
		configProvider:   configProvider,
		interval:         interval,
		poolRetryCount:   poolRetryCount,
		poolRetryTimeout: poolRetryTimeout,
		logger:           logger,

		batches: make(map[string][]batchItem),
	}
}

type sender struct {
	configProvider config.TechMetricsConfigProvider
	interval       time.Duration
	logger         zerolog.Logger

	poolRetryCount   int
	poolRetryTimeout time.Duration

	batchesMx sync.Mutex
	batches   map[string][]batchItem

	pool postgres.Pool
}

type batchItem struct {
	query     string
	arguments []interface{}
}

func (s *sender) Run(ctx context.Context) {
	defer func() {
		if s.pool != nil {
			s.pool.Close()
		}
	}()

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
}

func (s *sender) SendQueue(metricName string, timestamp time.Time, length int64) {
	query := fmt.Sprintf("INSERT INTO %s (time, length) VALUES($1, $2);", metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			timestamp.UTC(),
			length,
		},
	})
}

func (s *sender) SendAxePeriodical(metric AxePeriodicalMetric) {
	metricName := AxePeriodical
	query := fmt.Sprintf("INSERT INTO %s (time, interval, events) VALUES($1, $2, $3);", metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.Interval.Microseconds(),
			metric.Events,
		},
	})
}

func (s *sender) SendPBehaviorPeriodical(metric PbehaviorPeriodicalMetric) {
	metricName := PBehaviorPeriodical
	query := fmt.Sprintf("INSERT INTO %s (time, interval, events, entities, pbehaviors) VALUES($1, $2, $3, $4, $5);", metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.Interval.Microseconds(),
			metric.Events,
			metric.Entities,
			metric.Pbehaviors,
		},
	})
}

func (s *sender) SendCheEntityInfo(timestamp time.Time, name string) {
	metricName := CheInfos
	query := fmt.Sprintf("INSERT INTO %s (time, name) VALUES($1, $2);", metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			timestamp.UTC(),
			name,
		},
	})
}

func (s *sender) SendApiRequest(metric ApiRequestMetric) {
	metricName := ApiRequests
	query := fmt.Sprintf("INSERT INTO %s (time, method, url, interval) VALUES($1, $2, $3, $4);", metricName)
	url := metric.Url
	if len(url) > maxUrlLength {
		url = url[:maxUrlLength]
	}

	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.Method,
			url,
			metric.Interval.Microseconds(),
		},
	})
}

func (s *sender) SendSimpleEvent(metricName string, metric EventMetric) {
	query := fmt.Sprintf("INSERT INTO %s (time, type, interval) VALUES ($1, $2, $3)", metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.EventType,
			metric.Interval.Microseconds(),
		},
	})
}

func (s *sender) SendCheEvent(metric CheEventMetric) {
	metricName := CheEvent
	query := fmt.Sprintf(`
		INSERT INTO %s (time, type, interval, entity_type, is_new_entity, is_infos_updated, is_services_updated) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.EventType,
			metric.Interval.Microseconds(),
			metric.EntityType,
			metric.IsNewEntity,
			metric.IsInfosUpdated,
			metric.IsServicesUpdated,
		},
	})
}

func (s *sender) SendAxeEvent(metric AxeEventMetric) {
	metricName := AxeEvent
	query := fmt.Sprintf(`
		INSERT INTO %s (time, type, interval, entity_type, alarm_change_type) 
		VALUES ($1, $2, $3, $4, $5)`,
		metricName)
	s.addBatch(metricName, batchItem{
		query: query,
		arguments: []interface{}{
			metric.Timestamp.UTC(),
			metric.EventType,
			metric.Interval.Microseconds(),
			metric.EntityType,
			metric.AlarmChangeType,
		},
	})

}

func (s *sender) send(ctx context.Context) {
	if !s.configProvider.Get().Enabled {
		if s.pool != nil {
			s.pool.Close()
			s.pool = nil
		}

		s.cleanBatches()
		return
	}

	batches := s.flushBatches()

	if s.pool == nil {
		var err error
		s.pool, err = postgres.NewTechMetricsPool(ctx, s.poolRetryCount, s.poolRetryTimeout)
		if err != nil {
			s.logger.Err(err).Msg("cannot connect tech metrics Postgres")
			return
		}
	}

	batch := &pgx.Batch{}
	count := 0
	for _, items := range batches {
		for _, item := range items {
			batch.Queue(item.query, item.arguments...)
			count++

			if count >= canopsis.DefaultBulkSize {
				err := s.pool.SendBatch(ctx, batch)
				if err != nil {
					s.logger.Err(err).Msg("cannot send tech metrics")
					return
				}
				batch = &pgx.Batch{}
				count = 0
			}
		}
	}

	if count > 0 {
		err := s.pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msg("cannot send tech metrics")
		}
	}
}

func (s *sender) flushBatches() map[string][]batchItem {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	batches := make(map[string][]batchItem, len(s.batches))
	for metricName, items := range s.batches {
		if len(items) == 0 {
			continue
		}
		batches[metricName] = make([]batchItem, 0, len(items))
	}
	res := s.batches
	s.batches = batches

	return res
}

func (s *sender) cleanBatches() {
	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	if len(s.batches) > 0 {
		s.batches = make(map[string][]batchItem, 0)
	}
}

func (s *sender) addBatch(metricName string, item batchItem) {
	if !s.configProvider.Get().Enabled {
		return
	}

	s.batchesMx.Lock()
	defer s.batchesMx.Unlock()

	if _, ok := s.batches[metricName]; !ok {
		s.batches[metricName] = make([]batchItem, 0, 1)
	}

	s.batches[metricName] = append(s.batches[metricName], item)
}
