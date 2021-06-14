package main

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"github.com/rs/zerolog"
	"runtime/trace"
	"time"
)

type periodicalWorker struct {
	WatcherService 			watcher.Service
	PeriodicalInterval 		time.Duration
	Logger             		zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work() error {
	ctx, task := trace.NewTask(context.Background(), "watcher.PeriodicalWorker")
	defer task.End()

	w.Logger.Debug().Msg("Recompute watchers")
	err := w.WatcherService.ComputeAllWatchers(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing all watchers")
	}

	return nil
}
