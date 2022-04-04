package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"github.com/rs/zerolog"
	"runtime/trace"
	"time"
)

type recomputeAllPeriodicalWorker struct {
	EntityServiceService entityservice.Service
	PeriodicalInterval   time.Duration
	Logger               zerolog.Logger
}

func (w *recomputeAllPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *recomputeAllPeriodicalWorker) Work(parentCtx context.Context) {
	ctx, task := trace.NewTask(parentCtx, "service.PeriodicalWorker")
	defer task.End()

	w.Logger.Debug().Msg("Recompute entity services")
	err := w.EntityServiceService.ComputeAllServices(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing all entity services")
	}
}
