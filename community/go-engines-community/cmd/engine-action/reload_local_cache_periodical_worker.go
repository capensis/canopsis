package main

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"github.com/rs/zerolog"
)

type reloadLocalCachePeriodicalWorker struct {
	PeriodicalInterval    time.Duration
	ActionScenarioStorage action.ScenarioStorage
	Logger                zerolog.Logger
}

func (w *reloadLocalCachePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *reloadLocalCachePeriodicalWorker) Work(ctx context.Context) {
	err := w.ActionScenarioStorage.ReloadScenarios(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("failed to reload action scenarios")
	}
}
