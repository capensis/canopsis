package main

import (
	"context"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	mock_pbehavior "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/pbehavior"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestRecomputeMessageProcessor_Process_GivenPbehaviorID_ShouldRecomputePbehavior(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"

	mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any()).Do(func(_ []byte, event *rpc.PbehaviorRecomputeEvent) {
		*event = rpc.PbehaviorRecomputeEvent{Ids: []string{pbehaviorID}}
	})

	mockService.EXPECT().RecomputeByIds(gomock.Any(), gomock.Eq([]string{pbehaviorID}))

	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
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

	p := recomputeMessageProcessor{
		PbhService:          mockService,
		PbehaviorCollection: mockPbhDbCollection,
		EntityCollection:    mockEntityDbCollection,
		EventManager:        mockEventManager,
		Encoder:             mockEncoder,
		Decoder:             mockDecoder,
		Publisher:           mockPublisher,
		Exchange:            "",
		Queue:               "test-queue",
		Logger:              zerolog.Nop(),
	}

	_, err := p.Process(ctx, amqp.Delivery{Body: []byte("{\"_id\":\"" + pbehaviorID + "\"}")})
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestRecomputeMessageProcessor_Process_GivenEmptyPbehaviorID_ShouldRecomputeAllPbehaviors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)

	mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any())

	mockService.EXPECT().Recompute(gomock.Any())

	p := recomputeMessageProcessor{
		PbhService:          mockService,
		PbehaviorCollection: mockPbhDbCollection,
		EntityCollection:    mockEntityDbCollection,
		EventManager:        mockEventManager,
		Encoder:             mockEncoder,
		Decoder:             mockDecoder,
		Publisher:           mockPublisher,
		Exchange:            "",
		Queue:               "test-queue",
		Logger:              zerolog.Nop(),
	}

	_, err := p.Process(ctx, amqp.Delivery{Body: []byte("{\"_id\":\"\"}")})
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}

func TestRecomputeMessageProcessor_Process_GivenPbehaviorID_ShouldSendPbehaviorEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mock_pbehavior.NewMockService(ctrl)
	mockPbhDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEntityDbCollection := mock_mongo.NewMockDbCollection(ctrl)
	mockEventManager := mock_pbehavior.NewMockEventManager(ctrl)
	mockDecoder := mock_encoding.NewMockDecoder(ctrl)
	mockEncoder := mock_encoding.NewMockEncoder(ctrl)
	mockPublisher := mock_amqp.NewMockPublisher(ctrl)
	pbehaviorID := "test-pbehavior-id"
	queue := "test-queue"
	entity := types.Entity{
		ID:        "test-entity-id",
		Type:      types.EntityTypeResource,
		Component: "test-component",
	}
	resolveResult := pbehavior.ResolveResult{
		Type: pbehavior.Type{
			ID:   "test-type-id",
			Type: pbehavior.TypeMaintenance,
		},
		ID: pbehaviorID,
	}
	event := types.Event{EventType: types.EventTypePbhEnter}
	body := []byte("test-body")

	mockDecoder.EXPECT().Decode(gomock.Any(), gomock.Any()).Do(func(_ []byte, event *rpc.PbehaviorRecomputeEvent) {
		*event = rpc.PbehaviorRecomputeEvent{Ids: []string{pbehaviorID}}
	})

	mockResolver := mock_pbehavior.NewMockComputedEntityTypeResolver(ctrl)
	mockService.EXPECT().RecomputeByIds(gomock.Any(), gomock.Any()).Return(mockResolver, nil)
	mockResolver.EXPECT().Resolve(gomock.Any(), gomock.Any(), gomock.Any()).Return(resolveResult, nil)

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
	mockEntityCursor.EXPECT().Next(gomock.Any()).Return(true)
	mockEntityCursor.EXPECT().Next(gomock.Any()).Return(false)
	mockEntityCursor.EXPECT().Decode(gomock.Any()).Do(func(e *types.Entity) {
		*e = entity
	}).Return(nil)
	mockEntityCursor.EXPECT().Close(gomock.Any())
	mockEntityCursor2 := mock_mongo.NewMockCursor(ctrl)
	mockEntityDbCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockEntityCursor2, nil)
	mockEntityCursor2.EXPECT().Next(gomock.Any()).Return(false)
	mockEntityCursor2.EXPECT().Close(gomock.Any())

	mockEventManager.EXPECT().GetEvent(gomock.Eq(resolveResult), gomock.Any(), gomock.Any()).
		Return(event, nil)

	mockEncoder.EXPECT().Encode(gomock.Any()).Return(body, nil)

	mockPublisher.EXPECT().
		PublishWithContext(gomock.Any(), gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}))

	p := recomputeMessageProcessor{
		PbhService:          mockService,
		PbehaviorCollection: mockPbhDbCollection,
		EntityCollection:    mockEntityDbCollection,
		EventManager:        mockEventManager,
		Encoder:             mockEncoder,
		Decoder:             mockDecoder,
		Publisher:           mockPublisher,
		Exchange:            "",
		Queue:               "test-queue",
		Logger:              zerolog.Nop(),
	}

	_, err := p.Process(ctx, amqp.Delivery{Body: []byte("{\"_id\":\"" + pbehaviorID + "\"}")})
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}
}
