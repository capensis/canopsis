package entityinfodictionary

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	FindKeys(ctx context.Context, r ListKeysRequest) (AggregationResult, error)
	FindValues(ctx context.Context, r ListValuesRequest) (AggregationResult, error)
}

type store struct {
	db         mongo.DbClient
	collection mongo.DbCollection
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:         db,
		collection: db.Collection(mongo.EntityInfosDictionaryCollection),
	}
}

func (s *store) FindKeys(ctx context.Context, r ListKeysRequest) (AggregationResult, error) {
	res := AggregationResult{}

	var pipeline []bson.M

	searchQuery := common.GetSearchQuery(r.Search, []string{"_id"})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "_id", Value: 1}}},
		[]bson.M{{"$project": bson.M{"value": "$_id"}}},
	))
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
	}

	return res, err
}

func (s *store) FindValues(ctx context.Context, r ListValuesRequest) (AggregationResult, error) {
	res := AggregationResult{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": r.Key,
			},
		},
		{
			"$unwind": "$values",
		},
	}

	searchQuery := common.GetSearchQuery(r.Search, []string{"values"})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "values", Value: 1}}},
		[]bson.M{{"$project": bson.M{"value": "$values"}}},
	))
	if err != nil {
		return res, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
	}

	return res, err
}
