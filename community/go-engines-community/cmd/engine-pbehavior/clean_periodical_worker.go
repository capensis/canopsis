package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	libpbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

type cleanPeriodicalWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	PbehaviorCleaner          libpbehavior.Cleaner
	Logger                    zerolog.Logger
}

func (w *cleanPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *cleanPeriodicalWorker) Work(ctx context.Context) error {
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

	conf, err := w.LimitConfigAdapter.Get(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("fail to retrieve data storage config")
		return nil
	}
	// Skip if already executed today.
	dateFormat := "2006-01-02"
	if conf.History.Pbehavior != nil && conf.History.Pbehavior.Time.Format(dateFormat) == now.Format(dateFormat) {
		return nil
	}

	if conf.Config.Pbehavior.DeleteAfter == nil || !*conf.Config.Pbehavior.DeleteAfter.Enabled {
		return nil
	}

	d := conf.Config.Pbehavior.DeleteAfter.Duration()
	if d == 0 {
		return nil
	}

	deleted, err := w.PbehaviorCleaner.Clean(ctx, d)
	if err != nil {
		w.Logger.Err(err).Msg("cannot accumulate week statistics")
	} else if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("pbehaviors were deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryPbehavior(ctx, types.CpsTime{Time: now})
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}

	return nil
}
