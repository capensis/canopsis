package entity

import (
	"context"
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
	ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error)
	DeleteArchivedEntities(ctx context.Context) (int64, error)
	Toggle(ctx context.Context, id string, enabled bool) (bool, SimplifiedEntity, error)
	GetContextGraph(ctx context.Context, id string) (*ContextGraphResponse, error)
}

type store struct {
	db                 mongo.DbClient
	mainCollection     mongo.DbCollection
	archivedCollection mongo.DbCollection

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

func (s *store) ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error) {
	var totalArchived int64

	// do not cascade-archive connector dependencies
	archived, err := s.archiveEntitiesByType(ctx, types.EntityTypeConnector, false)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeComponent, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeResource, false)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	return totalArchived, nil
}

func (s *store) archiveEntitiesByType(ctx context.Context, eType string, archiveDeps bool) (int64, error) {
	cursor, err := s.mainCollection.Find(
		ctx,
		bson.M{
			"enabled": bson.M{"$in": bson.A{false, nil}},
			"type":    eType,
		},
	)
	if err != nil {
		return 0, err
	}

	totalArchived, err := s.processCursor(ctx, cursor, archiveDeps)
	if err != nil {
		return 0, err
	}

	return totalArchived, cursor.Close(ctx)
}

func (s *store) processCursor(ctx context.Context, cursor mongo.Cursor, archiveDeps bool) (int64, error) {
	archiveModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize/2) // per 1 archive model there are 2 context graph update models
	contextGraphModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	ids := make([]string, 0, canopsis.DefaultBulkSize/2)

	archiveBulkBytesSize := 0
	contextGraphBulkBytesSize := 0

	var totalArchived int64

	for cursor.Next(ctx) {
		var entity types.Entity

		err := cursor.Decode(&entity)
		if err != nil {
			return 0, err
		}

		if entity.Type == types.EntityTypeComponent && archiveDeps {
			archived, err := s.archiveComponentDependencies(ctx, entity.Depends)
			if err != nil {
				return 0, err
			}

			totalArchived += archived
		}

		newContextGraphModels := []mongodriver.WriteModel{
			mongodriver.NewUpdateManyModel().
				SetFilter(bson.M{"_id": bson.M{"$in": entity.Depends}}).
				SetUpdate(bson.M{"$pull": bson.M{"impact": entity.ID}}),
			mongodriver.NewUpdateManyModel().
				SetFilter(bson.M{"_id": bson.M{"$in": entity.Impacts}}).
				SetUpdate(bson.M{"$pull": bson.M{"depends": entity.ID}}),
		}

		b, err := bson.Marshal(struct {
			Arr []mongodriver.WriteModel
		}{Arr: newContextGraphModels})
		if err != nil {
			return 0, err
		}

		newContextGraphModelsLen := len(b)

		newArchiveModel := mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$set": entity}).
			SetUpsert(true)

		b, err = bson.Marshal(newArchiveModel)
		if err != nil {
			return 0, err
		}

		newArchiveModelLen := len(b)

		if contextGraphBulkBytesSize+newContextGraphModelsLen > canopsis.DefaultBulkBytesSize ||
			archiveBulkBytesSize+newArchiveModelLen > canopsis.DefaultBulkBytesSize {
			archived, err := s.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived

			contextGraphModels = contextGraphModels[:0]
			archiveModels = archiveModels[:0]
			ids = ids[:0]

			contextGraphBulkBytesSize = 0
			archiveBulkBytesSize = 0
		}

		archiveModels = append(archiveModels, newArchiveModel)
		contextGraphModels = append(contextGraphModels, newContextGraphModels...)
		ids = append(ids, entity.ID)

		contextGraphBulkBytesSize += newContextGraphModelsLen
		archiveBulkBytesSize += newArchiveModelLen

		if len(contextGraphModels) == canopsis.DefaultBulkSize {
			archived, err := s.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived

			contextGraphModels = contextGraphModels[:0]
			archiveModels = archiveModels[:0]
			ids = ids[:0]

			contextGraphBulkBytesSize = 0
			archiveBulkBytesSize = 0
		}
	}

	if len(archiveModels) != 0 {
		archived, err := s.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
		if err != nil {
			return 0, err
		}

		totalArchived += archived
	}

	return totalArchived, nil
}

func (s *store) archiveComponentDependencies(ctx context.Context, depIds []string) (int64, error) {
	cursor, err := s.mainCollection.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": depIds}},
	)
	if err != nil {
		return 0, err
	}

	archived, err := s.processCursor(ctx, cursor, false)
	if err != nil {
		return 0, err
	}

	return archived, cursor.Close(ctx)
}

func (s *store) bulkArchive(ctx context.Context, models, contextGraphModels []mongodriver.WriteModel, ids []string) (int64, error) {
	res, err := s.archivedCollection.BulkWrite(ctx, models)
	if err != nil {
		return 0, err
	}

	_, err = s.mainCollection.BulkWrite(ctx, contextGraphModels)
	if err != nil {
		return 0, err
	}

	_, err = s.mainCollection.DeleteMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
	)
	if err != nil {
		return 0, err
	}

	return res.UpsertedCount + res.ModifiedCount, err
}

func (s *store) DeleteArchivedEntities(ctx context.Context) (int64, error) {
	deleted, err := s.archivedCollection.DeleteMany(ctx, bson.M{})

	return deleted, err
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
				SetProjection(bson.M{"_id": 1, "enabled": 1, "type": 1}).
				SetReturnDocument(options.Before),
		).Decode(&oldSimplifiedEntity)
		if err != nil {
			return err
		}

		isToggled = oldSimplifiedEntity.Enabled != enabled
		return nil
	})

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
