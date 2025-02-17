package datastorage_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	libconfig "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_datastorage "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/datastorage"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestPeriodicalWorker_Work_GivenScheduledTimeOnToday_ShouldExecCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: make(map[string]datastorage.HistoryWithCount),
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockDbCollection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(2)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner1 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner1.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner2 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner2.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	gomock.InOrder(
		mockCleaner1.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner2.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test1", mockCleaner1)
	w.AddCleaner("test2", mockCleaner2)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenScheduledTimeOnTodayAndLastExecuteTimeOnToday_ShouldNotExecCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: map[string]datastorage.HistoryWithCount{
				"test": {
					Time: datetime.CpsTime{Time: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, location)},
				},
			},
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test", mockCleaner)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenScheduledTimeOnTodayAndDisabledCleaner_ShouldNotExecCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: make(map[string]datastorage.HistoryWithCount),
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner.EXPECT().IsEnabled(gomock.Any()).
		Return(false).
		AnyTimes()
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test", mockCleaner)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenNoScheduledTime_ShouldNotExecCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{}}).
		AnyTimes()
	mockCleaner := mock_datastorage.NewMockCleaner(ctrl)
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return nil, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test", mockCleaner)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenScheduledTimeOnAnotherDay_ShouldNotExecCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	weekday := now.Weekday()
	if weekday == time.Monday {
		weekday = time.Wednesday
	} else {
		weekday = time.Monday
	}
	scheduledTime := libconfig.ScheduledTime{
		Weekday: weekday,
		Hour:    now.Hour(),
	}
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner := mock_datastorage.NewMockCleaner(ctrl)
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return nil, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test", mockCleaner)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenInterruptedWorker_ShouldContinueFromTheOldestExecutedCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	twoWeekAgo := now.AddDate(0, 0, -14)
	threeWeekAgo := now.AddDate(0, 0, -21)
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: map[string]datastorage.HistoryWithCount{
				"test1": {
					Time: datetime.CpsTime{Time: time.Date(twoWeekAgo.Year(), twoWeekAgo.Month(), twoWeekAgo.Day(), twoWeekAgo.Hour(), 0, 0, 0, location)},
				},
				"test2": {
					Time: datetime.CpsTime{Time: time.Date(weekAgo.Year(), weekAgo.Month(), weekAgo.Day(), weekAgo.Hour(), 0, 0, 0, location)},
				},
				"test3": {
					Time: datetime.CpsTime{Time: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, location)},
				},
				"test4": {
					Time: datetime.CpsTime{Time: time.Date(threeWeekAgo.Year(), threeWeekAgo.Month(), threeWeekAgo.Day(), threeWeekAgo.Hour(), 0, 0, 0, location)},
				},
			},
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockDbCollection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(3)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner1 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner1.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner2 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner2.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner3 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner3.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner4 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner4.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	gomock.InOrder(
		mockCleaner4.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner1.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner2.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test1", mockCleaner1)
	w.AddCleaner("test2", mockCleaner2)
	w.AddCleaner("test3", mockCleaner3)
	w.AddCleaner("test4", mockCleaner4)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenInterruptedWorker_ShouldContinueFromNeverExecutedCleaner(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	weekAgo := now.AddDate(0, 0, -7)
	twoWeekAgo := now.AddDate(0, 0, -14)
	threeWeekAgo := now.AddDate(0, 0, -21)
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: map[string]datastorage.HistoryWithCount{
				"test1": {
					Time: datetime.CpsTime{Time: time.Date(twoWeekAgo.Year(), twoWeekAgo.Month(), twoWeekAgo.Day(), twoWeekAgo.Hour(), 0, 0, 0, location)},
				},
				"test2": {
					Time: datetime.CpsTime{Time: time.Date(weekAgo.Year(), weekAgo.Month(), weekAgo.Day(), weekAgo.Hour(), 0, 0, 0, location)},
				},
				"test3": {
					Time: datetime.CpsTime{Time: time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, location)},
				},
				"test4": {
					Time: datetime.CpsTime{Time: time.Date(threeWeekAgo.Year(), threeWeekAgo.Month(), threeWeekAgo.Day(), threeWeekAgo.Hour(), 0, 0, 0, location)},
				},
				"test6": {
					Time: datetime.CpsTime{Time: time.Date(threeWeekAgo.Year(), threeWeekAgo.Month(), threeWeekAgo.Day(), threeWeekAgo.Hour(), 0, 0, 0, location)},
				},
			},
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockDbCollection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(5)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{TimeToExecute: libconfig.ScheduledTimes{scheduledTime}}).
		AnyTimes()
	mockCleaner1 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner1.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner2 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner2.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner3 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner3.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner4 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner4.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner5 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner5.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner6 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner6.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	gomock.InOrder(
		mockCleaner5.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner6.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner1.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner2.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
		mockCleaner4.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()),
	)
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test1", mockCleaner1)
	w.AddCleaner("test2", mockCleaner2)
	w.AddCleaner("test3", mockCleaner3)
	w.AddCleaner("test4", mockCleaner4)
	w.AddCleaner("test5", mockCleaner5)
	w.AddCleaner("test6", mockCleaner6)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenCleanerErr_ShouldStopExec(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: make(map[string]datastorage.HistoryWithCount),
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockDbCollection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{
			TimeToExecute: libconfig.ScheduledTimes{scheduledTime},
			Timeout:       time.Hour,
		}).
		AnyTimes()
	mockCleaner1 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner1.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner1.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)
	mockCleaner2 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner2.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	mockCleaner2.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, _ mongo.DbClient, _ datastorage.Config, _ datetime.CpsTime, _ int) (datastorage.CleanResult, error) {
			<-ctx.Done()

			return datastorage.CleanResult{}, fmt.Errorf("test %w", ctx.Err())
		}).
		Times(1)
	mockCleaner3 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner3.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test1", mockCleaner1)
	w.AddCleaner("test2", mockCleaner2)
	w.AddCleaner("test3", mockCleaner3)
	w.Work(ctx)
}

