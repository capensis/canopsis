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
	query := fmt.Sprintf("INSERT INTO %s (time, type, interval) VALUES ($1, $2, $3)", FIFOEvents)

	batch := &pgx.Batch{}

	inserts := 0
	for _, metric := range metrics {
		batch.Queue(query, metric.Timestamp.UTC(), metric.EventType, metric.Interval)
		inserts++

		if inserts >= canopsis.DefaultBulkSize {
			err := s.pool.SendBatch(ctx, batch)
			if err != nil {
				s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvents)
				break
			}

			inserts = 0
			batch = &pgx.Batch{}
		}
	}

	if inserts > 0 {
		err := s.pool.SendBatch(ctx, batch)
		if err != nil {
			s.logger.Err(err).Msgf("failed to send %s metric: unable to send batch", FIFOEvents)
		}
	}
}
