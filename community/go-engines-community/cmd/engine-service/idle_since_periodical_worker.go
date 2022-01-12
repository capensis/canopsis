package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"github.com/rs/zerolog"
	"runtime/trace"
	"time"
)

type idleSincePeriodicalWorker struct {
	EntityServiceService entityservice.Service
	PeriodicalInterval   time.Duration
	Logger               zerolog.Logger
}

func (w *idleSincePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *idleSincePeriodicalWorker) Work(parenCtx context.Context) {
	ctx, task := trace.NewTask(parenCtx, "service.IdleSincePeriodicalWorker")
	defer task.End()

	w.Logger.Debug().Msg("Recompute idle_since")
	err := w.EntityServiceService.RecomputeIdleSince(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing idle_since")
	}
	w.Logger.Debug().Msg("Recompute idle_since finished")
}
