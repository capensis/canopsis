package permission

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	securitymodel "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, request ListRequest) (*AggregationResult, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.RightsMongoCollection),
		defaultSearchByFields: []string{"_id", "crecord_name", "description"},
		defaultSortBy:         "name",
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
	defaultSortBy         string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"crecord_type": securitymodel.LineTypeObject}},
		{"$addFields": bson.M{
			"name": "$crecord_name",
		}},
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	sortBy := s.defaultSortBy
	if r.SortBy != "" {
		sortBy = r.SortBy
	}

	project := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "view",
		}},
		{"$unwind": bson.M{"path": "$view", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewGroupMongoCollection,
			"localField":   "view.group_id",
			"foreignField": "_id",
			"as":           "view_group",
		}},
		{"$unwind": bson.M{"path": "$view_group", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.PlaylistMongoCollection,
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "playlist",
		}},
		{"$unwind": bson.M{"path": "$playlist", "preserveNullAndEmptyArrays": true}},
	}

	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		common.GetSortQuery(sortBy, r.Sort),
		project,
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
