package che

import (
	"context"
	libcontext "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"github.com/rs/zerolog"
	"time"
)

type reloadLocalCachePeriodicalWorker struct {
	EventFilterService eventfilter.Service
	EnrichmentCenter   libcontext.EnrichmentCenter
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *reloadLocalCachePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *reloadLocalCachePeriodicalWorker) Work(ctx context.Context) {
	err := w.EventFilterService.LoadRules(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load rules")
	}

	err = w.EnrichmentCenter.LoadServices(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load services")
	}
}