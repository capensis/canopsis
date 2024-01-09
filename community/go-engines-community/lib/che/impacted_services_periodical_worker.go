package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/contextgraph"
	"github.com/rs/zerolog"
)

type impactedServicesPeriodicalWorker struct {
	Manager            contextgraph.Manager
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *impactedServicesPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *impactedServicesPeriodicalWorker) Work(ctx context.Context) {
	err := w.Manager.UpdateImpactedServicesFromDependencies(ctx)
	if err != nil {
		w.Logger.Warn().Err(err).Msg("error while recomputing impacted services for connectors")
	}
}
