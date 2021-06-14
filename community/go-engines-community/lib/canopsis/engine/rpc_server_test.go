package engine_test

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
	"testing"
)

func TestRpcServer_Consume_GivenMessage_ShouldProcessIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 1,
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	d := amqp.Delivery{Body: body}
	msgs := make(chan amqp.Delivery, 1)
	msgs <- d
	close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)
	mockChannel.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any())

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Eq(d)).Return(body, nil)

	err := consumer.Consume(context.Background())
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
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 1,
		mockConnection,
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
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(resultBody, nil)
	mockChannel.EXPECT().Publish(
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

	err := consumer.Consume(context.Background())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestRpcServer_Consume_GivenErrorOnMessage_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 1,
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: body}
	defer close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)
	mockChannel.EXPECT().Publish(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Times(0)
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any())

	expectedErr := errors.New("test err")
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	err := consumer.Consume(context.Background())
	if err != expectedErr {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestRpcServer_Consume_GivenContextDone_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewRPCServer(
		name, queue,
		1, 1,
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	msgs := make(chan amqp.Delivery, 1)
	defer close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Times(0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := consumer.Consume(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}
