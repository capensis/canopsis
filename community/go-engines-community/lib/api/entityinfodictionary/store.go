package entityinfodictionary

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/che"
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

	searchQuery := common.GetSearchQuery(r.Search, []string{"_id.k"})
	if searchQuery != nil {
		pipeline = append(pipeline, bson.M{"$match": searchQuery})
	}

	pipeline = append(pipeline, []bson.M{
		{
			"$group": bson.M{"_id": "$_id.k", "type": bson.M{"$max": "$type"}},
		},
		{
			"$lookup": bson.M{
				"from":         mongo.EntityInfosPropertyCollection,
				"localField":   "_id",
				"foreignField": "key",
				"as":           "props",
			},
		},
		{
			"$unwind": bson.M{"path": "$props", "preserveNullAndEmptyArrays": true},
		},
		{
			"$set": bson.M{
				"type": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{bson.M{"$ifNull": bson.A{"$props", ""}}, ""}},
						bson.M{
							"$switch": bson.M{
								"branches": []bson.M{
									{"case": bson.M{"$eq": bson.A{"$type", che.TypeStringArray}}, "then": "string_array"},
									{"case": bson.M{"$eq": bson.A{"$type", che.TypeString}}, "then": "string"},
									{"case": bson.M{"$eq": bson.A{"$type", che.TypeNumber}}, "then": "number"},
									{"case": bson.M{"$eq": bson.A{"$type", che.TypeBoolean}}, "then": "boolean"},
								},
								"default": che.TypeString,
							},
						},
						"$props.type",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"props": 0,
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

	searchQuery := common.GetSearchQuery(r.Search, []string{"_id.v"})
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
