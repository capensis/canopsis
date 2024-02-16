package engine_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/bsm/redislock"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestLockedPeriodicalWorker_Work_GivenObtainedLock_ShouldRunWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := time.Minute
	lockKey := "test-lock"
	mockLock := mock_redis.NewMockLock(ctrl)
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Eq(lockKey), gomock.Any(), gomock.Any()).
		Return(mockLock, nil)
	mockWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	mockWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	mockWorker.EXPECT().Work(gomock.Any())

	worker := engine.NewLockedPeriodicalWorker(mockLockClient, lockKey, mockWorker, zerolog.Nop())
	worker.Work(ctx)
}

func TestLockedPeriodicalWorker_Work_GivenNotObtainedLock_ShouldNotRunWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interval := time.Minute
	lockKey := "test-lock"
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Eq(lockKey), gomock.Any(), gomock.Any()).
		Return(nil, redislock.ErrNotObtained)
	mockWorker := mock_engine.NewMockPeriodicalWorker(ctrl)
	mockWorker.EXPECT().GetInterval().Return(interval).AnyTimes()
	mockWorker.EXPECT().Work(gomock.Any()).Times(0)

	worker := engine.NewLockedPeriodicalWorker(mockLockClient, lockKey, mockWorker, zerolog.Nop())
	worker.Work(ctx)
}
