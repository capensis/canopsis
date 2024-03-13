package action

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"github.com/rs/zerolog"
)

type scenarioPeriodicalWorker struct {
	PeriodicalInterval time.Duration
	ActionService      action.Service
	Logger             zerolog.Logger
}

func (w *scenarioPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *scenarioPeriodicalWorker) Work(ctx context.Context) {
	err := w.ActionService.ProcessAbandonedExecutions(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("failed to process abandoned scenarios")
	}
}
