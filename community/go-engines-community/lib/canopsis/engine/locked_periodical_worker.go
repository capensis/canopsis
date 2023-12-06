package engine

import (
	"context"
	"errors"
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

func (w *lockedPeriodicalWorker) Work(ctx context.Context) {
	ttl := w.GetInterval()
	if ttl > ttlDiff {
		ttl -= ttlDiff
	}

	// Lock periodical, do not release lock to not allow another instance start periodical.
	_, err := w.lockClient.Obtain(ctx, w.lockKey, ttl, &redislock.Options{})
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			w.logger.Debug().Msg("lock already obtained")
			return
		}

		w.logger.Err(err).Msg("cannot obtain lock")
		return
	}

	w.worker.Work(ctx)
}
