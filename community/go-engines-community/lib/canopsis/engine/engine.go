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
		init:              init,
		deferFunc:         deferFunc,
		periodicalWorkers: make(map[string]PeriodicalWorker),
		logger:            logger,
	}
}

type engine struct {
	init              func(ctx context.Context) error
	deferFunc         func(ctx context.Context)
	consumers         []Consumer
	periodicalWorkers map[string]PeriodicalWorker
	routines          []Routine
	logger            zerolog.Logger
}

func (e *engine) AddConsumer(consumer Consumer) {
	e.consumers = append(e.consumers, consumer)
}

func (e *engine) AddPeriodicalWorker(name string, worker PeriodicalWorker) {
	if _, ok := e.periodicalWorkers[name]; ok {
		panic(fmt.Errorf("%q worker already exists", name))
	}

	e.periodicalWorkers[name] = worker
}

func (e *engine) AddRoutine(v Routine) {
	e.routines = append(e.routines, v)
}

func (e *engine) AddDeferFunc(deferFunc func(ctx context.Context)) {
	if deferFunc == nil {
		return
	}
	if e.deferFunc == nil {
		e.deferFunc = deferFunc
		return
	}

	prev := e.deferFunc
	e.deferFunc = func(ctx context.Context) {
		prev(ctx)
		deferFunc(ctx)
	}
}

func (e *engine) Run(ctx context.Context) error {
	e.logger.Info().
		Int("consumers", len(e.consumers)).
		Int("periodical workers", len(e.periodicalWorkers)).
		Int("routines", len(e.routines)).
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

	for k, v := range e.periodicalWorkers {
		name := k
		worker := v
		g.Go(func() error {
			return e.runPeriodicalWorker(ctx, name, worker)
		})
	}

	for _, r := range e.routines {
		routine := r
		g.Go(func() error {
			return routine(ctx)
		})
	}

	return g.Wait()
}

func (e *engine) runPeriodicalWorker(
	ctx context.Context,
	name string,
	worker PeriodicalWorker,
) (resErr error) {
	defer func() {
		if r := recover(); r != nil {
			var err error
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}

			e.logger.Err(err).Str("worker", name).Msgf("periodical worker recovered from panic\n%s\n", debug.Stack())
			resErr = fmt.Errorf("periodical worker %q recovered from panic: %w", name, err)
		}
	}()

	interval := worker.GetInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var start time.Time
	done := make(chan time.Duration, 1)
	defer close(done)

	var skip, prevSkip bool

	for {
		select {
		case <-ticker.C:
			prevSkip = skip
			skip = false
			var d time.Duration

			if !start.IsZero() {
				select {
				case d = <-done:
				default:
					skip = true
					e.logger.Error().
						Str("worker", name).
						Time("start", start).
						Str("spent time", time.Since(start).String()).
						Msg("previous run still in progress, skip periodical worker")
				}
			}

			if !skip {
				if prevSkip {
					e.logger.Info().
						Str("worker", name).
						Time("start", start).
						Str("spent time", d.String()).
						Msg("periodical worker continues to work properly")
				}

				start = time.Now()
				go func() {
					worker.Work(ctx)

					select {
					case <-ctx.Done():
						return
					default:
					}

					done <- time.Since(start)
				}()
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
