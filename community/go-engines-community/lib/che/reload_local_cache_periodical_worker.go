package che

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"github.com/rs/zerolog"
)

type reloadLocalCachePeriodicalWorker struct {
	EventFilterService eventfilter.Service
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
	LoadRules          bool
}

func (w *reloadLocalCachePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *reloadLocalCachePeriodicalWorker) Work(ctx context.Context) {
	if w.LoadRules {
		err := w.EventFilterService.LoadRules(ctx, []string{eventfilter.RuleTypeDrop, eventfilter.RuleTypeEnrichment, eventfilter.RuleTypeBreak})
		if err != nil {
			w.Logger.Error().Err(err).Msg("unable to load rules")
		}
	}
}
