package entityservice_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_v8 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/go-redis/redis/v8"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_entity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entity"
	mock_entityservice "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entityservice"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestService_Process_GivenEvent_ShouldUpdateServices(t *testing.T) {
	ctx := context.Background()
	serviceID := "test-service"
	services := []entityservice.ServiceData{
		{
			ID:             serviceID,
			OutputTemplate: "test-service-output",
		},
	}
	alarm := types.Alarm{
		ID:       "test-alarm",
		EntityID: "test-entity",
		Value: types.AlarmValue{
			State: &types.AlarmStep{Value: types.AlarmStateCritical},
		},
	}
	serviceIncrements := entityservice.AlarmCounters{
		State:             entityservice.StateCounters{Critical: 1, Major: -1},
		PbehaviorCounters: map[string]int64{},
	}
	newServiceCounters := entityservice.AlarmCounters{
		All:             10,
		Active:          10,
		NotAcknowledged: 10,
		State:           entityservice.StateCounters{Critical: 1, Major: 9},
	}
	eventBody := []byte("test-event")
	pubExchangeName, pubQueueName := "test-exchange", "test-queue"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAmqpPublisher := mock_amqp.NewMockPublisher(ctrl)
	mockAmqpPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(eventBody, nil)
	mockAdapter := mock_entityservice.NewMockAdapter(ctrl)
	mockAdapter.EXPECT().GetDependenciesCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	mockAdapter.EXPECT().UpdateCounters(gomock.Any(), gomock.Eq(serviceID), gomock.Eq(newServiceCounters)).Return(nil)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockCountersCache := mock_entityservice.NewMockCountersCache(ctrl)
	mockCountersCache.EXPECT().
		RemoveAndGet(gomock.Any(), gomock.Eq(fmt.Sprintf("%s&&%s", serviceID, alarm.EntityID))).
		Return(nil, nil)
	mockCountersCache.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, m map[string]entityservice.AlarmCounters) {
			if counters, ok := m[serviceID]; ok {
				if !reflect.DeepEqual(counters, serviceIncrements) {
					t.Errorf("expected %v counters but got %v", serviceIncrements, counters)
				}
			} else {
				t.Errorf("expected service counters but got nothing")
			}
		}).
		Return(map[string]entityservice.AlarmCounters{
			serviceID: newServiceCounters,
		}, nil)
	mockStorage := mock_entityservice.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(ctx).Return(services, nil)
	mockRedisClient := mock_v8.NewMockCmdable(ctrl)
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockServiceLock := mock_redis.NewMockLock(ctrl)
	mockServiceUpdateLock := mock_redis.NewMockLock(ctrl)
	firstCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceLock, nil)
	secondCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-update-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceUpdateLock, nil)
	gomock.InOrder(firstCall, secondCall)
	mockServiceLock.EXPECT().Release(gomock.Any())
	mockServiceUpdateLock.EXPECT().Release(gomock.Any())

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	service := entityservice.NewService(
		mockAmqpPublisher,
		pubExchangeName,
		pubQueueName,
		mockEncoder,
		mockAdapter,
		mockEntityAdapter,
		mock_alarm.NewMockAdapter(ctrl),
		mockCountersCache,
		mockStorage,
		mockLockClient,
		mockRedisClient,
		tplExecutor,
		zerolog.Nop(),
	)

	event := types.Event{
		Component: alarm.EntityID,
		Entity:    &types.Entity{ID: alarm.EntityID, Services: []string{serviceID}},
		Alarm:     &alarm,
		AlarmChange: &types.AlarmChange{
			Type:          types.AlarmChangeTypeStateIncrease,
			PreviousState: types.AlarmStateMajor,
		},
	}
	err := service.Process(context.Background(), event)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestService_Process_GivenEventAndCachedAlarmCounters_ShouldUpdateServices(t *testing.T) {
	serviceID := "test-service"
	services := []entityservice.ServiceData{
		{
			ID:             serviceID,
			OutputTemplate: "test-service-output",
		},
	}
	alarm := types.Alarm{
		ID:       "test-alarm",
		EntityID: "test-entity",
		Value: types.AlarmValue{
			State: &types.AlarmStep{Value: types.AlarmStateCritical},
		},
	}
	serviceIncrements := entityservice.AlarmCounters{
		Acknowledged:      -1,
		NotAcknowledged:   1,
		State:             entityservice.StateCounters{Critical: 1, Major: -1},
		PbehaviorCounters: map[string]int64{},
	}
	newServiceCounters := entityservice.AlarmCounters{
		All:             10,
		Active:          10,
		NotAcknowledged: 10,
		State:           entityservice.StateCounters{Critical: 1, Major: 9},
	}
	eventBody := []byte("test-event")
	pubExchangeName, pubQueueName := "test-exchange", "test-queue"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAmqpPublisher := mock_amqp.NewMockPublisher(ctrl)
	mockAmqpPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(eventBody, nil)
	mockAdapter := mock_entityservice.NewMockAdapter(ctrl)
	mockAdapter.EXPECT().GetDependenciesCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	mockAdapter.EXPECT().
		UpdateCounters(gomock.Any(), gomock.Eq(serviceID), gomock.Eq(newServiceCounters)).
		Return(nil)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockCountersCache := mock_entityservice.NewMockCountersCache(ctrl)
	mockCountersCache.EXPECT().
		RemoveAndGet(gomock.Any(), gomock.Eq(fmt.Sprintf("%s&&%s", serviceID, alarm.EntityID))).
		Return(&entityservice.AlarmCounters{
			All:          1,
			Active:       1,
			State:        entityservice.StateCounters{Major: 1},
			Acknowledged: 1,
		}, nil)
	mockCountersCache.EXPECT().
		Update(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, m map[string]entityservice.AlarmCounters) {
			if counters, ok := m[serviceID]; ok {
				if !reflect.DeepEqual(counters, serviceIncrements) {
					t.Errorf("expected %v counters but got %v", serviceIncrements, counters)
				}
			} else {
				t.Errorf("expected service counters but got nothing")
			}
		}).
		Return(map[string]entityservice.AlarmCounters{
			serviceID: newServiceCounters,
		}, nil)
	mockStorage := mock_entityservice.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return(services, nil)
	mockRedisClient := mock_v8.NewMockCmdable(ctrl)
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockServiceLock := mock_redis.NewMockLock(ctrl)
	mockServiceUpdateLock := mock_redis.NewMockLock(ctrl)
	firstCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceLock, nil)
	secondCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-update-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceUpdateLock, nil)
	gomock.InOrder(firstCall, secondCall)
	mockServiceLock.EXPECT().Release(gomock.Any())
	mockServiceUpdateLock.EXPECT().Release(gomock.Any())

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	service := entityservice.NewService(
		mockAmqpPublisher,
		pubExchangeName,
		pubQueueName,
		mockEncoder,
		mockAdapter,
		mockEntityAdapter,
		mock_alarm.NewMockAdapter(ctrl),
		mockCountersCache,
		mockStorage,
		mockLockClient,
		mockRedisClient,
		tplExecutor,
		zerolog.Nop(),
	)

	event := types.Event{
		Component:   alarm.EntityID,
		Entity:      &types.Entity{ID: alarm.EntityID, Services: []string{serviceID}},
		Alarm:       &alarm,
		AlarmChange: &types.AlarmChange{},
	}
	err := service.Process(context.Background(), event)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestService_Process_GivenEventAndLockedService_ShouldSkipEvent(t *testing.T) {
	serviceID := "test-service"
	services := []entityservice.ServiceData{
		{
			ID:             serviceID,
			OutputTemplate: "test-service-output",
		},
	}
	alarm := types.Alarm{
		ID: "test-alarm",
		Value: types.AlarmValue{
			State: &types.AlarmStep{Value: types.AlarmStateCritical},
		},
	}
	pubExchangeName, pubQueueName := "test-exchange", "test-queue"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAmqpPublisher := mock_amqp.NewMockPublisher(ctrl)
	mockAmqpPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(0)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	resendEventBody := []byte("test-body")
	mockEncoder.EXPECT().Encode(gomock.Any()).
		Do(func(event types.Event) {
			if event.EventType != types.EventTypeAlarmSkipped {
				t.Errorf("expected event %v but got %v", types.EventTypeAlarmSkipped, event.EventType)
			}
		}).
		Return(resendEventBody, nil)
	mockAdapter := mock_entityservice.NewMockAdapter(ctrl)
	mockAdapter.EXPECT().UpdateCounters(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockCountersCache := mock_entityservice.NewMockCountersCache(ctrl)
	mockCountersCache.EXPECT().RemoveAndGet(gomock.Any(), gomock.Any()).Times(0)
	mockCountersCache.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
	mockStorage := mock_entityservice.NewMockStorage(ctrl)
	mockStorage.EXPECT().GetAll(gomock.Any()).Return(services, nil)
	mockRedisClient := mock_v8.NewMockCmdable(ctrl)
	mockRedisClient.EXPECT().
		HSetNX(gomock.Any(), gomock.Eq(fmt.Sprintf("skipped-entities-%s", serviceID)), gomock.Any(), gomock.Eq(resendEventBody)).
		Return(redis.NewBoolResult(true, nil))
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockServiceLock := mock_redis.NewMockLock(ctrl)
	firstCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceLock, nil)
	secondCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-update-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(nil, redislock.ErrNotObtained)
	gomock.InOrder(firstCall, secondCall)
	mockServiceLock.EXPECT().Release(gomock.Any())

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	service := entityservice.NewService(
		mockAmqpPublisher,
		pubExchangeName,
		pubQueueName,
		mockEncoder,
		mockAdapter,
		mockEntityAdapter,
		mock_alarm.NewMockAdapter(ctrl),
		mockCountersCache,
		mockStorage,
		mockLockClient,
		mockRedisClient,
		tplExecutor,
		zerolog.Nop(),
	)

	event := types.Event{
		Entity: &types.Entity{Services: []string{serviceID}},
		Alarm:  &alarm,
		AlarmChange: &types.AlarmChange{
			Type:          types.AlarmChangeTypeStateIncrease,
			PreviousState: types.AlarmStateMajor,
		},
	}
	err := service.Process(context.Background(), event)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestService_UpdateService_GivenEvent_ShouldUpdateService(t *testing.T) {
	serviceID := "test-service"
	serviceData := entityservice.ServiceData{
		OutputTemplate: "test-output",
	}
	entity := types.Entity{
		ID: "test-entity",
	}
	alarm := types.Alarm{
		ID:       "test-alarm",
		EntityID: entity.ID,
		Value: types.AlarmValue{
			State: &types.AlarmStep{Value: types.AlarmStateCritical},
		},
	}
	serviceIncrements := entityservice.AlarmCounters{
		All:             1,
		Active:          1,
		NotAcknowledged: 1,
		State:           entityservice.StateCounters{Critical: 1},
	}
	newServiceCounters := entityservice.AlarmCounters{
		All:               1,
		Active:            1,
		NotAcknowledged:   1,
		State:             entityservice.StateCounters{Critical: 1},
		PbehaviorCounters: map[string]int64{},
	}
	eventBody := []byte("test-event")
	resendEventBody := []byte("test-resend-event")
	pubExchangeName, pubQueueName := "test-exchange", "test-queue"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAmqpPublisher := mock_amqp.NewMockPublisher(ctrl)
	mockAmqpPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		Times(2)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Any()).Return(eventBody, nil)
	mockAdapter := mock_entityservice.NewMockAdapter(ctrl)
	mockAdapter.EXPECT().GetDependenciesCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	mockEntityAdapter := mock_entity.NewMockAdapter(ctrl)
	mockEntityCursor := mock_mongo.NewMockCursor(ctrl)
	gomock.InOrder(
		mockEntityCursor.EXPECT().Next(gomock.Any()).Return(true),
		mockEntityCursor.EXPECT().Next(gomock.Any()).Return(false),
	)

	mockEntityCursor.EXPECT().
		Decode(gomock.Any()).
		Do(func(v *types.Entity) {
			*v = entity
		}).
		Return(nil)
	mockEntityCursor.EXPECT().Close(gomock.Any())
	mockAdapter.EXPECT().GetServiceDependencies(gomock.Any(), gomock.Eq(serviceID)).Return(mockEntityCursor, nil)
	mockAdapter.EXPECT().UpdateCounters(gomock.Any(), gomock.Eq(serviceID), gomock.Eq(newServiceCounters)).Return(nil)
	mockCountersCache := mock_entityservice.NewMockCountersCache(ctrl)
	firstReplaceCall := mockCountersCache.EXPECT().
		Replace(gomock.Any(), gomock.Eq(fmt.Sprintf("%s&&%s", serviceID, alarm.EntityID)), gomock.Eq(serviceIncrements)).
		Return(nil)
	secondReplaceCall := mockCountersCache.EXPECT().
		Replace(gomock.Any(), gomock.Eq(serviceID), gomock.Eq(newServiceCounters)).
		Return(nil)
	gomock.InOrder(firstReplaceCall, secondReplaceCall)
	mockStorage := mock_entityservice.NewMockStorage(ctrl)
	mockStorage.EXPECT().Reload(gomock.Any(), gomock.Any()).Return(&serviceData, false, false, false, nil)
	mockRedisClient := mock_v8.NewMockCmdable(ctrl)
	mockRedisClient.EXPECT().
		HGetAll(gomock.Any(), gomock.Eq(fmt.Sprintf("skipped-entities-%s", serviceID))).
		Return(redis.NewStringStringMapResult(map[string]string{"entity-id": string(resendEventBody)}, nil))
	mockRedisClient.EXPECT().
		Del(gomock.Any(), gomock.Eq(fmt.Sprintf("skipped-entities-%s", serviceID))).
		Return(redis.NewIntResult(0, nil))
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockServiceLock := mock_redis.NewMockLock(ctrl)
	mockServiceUpdateLock := mock_redis.NewMockLock(ctrl)
	firstCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-update-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceUpdateLock, nil)
	secondCall := mockLockClient.EXPECT().
		Obtain(gomock.Any(), gomock.Eq(fmt.Sprintf("lock-service-%s", serviceID)), gomock.Any(), gomock.Any()).
		Return(mockServiceLock, nil).
		Times(2)
	gomock.InOrder(firstCall, secondCall)
	mockServiceLock.EXPECT().Release(gomock.Any()).Return(nil).Times(2)
	mockServiceUpdateLock.EXPECT().Release(gomock.Any()).Return(nil)

	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsByIDs(gomock.Any(), gomock.All(
		gomock.Len(1),
		gomock.AssignableToTypeOf([]string{}),
	), gomock.All(
		gomock.AssignableToTypeOf(&[]types.Alarm{}),
		gomock.Not(gomock.Nil()),
	)).Do(func(_ context.Context, ids []string, alarms *[]types.Alarm) {
		if !reflect.DeepEqual(ids, []string{entity.ID}) {
			t.Errorf("expected %v but got %v", []string{entity.ID}, ids)
		}

		*alarms = []types.Alarm{alarm}
	}).Return(nil)

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	service := entityservice.NewService(
		mockAmqpPublisher,
		pubExchangeName,
		pubQueueName,
		mockEncoder,
		mockAdapter,
		mockEntityAdapter,
		mockAlarmAdapter,
		mockCountersCache,
		mockStorage,
		mockLockClient,
		mockRedisClient,
		tplExecutor,
		zerolog.Nop(),
	)

	event := types.Event{
		Component: serviceID,
	}
	err := service.UpdateService(context.Background(), event)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}
