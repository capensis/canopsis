package entity

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type Archiver interface {
	ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error)
	DeleteArchivedEntities(ctx context.Context) (int64, error)
}

type archiver struct {
	db                 mongo.DbClient
	mainCollection     mongo.DbCollection
	archivedCollection mongo.DbCollection
}

func NewArchiver(db mongo.DbClient) Archiver {
	return &archiver{
		db:                 db,
		mainCollection:     db.Collection(mongo.EntityMongoCollection),
		archivedCollection: db.Collection(mongo.ArchivedEntitiesMongoCollection),
	}
}

func (a *archiver) ArchiveDisabledEntities(ctx context.Context, archiveDeps bool) (int64, error) {
	var totalArchived int64

	// do not cascade-archive connector dependencies
	archived, err := a.archiveEntitiesByType(ctx, types.EntityTypeConnector, false)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = a.archiveEntitiesByType(ctx, types.EntityTypeComponent, archiveDeps)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	archived, err = a.archiveEntitiesByType(ctx, types.EntityTypeResource, false)
	if err != nil {
		return 0, err
	}

	totalArchived += archived

	return totalArchived, nil
}

func (a *archiver) archiveEntitiesByType(ctx context.Context, eType string, archiveDeps bool) (int64, error) {
	cursor, err := a.mainCollection.Find(
		ctx,
		bson.M{
			"enabled": bson.M{"$in": bson.A{false, nil}},
			"type":    eType,
		},
	)
	if err != nil {
		return 0, err
	}

	totalArchived, err := a.processCursor(ctx, cursor, archiveDeps)
	if err != nil {
		return 0, err
	}

	return totalArchived, cursor.Close(ctx)
}

func (a *archiver) processCursor(ctx context.Context, cursor mongo.Cursor, archiveDeps bool) (int64, error) {
	archiveModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	contextGraphModels := make([]mongodriver.WriteModel, 0, canopsis.DefaultBulkSize)
	ids := make([]string, 0, canopsis.DefaultBulkSize)
	archiveBulkBytesSize := 0
	var totalArchived int64

	for cursor.Next(ctx) {
		var entity types.Entity
		err := cursor.Decode(&entity)
		if err != nil {
			return 0, err
		}

		var newContextGraphModel mongodriver.WriteModel
		switch entity.Type {
		case types.EntityTypeComponent:
			if archiveDeps {
				archived, err := a.archiveComponentDependencies(ctx, entity.ID)
				if err != nil {
					return 0, err
				}

				totalArchived += archived
			}
		case types.EntityTypeConnector:
			newContextGraphModel = mongodriver.NewUpdateManyModel().
				SetFilter(bson.M{"connector": entity.ID}).
				SetUpdate(bson.M{"$unset": bson.M{"connector": ""}})
		}

		newArchiveModel := mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": entity.ID}).
			SetUpdate(bson.M{"$set": entity}).
			SetUpsert(true)

		b, err := bson.Marshal(newArchiveModel)
		if err != nil {
			return 0, err
		}

		newArchiveModelLen := len(b)
		archiveBulkBytesSize += newArchiveModelLen

		if archiveBulkBytesSize > canopsis.DefaultBulkBytesSize {
			archived, err := a.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived

			contextGraphModels = contextGraphModels[:0]
			archiveModels = archiveModels[:0]
			ids = ids[:0]

			archiveBulkBytesSize = newArchiveModelLen
		}

		archiveModels = append(archiveModels, newArchiveModel)
		if newContextGraphModel != nil {
			contextGraphModels = append(contextGraphModels, newContextGraphModel)
		}
		ids = append(ids, entity.ID)

		if len(archiveModels) == canopsis.DefaultBulkSize {
			archived, err := a.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
			if err != nil {
				return 0, err
			}

			totalArchived += archived

			contextGraphModels = contextGraphModels[:0]
			archiveModels = archiveModels[:0]
			ids = ids[:0]

			archiveBulkBytesSize = 0
		}
	}

	if len(archiveModels) != 0 {
		archived, err := a.bulkArchive(ctx, archiveModels, contextGraphModels, ids)
		if err != nil {
			return 0, err
		}

		totalArchived += archived
	}

	return totalArchived, nil
}

func (a *archiver) archiveComponentDependencies(ctx context.Context, id string) (int64, error) {
	cursor, err := a.mainCollection.Find(
		ctx,
		bson.M{"component": id},
	)
	if err != nil {
		return 0, err
	}

	archived, err := a.processCursor(ctx, cursor, false)
	if err != nil {
		_ = cursor.Close(ctx)
		return 0, err
	}

	return archived, cursor.Close(ctx)
}

func (a *archiver) bulkArchive(ctx context.Context, models, contextGraphModels []mongodriver.WriteModel, ids []string) (int64, error) {
	var count int64
	if len(models) > 0 {
		res, err := a.archivedCollection.BulkWrite(ctx, models)
		if err != nil {
			return 0, err
		}
		count = res.UpsertedCount + res.ModifiedCount
	}

	if len(contextGraphModels) > 0 {
		_, err := a.mainCollection.BulkWrite(ctx, contextGraphModels)
		if err != nil {
			return 0, err
		}
	}

	if len(ids) > 0 {
		_, err := a.mainCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

func (a *archiver) DeleteArchivedEntities(ctx context.Context) (int64, error) {
	var totalDeleted int64
	ids := make([]string, 0, canopsis.DefaultBulkSize)

	cursor, err := a.archivedCollection.Find(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	for cursor.Next(ctx) {
		var doc struct {
			ID string `bson:"_id"`
		}

		err = cursor.Decode(&doc)
		if err != nil {
			return 0, err
		}

		ids = append(ids, doc.ID)

		if len(ids) == canopsis.DefaultBulkSize {
			deleted, err := a.archivedCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
			if err != nil {
				return 0, err
			}

			ids = ids[:0]

			totalDeleted += deleted
		}
	}

	if len(ids) > 0 {
		deleted, err := a.archivedCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
		if err != nil {
			return 0, err
		}

		totalDeleted += deleted
	}

	return totalDeleted, nil
}
