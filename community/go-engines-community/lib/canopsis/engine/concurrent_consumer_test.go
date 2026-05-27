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

func TestConcurrentConsumer_Consume_GivenMessage_ShouldProcessIt(t *testing.T) {
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
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		"", "", "", "", 10, false,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	body := []byte("test-body")
	d := amqp.Delivery{Body: body}
	msgs := make(chan amqp.Delivery, 1)
	msgs <- d
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Eq(d))

	err := consumer.Consume(t.Context())
	if !errors.Is(err, libamqp.ErrChannelClosed) {
		t.Errorf("expected error %v but got %v", libamqp.ErrChannelClosed, err)
	}
}

func TestConcurrentConsumer_Consume_GivenProcessedMessage_ShouldPublishResultMessageToNextQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	nextExchange := "test-next-exchange"
	nextQueue := "test-next-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		nextExchange, nextQueue, "", "", 10, false,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	resultBody := []byte("test-result-body")
	msgs := make(chan amqp.Delivery, 1)
	msgs <- amqp.Delivery{Body: []byte("test-body")}
	close(msgs)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(resultBody, nil)
	mockPubCh.EXPECT().PublishWithContext(
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

	err := consumer.Consume(t.Context())
	if !errors.Is(err, libamqp.ErrChannelClosed) {
		t.Errorf("expected error %v but got %v", libamqp.ErrChannelClosed, err)
	}
}

func TestConcurrentConsumer_Consume_GivenProcessedMessageAndNoNextQueue_ShouldPublishMessageToFIFOQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	fifoExchange := "test-fifo-exchange"
	fifoQueue := "test-fifo-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		"", "", fifoExchange, fifoQueue, 10, false,
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
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockPubCh.EXPECT().PublishWithContext(
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

	err := consumer.Consume(t.Context())
	if !errors.Is(err, libamqp.ErrChannelClosed) {
		t.Errorf("expected error %v but got %v", libamqp.ErrChannelClosed, err)
	}
}

func TestConcurrentConsumer_Consume_GivenProcessedMessageAndNoNextMessage_ShouldPublishMessageToFIFOQueue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "test-consumer"
	queue := "test-queue"
	nextExchange := "next-fifo-exchange"
	nextQueue := "next-fifo-queue"
	fifoExchange := "test-fifo-exchange"
	fifoQueue := "test-fifo-queue"
	mockPubCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool := mock_amqp.NewMockChannelPool(ctrl)
	mockConsumeCh := mock_amqp.NewMockChannel(ctrl)
	mockChPool.EXPECT().Get(gomock.Any()).Return(mockConsumeCh, nil)
	mockChPool.EXPECT().Put(gomock.Any())
	mockMessageProcessor := mock_engine.NewMockMessageProcessor(ctrl)
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		nextExchange, nextQueue, fifoExchange, fifoQueue, 10, false,
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
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Do(func(_ uint64, _ bool) {
		close(notifyClose)
	})
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil, libamqp.ErrChannelClosed)

	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil, nil)
	mockPubCh.EXPECT().PublishWithContext(
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

	err := consumer.Consume(t.Context())
	if !errors.Is(err, libamqp.ErrChannelClosed) {
		t.Errorf("expected error %v but got %v", libamqp.ErrChannelClosed, err)
	}
}

func TestConcurrentConsumer_Consume_GivenErrorOnMessage_ShouldStopConsumer(t *testing.T) {
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
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		"", "", "", "", 10, false,
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

func TestConcurrentConsumer_Consume_GivenContextDone_ShouldStopConsumer(t *testing.T) {
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
	consumer := engine.NewConcurrentConsumer(
		name, queue,
		false,
		"", "", "", "", 10, false,
		mockPubCh,
		mockChPool,
		mockMessageProcessor,
		zerolog.Logger{},
	)
	msgs := make(chan amqp.Delivery, 1)
	notifyClose := make(chan *amqp.Error)
	mockConsumeCh.EXPECT().Ack(gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().Nack(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockConsumeCh.EXPECT().
		ConsumeWithCloseNotify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(msgs, notifyClose, nil)
	mockMessageProcessor.EXPECT().Process(gomock.Any(), gomock.Any()).Times(0)

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	close(msgs)

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
