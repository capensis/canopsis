package engine_test

import (
	"context"
	"errors"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func TestDefaultConsumer_Consume_GivenMessage_ShouldProcessIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		"", "", "", "",
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

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Eq(d))

	err := consumer.Consume(context.Background())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestDefaultConsumer_Consume_GivenProcessedMessage_ShouldPublishResultMessageToNextQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	nextExchange := "test-next-exchange"
	nextQueue := "test-next-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		nextExchange, nextQueue, "", "",
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	resultBody := []byte("test-result-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: []byte("test-body")}
	close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(resultBody, nil)
	mockChannel.EXPECT().PublishWithContext(
		gomock.Any(),
		gomock.Eq(nextExchange),
		gomock.Eq(nextQueue),
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         resultBody,
			DeliveryMode: amqp.Persistent,
		}),
	)

	err := consumer.Consume(context.Background())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestDefaultConsumer_Consume_GivenProcessedMessageAndNoNextQueue_ShouldPublishMessageToFIFOQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	fifoExchange := "test-fifo-exchange"
	fifoQueue := "test-fifo-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		"", "", fifoExchange, fifoQueue,
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: body}
	close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockChannel.EXPECT().PublishWithContext(
		gomock.Any(),
		gomock.Eq(fifoExchange),
		gomock.Eq(fifoQueue),
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}),
	)

	err := consumer.Consume(context.Background())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestDefaultConsumer_Consume_GivenProcessedMessageAndNoNextMessage_ShouldPublishMessageToFIFOQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	nextExchange := "next-fifo-exchange"
	nextQueue := "next-fifo-queue"
	fifoExchange := "test-fifo-exchange"
	fifoQueue := "test-fifo-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		nextExchange, nextQueue, fifoExchange, fifoQueue,
		mockConnection,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: body}
	close(msgs)
	mockConnection.EXPECT().Channel().Return(mockChannel, nil).AnyTimes()
	mockChannel.EXPECT().Qos(gomock.Any(), gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any())
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Close().AnyTimes()
	mockChannel.EXPECT().Consume(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(msgs, nil)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockChannel.EXPECT().PublishWithContext(
		gomock.Any(),
		gomock.Eq(fifoExchange),
		gomock.Eq(fifoQueue),
		gomock.Any(),
		gomock.Any(),
		gomock.Eq(amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}),
	)

	err := consumer.Consume(context.Background())
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestDefaultConsumer_Consume_GivenErrorOnMessage_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		"", "", "", "",
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
	mockChannel.EXPECT().PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any()).Times(0)
	mockChannel.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockChannel.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any())

	expectedErr := &testError{msg: "test error"}
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	err := consumer.Consume(context.Background())
	testErr := &testError{}
	if !errors.As(err, &testErr) || testErr.Error() != expectedErr.Error() {
		t.Errorf("expected error %v but got %v", expectedErr, err)
	}
}

func TestDefaultConsumer_Consume_GivenContextDone_ShouldStopConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	mockConnection := mock_amqp.NewMockConnection(ctrl)
	mockChannel := mock_amqp.NewMockChannel(ctrl)
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewDefaultConsumer(
		name, queue,
		1, 1, false,
		"", "", "", "",
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

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
