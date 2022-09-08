package scheduler_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/scheduler"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_v8 "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/github.com/go-redis/redis/v8"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestScheduler_ProcessEvent_GivenEventAndNoLock_ShouldPublishEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	event := types.Event{
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
		Resource:      "test-resource",
	}
	lockID := "test-resource/test-component"
	body := []byte("test-body")

	mockRedisLockStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisLockStorage.EXPECT().SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))
	mockRedisQueueStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockChannel.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	publishToQueue := "test-queue"
	lockTtl := 100
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Eq(event)).Return(body, nil)

	service := scheduler.NewSchedulerService(mockRedisLockStorage, mockRedisQueueStorage,
		mockChannel, publishToQueue, zerolog.Nop(), lockTtl, mockDecoder, mockEncoder)

	err := service.ProcessEvent(ctx, event)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestScheduler_ProcessEvent_GivenEventAndLock_ShouldAddEventToQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	event := types.Event{
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
		Resource:      "test-resource",
	}
	lockID := "test-resource/test-component"
	body := []byte("test-body")

	mockRedisLockStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisLockStorage.EXPECT().SetNX(gomock.Any(), gomock.Eq(lockID), gomock.Any(), gomock.Any()).
		Return(redis.NewBoolResult(false, nil))
	mockRedisQueueStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisQueueStorage.EXPECT().RPush(gomock.Any(), gomock.Eq(lockID), gomock.Any()).
		Return(redis.NewIntResult(1, nil))
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	publishToQueue := "test-queue"
	lockTtl := 100
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockEncoder.EXPECT().Encode(gomock.Eq(event)).Return(body, nil)

	service := scheduler.NewSchedulerService(mockRedisLockStorage, mockRedisQueueStorage,
		mockChannel, publishToQueue, zerolog.Nop(), lockTtl, mockDecoder, mockEncoder)

	err := service.ProcessEvent(ctx, event)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestScheduler_AckEvent_GivenEventAndNoNextEventInQueue_ShouldUnlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	event := types.Event{
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
		Resource:      "test-resource",
	}
	lockID := "test-resource/test-component"

	mockRedisLockStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisLockStorage.EXPECT().Expire(gomock.Any(), gomock.Eq(lockID), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))
	mockRedisLockStorage.EXPECT().Del(gomock.Any(), gomock.Eq(lockID)).
		Return(redis.NewIntResult(1, nil))
	mockRedisQueueStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisQueueStorage.EXPECT().LPop(gomock.Any(), gomock.Eq(lockID)).
		Return(redis.NewStringResult("", redis.Nil))
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	publishToQueue := "test-queue"
	lockTtl := 100
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)

	service := scheduler.NewSchedulerService(mockRedisLockStorage, mockRedisQueueStorage,
		mockChannel, publishToQueue, zerolog.Nop(), lockTtl, mockDecoder, mockEncoder)

	err := service.AckEvent(ctx, event)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	// Wait unlock in goroutine.
	time.Sleep(time.Millisecond * 10)
}

func TestScheduler_AckEvent_GivenEventAndNextEvent_ShouldPublishNextEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	event := types.Event{
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
		Resource:      "test-resource",
	}
	lockID := "test-resource/test-component"
	body := []byte("test-body")

	mockRedisLockStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisLockStorage.EXPECT().Expire(gomock.Any(), gomock.Eq(lockID), gomock.Any()).
		Return(redis.NewBoolResult(true, nil))
	mockRedisQueueStorage := mock_v8.NewMockUniversalClient(ctrl)
	mockRedisQueueStorage.EXPECT().LPop(gomock.Any(), gomock.Eq(lockID)).
		Return(redis.NewStringResult(string(body), nil))
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockChannel.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)
	publishToQueue := "test-queue"
	lockTtl := 100
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)

	service := scheduler.NewSchedulerService(mockRedisLockStorage, mockRedisQueueStorage,
		mockChannel, publishToQueue, zerolog.Nop(), lockTtl, mockDecoder, mockEncoder)

	err := service.AckEvent(ctx, event)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}
