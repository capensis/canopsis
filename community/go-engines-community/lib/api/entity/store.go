package entity

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
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

		cursor, err := s.mainCollection.Aggregate(ctx, []bson.M{
			{"$match": bson.M{"_id": id}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "component",
				"as":                      "resources",
				"restrictSearchWithMatch": bson.M{"type": types.EntityTypeResource},
				"maxDepth":                0,
			}},
			{"$project": bson.M{
				"_id":       1,
				"enabled":   1,
				"type":      1,
				"resources": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		})
		if err != nil {
			return err
		}
		if cursor.Next(ctx) {
			err = cursor.Decode(&oldSimplifiedEntity)
			if err != nil {
				return err
			}
		} else {
			return nil
		}

		_, err = s.mainCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"enabled": enabled}})
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

	if oldSimplifiedEntity.ID == "" {
		return false, SimplifiedEntity{}, nil
	}

	if isToggled && !enabled && oldSimplifiedEntity.Type == types.EntityTypeComponent {
		depLen := len(oldSimplifiedEntity.Resources)
		from := 0

		for to := canopsis.DefaultBulkSize; to <= depLen; to += canopsis.DefaultBulkSize {
			_, err = s.mainCollection.UpdateMany(
				ctx,
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Resources[from:to]}},
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
				bson.M{"_id": bson.M{"$in": oldSimplifiedEntity.Resources[from:depLen]}},
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
	entity := Entity{}
	err := s.mainCollection.
		FindOne(ctx, bson.M{"_id": id}, options.FindOne().SetProjection(bson.M{"type": 1})).
		Decode(&entity)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	pipeline := []bson.M{{"$match": bson.M{"_id": id}}}
	switch entity.Type {
	case types.EntityTypeResource:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$component",
				"connectFromField":        "component",
				"connectToField":          "_id",
				"as":                      "component",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$connector",
				"connectFromField":        "connector",
				"connectToField":          "_id",
				"as":                      "connector",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$component", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$connector", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeComponent:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "component",
				"as":               "resources",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeResource,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$connector",
				"connectFromField":        "connector",
				"connectToField":          "_id",
				"as":                      "connector",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$connector", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeConnector:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "connector",
				"as":               "resources",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeResource,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":             mongo.EntityMongoCollection,
				"startWith":        "$_id",
				"connectFromField": "_id",
				"connectToField":   "connector",
				"as":               "components",
				"restrictSearchWithMatch": bson.M{
					"type":         types.EntityTypeComponent,
					"soft_deleted": bson.M{"$exists": false},
				},
				"maxDepth": 0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "services",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact": bson.M{"$concatArrays": bson.A{
					bson.M{"$map": bson.M{"input": "$resources", "in": "$$this._id"}},
					bson.M{"$map": bson.M{"input": "$services", "in": "$$this._id"}},
				}},
				"depends": bson.M{"$map": bson.M{"input": "$components", "in": "$$this._id"}},
			}},
		}...)
	case types.EntityTypeService:
		pipeline = append(pipeline, []bson.M{
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$_id",
				"connectFromField":        "_id",
				"connectToField":          "services",
				"as":                      "depends",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$graphLookup": bson.M{
				"from":                    mongo.EntityMongoCollection,
				"startWith":               "$services",
				"connectFromField":        "services",
				"connectToField":          "_id",
				"as":                      "impact",
				"restrictSearchWithMatch": bson.M{"soft_deleted": bson.M{"$exists": false}},
				"maxDepth":                0,
			}},
			{"$addFields": bson.M{
				"impact":  bson.M{"$map": bson.M{"input": "$impact", "in": "$$this._id"}},
				"depends": bson.M{"$map": bson.M{"input": "$depends", "in": "$$this._id"}},
			}},
		}...)
	}

	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"impact":  1,
		"depends": 1,
	}})
	cursor, err := s.mainCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		res := ContextGraphResponse{}
		err = cursor.Decode(&res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	return nil, nil
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
