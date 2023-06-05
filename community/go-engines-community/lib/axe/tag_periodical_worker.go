package axe

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarmtag"
	"github.com/rs/zerolog"
)

type tagPeriodicalWorker struct {
	PeriodicalInterval time.Duration
	TagUpdater         alarmtag.Updater
	Logger             zerolog.Logger
}

func (w *tagPeriodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *tagPeriodicalWorker) Work(ctx context.Context) {
	err := w.TagUpdater.Update(ctx)
	if err != nil {
		w.Logger.Err(err).Msg("cannot update alarm tags")
	}
}
