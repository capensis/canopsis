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

func (w *cleanPeriodicalWorker) Work(ctx context.Context) {
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
		w.Logger.Err(err).Msg("fail to retrieve data storage config")
		return
	}
	// Skip if already executed today.
	if conf.History.Pbehavior != nil && conf.History.Pbehavior.EqualDay(now) {
		return
	}

	d := conf.Config.Pbehavior.DeleteAfter
	if d == nil || !*d.Enabled || d.Value == 0 {
		return
	}

	deleted, err := w.PbehaviorCleaner.Clean(ctx, d.SubFrom(now))
	if err != nil {
		w.Logger.Err(err).Msg("cannot accumulate week statistics")
	} else if deleted > 0 {
		w.Logger.Info().Int64("count", deleted).Msg("pbehaviors were deleted")
	}

	err = w.LimitConfigAdapter.UpdateHistoryPbehavior(ctx, now)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update config history")
	}
}
