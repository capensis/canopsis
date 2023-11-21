package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
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

	var lastExecuted datetime.CpsTime
	if conf.History.Pbehavior != nil {
		lastExecuted = *conf.History.Pbehavior
	}

	ok := datastorage.CanRun(lastExecuted, w.DataStorageConfigProvider.Get().TimeToExecute, w.TimezoneConfigProvider.Get().Location)
	if !ok {
		return
	}

	d := conf.Config.Pbehavior.DeleteAfter
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
	now := datetime.NewCpsTime()
	cleaner := pbehavior.NewCleaner(mongoClient, datastorage.BulkSize, w.Logger)
	maxUpdates := int64(w.DataStorageConfigProvider.Get().MaxUpdates)
	deleted, err := cleaner.Clean(ctx, d.SubFrom(now), maxUpdates)
	if err != nil {
		w.Logger.Err(err).Msg("cannot delete pbehaviors")
		return
	}
	if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("pbehaviors were deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryPbehavior(ctx, now)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}
}
