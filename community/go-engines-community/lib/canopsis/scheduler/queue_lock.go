package scheduler

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/keymutex"
	redismod "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// QueueLock interface is used to implement a lock to consistently process items
// for the same resource. Base implementation uses redis to set lock and store next items
// while current item is processing. Im-memory mutex is used to synchronize access to redis
// so it cannot be used in multi-instance app.
type QueueLock interface {
	// LockOrPush tries to lock lockID and pushes item to queue by lockID if fails.
	// Return true if locks or false if error or item is added to queue.
	LockOrPush(ctx context.Context, lockID string, item []byte) (bool, error)
	// PopOrUnlock tries to extend lock lockID and pops item from queue by lockID.
	// It unlocks lockID if either fails.
	PopOrUnlock(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error)
	// LockAndPop tries to lock lockID and pops item from queue by lockID.
	LockAndPop(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error)
}

const defaultLockValue = 1
const defaultLockExpirationTime = time.Second * 10

type baseQueueLock struct {
	// lockClient is used to set lock.
	lockClient         redismod.Cmdable
	lockExpirationTime time.Duration
	// queueClient is used to set queue.
	queueClient redismod.Cmdable
	// mutex is used to synchronize operations on lockClient and queueClient.
	mutex  keymutex.KeyMutex
	logger zerolog.Logger
}

// NewQueueLock creates lock.
func NewQueueLock(
	lockClient redismod.Cmdable,
	lockExpirationTime time.Duration,
	queueClient redismod.Cmdable,
	logger zerolog.Logger,
) QueueLock {
	if lockExpirationTime.Seconds() == 0 {
		lockExpirationTime = defaultLockExpirationTime
	}

	return &baseQueueLock{
		lockClient:         lockClient,
		lockExpirationTime: lockExpirationTime,
		queueClient:        queueClient,
		logger:             logger,
		mutex:              keymutex.New(),
	}
}

func (s *baseQueueLock) LockOrPush(ctx context.Context, lockID string, item []byte) (bool, error) {
	s.mutex.Lock(lockID)
	defer func() {
		err := s.mutex.Unlock(lockID)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	locked, err := s.lock(ctx, lockID)

	if err != nil {
		return false, err
	}

	if !locked {
		return false, s.push(ctx, lockID, item)
	}

	return true, nil
}

func (s *baseQueueLock) PopOrUnlock(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error) {
	s.mutex.Lock(lockID)
	unlock := false

	defer func() {
		if !unlock {
			err := s.mutex.Unlock(lockID)
			if err != nil {
				s.logger.Err(err).Msg("cannot unlock mutex")
			}
		}
	}()

	extended, err := s.extendLock(ctx, lockID)

	if !extended || err != nil {
		return nil, err
	}

	nextItem, err := s.pop(ctx, lockID)
	if err != nil {
		return nil, err
	}

	if nextItem == nil {
		// Unlock in another goroutine for performance.
		if asyncUnlock {
			unlock = true
			go func() {
				defer func() {
					err := s.mutex.Unlock(lockID)
					if err != nil {
						s.logger.Err(err).Msg("cannot unlock mutex")
					}
				}()

				err := s.unlock(ctx, lockID)
				if err != nil {
					s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
				}
			}()
		} else {
			err := s.unlock(ctx, lockID)
			if err != nil {
				s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
			}
		}
	}

	return nextItem, nil
}

func (s *baseQueueLock) LockAndPop(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error) {
	s.mutex.Lock(lockID)
	unlock := false

	defer func() {
		if !unlock {
			err := s.mutex.Unlock(lockID)
			if err != nil {
				s.logger.Err(err).Msg("cannot unlock mutex")
			}
		}
	}()

	locked, err := s.lock(ctx, lockID)

	if !locked || err != nil {
		return nil, err
	}

	nextItem, err := s.pop(ctx, lockID)
	if err != nil {
		return nil, err
	}

	if nextItem == nil {
		// Unlock in another goroutine for performance.
		if asyncUnlock {
			unlock = true
			go func() {
				defer func() {
					err := s.mutex.Unlock(lockID)
					if err != nil {
						s.logger.Err(err).Msg("cannot unlock mutex")
					}
				}()

				err := s.unlock(ctx, lockID)
				if err != nil {
					s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
				}
			}()
		} else {
			err := s.unlock(ctx, lockID)
			if err != nil {
				s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
			}
		}
	}

	return nextItem, nil
}

func (s *baseQueueLock) lock(ctx context.Context, lockID string) (bool, error) {
	result := s.lockClient.SetNX(ctx, lockID, defaultLockValue, s.lockExpirationTime)

	return result.Val(), result.Err()
}

func (s *baseQueueLock) extendLock(ctx context.Context, lockID string) (bool, error) {
	result := s.lockClient.Expire(ctx, lockID, s.lockExpirationTime)

	if err := result.Err(); err != nil {
		s.logger.Err(err).
			Str(lockID, "lockID").
			Msg("error on updating expiration time of the queue lock")

		return false, err
	}

	return result.Val(), nil
}

func (s *baseQueueLock) unlock(ctx context.Context, lockID ...string) error {
	result := s.lockClient.Del(ctx, lockID...)

	return result.Err()
}

func (s *baseQueueLock) push(ctx context.Context, lockID string, item []byte) error {
	result := s.queueClient.RPush(ctx, lockID, item)

	if err := result.Err(); err != nil {
		s.logger.Err(err).
			Str("lockID", lockID).
			Msg("error on pushing item to redis queue")

		return err
	}

	return nil
}

func (s *baseQueueLock) pop(ctx context.Context, lockID string) ([]byte, error) {
	result := s.queueClient.LPop(ctx, lockID)

	if err := result.Err(); err != nil {
		if errors.Is(err, redismod.Nil) {
			return nil, nil
		}

		s.logger.Err(err).
			Str("lockID", lockID).
			Msg("error on popping item from redis queue")

		return nil, err
	}

	return []byte(result.Val()), nil
}
