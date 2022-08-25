package metrics

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"

	"github.com/jackc/pgx/v4"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

type TechSender interface {
	SendFifoQueue(ctx context.Context, timestamp time.Time, length int64)
	SendFifoEventBatch(ctx context.Context, metrics []FifoEventMetric)
	SendCheEventBatch(ctx context.Context, metrics []CheEventMetric)
	SendAxeEventBatch(ctx context.Context, metrics []AxeEventMetric)
	SendAxePeriodical(ctx context.Context, timestamp time.Time, length int64)
	SendPBehaviorPeriodical(ctx context.Context, timestamp time.Time, length int64)
	SendCheEntityInfo(ctx context.Context, timestamp time.Time, name string)
	SendApiRequest(ctx context.Context, timestamp time.Time, interval int64)
}

type techSender struct {
	poolProvider postgres.PoolProvider
	logger       zerolog.Logger
}

func NewTechMetricsSender(
	pool postgres.PoolProvider,
	logger zerolog.Logger,
) TechSender {
	return &techSender{
		poolProvider: pool,
		logger:       logger,
	}
}

func (s *techSender) SendFifoQueue(ctx context.Context, timestamp time.Time, length int64) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, length) VALUES($1, $2);", FIFOQueue)
	_, err := pool.Exec(ctx, query, timestamp.UTC(), length)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", FIFOQueue)
	}
}

func (s *techSender) SendFifoEventBatch(ctx context.Context, metrics []FifoEventMetric) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, type, interval) VALUES ($1, $2, $3)", FIFOEvent)

	batch := &pgx.Batch{}

	inserts := 0
	for _, metric := range metrics {
		batch.Queue(query, metric.Timestamp.UTC(), metric.EventType, metric.Interval)
		inserts++

		if inserts >= canopsis.DefaultBulkSize {
			err := pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvent)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvent)
		}
	}
}

func (s *techSender) SendCheEventBatch(ctx context.Context, metrics []CheEventMetric) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (time, type, interval, entity_type, is_new_entity, is_infos_updated, is_services_updated) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		CheEvent)

	batch := &pgx.Batch{}

	inserts := 0
	for _, metric := range metrics {
		batch.Queue(query, metric.Timestamp.UTC(), metric.EventType, metric.Interval, metric.EntityType, metric.IsNewEntity, metric.IsInfosUpdated, metric.IsServicesUpdated)
		inserts++

		if inserts >= canopsis.DefaultBulkSize {
			err := pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", CheEvent)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", CheEvent)
		}
	}
}

func (s *techSender) SendAxeEventBatch(ctx context.Context, metrics []AxeEventMetric) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (time, type, interval, entity_type, alarm_change_type) 
		VALUES ($1, $2, $3, $4, $5)`,
		AxeEvent)

	batch := &pgx.Batch{}

	inserts := 0
	for _, metric := range metrics {
		batch.Queue(query, metric.Timestamp.UTC(), metric.EventType, metric.Interval, metric.EntityType, metric.AlarmChangeType)
		inserts++

		if inserts >= canopsis.DefaultBulkSize {
			err := pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", AxeEvent)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", AxeEvent)
		}
	}
}

func (s *techSender) SendAxePeriodical(ctx context.Context, timestamp time.Time, length int64) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, interval) VALUES($1, $2);", AxePeriodical)
	_, err := pool.Exec(ctx, query, timestamp.UTC(), length)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", AxePeriodical)
	}
}

func (s *techSender) SendPBehaviorPeriodical(ctx context.Context, timestamp time.Time, length int64) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, interval) VALUES($1, $2);", PBehaviorPeriodical)
	_, err := pool.Exec(ctx, query, timestamp.UTC(), length)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", PBehaviorPeriodical)
	}
}

func (s *techSender) SendCheEntityInfo(ctx context.Context, timestamp time.Time, name string) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, name) VALUES($1, $2);", CheInfos)
	_, err := pool.Exec(ctx, query, timestamp.UTC(), name)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", CheInfos)
	}
}

func (s *techSender) SendApiRequest(ctx context.Context, timestamp time.Time, interval int64) {
	pool := s.poolProvider.GetPool()
	if pool == nil {
		return
	}

	query := fmt.Sprintf("INSERT INTO %s (time, interval) VALUES($1, $2);", ApiRequests)
	_, err := pool.Exec(ctx, query, timestamp.UTC(), interval)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", ApiRequests)
	}
}
