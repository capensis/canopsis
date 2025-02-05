package permission

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type Store interface {
	Find(ctx context.Context, request ListRequest) (*AggregationResult, error)
}

func NewStore(dbClient mongo.DbClient) Store {
	return &store{
		dbClient:              dbClient,
		dbCollection:          dbClient.Collection(mongo.PermissionCollection),
		defaultSearchByFields: []string{"_id", "name", "description"},
	}
}

type store struct {
	dbClient              mongo.DbClient
	dbCollection          mongo.DbCollection
	defaultSearchByFields []string
}

func (s *store) Find(ctx context.Context, r ListRequest) (*AggregationResult, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"hidden": bson.M{"$in": bson.A{false, nil}}}},
	}

	filter := common.GetSearchQuery(r.Search, s.defaultSearchByFields)
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.M{"$match": filter})
	}

	pipeline = append(pipeline, []bson.M{
		{"$unwind": bson.M{
			"path":                       "$groups",
			"preserveNullAndEmptyArrays": true,
			"includeArrayIndex":          "group_index",
		}},
		{"$lookup": bson.M{
			"from":         mongo.PermissionGroupCollection,
			"localField":   "groups",
			"foreignField": "_id",
			"as":           "group_perm",
		}},
		{"$unwind": bson.M{"path": "$group_perm", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewGroupMongoCollection,
			"localField":   "groups",
			"foreignField": "_id",
			"as":           "group_viewgroup",
		}},
		{"$unwind": bson.M{"path": "$group_viewgroup", "preserveNullAndEmptyArrays": true}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "groups",
			"foreignField": "_id",
			"as":           "group_view",
		}},
		{"$unwind": bson.M{"path": "$group_view", "preserveNullAndEmptyArrays": true}},
		{"$sort": bson.M{"group_index": 1}},
		{"$group": bson.M{
			"_id": "$_id",
			"doc": bson.M{"$first": "$$ROOT"},
			"groups": bson.M{"$push": bson.M{"$switch": bson.M{
				"branches": []bson.M{
					{"case": "$group_perm", "then": "$group_perm"},
					{"case": "$group_viewgroup", "then": "$group_viewgroup"},
					{"case": "$group_view", "then": "$group_view"},
				},
				"default": "$$REMOVE",
			}}},
		}},
		{"$lookup": bson.M{
			"from":         mongo.ViewMongoCollection,
			"localField":   "doc.view",
			"foreignField": "_id",
			"as":           "view",
		}},
		{"$unwind": bson.M{"path": "$view", "preserveNullAndEmptyArrays": true}},
		{"$replaceRoot": bson.M{"newRoot": bson.M{"$mergeObjects": bson.A{
			"$doc",
			bson.M{"groups": "$groups"},
			bson.M{"position": bson.M{"$cond": bson.M{
				"if":   "$view",
				"then": "$view.position",
				"else": 0,
			}}},
		}}}},
	}...)
	cursor, err := s.dbCollection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{
			{Key: "position", Value: 1},
			{Key: "name", Value: 1},
			{Key: "_id", Value: 1},
		}},
	))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	res := AggregationResult{}
	if cursor.Next(ctx) {
		err := cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}
