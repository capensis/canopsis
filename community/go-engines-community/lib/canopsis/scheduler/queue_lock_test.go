package scheduler

import (
	"bytes"
	mock_v7 "git.canopsis.net/canopsis/go-engines/mocks/github.com/go-redis/redis/v7"
	"github.com/go-redis/redis/v7"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestQueueLock_LockOrPush_GivenLockIsNotSet_ShouldSetLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Eq(lockID), gomock.Eq(defaultLockValue), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	locked, err := queueLock.LockOrPush(lockID, item)

	if !locked {
		t.Error("expected returns true")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockOrPush_GivenLockIsNotSet_ShouldNotAddItemToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Any()).
		Times(0)

	_, _ = queueLock.LockOrPush(lockID, item)
}

func TestQueueLock_LockOrPush_GivenLockIsSet_ShouldAddItemToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Eq(lockID), gomock.Eq(item)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	locked, err := queueLock.LockOrPush(lockID, item)

	if locked {
		t.Error("expected returns false")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockMultipleOrPush_GivenLockIsNotSet_ShouldSetLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockList := map[string]interface{}{"testlock1": defaultLockValue, "testlock2": defaultLockValue, "testlock3": defaultLockValue}
	lockID := "testlock1"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Eq(lockList)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Eq(lockIDList[0]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Eq(lockIDList[1]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Eq(lockIDList[2]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	locked, err := queueLock.LockMultipleOrPush(lockIDList, lockID, item)

	if !locked {
		t.Error("expected returns true")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockMultipleOrPush_GivenLockIsNotSet_ShouldNotAddItemToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Any()).
		Times(0)

	_, _ = queueLock.LockMultipleOrPush(lockIDList, lockID, item)
}

func TestQueueLock_LockMultipleOrPush_GivenLockIsSet_ShouldAddItemToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		Times(0)

	queueClient.
		EXPECT().
		RPush(gomock.Eq(lockID), gomock.Eq(item)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	locked, err := queueLock.LockMultipleOrPush(lockIDList, lockID, item)

	if locked {
		t.Error("expected returns false")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsNotEmpty_ShouldReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	expectedItem := make([]byte, 1)

	lockClient.
		EXPECT().
		Expire(gomock.Eq(lockID), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult(string(expectedItem), nil))

	item, err := queueLock.PopOrUnlock(lockID)

	if 0 != bytes.Compare(item, expectedItem) {
		t.Errorf("expected item: %v but got %v", expectedItem, item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsNotSet_ShouldNotReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any()).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any()).
		Times(0)

	item, err := queueLock.PopOrUnlock(lockID)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsNotSet_ShouldDeleteLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	lockClient.
		EXPECT().
		Del(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	_, _ = queueLock.PopOrUnlock(lockID)
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsEmpty_ShouldNotReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any()).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult("", redis.Nil))

	item, err := queueLock.PopOrUnlock(lockID)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsEmpty_ShouldDeleteLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any()).
		Return(redis.NewStringResult("", redis.Nil))

	_, _ = queueLock.PopOrUnlock(lockID)
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsNotEmpty_ShouldReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	expectedItem := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Eq(lockID), gomock.Eq(defaultLockValue), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult(string(expectedItem), nil))

	item, err := queueLock.LockAndPop(lockID)

	if 0 != bytes.Compare(item, expectedItem) {
		t.Errorf("expected item: %v but got %v", expectedItem, item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsSet_ShouldNotReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any()).
		Times(0)

	item, err := queueLock.LockAndPop(lockID)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsEmpty_ShouldNotReturnNextItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any()).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult("", redis.Nil))

	item, err := queueLock.LockAndPop(lockID)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsEmpty_ShouldDeleteLock(t *testing.T) {
	ctrl := gomock.NewController(t)
	lockClient := mock_v7.NewMockCmdable(ctrl)
	queueClient := mock_v7.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any()).
		Return(redis.NewStringResult("", redis.Nil))

	_, _ = queueLock.LockAndPop(lockID)
}
