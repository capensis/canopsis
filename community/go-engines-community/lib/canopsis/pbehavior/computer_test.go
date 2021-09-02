package pbehavior_test

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/timespan"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	mock_redis "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/redis"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestCancelableComputer_Compute_GivenPbehaviorID_ShouldRecomputePbehavior(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockLock := mock_redis.NewMockLock(ctrl)
	mockStore := mock_redis.NewMockStore(ctrl)
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"

	mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockLock, nil)
	mockLock.EXPECT().Release(gomock.Any())

	mockStore.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(true, nil)
	mockStore.EXPECT().Save(gomock.Any(), gomock.Any())

	mockService.EXPECT().Recompute(gomock.Any(), gomock.Eq(pbehaviorID))

	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(mongo.PbehaviorMongoCollection).Return(mockPbhDbCollection)
	mockSingleResultHelper := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockPbhDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSingleResultHelper)
	mockSingleResultHelper.EXPECT().Decode(gomock.Any()).Do(func(pbh *pbehavior.PBehavior) {
		pbh.Filter = "{\"name\":\"test-name\"}"
	}).Return(nil)
	mockDbClient.EXPECT().Collection(mongo.EntityMongoCollection).Return(mockEntityDbCollection)
	mockEntityDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockCursor, nil)
	mockCursor.EXPECT().Next(gomock.Any()).Return(false)

	computer := pbehavior.NewCancelableComputer(mockLockClient, mockStore, mockService,
		mockDbClient, mockPublisher, mockEventManager, mockEncoder, "test-queue",
		zerolog.Logger{})

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{
		PbehaviorID:   pbehaviorID,
		OperationType: pbehavior.OperationCreate,
	}
	defer close(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	computer.Compute(ctx, ch)
}

func TestCancelableComputer_Compute_GivenEmptyPbehaviorID_ShouldRecomputeAllPbehaviors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockLock := mock_redis.NewMockLock(ctrl)
	mockStore := mock_redis.NewMockStore(ctrl)
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	span := timespan.New(time.Now(), time.Now().Add(time.Hour))

	mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockLock, nil)
	mockLock.EXPECT().Release(gomock.Any())

	mockStore.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(true, nil)
	mockStore.EXPECT().Save(gomock.Any(), gomock.Any())

	mockService.EXPECT().GetSpan().Return(span)
	mockService.EXPECT().Compute(gomock.Any(), gomock.Eq(span))

	computer := pbehavior.NewCancelableComputer(mockLockClient, mockStore, mockService,
		mockDbClient, mockPublisher, mockEventManager, mockEncoder, "test-queue",
		zerolog.Logger{})

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{}
	defer close(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	computer.Compute(ctx, ch)
}

func TestCancelableComputer_Compute_GivenPbehaviorIDAndOperationType_ShouldSendPbehaviorEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLockClient := mock_redis.NewMockLockClient(ctrl)
	mockLock := mock_redis.NewMockLock(ctrl)
	mockStore := mock_redis.NewMockStore(ctrl)
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"
	queue := "test-queue"
	alarm := types.AlarmWithEntity{
		Alarm:  types.Alarm{},
		Entity: types.Entity{},
	}
	resolveResult := pbehavior.ResolveResult{
		ResolvedType: &pbehavior.Type{
			ID:   "test-type-id",
			Type: pbehavior.TypeMaintenance,
		},
		ResolvedPbhID: pbehaviorID,
	}
	event := types.Event{EventType: types.EventTypePbhEnter}
	body := []byte("test-body")

	mockLockClient.EXPECT().Obtain(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockLock, nil)
	mockLock.EXPECT().Release(gomock.Any())

	mockStore.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(true, nil)
	mockStore.EXPECT().Save(gomock.Any(), gomock.Any())

	mockService.EXPECT().Recompute(gomock.Any(), gomock.Any())
	mockService.EXPECT().Resolve(gomock.Any(), gomock.Any(), gomock.Any()).Return(resolveResult, nil)

	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.PbehaviorMongoCollection:
			return mockPbhDbCollection
		case mongo.EntityMongoCollection:
			return mockEntityDbCollection
		}

		t.Errorf("uknown collection")
		return nil
	}).AnyTimes()
	mockPbhSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockPbhDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPbhSingleResult)
	mockPbhSingleResult.EXPECT().Decode(gomock.Any()).Do(func(pbh *pbehavior.PBehavior) {
		*pbh = pbehavior.PBehavior{Filter: "{\"name\":\"test-name\"}"}
	}).Return(nil)
	mockEntityCursor := mock_mongo.NewMockCursor(ctrl)
	mockEntityDbCollection.EXPECT().Aggregate(gomock.Any(), gomock.Any()).Return(mockEntityCursor, nil)
	firstCall := mockEntityCursor.EXPECT().Next(gomock.Any()).Return(true)
	secondCall := mockEntityCursor.EXPECT().Next(gomock.Any()).Return(false)
	gomock.InOrder(firstCall, secondCall)
	mockEntityCursor.EXPECT().Decode(gomock.Any()).Do(func(a *types.AlarmWithEntity) {
		*a = alarm
	}).Return(nil)

	mockEventManager.EXPECT().GetEvent(gomock.Eq(resolveResult), gomock.Any(), gomock.Any()).
		Return(event)

	mockEncoder.EXPECT().Encode(gomock.Eq(event)).Return(body, nil)

	mockPublisher.EXPECT().
		Publish(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}))

	computer := pbehavior.NewCancelableComputer(mockLockClient, mockStore, mockService,
		mockDbClient, mockPublisher, mockEventManager, mockEncoder, queue,
		zerolog.Logger{})

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{PbehaviorID: pbehaviorID, OperationType: pbehavior.OperationCreate}
	defer close(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	computer.Compute(ctx, ch)
}
