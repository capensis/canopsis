package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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

func (w *resolvedArchiverWorker) Work(ctx context.Context) {
	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot retrieve data storage config")
		return
	}

	var lastExecuted datetime.CpsTime
	if conf.History.Alarm != nil {
		lastExecuted = conf.History.Alarm.Time
	}

	ok := datastorage.CanRun(lastExecuted, w.DataStorageConfigProvider.Get().TimeToExecute, w.TimezoneConfigProvider.Get().Location)
	if !ok {
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
	cleaner := alarm.NewCleaner(mongoClient, datastorage.BulkSize)
	maxUpdates := int64(w.DataStorageConfigProvider.Get().MaxUpdates)
	var archived, deleted int64
	updated := false
	archiveAfter := conf.Config.Alarm.ArchiveAfter
	if datetime.IsDurationEnabledAndValid(archiveAfter) {
		updated = true
		archived, err = cleaner.ArchiveResolvedAlarms(ctx, archiveAfter.SubFrom(now), maxUpdates)
		if err != nil {
			w.Logger.Err(err).Msg("cannot archive resolved alarms")
			return
		}

		if archived > 0 {
			w.Logger.Info().Int64("alarm_number", archived).Msg("resolved alarm archiving")
		}
	}

	deleteAfter := conf.Config.Alarm.DeleteAfter
	if datetime.IsDurationEnabledAndValid(deleteAfter) {
		updated = true
		deleted, err = cleaner.DeleteArchivedResolvedAlarms(ctx, deleteAfter.SubFrom(now), maxUpdates)
		if err != nil {
			w.Logger.Err(err).Msg("cannot delete resolved alarms")
			return
		}

		if deleted > 0 {
			w.Logger.Info().Int64("alarm_number", deleted).Msg("resolved alarm removing")
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
}
