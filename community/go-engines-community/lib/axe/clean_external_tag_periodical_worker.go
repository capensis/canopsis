package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cleanExternalTagPeriodicalWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	Logger                    zerolog.Logger
}

func (w *cleanExternalTagPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *cleanExternalTagPeriodicalWorker) Work(ctx context.Context) {
	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot retrieve data storage config")

		return
	}

	var lastExecuted datetime.CpsTime
	if conf.History.AlarmExternalTag != nil {
		lastExecuted = conf.History.AlarmExternalTag.Time
	}

	dataStorageConf := w.DataStorageConfigProvider.Get()
	if !datastorage.CanRun(lastExecuted, dataStorageConf.TimeToExecute, w.TimezoneConfigProvider.Get().Location) {
		return
	}

	mongoClient, err := mongo.NewClientWithOptions(ctx, 0, 0, 0, dataStorageConf.MongoClientTimeout, w.Logger)
	if err != nil {
		w.Logger.Err(err).Msg("cannot connect to mongo")

		return
	}

	defer func() {
		err = mongoClient.Disconnect(ctx)
		if err != nil {
			w.Logger.Err(err).Msg("cannot disconnect from mongo")
		}
	}()

	deleteAfter := conf.Config.AlarmExternalTag.DeleteAfter
	if !datetime.IsDurationEnabledAndValid(deleteAfter) {
		return
	}

	now := datetime.NewCpsTime()
	dbCollection := mongoClient.Collection(mongo.AlarmTagCollection)
	colorDbCollection := mongoClient.Collection(mongo.AlarmTagColorCollection)
	deleted, err := w.delete(ctx, deleteAfter.SubFrom(now), dataStorageConf.MaxUpdates, datastorage.BulkSize, dbCollection)
	if err != nil {
		w.Logger.Err(err).Msg("cannot delete alarm tags")

		return
	}

	if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("alarm external tags are deleted")
	}

	deletedColors, err := w.deleteColors(ctx, dataStorageConf.MaxUpdates, datastorage.BulkSize, colorDbCollection)
	if err != nil {
		w.Logger.Err(err).Msg("cannot delete alarm tag colors")

		return
	}

	if deletedColors > 0 {
		w.Logger.Info().Int64("count", deletedColors).Msg("alarm tag colors are deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryAlarmExternalTag(ctx, datastorage.HistoryWithCount{
		Time:    now,
		Deleted: deleted,
	})
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}
}

func (w *cleanExternalTagPeriodicalWorker) delete(
	ctx context.Context,
	before datetime.CpsTime,
	limit int,
	bulkSize int,
	dbCollection mongo.DbCollection,
) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}

	cursor, err := dbCollection.Find(ctx, bson.M{
		"type":            alarmtag.TypeExternal,
		"last_event_date": bson.M{"$lte": before},
	}, opts)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, bulkSize)
	var deleted int64
	for cursor.Next(ctx) {
		var tag alarmtag.AlarmTag
		err := cursor.Decode(&tag)
		if err != nil {
			return 0, err
		}

		ids = append(ids, tag.ID)
		if len(ids) >= bulkSize {
			res, err := dbCollection.DeleteMany(
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
		res, err := dbCollection.DeleteMany(
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
func (w *cleanExternalTagPeriodicalWorker) deleteColors(
	ctx context.Context,
	limit int,
	bulkSize int,
	dbCollection mongo.DbCollection,
) (int64, error) {
	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from":         mongo.AlarmTagCollection,
			"localField":   "_id",
			"foreignField": "label",
			"as":           "tags",
			"pipeline": []bson.M{
				{"$limit": 1},
				{"$project": bson.M{
					"_id": 1,
				}},
			},
		}},
		{"$match": bson.M{
			"tags": bson.A{},
		}},
		{"$project": bson.M{
			"_id": 1,
		}},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": limit,
		})
	}

	cursor, err := dbCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)
	ids := make([]string, 0, bulkSize)
	var deleted int64
	for cursor.Next(ctx) {
		var color struct {
			ID string `bson:"_id"`
		}
		err := cursor.Decode(&color)
		if err != nil {
			return 0, err
		}

		ids = append(ids, color.ID)
		if len(ids) >= bulkSize {
			res, err := dbCollection.DeleteMany(
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
		res, err := dbCollection.DeleteMany(
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
