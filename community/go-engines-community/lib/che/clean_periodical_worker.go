package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	libtime "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/time"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cleanPeriodicalWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	Logger                    zerolog.Logger
}

func (w *cleanPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *cleanPeriodicalWorker) Work(ctx context.Context) {
	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("fail to retrieve data storage config")
		return
	}

	var lastExecuted libtime.CpsTime
	if conf.History.EventFilterFailure != nil {
		lastExecuted = *conf.History.EventFilterFailure
	}

	ok := datastorage.CanRun(lastExecuted, w.DataStorageConfigProvider.Get().TimeToExecute, w.TimezoneConfigProvider.Get().Location)
	if !ok {
		return
	}

	d := conf.Config.EventFilterFailure.DeleteAfter
	if d == nil || !*d.Enabled || d.Value == 0 {
		return
	}

	mongoClient, err := mongo.NewClientWithOptions(ctx, 0, 0, 0,
		w.DataStorageConfigProvider.Get().MongoClientTimeout, w.Logger)
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
	now := libtime.NewCpsTime()
	maxUpdates := int64(w.DataStorageConfigProvider.Get().MaxUpdates)
	deleted, err := w.delete(ctx, mongoClient.Collection(mongo.EventFilterFailureCollection), d.SubFrom(now), maxUpdates)
	if err != nil {
		w.Logger.Err(err).Msg("cannot delete event filter failures")
		return
	}
	if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("event filter failures were deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryEventFilterFailure(ctx, now)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}
}

func (w *cleanPeriodicalWorker) delete(ctx context.Context, dbCollection mongo.DbCollection, before libtime.CpsTime, limit int64) (int64, error) {
	opts := options.Find().SetProjection(bson.M{"_id": 1})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := dbCollection.Find(ctx, bson.M{
		"t": bson.M{"$lt": before},
	}, opts)
	if err != nil {
		return 0, err
	}

	defer cursor.Close(ctx)

	ids := make([]string, 0, datastorage.BulkSize)
	var deleted int64

	for cursor.Next(ctx) {
		var item struct {
			ID string `bson:"_id"`
		}
		err := cursor.Decode(&item)
		if err != nil {
			return 0, err
		}

		ids = append(ids, item.ID)

		if len(ids) >= datastorage.BulkSize {
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
