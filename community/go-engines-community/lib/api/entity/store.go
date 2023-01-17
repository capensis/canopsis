package entity

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store interface {
	Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error)
	Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error)
	GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error)
}

type store struct {
	db                     mongo.DbClient
	mainCollection         mongo.DbCollection
	archivedCollection     mongo.DbCollection
	timezoneConfigProvider config.TimezoneConfigProvider
}

func NewStore(db mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Store {
	return &store{
		db:                     db,
		mainCollection:         db.Collection(mongo.EntityMongoCollection),
		archivedCollection:     db.Collection(mongo.ArchivedEntitiesMongoCollection),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error) {
	location := s.timezoneConfigProvider.Get().Location
	now := types.CpsTime{Time: time.Now().In(location)}

	pipeline, err := s.getQueryBuilder().CreateListAggregationPipeline(ctx, r, now)
	if err != nil {
		return nil, err
	}

	cursor, err := s.mainCollection.Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))

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

	s.fillConnectorType(&res)

	return &res, nil
}

func (s *store) Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error) {
	var isToggled bool
	var oldSimplifiedEntity SimplifiedEntity

	err := s.db.WithTransaction(ctx, func(ctx context.Context) error {
		isToggled = false
		oldSimplifiedEntity = SimplifiedEntity{}

		err := s.mainCollection.FindOneAndUpdate(
			ctx,
			bson.M{"_id": id},
			bson.M{"$set": bson.M{"enabled": enabled}},
			options.
				FindOneAndUpdate().
				SetProjection(bson.M{"_id": 1, "enabled": 1, "type": 1, "depends": 1}).
				SetReturnDocument(options.Before),
		).Decode(&oldSimplifiedEntity)
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

	if isToggled && !enabled && oldSimplifiedEntity.Type == types.EntityTypeComponent {
		depLen := len(oldSimplifiedEntity.Depends)
		from := 0

		for to := canopsis.DefaultBulkSize; to <= depLen; to += canopsis.DefaultBulkSize {
			_, err = s.mainCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Depends[from:to]}},
				bson.M{"$set": bson.M{"enabled": enabled}},
			)
			if err != nil {
				return isToggled, oldSimplifiedEntity, err
			}

			from = to
		}

		if from < depLen {
			_, err = s.mainCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Depends[from:depLen]}},
				bson.M{"$set": bson.M{"enabled": enabled}},
			)
			if err != nil {
				return isToggled, oldSimplifiedEntity, err
			}
		}
	}

	return isToggled, oldSimplifiedEntity, err
}

func (s *store) GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error) {
	res := ContextGraphResponse{}

	cursor, err := s.mainCollection.Aggregate(ctx, []bson.M{
		{
			"$match": bson.M{"_id": id},
		},
		{
			"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$impact",
				"connectFromField":        "impact",
				"connectToField":          "_id",
				"as":                      "impact",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			},
		},
		{
			"$addFields": bson.M{
				"impact": bson.M{"$map": bson.M{"input": "$impact", "as": "each", "in": "$$each._id"}},
			},
		},
		{
			"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$depends",
				"connectFromField":        "depends",
				"connectToField":          "_id",
				"as":                      "depends",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			},
		},
		{
			"$addFields": bson.M{
				"depends": bson.M{"$map": bson.M{"input": "$depends", "as": "each", "in": "$$each._id"}},
			},
		},
		{
			"$project": bson.M{
				"impact":  1,
				"depends": 1,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return &res, nil
}

func (s *store) fillConnectorType(result *AggregationResult) {
	if result == nil {
		return
	}
	for i := range result.Data {
		result.Data[i].fillConnectorType()
	}
}

func (s *store) getQueryBuilder() *MongoQueryBuilder {
	return NewMongoQueryBuilder(s.db)
}
