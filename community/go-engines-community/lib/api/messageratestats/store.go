package messageratestats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/postgres"
	"github.com/jackc/pgx/v5"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(context.Context, ListRequest) ([]StatsResponse, error)
	// GetDeletedBeforeForHours gets the lower bound time value for hourly request.
	GetDeletedBeforeForHours(ctx context.Context) (*datetime.CpsTime, error)
}

type store struct {
	db             mongo.DbClient
	pgPoolProvider postgres.PoolProvider
}

// NewStore creates new store.
func NewStore(db mongo.DbClient, pgPoolProvider postgres.PoolProvider) Store {
	return &store{
		db:             db,
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
		return nil, err
	}

	search, args := s.getSearchQuery(r)

	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('1 minute', time), count(*) FROM "+metrics.MessageRate+
		search+" GROUP BY time_bucket_gapfill('1 minute', time)", args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	rates := make([]StatsResponse, 0)
	for rows.Next() {
		var t time.Time
		var rateColumn *int64
		err := rows.Scan(&t, &rateColumn)
		if err != nil {
			return nil, err
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

	return rates, nil
}

func (s *store) findHourStats(ctx context.Context, r ListRequest) ([]StatsResponse, error) {
	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, err
	}

	search, args := s.getSearchQuery(r)

	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('1 hour', time), sum(count) FROM "+metrics.MessageRateHourly+
		search+" GROUP BY time_bucket_gapfill('1 hour', time)", args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	rates := make([]StatsResponse, 0)
	for rows.Next() {
		var t time.Time
		var rateColumn *int64
		err := rows.Scan(&t, &rateColumn)
		if err != nil {
			return nil, err
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

	collection := s.db.Collection(mongo.MessageRateStatsHourCollectionName)
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$gte": r.From, "$lte": r.To}},
		options.Find().SetSort(bson.M{"_id": 1}))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	from := r.From.Unix()
	interval := int64(time.Hour.Seconds())
	for cursor.Next(ctx) {
		var rate StatsResponse
		err = cursor.Decode(&rate)
		if err != nil {
			return nil, err
		}

		i := int((rate.ID - from) / interval)
		if i < 0 || i >= len(rates) {
			return nil, errors.New("invalid postgres query, rates must contain gaps")
		}

		rates[i].Rate += rate.Rate
	}

	return rates, nil
}

func (s *store) GetDeletedBeforeForHours(ctx context.Context) (*datetime.CpsTime, error) {
	res := struct {
		Time datetime.CpsTime `bson:"_id"`
	}{}
	collection := s.db.Collection(mongo.MessageRateStatsHourCollectionName)
	err := collection.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"_id": 1})).Decode(&res)
	if err == nil {
		return &res.Time, nil
	}

	if !errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, err
	}

	pgPool, err := s.pgPoolProvider.Get(ctx)
	if err != nil {
		return nil, err
	}

	var t time.Time
	err = pgPool.QueryRow(ctx, "SELECT min(time) FROM "+metrics.MessageRateHourly).Scan(&t)
	if err != nil {
		return nil, err
	}

	return &datetime.CpsTime{Time: t}, nil
}

func (s *store) getSearchQuery(r ListRequest) (string, pgx.NamedArgs) {
	var start, end time.Time

	if r.From.IsZero() || r.To.IsZero() {
		nowTrunc := time.Now().Truncate(time.Minute).UTC()

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
