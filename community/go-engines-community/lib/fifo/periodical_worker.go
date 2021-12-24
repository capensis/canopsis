package fifo

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"time"
)

type periodicalWorker struct {
	RuleService        eventfilter.Service
	PeriodicalInterval time.Duration
}

func (w *periodicalWorker) GetInterval() time.Duration {
	return w.PeriodicalInterval
}

func (w *periodicalWorker) Work(ctx context.Context) error {
	err := w.RuleService.LoadRules(ctx, []string{eventfilter.RuleTypeChangeEntity})
	if err != nil {
		return fmt.Errorf("unable to load rules: %w", err)
	}

	return nil
}
