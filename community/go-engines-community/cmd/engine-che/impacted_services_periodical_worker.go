package main

import (
	"context"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"github.com/rs/zerolog"
	"time"
)

type impactedServicesPeriodicalWorker struct {
	EnrichmentCenter   libcontext.EnrichmentCenter
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *impactedServicesPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *impactedServicesPeriodicalWorker) Work(ctx context.Context) error {
	w.Logger.Debug().Msg("Recompute impacted services for connectors")
	err := w.EnrichmentCenter.UpdateImpactedServices(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
	}
	w.Logger.Debug().Msg("Recompute impacted services for connectors finished")

	return nil
}
