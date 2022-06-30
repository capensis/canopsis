package entityinfodictionary

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (AggregationResult, error)
}

type store struct {
	db         mongo.DbClient
	collection mongo.DbCollection
}

func NewStore(db mongo.DbClient) Store {
	return &store{
		db:         db,
		collection: db.Collection(mongo.InfosDictionaryCollection),
	}
}

func (s *store) Find(ctx context.Context, r ListRequest) (AggregationResult, error) {
	res := AggregationResult{}

	var pipeline []bson.M

	searchField := "_id"
	if r.Key != "" {
		pipeline = []bson.M{
			{
				"$match": bson.M{
					"_id": r.Key,
				},
			},
			{
				"$unwind": "$values",
			},
		}

		searchField = "values"
	}

	searchQuery := common.GetSearchQuery(r.Search, []string{searchField})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{{Key: searchField, Value: 1}}},
		[]bson.M{{"$project": bson.M{"value": "$" + searchField}}},
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
