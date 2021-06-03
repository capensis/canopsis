package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"time"
)

const PeriodicalLockKey = "action-periodical-lock-key"

type periodicalWorker struct {
	PeriodicalInterval    time.Duration
	LockerClient          redis.LockClient
	ActionService         action.Service
	ActionScenarioStorage action.ScenarioStorage
	Logger                zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) error {
	err := w.ActionScenarioStorage.ReloadScenarios()
	if err != nil {
		w.Logger.Error().Err(err).Msg("Periodical process: failed to reload actions")
		return nil
	}

	_, err = w.LockerClient.Obtain(ctx, PeriodicalLockKey, w.GetInterval(), nil)
	if err == redislock.ErrNotObtained {
		w.Logger.Debug().Msg("Could not obtain lock! Skip periodical process")
		return nil
	} else if err != nil {
		w.Logger.Error().Err(err).Msg("Obtain redis lock: unexpected error")
		return nil
	}

	err = w.ActionService.ProcessAbandonedExecutions(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("Periodical process: failed to process abandoned scenarios.")
		return nil
	}

	return nil
}