package engine

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"runtime/debug"
	"sync"
	"time"
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

func (e *engine) Run(parentCtx context.Context) error {
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

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if e.init != nil {
		err := e.init(ctx)
		if err != nil {
			return fmt.Errorf("cannot init engine: %w", err)
		}
	}

	wg := &sync.WaitGroup{}
	exitCh := make(chan error, len(e.consumers)+len(e.periodicalWorkers))
	defer close(exitCh)

	for _, c := range e.consumers {
		wg.Add(1)
		go func(consumer Consumer) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					var err error
					var ok bool
					if err, ok = r.(error); !ok {
						err = fmt.Errorf("%v", r)
					}

					e.logger.Err(err).Msgf("consumer recovered from panic\n%s\n", debug.Stack())
					exitCh <- fmt.Errorf("consumer recovered from panic: %w", err)
				}
			}()

			err := consumer.Consume(ctx)
			if err != nil {
				exitCh <- fmt.Errorf("consumer failed: %w", err)
			}
		}(c)
	}

	for _, w := range e.periodicalWorkers {
		wg.Add(1)
		go e.runPeriodicalWorker(ctx, wg, w, exitCh)
	}

	// Wait context done or error from goroutines
	var exitErr error
	select {
	case <-ctx.Done():
	case exitErr = <-exitCh:
		cancel()
	}

	// Wait goroutines finish
	wg.Wait()

	return exitErr
}

func (e *engine) runPeriodicalWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	worker PeriodicalWorker,
	exitCh chan<- error,
) {
	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			var err error
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("%v", r)
			}

			e.logger.Err(err).Msgf("periodical worker recovered from panic\n%s\n", debug.Stack())
			exitCh <- fmt.Errorf("periodical worker recovered from panic: %w", err)
		}
	}()

	interval := worker.GetInterval()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := worker.Work(ctx)
			if err != nil {
				exitCh <- fmt.Errorf("periodical worker failed: %w", err)
				return
			}

			newInterval := worker.GetInterval()
			if newInterval != interval {
				ticker.Stop()
				interval = newInterval
				ticker = time.NewTicker(interval)
			}
		case <-ctx.Done():
			return
		}
	}
}
