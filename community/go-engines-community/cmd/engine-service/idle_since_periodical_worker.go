package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"runtime/trace"
	"time"
)

type idleSincePeriodicalWorker struct {
	LockClient           redis.LockClient
	EntityServiceService entityservice.Service
	PeriodicalInterval   time.Duration
	Logger               zerolog.Logger
}

func (w *idleSincePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *idleSincePeriodicalWorker) Work(parenCtx context.Context) error {
	ctx, task := trace.NewTask(parenCtx, "service.IdleSincePeriodicalWorker")
	defer task.End()

	// Lock periodical, do not release lock to not allow another instance start periodical.
	_, err := w.LockClient.Obtain(ctx, periodicalIdleSinceLock, w.PeriodicalInterval, &redislock.Options{})
	if err != nil {
		if err == redislock.ErrNotObtained {
			return nil
		}

		w.Logger.Err(err).Msg("cannot obtain lock")

		return nil
	}

	w.Logger.Debug().Msg("Recompute idle_since")
	err = w.EntityServiceService.RecomputeIdleSince(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing idle_since")
	}
	w.Logger.Debug().Msg("Recompute idle_since finished")

	return nil
}
