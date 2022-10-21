package engine

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const shutdownTimout = 5 * time.Second

func New(
	init func(ctx context.Context) error,
	deferFunc func(ctx context.Context),
	logger zerolog.Logger,
) Engine {
	return &engine{
		init:      init,
		deferFunc: deferFunc,
		logger:    logger,
	}
}

type engine struct {
	init              func(ctx context.Context) error
	deferFunc         func(ctx context.Context)
	consumers         []Consumer
	periodicalWorkers []PeriodicalWorker
	logger            zerolog.Logger
}

func (e *engine) AddConsumer(consumer Consumer) {
	e.consumers = append(e.consumers, consumer)
}

func (e *engine) AddPeriodicalWorker(worker PeriodicalWorker) {
	e.periodicalWorkers = append(e.periodicalWorkers, worker)
}

func (e *engine) Run(ctx context.Context) error {
	e.logger.Info().
		Int("consumers", len(e.consumers)).
		Int("periodical workers", len(e.periodicalWorkers)).
		Msg("engine started")
	defer e.logger.Info().Msg("engine stopped")
	defer func() {
		if e.deferFunc != nil {
			deferCtx, deferCancel := context.WithTimeout(context.Background(), shutdownTimout)
			defer deferCancel()
			e.deferFunc(deferCtx)
		}
	}()

	if e.init != nil {
		err := e.init(ctx)
		if err != nil {
			return fmt.Errorf("cannot init engine: %w", err)
		}
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, c := range e.consumers {
		consumer := c

		g.Go(func() (resErr error) {
			defer func() {
				if r := recover(); r != nil {
					var err error
					var ok bool
					if err, ok = r.(error); !ok {
						err = fmt.Errorf("%v", r)
					}

					e.logger.Err(err).Msgf("consumer recovered from panic\n%s\n", debug.Stack())
					resErr = fmt.Errorf("consumer recovered from panic: %w", err)
				}
			}()

			err := consumer.Consume(ctx)
			if err != nil {
				return fmt.Errorf("consumer failed: %w", err)
			}

			return nil
		})
	}

	for _, w := range e.periodicalWorkers {
		worker := w
		g.Go(func() error {
			return e.runPeriodicalWorker(ctx, worker)
		})
	}

	return g.Wait()
}

func (e *engine) runPeriodicalWorker(
	ctx context.Context,
	worker PeriodicalWorker,
) (resErr error) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}

			e.logger.Err(err).Msgf("periodical worker recovered from panic\n%s\n", debug.Stack())
			resErr = fmt.Errorf("periodical worker recovered from panic: %w", err)
		}
	}()

	interval := worker.GetInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			err := worker.Work(ctx)
			d := time.Since(start)
			if d > worker.GetInterval() {
				e.logger.Error().
					Time("start", start).
					Str("spent time", d.String()).
					Msgf("periodical worker %T run too long", worker)
			}
			if err != nil {
				return fmt.Errorf("periodical worker failed: %w", err)
			}

			newInterval := worker.GetInterval()
			if newInterval != interval {
				ticker.Stop()
				interval = newInterval
				ticker = time.NewTicker(interval)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
