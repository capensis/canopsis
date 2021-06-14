package main

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"runtime/trace"
	"time"
)

type periodicalWorker struct {
	LockClient           redis.LockClient
	EntityServiceService entityservice.Service
	PeriodicalInterval   time.Duration
	Logger               zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) error {
	ctx, task := trace.NewTask(ctx, "service.PeriodicalWorker")
	defer task.End()

	// Lock periodical, do not release lock to not allow another instance start periodical.
	_, err := w.LockClient.Obtain(ctx, periodicalLock, w.PeriodicalInterval, &redislock.Options{})
	if err != nil {
		if err == redislock.ErrNotObtained {
			return nil
		}

		w.Logger.Err(err).Msg("cannot obtain lock")

		return nil
	}

	w.Logger.Debug().Msg("Recompute entity services")
	err = w.EntityServiceService.ComputeAllServices(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing all entity services")
	}

	return nil
}
