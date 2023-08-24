package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/axe/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmstatus"
	"github.com/rs/zerolog"
)

type reloadLocalCachePeriodicalWorker struct {
	PeriodicalInterval     time.Duration
	AlarmStatusService     alarmstatus.Service
	AutoInstructionMatcher event.AutoInstructionMatcher
	Logger                 zerolog.Logger
}

func (w *reloadLocalCachePeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *reloadLocalCachePeriodicalWorker) Work(ctx context.Context) {
	err := w.AutoInstructionMatcher.Load(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot load auto instructions")
	}

	err = w.AlarmStatusService.Load(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot load alarm status rules")
	}
}
