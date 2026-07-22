package engine_test

import (
	"context"
	"errors"
	"testing"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestRpcServer_Consume_GivenMessage_ShouldProcessIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 100,
		10,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	replyTo := "test-reply"
	d := amqp.Delivery{
		Body:    body,
		ReplyTo: replyTo,
	}
	msgs := make(chan amqp.Delivery, 1)
	msgs <- d
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false)).Times(2)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)
	mockPubCh.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Eq(replyTo), gomock.Any(),
		gomock.Any(), gomock.Any())

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Eq(d)).Return(body, nil)

	err := consumer.Consume(t.Context())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestRpcServer_Consume_GivenProcessedMessage_ShouldPublishResultMessageToBackQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	backQueue := "test-back-queue"
	corrId := "test-corr-id"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 100,
		10,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	resultBody := []byte("test-result-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{
		Body:          []byte("test-body"),
		ReplyTo:       backQueue,
		CorrelationId: corrId,
	}
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false)).Times(2)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(resultBody, nil)
	mockPubCh.EXPECT().PublishWithContext(gomock.Any(),
		gomock.Eq(""),
		gomock.Eq(backQueue),
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			Body:          resultBody,
			DeliveryMode:  amqp.Persistent,
		}),
	)

	err := consumer.Consume(t.Context())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestRpcServer_Consume_GivenErrorOnMessage_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 100,
		10,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: body}
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockPubCh.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any())

	expectedErr := &testError{msg: "test error"}
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	err := consumer.Consume(t.Context())
	testErr := &testError{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestRpcServer_Consume_GivenContextDone_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 100,
		10,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	msgs := make(chan amqp.Delivery, 1)
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Times(0)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := consumer.Consume(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}
