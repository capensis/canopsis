package engine

import (
	"context"
	"errors"
	"sync"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/bsm/redislock"
	"github.com/rs/zerolog"
)

const ttlDiff = 100 * time.Millisecond

func NewLockedPeriodicalWorker(
	lockClient redis.LockClient,
	lockKey string,
	worker PeriodicalWorker,
	logger zerolog.Logger,
) PeriodicalWorker {
	return &lockedPeriodicalWorker{
		lockClient: lockClient,
		lockKey:    lockKey,
		worker:     worker,
		logger:     logger,
	}
}

type lockedPeriodicalWorker struct {
	lockClient redis.LockClient
	lockKey    string
	worker     PeriodicalWorker
	logger     zerolog.Logger
}

func (w *lockedPeriodicalWorker) GetInterval() time.Duration {
	return w.worker.GetInterval()
}

func (w *lockedPeriodicalWorker) Work(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	ttl := w.GetInterval()
	if ttl > ttlDiff {
		ttl -= ttlDiff
	}

	lockEnd := time.Now().Add(ttl)
	lockOpts := &redislock.Options{}
	// Lock periodical, do not release lock to not allow another instance start periodical.
	l, err := w.lockClient.Obtain(ctx, w.lockKey, ttl, lockOpts)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			w.logger.Debug().Str("key", w.lockKey).Msg("lock already obtained")

			return
		}

		w.logger.Err(err).Str("key", w.lockKey).Msg("cannot obtain lock")

		return
	}

	lockRefreshed := false
	defer func() {
		if !lockRefreshed {
			return
		}

		d := time.Until(lockEnd)
		if d > 0 {
			err = l.Refresh(context.WithoutCancel(parentCtx), d, lockOpts)
			if err != nil {
				w.logger.Err(err).Str("key", w.lockKey).Msg("cannot refresh lock")
			}
		} else {
			err = l.Release(context.WithoutCancel(parentCtx))
			if err != nil {
				w.logger.Err(err).Str("key", w.lockKey).Msg("cannot release lock")
			}
		}
	}()

	wg := sync.WaitGroup{}

	done := make(chan struct{})
	wg.Go(func() {
		defer close(done)
		w.worker.Work(ctx)
	})

	wg.Go(func() {
		ticker := time.NewTicker(ttl / 2)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err = l.Refresh(ctx, ttl, lockOpts)
				if err != nil {
					w.logger.Err(err).Str("key", w.lockKey).Msg("cannot refresh lock, stop periodical worker")
					cancel()

					return
				}

				lockRefreshed = true
			}
		}
	})

	wg.Wait()
}
