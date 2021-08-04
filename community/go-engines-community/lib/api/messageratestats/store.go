package messageratestats

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(context.Context, ListRequest) (*AggregationResult, error)
	// get the lower bound time value for hourly /message-rate-stats request.
	GetDeletedBeforeForHours(ctx context.Context) (*types.CpsTime, error)
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

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	filter := bson.M{"$and": []bson.M{
		{"_id": bson.M{"$gte": r.From}},
		{"_id": bson.M{"$lte": r.To}},
	}}

	sortBy := "_id"
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	collectionNames := map[string]string{
		IntervalMinute: mongo.MessageRateStatsMinuteCollectionName,
		IntervalHour:   mongo.MessageRateStatsHourCollectionName,
	}
	collectionName, ok := collectionNames[r.Interval]
	if !ok {
		return nil, fmt.Errorf("unknown interval %v", r.Interval)
	}

	collection := s.db.Collection(collectionName)
	cursor, err := collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		[]bson.M{{"$match": filter}},
		common.GetSortQuery(sortBy, r.Sort),
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	res := AggregationResult{}

	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}

func (s *store) GetDeletedBeforeForHours(ctx context.Context) (*types.CpsTime, error) {
	res := struct {
		Time types.CpsTime `bson:"_id"`
	}{}

	err := s.db.Collection(mongo.MessageRateStatsHourCollectionName).FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"_id": 1})).Decode(&res)

	return &res.Time, err
}
