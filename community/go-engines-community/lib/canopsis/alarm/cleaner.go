package alarm

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewCleaner(logger zerolog.Logger) datastorage.Cleaner {
	return &cleaner{
		logger:   logger,
		bulkSize: datastorage.BulkSize,
	}
}

type cleaner struct {
	bulkSize int
	logger   zerolog.Logger
}

func (c *cleaner) IsEnabled(conf datastorage.Config) bool {
	return datetime.IsDurationEnabledAndValid(conf.Alarm.ArchiveAfter) ||
		datetime.IsDurationEnabledAndValid(conf.Alarm.DeleteAfter)
}

func (c *cleaner) Clean(ctx context.Context, dbClient mongo.DbClient, conf datastorage.Config, t datetime.CpsTime, limit int) (datastorage.CleanResult, error) {
	res := datastorage.CleanResult{}
	if !c.IsEnabled(conf) {
		return res, nil
	}

	defer func() {
		if res.Archived > 0 {
			c.logger.Info().Int64("alarm_number", res.Archived).Msg("resolved alarm archiving")
		}

		if res.Deleted > 0 {
			c.logger.Info().Int64("alarm_number", res.Deleted).Msg("resolved alarm removing")
		}
	}()

	resolvedDbCollection := dbClient.Collection(mongo.ResolvedAlarmMongoCollection)
	archivedDbCollection := dbClient.Collection(mongo.ArchivedAlarmMongoCollection)
	var err error
	archiveAfter := conf.Alarm.ArchiveAfter
	if datetime.IsDurationEnabledAndValid(archiveAfter) {
		res.Archived, err = c.archiveResolvedAlarms(ctx, resolvedDbCollection, archivedDbCollection, archiveAfter.SubFrom(t), limit)
		if err != nil {
			return res, fmt.Errorf("cannot archive resolved alarms: %w", err)
		}
	}

	deleteAfter := conf.Alarm.DeleteAfter
	if datetime.IsDurationEnabledAndValid(deleteAfter) {
		res.Deleted, err = c.deleteArchivedResolvedAlarms(ctx, archivedDbCollection, deleteAfter.SubFrom(t), limit)
		if err != nil {
			return res, fmt.Errorf("cannot delete resolved alarms: %w", err)
		}
	}

	return res, nil
}

// archiveResolvedAlarms archives alarm to archived alarm collection.
func (c *cleaner) archiveResolvedAlarms(ctx context.Context, resolvedDbCollection, archivedDbCollection mongo.DbCollection, before datetime.CpsTime, limit int) (int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := resolvedDbCollection.Find(ctx, bson.M{
		"v.resolved": bson.M{"$lte": before},
	}, opts)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	writeModels := make([]mongodriver.WriteModel, 0, c.bulkSize)
	archivedIds := make([]string, 0, c.bulkSize)
	bulkBytesSize := 0
	var archived int64

	for cursor.Next(ctx) {
		var alarm types.Alarm
		err := cursor.Decode(&alarm)
		if err != nil {
			return archived, err
		}

		writeModel := mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": alarm.ID}).
			SetUpdate(bson.M{"$set": alarm}).
			SetUpsert(true)
		b, err := bson.Marshal(writeModel)
		if err != nil {
			return archived, err
		}
		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			res, err := archivedDbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return archived, err
			}

			archived += res.UpsertedCount
			_, err = resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
			if err != nil {
				return archived, err
			}

			bulkBytesSize = 0
			writeModels = writeModels[:0]
			archivedIds = archivedIds[:0]
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, writeModel)
		archivedIds = append(archivedIds, alarm.ID)

		if len(writeModels) >= c.bulkSize {
			res, err := archivedDbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return archived, err
			}

			archived += res.UpsertedCount
			_, err = resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
			if err != nil {
				return archived, err
			}

			bulkBytesSize = 0
			writeModels = writeModels[:0]
			archivedIds = archivedIds[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return archived, err
	}

	if len(writeModels) > 0 {
		res, err := archivedDbCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return archived, err
		}

		archived += res.UpsertedCount

		_, err = resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
		if err != nil {
			return archived, err
		}
	}

	return archived, nil
}

// deleteArchivedResolvedAlarms deletes resolved alarms from archived collection after some time.
func (c *cleaner) deleteArchivedResolvedAlarms(ctx context.Context, archivedDbCollection mongo.DbCollection, before datetime.CpsTime, limit int) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := archivedDbCollection.Find(ctx, bson.M{
		"v.resolved": bson.M{"$lte": before},
	}, opts)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	ids := make([]string, 0, c.bulkSize)
	var deleted int64

	for cursor.Next(ctx) {
		var alarm types.Alarm
		err := cursor.Decode(&alarm)
		if err != nil {
			return deleted, err
		}

		ids = append(ids, alarm.ID)

		if len(ids) >= c.bulkSize {
			res, err := archivedDbCollection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return deleted, err
			}

			deleted += res
			ids = ids[:0]
		}
	}

	if err = cursor.Err(); err != nil {
		return deleted, err
	}

	if len(ids) > 0 {
		res, err := archivedDbCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return deleted, err
		}

		deleted += res
	}

	return deleted, nil
}
