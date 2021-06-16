package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"time"
)

const axeResolvedArchiverPeriodicalLockKey = "axe-resolved-archiver-periodical-lock-key"

type resolvedArchiverWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	AlarmAdapter              alarm.Adapter
	LockerClient              redis.LockClient
	Logger                    zerolog.Logger
}

func (w *resolvedArchiverWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *resolvedArchiverWorker) Work(ctx context.Context) error {
	schedule := w.DataStorageConfigProvider.Get().TimeToExecute
	// Skip if schedule is not defined.
	if schedule == nil {
		return nil
	}
	// Check now = schedule.
	location := w.TimezoneConfigProvider.Get().Location
	now := time.Now().In(location)
	if now.Weekday() != schedule.Weekday || now.Hour() != schedule.Hour {
		return nil
	}

	// Do actions under lock.
	_, err := w.LockerClient.Obtain(ctx, axeResolvedArchiverPeriodicalLockKey, w.GetInterval(), nil)
	if err == redislock.ErrNotObtained {
		w.Logger.Debug().Msg("skip periodical process")
		return nil
	} else if err != nil {
		w.Logger.Error().Err(err).Msg("skip periodical process")
		return nil
	}

	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("fail to retrieve data storage config")
		return nil
	}
	//Skip if already executed today.
	dateFormat := "2006-01-02"
	if conf.History.Alarm != nil && conf.History.Alarm.Time.Time.Format(dateFormat) == now.Format(dateFormat) {
		return nil
	}

	var archived, deleted int64

	updated := false
	if conf.Config.Alarm.ArchiveAfter != nil && *conf.Config.Alarm.ArchiveAfter.Enabled {
		d := conf.Config.Alarm.ArchiveAfter.Duration()
		if d > 0 {
			updated = true
			archived, err = w.AlarmAdapter.ArchiveResolvedAlarms(ctx, d)
			if err != nil {
				w.Logger.Err(err).Msg("cannot archive resolved alarms")
				return err
			}

			w.Logger.Info().Int64("alarm number", archived).Msg("resolved alarm archiving")
		}
	}

	if conf.Config.Alarm.DeleteAfter != nil && *conf.Config.Alarm.DeleteAfter.Enabled {
		d := conf.Config.Alarm.DeleteAfter.Duration()
		if d > 0 {
			updated = true
			deleted, err = w.AlarmAdapter.DeleteArchivedResolvedAlarms(ctx, d)
			if err != nil {
				w.Logger.Err(err).Msg("cannot delete resolved alarms")
			}

			w.Logger.Info().Int64("alarm number", deleted).Msg("resolved alarm removing")
		}
	}

	if updated {
		err := w.LimitConfigAdapter.UpdateHistoryAlarm(ctx, datastorage.AlarmHistory{
			Time:           types.CpsTime{Time: now},
			AlarmsArchived: archived,
			AlarmsDeleted:  deleted,
		})
		if err != nil {
			w.Logger.Err(err).Msg("cannot update config history")
		}
	}

	return nil
}
