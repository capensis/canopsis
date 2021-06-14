package engine

import (
	"context"
	"github.com/rs/zerolog"
	"time"
)

func NewRunInfoPeriodicalWorker(
	periodicalInterval time.Duration,
	manager RunInfoManager,
	info RunInfo,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &runInfoPeriodicalWorker{
		periodicalInterval: periodicalInterval,
		manager:            manager,
		info:               info,
		logger:             logger,
	}
}

type runInfoPeriodicalWorker struct {
	periodicalInterval time.Duration
	manager            RunInfoManager
	info               RunInfo
	logger             zerolog.Logger
}

func (w *runInfoPeriodicalWorker) GetInterval() time.Duration {
	return w.periodicalInterval
}

func (w *runInfoPeriodicalWorker) Work(ctx context.Context) error {
	err := w.manager.Save(ctx, w.info, w.GetInterval())
	if err != nil {
		w.logger.Error().Err(err).Msg("cannot save run info")
	}

	return nil
}
