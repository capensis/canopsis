package axe

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

type resolvedArchiverWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	AlarmAdapter              alarm.Adapter
	Logger                    zerolog.Logger
}

func (w *resolvedArchiverWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *resolvedArchiverWorker) Work(ctx context.Context) {
	schedule := w.DataStorageConfigProvider.Get().TimeToExecute
	// Skip if schedule is not defined.
	if schedule == nil {
		return
	}
	// Check now = schedule.
	location := w.TimezoneConfigProvider.Get().Location
	now := types.NewCpsTime().In(location)
	if now.Weekday() != schedule.Weekday || now.Hour() != schedule.Hour {
		return
	}

	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot retrieve data storage config")
		return
	}
	//Skip if already executed today.
	if conf.History.Alarm != nil && conf.History.Alarm.Time.EqualDay(now) {
		return
	}

	var archived, deleted int64

	updated := false
	archiveAfter := conf.Config.Alarm.ArchiveAfter
	if archiveAfter != nil && *archiveAfter.Enabled && archiveAfter.Value > 0 {
		updated = true
		archived, err = w.AlarmAdapter.ArchiveResolvedAlarms(ctx, archiveAfter.SubFrom(now))
		if err != nil {
			w.Logger.Err(err).Msg("cannot archive resolved alarms")
			return
		}

		if archived > 0 {
			w.Logger.Info().Int64("alarm number", archived).Msg("resolved alarm archiving")
		}
	}

	deleteAfter := conf.Config.Alarm.DeleteAfter
	if deleteAfter != nil && *deleteAfter.Enabled && deleteAfter.Value > 0 {
		updated = true
		deleted, err = w.AlarmAdapter.DeleteArchivedResolvedAlarms(ctx, deleteAfter.SubFrom(now))
		if err != nil {
			w.Logger.Err(err).Msg("cannot delete resolved alarms")
			return
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
}
