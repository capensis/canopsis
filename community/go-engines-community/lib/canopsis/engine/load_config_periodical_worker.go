package engine

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/rs/zerolog"
	"time"
)

func NewLoadConfigPeriodicalWorker(
	periodicalInterval time.Duration,
	adapter config.Adapter,
	updater config.Updater,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &loadConfigPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		adapter:            adapter,
		updater:            updater,
		logger:             logger,
	}
}

type loadConfigPeriodicalWorker struct {
	periodicalInterval time.Duration
	adapter            config.Adapter
	logger             zerolog.Logger
	updater            config.Updater
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

	w.updater.Update(cfg)
}
