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

	bucket := "1 minute"
	table := metrics.MessageRate
	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('"+bucket+"', time), count(*) FROM "+table+
		" WHERE time >= $1 AND time <= $2 GROUP BY time_bucket_gapfill('"+bucket+"', time)", r.From.UTC(), r.To.UTC())
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

	bucket := "1 hour"
	table := metrics.MessageRateHourly
	rows, err := pgPool.Query(ctx, "SELECT time_bucket_gapfill('"+bucket+"', time), sum(count) FROM "+table+
		" WHERE time >= $1 AND time <= $2 GROUP BY time_bucket_gapfill('"+bucket+"', time)", r.From.UTC(), r.To.UTC())
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
			return nil, fmt.Errorf("invalid postgres query, rates must contain gaps")
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
