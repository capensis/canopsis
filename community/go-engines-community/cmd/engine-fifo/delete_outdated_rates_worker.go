package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/ratelimit"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog"
	"time"
)

type deleteOutdatedRatesWorker struct {
	PeriodicalInterval        time.Duration
	TimezoneConfigProvider    config.TimezoneConfigProvider
	DataStorageConfigProvider config.DataStorageConfigProvider
	LimitConfigAdapter        datastorage.Adapter
	RateLimitAdapter          ratelimit.Adapter
	Logger                    zerolog.Logger
}

func (w *deleteOutdatedRatesWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *deleteOutdatedRatesWorker) Work(ctx context.Context) error {
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
	if conf.History.HealthCheck != nil && conf.History.HealthCheck.Time.Format(dateFormat) == now.Format(dateFormat) {
		return nil
	}

	updated := false
	d := conf.Config.HealthCheck.DeleteAfter
	if d != nil && *d.Enabled && d.Value > 0 {
		updated = true
		before := d.SubFrom(now)
		deleted, err := w.RateLimitAdapter.DeleteBefore(ctx, types.CpsTime{Time: before})
		if err != nil {
			w.Logger.Err(err).Msg("cannot delete message rates")
		} else if deleted > 0 {
			w.Logger.Info().Int64("count", deleted).Msg("message rates were deleted")
		}
	}

	if updated {
		err := w.LimitConfigAdapter.UpdateHistoryHealthCheck(ctx, types.CpsTime{Time: now})
		if err != nil {
			w.Logger.Err(err).Msg("cannot update config history")
		}
	}

	return nil
}
