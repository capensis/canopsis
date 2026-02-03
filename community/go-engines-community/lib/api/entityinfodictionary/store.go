package entityinfodictionary

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/mongoquery"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	searchQuery := mongoquery.GetSearchQuery(r.Search, []string{"_id.k"})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	pipeline = append(pipeline, []bson.M{
		{
			"$unionWith": bson.M{
				"coll": mongo.EntityInfosPropertyCollection,
				"pipeline": []bson.M{
					{
						"$match": bson.M{
							"name": bson.Regex{
								Pattern: fmt.Sprintf(".*%s.*", r.Search),
								Options: "i",
							},
						},
					},
					{
						"$project": bson.M{
							"_id.k":    "$name",
							"proptype": "$type",
						},
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id":      "$_id.k",
				"type":     bson.M{"$max": "$type"},
				"proptype": bson.M{"$max": "$proptype"},
			},
		},
		{
			"$set": bson.M{
				"type": bson.M{
					"$cond": bson.A{
						bson.M{"$or": bson.A{
							bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$proptype", ""}}, ""}},
							bson.M{"$gt": bson.A{"$proptype", types.EntityInfoTypeStringArray}},
							bson.M{"$lt": bson.A{"$proptype", types.EntityInfoTypeBoolean}},
						}},
						bson.M{
							"$cond": bson.A{
								bson.M{"$or": bson.A{
									bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$type", ""}}, ""}},
									bson.M{"$gt": bson.A{"$type", types.EntityInfoTypeStringArray}},
									bson.M{"$lt": bson.A{"$type", types.EntityInfoTypeBoolean}},
								}},
								types.EntityInfoTypeString,
								"$type",
							},
						},
						"$proptype",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"proptype": 0,
			},
		},
	}...)

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "_id", Value: 1}}},
		[]bson.M{{"$project": bson.M{"value": "$_id", "type": "$type"}}},
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
				"_id.k": r.Key,
				"_id.v": bson.M{"$ne": ""},
			},
		},
	}

	searchQuery := mongoquery.GetSearchQuery(r.Search, []string{"_id.v"})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	cursor, err := s.collection.Aggregate(ctx, pagination.CreateAggregationPipeline(
		r.Query,
		pipeline,
		bson.M{"$sort": bson.D{{Key: "_id.v", Value: 1}}},
		[]bson.M{{"$project": bson.M{"value": "$_id.v"}}},
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
