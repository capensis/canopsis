package entity

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const bulkMaxSize = 10000

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

	queryBuilder           *MongoQueryBuilder
	timezoneConfigProvider config.TimezoneConfigProvider
}

func NewStore(db mongo.DbClient, timezoneConfigProvider config.TimezoneConfigProvider) Store {
	return &store{
		db:                     db,
		mainCollection:         db.Collection(mongo.EntityMongoCollection),
		archivedCollection:     db.Collection(mongo.ArchivedEntitiesMongoCollection),
		queryBuilder:           NewMongoQueryBuilder(db),
		timezoneConfigProvider: timezoneConfigProvider,
	}
}

func (s *store) Find(ctx context.Context, r ListRequestWithPagination) (*AggregationResult, error) {
	location := s.timezoneConfigProvider.Get().Location
	now := types.CpsTime{Time: time.Now().In(location)}

	pipeline, err := s.queryBuilder.CreateListAggregationPipeline(ctx, r, now)
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

	archived, err := s.archiveEntitiesByType(ctx, types.EntityTypeConnector, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeComponent, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = s.archiveEntitiesByType(ctx, types.EntityTypeResource, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	return totalArchived, nil
}

func (s *store) archiveEntitiesByType(ctx context.Context, eType string, archiveDeps bool) (int64, error) {
	models := make([]mongodriver.WriteModel, 0, bulkMaxSize)
	ids := make([]string, 0, bulkMaxSize)

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

	totalArchived, err := s.processCursor(ctx, cursor, &models, &ids, archiveDeps)
	if err != nil {
		return 0, err
	}

	if len(models) != 0 {
		archived, err := s.bulkArchive(ctx, &models, &ids)
		if err != nil {
			return 0, err
		}

		totalArchived += archived
	}

	return totalArchived, cursor.Close(ctx)
}

func (s *store) processCursor(ctx context.Context, cursor mongo.Cursor, models *[]mongodriver.WriteModel, ids *[]string, archiveDeps bool) (int64, error) {
	var totalArchived int64

	for cursor.Next(ctx) {
		var entity types.Entity

		err := cursor.Decode(&entity)
		if err != nil {
			return 0, err
		}

		if archiveDeps {
			switch entity.Type {
			case types.EntityTypeConnector:
				archived, err := s.archiveDependencies(ctx, append(entity.Impacts, entity.Depends...), models, ids)
				if err != nil {
					return 0, err
				}

				totalArchived += archived
			case types.EntityTypeComponent:
				archived, err := s.archiveDependencies(ctx, entity.Depends, models, ids)
				if err != nil {
					return 0, err
				}

				totalArchived += archived
			case types.EntityTypeResource:
			default:
				continue
			}
		}

		*ids = append(*ids, entity.ID)
		*models = append(
			*models,
			mongodriver.NewUpdateOneModel().
				SetFilter(bson.M{"_id": entity.ID}).
				SetUpdate(bson.M{"$set": entity}).
				SetUpsert(true),
		)

		if len(*models) == bulkMaxSize {
			archived, err := s.bulkArchive(ctx, models, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived
		}
	}

	return totalArchived, nil
}

func (s *store) archiveDependencies(ctx context.Context, depIds []string, models *[]mongodriver.WriteModel, ids *[]string) (int64, error) {
	cursor, err := s.mainCollection.Find(
		ctx,
		bson.M{"_id": bson.M{"$in": depIds}},
	)
	if err != nil {
		return 0, err
	}

	// do not archive dependencies' dependencies
	archived, err := s.processCursor(ctx, cursor, models, ids, false)
	if err != nil {
		return 0, err
	}

	return archived, cursor.Close(ctx)
}

func (s *store) bulkArchive(ctx context.Context, models *[]mongodriver.WriteModel, ids *[]string) (int64, error) {
	res, err := s.archivedCollection.BulkWrite(ctx, *models)
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

	*models = (*models)[:0]
	*ids = (*ids)[:0]

	return res.UpsertedCount, err
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
	err := s.mainCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&res)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
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
