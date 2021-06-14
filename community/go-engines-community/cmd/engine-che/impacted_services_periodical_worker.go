package main

import (
	"context"
	libcontext "git.canopsis.net/canopsis/go-engines/lib/canopsis/context"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
	"time"
)

const impactedServicesWorkerLock = "impacted-services-worker-lock"

type impactedServicesPeriodicalWorker struct {
	LockClient         redis.LockClient
	EnrichmentCenter   libcontext.EnrichmentCenter
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *impactedServicesPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *impactedServicesPeriodicalWorker) Work(ctx context.Context) error {
	// Lock periodical, do not release lock to not allow another instance start periodical.
	_, err := w.LockClient.Obtain(ctx, impactedServicesWorkerLock, w.PeriodicalInterval, &redislock.Options{})
	if err != nil {
		if err == redislock.ErrNotObtained {
			return nil
		}

		w.Logger.Err(err).Msg("cannot obtain lock")

		return nil
	}

	w.Logger.Debug().Msg("Recompute impacted services for connectors")
	err = w.EnrichmentCenter.UpdateImpactedServices(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
	}
	w.Logger.Debug().Msg("Recompute impacted services for connectors finished")

	return nil
}
