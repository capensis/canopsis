package scheduler_test

import (
	"bytes"
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	mock_v8 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/go-redis/redis/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

func TestQueueLock_LockOrPush_GivenLockIsNotSet_ShouldSetLock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	lockExpirationTime := time.Second
	queueLock := scheduler.NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Eq(1), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	locked, err := queueLock.LockOrPush(ctx, lockID, item)

	if !locked {
		t.Error("expected returns true")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockOrPush_GivenLockIsNotSet_ShouldNotAddItemToQueue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Any(), gomock.Any()).
		Times(0)

	_, _ = queueLock.LockOrPush(ctx, lockID, item)
}

func TestQueueLock_LockOrPush_GivenLockIsSet_ShouldAddItemToQueue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Any(), gomock.Eq(lockID), gomock.Eq(item)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	locked, err := queueLock.LockOrPush(ctx, lockID, item)

	if locked {
		t.Error("expected returns false")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsNotEmpty_ShouldReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	lockExpirationTime := time.Second
	queueLock := scheduler.NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	expectedItem := make([]byte, 1)

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Eq(lockID), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult(string(expectedItem), nil))

	item, err := queueLock.PopOrUnlock(ctx, lockID, false)

	if !bytes.Equal(item, expectedItem) {
		t.Errorf("expected item: %v but got %v", expectedItem, item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsNotSet_ShouldNotReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any()).
		Times(0)

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Any()).
		Times(0)

	item, err := queueLock.PopOrUnlock(ctx, lockID, false)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsEmpty_ShouldNotReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any(), gomock.Any()).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult("", redis.Nil))

	item, err := queueLock.PopOrUnlock(ctx, lockID, false)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_PopOrUnlock_GivenLockIsSetAndQueueIsEmpty_ShouldDeleteLock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Any()).
		Return(redis.NewStringResult("", redis.Nil))

	_, _ = queueLock.PopOrUnlock(ctx, lockID, false)
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsNotEmpty_ShouldReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	lockExpirationTime := time.Second
	queueLock := scheduler.NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		logger,
	)
	lockID := "testlock"
	expectedItem := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Eq(1), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult(string(expectedItem), nil))

	item, err := queueLock.LockAndPop(ctx, lockID, false)

	if !bytes.Equal(item, expectedItem) {
		t.Errorf("expected item: %v but got %v", expectedItem, item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsSet_ShouldNotReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Any()).
		Times(0)

	item, err := queueLock.LockAndPop(ctx, lockID, false)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsEmpty_ShouldNotReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any(), gomock.Any()).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult("", redis.Nil))

	item, err := queueLock.LockAndPop(ctx, lockID, false)

	if item != nil {
		t.Errorf("expected item: nil but got %v", item)
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsEmpty_ShouldDeleteLock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.Nop()
	queueLock := scheduler.NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		logger,
	)
	lockID := "testlock"

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Del(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Any()).
		Return(redis.NewStringResult("", redis.Nil))

	_, _ = queueLock.LockAndPop(ctx, lockID, false)
}