func TestPeriodicalWorker_Work_GivenCleanerDeadlineErr_ShouldStopExecAndUpdateHistory(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	location := now.Location()
	scheduledTime := libconfig.ScheduledTime{
		Weekday: now.Weekday(),
		Hour:    now.Hour(),
	}
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockDbClient.EXPECT().Disconnect(gomock.Any())
	mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Eq(mongo.ConfigurationMongoCollection)).Return(mockDbCollection)
	mockSingleRes := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockSingleRes.EXPECT().Decode(gomock.Any()).Do(func(v *datastorage.DataStorage) {
		*v = datastorage.DataStorage{
			History: make(map[string]datastorage.HistoryWithCount),
		}
	})
	mockDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any()).Return(mockSingleRes)
	mockDbCollection.EXPECT().
		UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(2)
	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().
		Return(libconfig.TimezoneConfig{Location: location}).
		AnyTimes()
	mockScheduleConfigProvider := mock_config.NewMockDataStorageConfigProvider(ctrl)
	mockScheduleConfigProvider.EXPECT().Get().
		Return(libconfig.DataStorageConfig{
			TimeToExecute: libconfig.ScheduledTimes{scheduledTime},
			Timeout:       time.Millisecond,
		}).
		AnyTimes()
	mockCleaner1 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner1.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner1.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)
	mockCleaner2 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner2.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	mockCleaner2.EXPECT().Clean(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, _ mongo.DbClient, _ datastorage.Config, _ datetime.CpsTime, _ int) (datastorage.CleanResult, error) {
			<-ctx.Done()

			return datastorage.CleanResult{}, fmt.Errorf("test %w", ctx.Err())
		}).
		Times(1)
	mockCleaner3 := mock_datastorage.NewMockCleaner(ctrl)
	mockCleaner3.EXPECT().IsEnabled(gomock.Any()).
		Return(true).
		AnyTimes()
	w := datastorage.NewPeriodicalWorker(func(_ context.Context, _ time.Duration) (mongo.DbClient, error) {
		return mockDbClient, nil
	}, time.Hour, mockTimezoneConfigProvider, mockScheduleConfigProvider, zerolog.Nop())
	w.AddCleaner("test1", mockCleaner1)
	w.AddCleaner("test2", mockCleaner2)
	w.AddCleaner("test3", mockCleaner3)
	w.Work(ctx)
}
