package scheduler

import (
	"bytes"
	"context"
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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		false,
		logger,
	)
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Eq(defaultLockValue), gomock.Eq(lockExpirationTime)).
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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

func TestQueueLock_LockMultipleOrPush_GivenLockIsNotSet_ShouldSetLock(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		false,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockList := map[string]interface{}{"testlock1": defaultLockValue, "testlock2": defaultLockValue, "testlock3": defaultLockValue}
	lockID := "testlock1"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Any(), gomock.Eq(lockList)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Eq(lockIDList[0]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Eq(lockIDList[1]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Eq(lockIDList[2]), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	locked, err := queueLock.LockMultipleOrPush(ctx, lockIDList, lockID, item)

	if !locked {
		t.Error("expected returns true")
	}

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestQueueLock_LockMultipleOrPush_GivenLockIsNotSet_ShouldNotAddItemToQueue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		RPush(gomock.Any(), gomock.Any()).
		Times(0)

	_, _ = queueLock.LockMultipleOrPush(ctx, lockIDList, lockID, item)
}

func TestQueueLock_LockMultipleOrPush_GivenLockIsSet_ShouldAddItemToQueue(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
		logger,
	)
	lockIDList := []string{"testlock1", "testlock2", "testlock3"}
	lockID := "testlock"
	item := make([]byte, 1)

	lockClient.
		EXPECT().
		MSetNX(gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))

	lockClient.
		EXPECT().
		Expire(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)

	queueClient.
		EXPECT().
		RPush(gomock.Any(), gomock.Eq(lockID), gomock.Eq(item)).
		Times(1).
		Return(redis.NewIntResult(1, nil))

	locked, err := queueLock.LockMultipleOrPush(ctx, lockIDList, lockID, item)

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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		false,
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

	item, err := queueLock.PopOrUnlock(ctx, lockID)

	if 0 != bytes.Compare(item, expectedItem) {
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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	item, err := queueLock.PopOrUnlock(ctx, lockID)

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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	item, err := queueLock.PopOrUnlock(ctx, lockID)

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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	_, _ = queueLock.PopOrUnlock(ctx, lockID)
}

func TestQueueLock_LockAndPop_GivenLockIsNotSetAndQueueIsNotEmpty_ShouldReturnNextItem(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	lockClient := mock_v8.NewMockCmdable(ctrl)
	queueClient := mock_v8.NewMockCmdable(ctrl)
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	lockExpirationTime := time.Second
	queueLock := NewQueueLock(
		lockClient,
		lockExpirationTime,
		queueClient,
		false,
		logger,
	)
	lockID := "testlock"
	expectedItem := make([]byte, 1)

	lockClient.
		EXPECT().
		SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Eq(defaultLockValue), gomock.Eq(lockExpirationTime)).
		Times(1).
		Return(redis.NewBoolResult(true, nil))

	queueClient.
		EXPECT().
		LPop(gomock.Any(), gomock.Eq(lockID)).
		Times(1).
		Return(redis.NewStringResult(string(expectedItem), nil))

	item, err := queueLock.LockAndPop(ctx, lockID)

	if 0 != bytes.Compare(item, expectedItem) {
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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	item, err := queueLock.LockAndPop(ctx, lockID)

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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	item, err := queueLock.LockAndPop(ctx, lockID)

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
	logger := zerolog.New(bytes.NewBuffer(make([]byte, 0)))
	queueLock := NewQueueLock(
		lockClient,
		time.Second,
		queueClient,
		false,
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

	_, _ = queueLock.LockAndPop(ctx, lockID)
}
