package main

import (
	"testing"
	"time"

	mock_alarm "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/alarm"
	mock_idlealarm "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/idlealarm"
	mock_redis "git.canopsis.net/canopsis/go-engines/mocks/lib/redis"
	"github.com/bsm/redislock"
	"github.com/golang/mock/gomock"
)

func TestPeriodicalWorker_Work_GivenObtainedLock_ShouldDoWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockAlarmService := mock_alarm.NewMockService(ctrl)
	mockIdleAlarmService := mock_idlealarm.NewMockService(ctrl)
	interval := time.Minute
	worker := periodicalWorker{
		PeriodicalInterval: interval,
		LockerClient:       mockLockClient,
		AlarmService:       mockAlarmService,
		IdleAlarmService:   mockIdleAlarmService,
	}

	mockLockClient.EXPECT().
		Obtain(gomock.Eq(PeriodicalLockKey), gomock.Eq(interval), gomock.Any()).
		Return(nil, nil)

	mockAlarmService.EXPECT().ResolveAlarms(gomock.Any(), gomock.Any())
	mockAlarmService.EXPECT().ResolveSnoozes(gomock.Any(), false)
	mockAlarmService.EXPECT().ResolveCancels(gomock.Any(), gomock.Any())
	mockAlarmService.EXPECT().ResolveDone(gomock.Any())
	mockAlarmService.EXPECT().UpdateFlappingAlarms(gomock.Any())
	mockIdleAlarmService.EXPECT().Process()

	_ = worker.Work()
}

func TestPeriodicalWorker_Work_GivenNotObtainedLock_ShouldDoNotAnything(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockService := mock_alarm.NewMockService(ctrl)
	interval := time.Minute
	worker := periodicalWorker{
		PeriodicalInterval: interval,
		LockerClient:       mockLockClient,
		AlarmService:       mockService,
	}

	mockLockClient.EXPECT().
		Obtain(gomock.Eq(PeriodicalLockKey), gomock.Eq(interval), gomock.Any()).
		Return(nil, redislock.ErrNotObtained)

	mockService.EXPECT().ResolveAlarms(gomock.Any(), gomock.Any()).Times(0)

	_ = worker.Work()
}
