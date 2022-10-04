package alarmtag

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, r ListRequest) (*AggregationResult, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		collection: dbClient.Collection(mongo.AlarmTagCollection),

		defaultSearchByFields: []string{"value"},
		defaultSortBy:         "value",
	}
}

type store struct {
	collection mongo.DbCollection

	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	var pipeline []bson.M
	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
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
