package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/neweventfilter"
	"github.com/rs/zerolog"
	"time"
)

type periodicalWorker struct {
	RuleService neweventfilter.EventFilterService
	PeriodicalInterval time.Duration
	Logger      zerolog.Logger
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) error {
	w.Logger.Debug().Msg("Loading event filter rules")

	err := w.RuleService.LoadRules(ctx)
	if err != nil {
		w.Logger.Error().Err(err).Msg("unable to load rules")
		return err
	}

	return nil
}
