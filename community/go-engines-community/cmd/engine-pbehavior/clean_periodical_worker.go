package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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

func (w *cleanPeriodicalWorker) Work(ctx context.Context) error {
	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("fail to retrieve data storage config")
		return nil
	}

	var lastExecuted types.CpsTime
	if conf.History.Pbehavior != nil {
		lastExecuted = *conf.History.Pbehavior
	}

	ok := datastorage.CanRun(lastExecuted, w.DataStorageConfigProvider.Get().TimeToExecute, w.TimezoneConfigProvider.Get().Location)
	if !ok {
		return nil
	}

	d := conf.Config.Pbehavior.DeleteAfter
	if d == nil || !*d.Enabled || d.Seconds == 0 {
		return nil
	}

	mongoClient, err := mongo.NewClientWithOptions(ctx, 0, 0, 0,
		w.DataStorageConfigProvider.Get().MongoClientTimeout)
	if err != nil {
		w.Logger.Err(err).Msg("cannot connect to mongo")
		return nil
	}
	defer func() {
		err = mongoClient.Disconnect(ctx)
		if err != nil {
			w.Logger.Err(err).Msg("cannot disconnect from mongo")
		}
	}()
	now := types.CpsTime{Time: time.Now()}
	before := types.CpsTime{Time: now.Add(-d.Duration())}
	cleaner := pbehavior.NewCleaner(mongoClient, datastorage.BulkSize, w.Logger)
	maxUpdates := int64(w.DataStorageConfigProvider.Get().MaxUpdates)
	deleted, err := cleaner.Clean(ctx, before, maxUpdates)
	if err != nil {
		w.Logger.Err(err).Msg("cannot accumulate week statistics")
	} else if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("pbehaviors were deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryPbehavior(ctx, now)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}

	return nil
}
