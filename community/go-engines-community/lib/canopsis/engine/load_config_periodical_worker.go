package engine

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/rs/zerolog"
)

func NewLoadConfigPeriodicalWorker(
	periodicalInterval time.Duration,
	adapter config.Adapter,
	logger zerolog.Logger,
	updaters ...config.Updater,
) PeriodicalWorker {
	return &loadConfigPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		adapter:            adapter,
		updaters:           updaters,
		logger:             logger,
	}
}

type loadConfigPeriodicalWorker struct {
	periodicalInterval time.Duration
	adapter            config.Adapter
	logger             zerolog.Logger
	updaters           []config.Updater
}

func (w *loadConfigPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *loadConfigPeriodicalWorker) Work(ctx context.Context) {
	cfg, err := w.adapter.GetConfig(ctx)
	if err != nil {
		w.logger.Err(err).Msgf("cannot load config")
		return
	}

	for _, updater := range w.updaters {
		updater.Update(cfg)
	}
}
