package fifo

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"github.com/rs/zerolog"
)

type periodicalWorker struct {
	RuleService        eventfilter.Service
	PeriodicalInterval time.Duration
	Logger             zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) {
	err := w.RuleService.LoadRules(ctx, []string{eventfilter.RuleTypeChangeEntity})
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load rules")
	}
}
