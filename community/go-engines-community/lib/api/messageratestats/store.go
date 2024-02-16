package messageratestats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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
	db mongo.DbClient
}

// NewStore creates new store.
func NewStore(db mongo.DbClient) Store {
	return &store{
		db: db,
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) ([]StatsResponse, error) {
	collectionName := ""
	var interval int64
	switch r.Interval {
	case IntervalMinute:
		collectionName = mongo.MessageRateStatsMinuteCollectionName
		interval = int64(time.Minute.Seconds())
	case IntervalHour:
		collectionName = mongo.MessageRateStatsHourCollectionName
		interval = int64(time.Hour.Seconds())
	default:
		return nil, fmt.Errorf("unknown interval %v", r.Interval)
	}

	l := 1 + (r.To.Unix()-r.From.Unix())/interval
	rates := make([]StatsResponse, l)
	from := r.From.Unix()

	for i := int64(0); i < l; i++ {
		rates[i].ID = i*interval + from
		rates[i].Received = 0
	}

	collection := s.db.Collection(collectionName)
	cursor, err := collection.Find(ctx, bson.M{"_id": bson.M{"$gte": r.From, "$lte": r.To}})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var rate StatsResponse
		err = cursor.Decode(&rate)
		if err != nil {
			return nil, err
		}

		rates[(rate.ID-from)/interval].Received = rate.Received
	}

	return rates, nil
}

func (s *store) GetDeletedBeforeForHours(ctx context.Context) (*datetime.CpsTime, error) {
	res := struct {
		Time datetime.CpsTime `bson:"_id"`
	}{}

	err := s.db.Collection(mongo.MessageRateStatsHourCollectionName).FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"_id": 1})).Decode(&res)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, nil
	}

	return &res.Time, err
}
