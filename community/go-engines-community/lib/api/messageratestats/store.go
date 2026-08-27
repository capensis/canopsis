package messageratestats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v5"
)

type Store interface {
	Find(context.Context, ListRequest) ([]StatsResponse, error)
	// GetDeletedBeforeForHours gets the lower bound time value for hourly request.
	GetDeletedBeforeForHours(ctx context.Context) (*datetime.CpsTime, error)
}

type store struct {
	pgPoolProvider postgres.PoolProvider
}

// NewStore creates new store.
func NewStore(pgPoolProvider postgres.PoolProvider) Store {
	return &store{
		pgPoolProvider: pgPoolProvider,
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) ([]StatsResponse, error) {
	switch r.Interval {
	case IntervalMinute:
		return s.findMinuteStats(ctx, r)
	case IntervalHour:
		return s.findHourStats(ctx, r)
	default:
		return nil, fmt.Errorf("unknown interval %v", r.Interval)
	}
}

func (s *store) findMinuteStats(ctx context.Context, r ListRequest) ([]StatsResponse, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	search, args := s.getSearchQuery(r)

	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('1 minute', time), count(*) FROM "+metrics.MessageRate+
		search+" GROUP BY time_bucket_gapfill('1 minute', time)", args)
	if err != nil {
		return nil, fmt.Errorf("failed to find minute stats: %w", err)
	}

	defer rows.Close()

	rates := make([]StatsResponse, 0)
	for rows.Next() {
		var t time.Time
		var rateColumn *int64
		err := rows.Scan(&t, &rateColumn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan minute stats: %w", err)
		}

		var rate int64
		if rateColumn != nil {
			rate = *rateColumn
		}

		rates = append(rates, StatsResponse{
			ID:   t.Unix(),
			Rate: rate,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch minute stats: %w", err)
	}

	return rates, nil
}

func (s *store) findHourStats(ctx context.Context, r ListRequest) ([]StatsResponse, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	search, args := s.getSearchQuery(r)
	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('1 hour', time), sum(count) FROM "+metrics.MessageRateHourly+
		search+" GROUP BY time_bucket_gapfill('1 hour', time)", args)
	if err != nil {
		return nil, fmt.Errorf("failed to find hour stats: %w", err)
	}

	defer rows.Close()
	rates := make([]StatsResponse, 0)
	for rows.Next() {
		var t time.Time
		var rateColumn *int64
		err := rows.Scan(&t, &rateColumn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hour stats: %w", err)
		}

		var rate int64
		if rateColumn != nil {
			rate = *rateColumn
		}

		rates = append(rates, StatsResponse{
			ID:   t.Unix(),
			Rate: rate,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch hour stats: %w", err)
	}

	return rates, nil
}

func (s *store) GetDeletedBeforeForHours(ctx context.Context) (*datetime.CpsTime, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	var t *time.Time
	err = pgPool.QueryRow(ctx, "SELECT min(time) FROM "+metrics.MessageRateHourly).Scan(&t)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find min stats date: %w", err)
	}

	if t == nil {
		return nil, nil
	}

	return &datetime.CpsTime{Time: *t}, nil
}

func (s *store) getSearchQuery(r ListRequest) (string, pgx.NamedArgs) {
	var start, end time.Time

	if r.From.IsZero() || r.To.IsZero() {
		nowTrunc := time.Now().Truncate(time.Minute).UTC()

		// add one minute to include current minute to response
		end = nowTrunc.Add(time.Minute)
		start = end.Add(-time.Hour)
	} else {
		start = r.From.UTC()
		end = r.To.UTC()
	}

	search := " WHERE time >= @start AND time <= @end "

	if len(r.EventTypes) > 0 {
		search += "AND event_type = ANY(@event_types) "
	}

	if len(r.ConnectorNames) > 0 {
		search += " AND connector_name = ANY(@connector_names) "
	}

	return search, pgx.NamedArgs{
		"start":           start,
		"end":             end,
		"event_types":     r.EventTypes,
		"connector_names": r.ConnectorNames,
	}
}
