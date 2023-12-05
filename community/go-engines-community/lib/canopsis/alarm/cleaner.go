package alarm

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Cleaner interface {
	// ArchiveResolvedAlarms archives alarm to archived alarm collection.
	ArchiveResolvedAlarms(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error)

	// DeleteArchivedResolvedAlarms deletes resolved alarms from archived collection after some time.
	DeleteArchivedResolvedAlarms(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error)
}

func NewCleaner(dbClient mongo.DbClient, bulkSize int) Cleaner {
	return &cleaner{
		resolvedDbCollection: dbClient.Collection(mongo.ResolvedAlarmMongoCollection),
		archivedDbCollection: dbClient.Collection(mongo.ArchivedAlarmMongoCollection),
		bulkSize:             bulkSize,
	}
}

type cleaner struct {
	resolvedDbCollection mongo.DbCollection
	archivedDbCollection mongo.DbCollection
	bulkSize             int
}

func (c *cleaner) ArchiveResolvedAlarms(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.resolvedDbCollection.Find(ctx, bson.M{
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
			return 0, err
		}

		writeModel := mongodriver.NewUpdateOneModel().
			SetFilter(bson.M{"_id": alarm.ID}).
			SetUpdate(bson.M{"$set": alarm}).
			SetUpsert(true)
		b, err := bson.Marshal(writeModel)
		if err != nil {
			return 0, err
		}
		newModelLen := len(b)
		if bulkBytesSize+newModelLen > canopsis.DefaultBulkBytesSize {
			res, err := c.archivedDbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return 0, err
			}

			archived += res.UpsertedCount

			_, err = c.resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
			if err != nil {
				return 0, err
			}

			bulkBytesSize = 0
			writeModels = writeModels[:0]
			archivedIds = archivedIds[:0]
		}

		bulkBytesSize += newModelLen
		writeModels = append(writeModels, writeModel)
		archivedIds = append(archivedIds, alarm.ID)

		if len(writeModels) >= c.bulkSize {
			res, err := c.archivedDbCollection.BulkWrite(ctx, writeModels)
			if err != nil {
				return 0, err
			}

			archived += res.UpsertedCount

			_, err = c.resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
			if err != nil {
				return 0, err
			}

			bulkBytesSize = 0
			writeModels = writeModels[:0]
			archivedIds = archivedIds[:0]
		}
	}

	if len(writeModels) > 0 {
		res, err := c.archivedDbCollection.BulkWrite(ctx, writeModels)
		if err != nil {
			return 0, err
		}

		archived += res.UpsertedCount

		_, err = c.resolvedDbCollection.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": archivedIds}})
		if err != nil {
			return 0, err
		}
	}

	return archived, nil
}

func (c *cleaner) DeleteArchivedResolvedAlarms(ctx context.Context, before datetime.CpsTime, limit int64) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := c.archivedDbCollection.Find(ctx, bson.M{
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
			return 0, err
		}

		ids = append(ids, alarm.ID)

		if len(ids) >= c.bulkSize {
			res, err := c.archivedDbCollection.DeleteMany(
				ctx,
				bson.M{"_id": bson.M{"$in": ids}},
			)
			if err != nil {
				return 0, err
			}

			deleted += res
			ids = ids[:0]
		}
	}

	if len(ids) > 0 {
		res, err := c.archivedDbCollection.DeleteMany(
			ctx,
			bson.M{"_id": bson.M{"$in": ids}},
		)
		if err != nil {
			return 0, err
		}

		deleted += res
	}

	return deleted, nil
}
