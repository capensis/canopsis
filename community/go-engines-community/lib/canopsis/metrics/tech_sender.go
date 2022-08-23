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
	SendFifoRate(ctx context.Context, timestamp time.Time, length int64)
	SendFifoEventBatch(ctx context.Context, metrics []FifoEventMetric)
	SendCheEventBatch(ctx context.Context, metrics []CheEventMetric)
}

type techSender struct {
	pool   postgres.Pool
	logger zerolog.Logger
}

func NewTechMetricsSender(
	pool postgres.Pool,
	logger zerolog.Logger,
) TechSender {
	return &techSender{
		pool:   pool,
		logger: logger,
	}
}

func (s *techSender) SendFifoRate(ctx context.Context, timestamp time.Time, length int64) {
	query := fmt.Sprintf("INSERT INTO %s (time, length) VALUES($1, $2);", FIFOQueue)
	_, err := s.pool.Exec(ctx, query, timestamp.UTC(), length)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", FIFOQueue)
	}
}

func (s *techSender) SendFifoEventBatch(ctx context.Context, metrics []FifoEventMetric) {
	query := fmt.Sprintf("INSERT INTO %s (time, type, interval) VALUES ($1, $2, $3)", FIFOEvent)

	batch := &pgx.Batch{}

	inserts := 0
	for _, metric := range metrics {
		batch.Queue(query, metric.Timestamp.UTC(), metric.EventType, metric.Interval)
		inserts++

		if inserts >= canopsis.DefaultBulkSize {
			err := s.pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvent)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := s.pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvent)
		}
	}
}

func (s *techSender) SendCheEventBatch(ctx context.Context, metrics []CheEventMetric) {
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
			err := s.pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", CheEvent)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := s.pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", CheEvent)
		}
	}
}
