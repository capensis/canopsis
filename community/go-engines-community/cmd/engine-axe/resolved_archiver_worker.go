package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/rs/zerolog"
)

type resolvedArchiverWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	Logger                    zerolog.Logger
}

func (w *resolvedArchiverWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *resolvedArchiverWorker) Work(ctx context.Context) error {
	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot retrieve data storage config")
		return nil
	}

	var lastExecuted types.CpsTime
	if conf.History.Alarm != nil {
		lastExecuted = conf.History.Alarm.Time
	}

	ok := datastorage.CanRun(lastExecuted, w.DataStorageConfigProvider.Get().TimeToExecute, w.TimezoneConfigProvider.Get().Location)
	if !ok {
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
	cleaner := alarm.NewCleaner(mongoClient, datastorage.BulkSize)
	maxUpdates := int64(w.DataStorageConfigProvider.Get().MaxUpdates)
	var archived, deleted int64
	updated := false
	archiveAfter := conf.Config.Alarm.ArchiveAfter
	if archiveAfter != nil && *archiveAfter.Enabled && archiveAfter.Seconds > 0 {
		updated = true
		before := types.CpsTime{Time: now.Add(-archiveAfter.Duration())}
		archived, err = cleaner.ArchiveResolvedAlarms(ctx, before, maxUpdates)
		if err != nil {
			w.Logger.Err(err).Msg("cannot archive resolved alarms")
			return nil
		}

		if archived > 0 {
			w.Logger.Info().Int64("alarm number", archived).Msg("resolved alarm archiving")
		}
	}

	deleteAfter := conf.Config.Alarm.DeleteAfter
	if deleteAfter != nil && *deleteAfter.Enabled && deleteAfter.Seconds > 0 {
		updated = true
		before := types.CpsTime{Time: now.Add(-deleteAfter.Duration())}
		deleted, err = cleaner.DeleteArchivedResolvedAlarms(ctx, before, maxUpdates)
		if err != nil {
			w.Logger.Err(err).Msg("cannot delete resolved alarms")
			return nil
		}

		if deleted > 0 {
			w.Logger.Info().Int64("alarm number", deleted).Msg("resolved alarm removing")
		}
	}

	if updated {
		err := w.LimitConfigAdapter.UpdateHistoryAlarm(ctx, datastorage.HistoryWithCount{
			Time:     now,
			Archived: archived,
			Deleted:  deleted,
		})
		if err != nil {
			w.Logger.Err(err).Msg("cannot update config history")
		}
	}

	return nil
}
