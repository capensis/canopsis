package scheduler

import (
	"context"
	"sort"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/keymutex"
	redismod "github.com/go-redis/redis/v8"
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
	// LockMultipleOrPush tries to lock all lockIDList and lockID
	// and pushes item to the end of queue by lockID if fails.
	LockMultipleOrPush(ctx context.Context, lockIDList []string, lockID string, item []byte) (bool, error)
	// ExtendAndPopMultiple tries to expire lockID and pops item from lockID queue.
	// If next item exists it tries to lock lockIDList.
	// Arg getLockIDList retrieves lockIDList from next item.
	ExtendAndPopMultiple(ctx context.Context, lockID string, getLockIDList func([]byte) ([]string, error)) ([]byte, error)
	// ExtendAndPopRelatedOrMultiple tries to expire lockID and lockIDList and pops item from lockIDList queues.
	// If at least one next item from lockIDList exists it returns all next events from lockIDList
	// and unlocks lockIDList without events.
	// If there aren't next items it pops next item from lockID. If lockID queue is empty
	// it unlocks all locks.
	ExtendAndPopRelatedOrMultiple(ctx context.Context, lockIDList []string, lockID string) ([][]byte, error)
	// PopOrUnlock tries to extend lock lockID and pops item from queue by lockID.
	// It unlocks lockID if either fails.
	PopOrUnlock(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error)
	// LockAndPop tries to lock lockID and pops item from queue by lockID.
	LockAndPop(ctx context.Context, lockID string, asyncUnlock bool) ([]byte, error)
	// IsLocked returns true if lock lockID is set.
	IsLocked(ctx context.Context, lockID string) bool
	// IsEmpty returns true if queue lockID is empty.
	IsEmpty(ctx context.Context, lockID string) bool
	Unlock(ctx context.Context, lockID ...string) error
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

func (s *baseQueueLock) LockMultipleOrPush(
	ctx context.Context,
	lockIDList []string,
	lockID string,
	item []byte,
) (bool, error) {
	allLockIDList := append([]string{lockID}, lockIDList...)
	// Sort to prevent deadlock
	sort.Strings(allLockIDList)
	s.mutex.LockMultiple(allLockIDList...)

	defer func() {
		err := s.mutex.UnlockMultiple(allLockIDList...)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	/**
	The LockMultipleOrPush function is typically used in the metaalarm context,
	if we try to lock by all keys(metaalarm + children) and fail, then we place metaalarm event in the queue,
	but in that case we don't have metaalarm key in the lock storage, so event won't be released by ttl expiration.

	So we need to lock by metaalarm id first and then try to lock by children. With that way we'll have metaalarm key in the lock storage,
	so an event will be released from the queue by the ttl expiration if something goes wrong.
	*/
	locked, err := s.lock(ctx, lockID)
	if err != nil {
		return false, err
	}
	if !locked {
		return false, s.push(ctx, lockID, item)
	}

	if len(lockIDList) == 0 {
		return true, nil
	}

	locked, err = s.lockMultiple(ctx, lockIDList)
	if err != nil {
		return false, err
	}

	if !locked {
		return false, s.push(ctx, lockID, item)
	}

	return true, nil
}

func (s *baseQueueLock) ExtendAndPopMultiple(
	ctx context.Context,
	lockID string,
	f func([]byte) ([]string, error),
) (res []byte, resErr error) {
	s.mutex.Lock(lockID)
	var extended bool
	var err error

	defer func() {
		err := s.mutex.Unlock(lockID)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	/**
	The ExtendAndPopMultiple function is typically used in the metaalarm context,
	since metaalarm leaves a lock after itself, we should try to extend it.
	If success, then there is an event in the queue. We can try to pop it and lock children.
	*/
	extended, err = s.extendLock(ctx, lockID)
	if !extended || err != nil {
		return nil, err
	}

	nextItem, err := s.pop(ctx, lockID)
	if nextItem == nil || err != nil {
		return nil, err
	}

	lockIDList, err := f(nextItem)
	if len(lockIDList) == 0 || err != nil {
		return nil, err
	}

	// Sort to prevent deadlock
	sort.Strings(lockIDList)
	s.mutex.LockMultiple(lockIDList...)

	defer func() {
		err := s.mutex.UnlockMultiple(lockIDList...)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	lockedMultiple, err := s.lockMultiple(ctx, lockIDList)
	if err != nil {
		return nil, err
	}

	if !lockedMultiple {
		err := s.unshift(ctx, lockID, nextItem)
		return nil, err
	}

	return nextItem, nil
}

func (s *baseQueueLock) ExtendAndPopRelatedOrMultiple(
	ctx context.Context,
	lockIDList []string,
	lockID string,
) (res [][]byte, resErr error) {
	s.mutex.Lock(lockID)
	var extended bool
	var err error

	defer func() {
		err := s.mutex.Unlock(lockID)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	/**
	The ExtendAndPopRelatedOrMultiple function is typically used in the metaalarm context,
	since metaalarm leaves a lock after itself, we should try to extend it.
	If success, then there is an event in the queue. We can try to pop it and lock children.
	*/
	extended, err = s.extendLock(ctx, lockID)
	if !extended || err != nil {
		return nil, err
	}

	// Sort to prevent deadlock
	sort.Strings(lockIDList)
	s.mutex.LockMultiple(lockIDList...)

	defer func() {
		err := s.mutex.UnlockMultiple(lockIDList...)
		if err != nil {
			s.logger.Err(err).Msg("cannot unlock mutex")
		}
	}()

	events := make([][]byte, 0)
	noEvents := make([]string, 0)
	for _, relatedLockID := range lockIDList {
		extended, err = s.extendLock(ctx, relatedLockID)
		if err != nil {
			return nil, err
		}
		if extended {
			event, err := s.pop(ctx, relatedLockID)
			if err != nil {
				return nil, err
			}
			if event == nil {
				noEvents = append(noEvents, relatedLockID)
			} else {
				events = append(events, event)
			}
		}
	}

	if len(events) > 0 {
		if len(noEvents) > 0 {
			err = s.Unlock(ctx, noEvents...)
			if err != nil {
				s.logger.Err(err).Strs("lockID", noEvents).Msg("error on unlocking queue lock")
			}
		}

		return events, nil
	}

	nextItem, err := s.pop(ctx, lockID)
	if err != nil {
		return nil, err
	}

	if nextItem == nil {
		allLockIDList := append([]string{lockID}, lockIDList...)
		err = s.Unlock(ctx, allLockIDList...)
		if err != nil {
			s.logger.Err(err).Strs("lockID", allLockIDList).Msg("error on unlocking queue lock")
		}

		return nil, nil
	}

	return [][]byte{nextItem}, nil
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

				err := s.Unlock(ctx, lockID)
				if err != nil {
					s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
				}
			}()
		} else {
			err := s.Unlock(ctx, lockID)
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

				err := s.Unlock(ctx, lockID)
				if err != nil {
					s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
				}
			}()
		} else {
			err := s.Unlock(ctx, lockID)
			if err != nil {
				s.logger.Err(err).Str(lockID, "lockID").Msg("error on unlocking queue lock")
			}
		}
	}

	return nextItem, nil
}

func (s *baseQueueLock) IsLocked(ctx context.Context, lockID string) bool {
	result := s.lockClient.Exists(ctx, lockID)

	return result.Val() > 0
}

func (s *baseQueueLock) IsEmpty(ctx context.Context, lockID string) bool {
	result := s.queueClient.Exists(ctx, lockID)

	return result.Val() == 0
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

func (s *baseQueueLock) Unlock(ctx context.Context, lockID ...string) error {
	result := s.lockClient.Del(ctx, lockID...)

	return result.Err()
}

func (s *baseQueueLock) lockMultiple(ctx context.Context, lockIDList []string) (bool, error) {
	values := make(map[string]interface{}, len(lockIDList))

	for _, v := range lockIDList {
		values[v] = defaultLockValue
	}

	result := s.lockClient.MSetNX(ctx, values)

	if err := result.Err(); err != nil {
		return false, err
	}

	if !result.Val() {
		return false, nil
	}

	// Set timeout for each lock since MSETNX sets lock without expiration
	for _, v := range lockIDList {
		for i := 0; i < multipleLockSetExpireRetries; i++ {
			res := s.lockClient.Expire(ctx, v, s.lockExpirationTime)
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

func (s *baseQueueLock) unshift(ctx context.Context, lockID string, item []byte) error {
	result := s.queueClient.LPush(ctx, lockID, item)

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
