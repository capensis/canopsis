package main

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_idlealarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/idlealarm"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/bsm/redislock"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestPeriodicalWorker_Work_GivenObtainedLock_ShouldDoWork(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockAlarmService := mock_alarm.NewMockService(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockIdleAlarmService := mock_idlealarm.NewMockService(ctrl)
	mockAlarmConfigProvider := mock_config.NewMockAlarmConfigProvider(ctrl)
	interval := time.Minute
	worker := periodicalWorker{
		PeriodicalInterval:  interval,
		LockerClient:        mockLockClient,
		AlarmService:        mockAlarmService,
		AlarmAdapter:        mockAlarmAdapter,
		IdleAlarmService:    mockIdleAlarmService,
		AlarmConfigProvider: mockAlarmConfigProvider,
	}

	mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(PeriodicalLockKey), gomock.Eq(interval), gomock.Any()).
		Return(nil, nil)
	alarmConfig := config.AlarmConfig{
		FlappingFreqLimit:        1,
		FlappingInterval:         time.Second,
		StealthyInterval:         time.Second,
		BaggotTime:               time.Second,
		EnableLastEventDate:      true,
		CancelAutosolveDelay:     time.Second,
		OutputLength:             10,
		TimeToKeepResolvedAlarms: time.Second,
	}
	mockAlarmConfigProvider.EXPECT().Get().Return(alarmConfig)
	mockAlarmAdapter.EXPECT().DeleteResolvedAlarms(gomock.Any(), gomock.Any())
	mockAlarmService.EXPECT().ResolveAlarms(gomock.Any(), gomock.Eq(alarmConfig))
	mockAlarmService.EXPECT().ResolveSnoozes(gomock.Any(), gomock.Eq(alarmConfig))
	mockAlarmService.EXPECT().ResolveCancels(gomock.Any(), gomock.Eq(alarmConfig))
	mockAlarmService.EXPECT().ResolveDone(gomock.Any())
	mockAlarmService.EXPECT().UpdateFlappingAlarms(gomock.Any(), gomock.Eq(alarmConfig))
	mockIdleAlarmService.EXPECT().Process(gomock.Any())

	_ = worker.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenNotObtainedLock_ShouldDoNotAnything(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockService := mock_alarm.NewMockService(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)

	interval := time.Minute
	worker := periodicalWorker{
		PeriodicalInterval: interval,
		LockerClient:       mockLockClient,
		AlarmService:       mockService,
		AlarmAdapter:       mockAlarmAdapter,
	}

	mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(PeriodicalLockKey), gomock.Eq(interval), gomock.Any()).
		Return(nil, redislock.ErrNotObtained)

	mockService.EXPECT().ResolveAlarms(gomock.Any(), gomock.Any()).Times(0)

	_ = worker.Work(ctx)
}
