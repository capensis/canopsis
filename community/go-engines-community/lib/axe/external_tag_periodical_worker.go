package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"github.com/rs/zerolog"
)

type externalTagPeriodicalWorker struct {
	PeriodicalInterval time.Duration
	ExternalTagUpdater alarmtag.ExternalUpdater
	Logger             zerolog.Logger
}

func (w *externalTagPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *externalTagPeriodicalWorker) Work(ctx context.Context) {
	err := w.ExternalTagUpdater.Update(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update alarm tags")
	}
}
