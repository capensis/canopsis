package pbehavior_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestCancelableComputer_Compute_GivenPbehaviorID_ShouldRecomputePbehavior(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"

	mockService.EXPECT().RecomputeByIds(gomock.Any(), gomock.Eq([]string{pbehaviorID}))

	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockAlarmDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collection string) mongo.DbCollection {
		switch collection {
		case mongo.AlarmMongoCollection:
			return mockAlarmDbCollection
		case mongo.EntityMongoCollection:
			return mockEntityDbCollection
		case mongo.PbehaviorMongoCollection:
			return mockPbhDbCollection
		}
		t.Errorf("uknown collection")
		return nil
	}).AnyTimes()
	mockSingleResultHelper := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockCursor := mock_mongo.NewMockCursor(ctrl)
	mockPbhDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(mockSingleResultHelper)
	mockSingleResultHelper.EXPECT().Decode(gomock.Any()).Do(func(pbh *pbehavior.PBehavior) {
		pbh.EntityPattern = pattern.Entity{{
			{
				Field:     "name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
			},
		}}
	}).Return(nil)
	mockEntityDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockCursor, nil).Times(2)
	mockCursor.EXPECT().Next(gomock.Any()).Return(false).Times(2)
	mockCursor.EXPECT().Close(gomock.Any()).Times(2)

	computer := pbehavior.NewCancelableComputer(mockService, mockDbClient, mockPublisher, mockEventManager, mockDecoder,
		mockEncoder, "test-queue", zerolog.Nop())

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{
		PbehaviorIds: []string{pbehaviorID},
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
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockAlarmDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collection string) mongo.DbCollection {
		switch collection {
		case mongo.AlarmMongoCollection:
			return mockAlarmDbCollection
		case mongo.EntityMongoCollection:
			return mockEntityDbCollection
		case mongo.PbehaviorMongoCollection:
			return mockPbhDbCollection
		}
		t.Errorf("uknown collection")
		return nil
	}).AnyTimes()
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)

	mockService.EXPECT().Recompute(gomock.Any())

	computer := pbehavior.NewCancelableComputer(mockService, mockDbClient, mockPublisher, mockEventManager, mockDecoder,
		mockEncoder, "test-queue", zerolog.Nop())

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{}
	defer close(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	computer.Compute(ctx, ch)
}

func TestCancelableComputer_Compute_GivenPbehaviorID_ShouldSendPbehaviorEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockDbClient := mock_mongo.NewMockDbClient(ctrl)
	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockAlarmDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"
	queue := "test-queue"
	alarm := types.Alarm{Value: types.AlarmValue{
		Connector:     "test-connector",
		ConnectorName: "test-connector-name",
		Component:     "test-component",
		Resource:      "test-resource",
	}}
	entity := types.Entity{ID: "test-entity-id"}
	resolveResult := pbehavior.ResolveResult{
		ResolvedType: pbehavior.Type{
			ID:   "test-type-id",
			Type: pbehavior.TypeMaintenance,
		},
		ResolvedPbhID: pbehaviorID,
	}
	event := types.Event{EventType: types.EventTypePbhEnter}
	body := []byte("test-body")
	mockResolver := mock_pbehavior.NewMockComputedEntityTypeResolver(ctrl)
	mockService.EXPECT().RecomputeByIds(gomock.Any(), gomock.Any()).Return(mockResolver, nil)
	mockResolver.EXPECT().Resolve(gomock.Any(), gomock.Any(), gomock.Any()).Return(resolveResult, nil)

	mockDbClient.EXPECT().Collection(gomock.Any()).DoAndReturn(func(collectionName string) mongo.DbCollection {
		switch collectionName {
		case mongo.PbehaviorMongoCollection:
			return mockPbhDbCollection
		case mongo.AlarmMongoCollection:
			return mockAlarmDbCollection
		case mongo.EntityMongoCollection:
			return mockEntityDbCollection
		}

		t.Errorf("uknown collection")
		return nil
	}).AnyTimes()
	mockPbhSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockPbhDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockPbhSingleResult)
	mockPbhSingleResult.EXPECT().Decode(gomock.Any()).Do(func(pbh *pbehavior.PBehavior) {
		*pbh = pbehavior.PBehavior{}
		pbh.EntityPattern = pattern.Entity{{
			{
				Field:     "name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
			},
		}}
	}).Return(nil)
	mockEntityCursor := mock_mongo.NewMockCursor(ctrl)
	mockEntityDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockEntityCursor, nil)
	gomock.InOrder(
		mockEntityCursor.EXPECT().Next(gomock.Any()).Return(true),
		mockEntityCursor.EXPECT().Next(gomock.Any()).Return(false),
	)
	mockEntityCursor.EXPECT().Decode(gomock.Any()).Do(func(e *types.Entity) {
		*e = entity
	}).Return(nil)
	mockEntityCursor.EXPECT().Close(gomock.Any())
	mockEntityCursor2 := mock_mongo.NewMockCursor(ctrl)
	mockEntityDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockEntityCursor2, nil)
	mockEntityCursor2.EXPECT().Next(gomock.Any()).Return(false)
	mockEntityCursor2.EXPECT().Close(gomock.Any())
	mockAlarmSingleResult := mock_mongo.NewMockSingleResultHelper(ctrl)
	mockAlarmDbCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAlarmSingleResult)
	mockAlarmSingleResult.EXPECT().Decode(gomock.Any()).Do(func(a *types.Alarm) {
		*a = alarm
	}).Return(nil)

	mockEventManager.EXPECT().GetEventType(gomock.Eq(resolveResult), gomock.Any()).
		Return(event.EventType, event.Output)

	mockEncoder.EXPECT().Encode(gomock.Any()).Return(body, nil)

	mockPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}))

	computer := pbehavior.NewCancelableComputer(mockService, mockDbClient, mockPublisher, mockEventManager, mockDecoder,
		mockEncoder, queue, zerolog.Nop())

	ch := make(chan pbehavior.ComputeTask, 1)
	ch <- pbehavior.ComputeTask{PbehaviorIds: []string{pbehaviorID}}
	defer close(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	computer.Compute(ctx, ch)
}
