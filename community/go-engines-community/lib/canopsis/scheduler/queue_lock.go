package scheduler

import (
	"sort"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/keymutex"
	redismod "github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
)

// QueueLock interface is used to implement a lock to consistently process items
// for the same resource. Base implementation uses redis to set lock and store next items
// while current item is processing. Im-memory mutex is used to synchronize access to redis
// so it cannot be used in multi-instance app.
type QueueLock interface {
	// LockOrPush tries to lock lockID and pushes item to queue by lockID if fails.
	// Return true if locks or false if error or item is added to queue.
	LockOrPush(lockID string, item []byte) (bool, error)
	// LockMultipleOrPush tries to lock all lockIDList
	// and pushes item to queue by lockID if fails.
	LockMultipleOrPush(lockIDList []string, lockID string, item []byte) (bool, error)
	// PopOrUnlock tries to extend lock lockID and pops item from queue by lockID.
	// It unlocks lockID if either fails.
	PopOrUnlock(lockID string) ([]byte, error)
	// LockAndPop tries to lock lockID and pops item from queue by lockID.
	LockAndPop(lockID string) ([]byte, error)
	// IsLocked returns true if lock lockID is set.
	IsLocked(lockID string) bool
	// IsEmpty returns true if queue lockID is empty.
	IsEmpty(lockID string) bool
}

const defaultLockValue = 1
const defaultLockExpirationTime = time.Second * 10
const multipleLockSetExpireRetries = 3
const multipleLockSetExpireRetryTimeout = 10 * time.Millisecond

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

func (s *baseQueueLock) LockOrPush(lockID string, item []byte) (bool, error) {
	s.mutex.Lock(lockID)
	defer s.mutex.Unlock(lockID)

	locked, err := s.lock(lockID)

	if err != nil {
		return false, err
	}

	if !locked {
		return false, s.push(lockID, item)
	}

	return true, nil
}

func (s *baseQueueLock) LockMultipleOrPush(
	lockIDList []string,
	lockID string,
	item []byte,
) (bool, error) {
	sort.Strings(lockIDList)
	s.mutex.LockMultiple(lockIDList...)

	defer func() {
		s.mutex.UnlockMultiple(lockIDList...)
	}()

	locked, err := s.lockMultiple(lockIDList)

	if err != nil {
		return false, err
	}

	if !locked {
		return false, s.push(lockID, item)
	}

	return true, nil
}

func (s *baseQueueLock) PopOrUnlock(lockID string) ([]byte, error) {
	s.mutex.Lock(lockID)
	unlock := false

	defer func() {
		if !unlock {
			s.mutex.Unlock(lockID)
		}
	}()

	extended, err := s.extendLock(lockID)

	if !extended || err != nil {
		return nil, err
	}

	nextItem, err := s.pop(lockID)

	if err != nil {
		return nil, err
	}

	if nextItem == nil {
		unlock = true
		// Unlock in another goroutine for performance.
		go func() {
			defer s.mutex.Unlock(lockID)

			s.unlock(lockID)
		}()
	}

	return nextItem, nil
}

func (s *baseQueueLock) LockAndPop(lockID string) ([]byte, error) {
	s.mutex.Lock(lockID)
	unlock := false

	defer func() {
		if !unlock {
			s.mutex.Unlock(lockID)
		}
	}()

	locked, err := s.lock(lockID)

	if !locked || err != nil {
		return nil, err
	}

	nextItem, err := s.pop(lockID)

	if err != nil {
		return nil, err
	}

	if nextItem == nil {
		unlock = true
		// Unlock in another goroutine for performance.
		go func() {
			defer s.mutex.Unlock(lockID)

			s.unlock(lockID)
		}()
	}

	return nextItem, nil
}

func (s *baseQueueLock) IsLocked(lockID string) bool {
	result := s.lockClient.Exists(lockID)

	return result.Val() > 0
}

func (s *baseQueueLock) IsEmpty(lockID string) bool {
	result := s.queueClient.Exists(lockID)

	return result.Val() == 0
}

func (s *baseQueueLock) lock(lockID string) (bool, error) {
	result := s.lockClient.SetNX(lockID, defaultLockValue, s.lockExpirationTime)

	return result.Val(), result.Err()
}

func (s *baseQueueLock) extendLock(lockID string) (bool, error) {
	result := s.lockClient.Expire(lockID, s.lockExpirationTime)

	if err := result.Err(); err != nil {
		s.logger.Err(err).
			Str(lockID, "lockID").
			Msg("error on updating expiration time of the queue lock")

		return false, err
	}

	return result.Val(), nil
}

func (s *baseQueueLock) unlock(lockID string) {
	result := s.lockClient.Del(lockID)

	if err := result.Err(); err != nil {
		s.logger.Err(err).
			Str(lockID, "lockID").
			Msg("error on unlocking queue lock")
	}
}

func (s *baseQueueLock) lockMultiple(lockIDList []string) (bool, error) {
	values := make(map[string]interface{}, len(lockIDList))

	for _, v := range lockIDList {
		values[v] = defaultLockValue
	}

	result := s.lockClient.MSetNX(values)

	if err := result.Err(); err != nil {
		return false, err
	}

	if !result.Val() {
		return false, nil
	}

	// Set timeout for each lock since MSETNX sets lock without expiration
	for _, v := range lockIDList {
		for i := 0; i < multipleLockSetExpireRetries; i++ {
			res := s.lockClient.Expire(v, s.lockExpirationTime)
			err := res.Err()

			if err == nil {
				break
			}

			if i == multipleLockSetExpireRetries-1 {
				return false, err
			}

			time.Sleep(multipleLockSetExpireRetryTimeout)
		}
	}

	return true, nil
}

func (s *baseQueueLock) push(lockID string, item []byte) error {
	result := s.queueClient.RPush(lockID, item)

	if err := result.Err(); err != nil {
		s.logger.Err(err).
			Str("lockID", lockID).
			Msg("error on pushing item to redis queue")

		return err
	}

	return nil
}

func (s *baseQueueLock) pop(lockID string) ([]byte, error) {
	result := s.queueClient.LPop(lockID)

	if err := result.Err(); err != nil {
		if err == redismod.Nil {
			return nil, nil
		}

		s.logger.Err(err).
			Str("lockID", lockID).
			Msg("error on popping item from redis queue")

		return nil, err
	}

	return []byte(result.Val()), nil
}
