package metrics

import (
	"context"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/rs/zerolog"
)

type TechSender interface {
	SendFifoRate(ctx context.Context, timestamp time.Time, length int64)
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
	fmt.Printf("events %d\n", length)

	query := fmt.Sprintf("INSERT INTO %s (time, length) VALUES($1, $2);", FIFOQueue)
	_, err := s.pool.Exec(ctx, query, timestamp.UTC(), length)
	if err != nil {
		s.logger.Err(err).Msgf("failed to send %s metric: unable to execute insert", FIFOQueue)
	}
}
