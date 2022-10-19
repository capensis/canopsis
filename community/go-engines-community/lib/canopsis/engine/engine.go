package engine

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

func New(
	init func(ctx context.Context) error,
	deferFunc func(),
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
	deferFunc         func()
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
	e.logger.Info().Msg("engine started")
	defer e.logger.Info().Msg("engine stopped")
	defer func() {
		if e.deferFunc != nil {
			e.deferFunc()
		}
	}()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if e.init != nil {
		err := e.init(ctx)
		if err != nil {
			e.logger.Err(err).Msg("cannot init engine")
			return err
		}
	}

	wg := &sync.WaitGroup{}
	exitCh := make(chan error, len(e.consumers)+len(e.periodicalWorkers))
	defer close(exitCh)

	for _, c := range e.consumers {
		wg.Add(1)
		go func(consumer Consumer) {
			e.logger.Debug().Msg("consumer started")
			defer e.logger.Debug().Msg("consumer stopped")

			defer wg.Done()

			if r := recover(); r != nil {
				e.logger.Error().Msgf("worker recovered from panic: %v", r)
				debug.PrintStack()
				exitCh <- fmt.Errorf("from panic: %v", r)
			}

			err := consumer.Consume(ctx)
			if err != nil {
				exitCh <- err
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
	e.logger.Debug().Msg("periodical process started")
	defer e.logger.Debug().Msg("periodical process stopped")

	defer wg.Done()
	defer func() {
		if r := recover(); r != nil {
			e.logger.Error().Msgf("periodical recovered from panic: %v", r)
			debug.PrintStack()
			exitCh <- fmt.Errorf("from panic: %v", r)
		}
	}()

	ticker := time.NewTicker(worker.GetInterval())
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
				e.logger.Err(err).Msg("periodical process has been failed")
				exitCh <- err
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
