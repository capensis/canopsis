package amqp_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestConsumeWithReconnect_GivenSuccessfulConsume_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		workers := 3
		ch := mock_amqp.NewMockChannel(ctrl)
		d := make(chan amqp091.Delivery)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, workers, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenInitialConsumeError_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		cerr := errors.New("test")
		queue := "amq.queue"
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil, cerr)

		opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
		gotErr := amqp.ConsumeWithReconnect(t.Context(), ch, opts, 1, newNoopProcess, zerolog.Nop())

		if !errors.Is(gotErr, cerr) {
			t.Fatalf("expected %+v, got %+v", cerr, gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenNotifyClose_ShouldReConsume(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		notifyClose1 := make(chan *amqp091.Error)
		notifyClose2 := make(chan *amqp091.Error)
		queue := "amq.queue"
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(
				ctx context.Context,
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args amqp091.Table,
			) (<-chan amqp091.Delivery, <-chan *amqp091.Error, error) {
				return make(chan amqp091.Delivery), notifyClose1, nil
			})
		d := make(chan amqp091.Delivery)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(
				ctx context.Context,
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args amqp091.Table,
			) (<-chan amqp091.Delivery, <-chan *amqp091.Error, error) {
				return d, notifyClose2, nil
			})

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		close(notifyClose1)
		synctest.Wait()

		cancel()
		close(d)

		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenProcessReturnsAck_ShouldAckDelivery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery, 1)
		msg := amqp091.Delivery{DeliveryTag: 1}
		d <- msg
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		ch.EXPECT().Ack(gomock.Eq(msg.DeliveryTag), gomock.Any()).Return(nil)

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenProcessReturnsNack_ShouldNackDelivery(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery, 1)
		msg := amqp091.Delivery{DeliveryTag: 1}
		d <- msg
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		ch.EXPECT().Nack(gomock.Eq(msg.DeliveryTag), gomock.Any(), gomock.Any()).Return(nil)

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		process := func(context.Context, amqp091.Delivery) (amqp.AckAction, error) {
			return amqp.Nack, nil
		}

		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, 1, process, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenProcessReturnsNackAndError_ShouldNackDeliveryAndReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery, 1)
		msg := amqp091.Delivery{DeliveryTag: 1}
		d <- msg
		close(d)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		ch.EXPECT().Nack(gomock.Eq(msg.DeliveryTag), gomock.Any(), gomock.Any()).Return(nil)

		var gotErr error
		done := make(chan struct{})
		processErr := errors.New("test error")
		process := func(context.Context, amqp091.Delivery) (amqp.AckAction, error) {
			return amqp.Nack, processErr
		}

		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(t.Context(), ch, opts, 1, process, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if !errors.Is(gotErr, processErr) {
			t.Fatalf("expected %+v, got %+v", processErr, gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenAckBrokerError_ShouldNotReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery, 1)
		d <- amqp091.Delivery{DeliveryTag: 1}
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		ch.EXPECT().Ack(gomock.Any(), gomock.Any()).Return(errors.New("broker ack failed"))

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("ack errors should be logged, not returned; got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenReconnectSetupError_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery)
		notifyClose := make(chan *amqp091.Error)
		consumeErr := errors.New("reconnect failed")
		gomock.InOrder(
			ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false)),
			ch.EXPECT().
				ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(d, notifyClose, nil),
			ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false)),
			ch.EXPECT().
				ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil, consumeErr),
		)

		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(t.Context(), ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		close(notifyClose)
		synctest.Wait()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after reconnect setup error")
		}

		if !errors.Is(gotErr, consumeErr) {
			t.Fatalf("expected %+v, got %+v", consumeErr, gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenCtxCancelDuringProcess_ShouldDrainAck(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		workers := 5
		d := make(chan amqp091.Delivery, workers)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		for i := 0; i < workers; i++ {
			msg := amqp091.Delivery{DeliveryTag: uint64(i + 1)}
			d <- msg
			//Ack must be applied during the drain phase, after ctx is canceled.
			ch.EXPECT().Ack(gomock.Eq(msg.DeliveryTag), gomock.Any()).Return(nil)
		}

		processStarted := make(chan struct{}, workers)
		processDone := make(chan struct{})
		slowProcess := func(_ context.Context, _ amqp091.Delivery) (amqp.AckAction, error) {
			processStarted <- struct{}{}
			<-processDone

			return amqp.Ack, nil
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, workers, slowProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()

		for i := 0; i < workers; i++ {
			<-processStarted
		}

		synctest.Wait()

		cancel()
		close(d)

		synctest.Wait()

		close(processDone)

		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenProcessPanic_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		ch := mock_amqp.NewMockChannel(ctrl)

		d := make(chan amqp091.Delivery, 1)
		d <- amqp091.Delivery{DeliveryTag: 1}
		close(d)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)

		panicErr := errors.New("test panic")
		panicProcess := func(_ context.Context, _ amqp091.Delivery) (amqp.AckAction, error) {
			panic(panicErr)
		}

		opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
		gotErr := amqp.ConsumeWithReconnect(t.Context(), ch, opts, 1, panicProcess, zerolog.Nop())

		if !errors.Is(gotErr, panicErr) {
			t.Fatalf("expected %+v, got %+v", panicErr, gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenCtxCanceled_ShouldDrainDeliveries(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.queue"
		workers := 3
		ch := mock_amqp.NewMockChannel(ctrl)
		msgCount := 3 * workers
		d := make(chan amqp091.Delivery, msgCount)
		for i := 0; i < msgCount; i++ {
			msg := amqp091.Delivery{DeliveryTag: uint64(i + 1)}
			d <- msg

			if i < workers {
				ch.EXPECT().Ack(gomock.Eq(msg.DeliveryTag), gomock.Any()).Return(nil)
			} else {
				ch.EXPECT().Nack(gomock.Eq(msg.DeliveryTag), gomock.Any(), gomock.Any()).Return(nil)
			}
		}

		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)

		processStarted := make(chan struct{}, workers)
		processDone := make(chan struct{})
		var processExecutions atomic.Int32
		slowProcess := func(_ context.Context, _ amqp091.Delivery) (amqp.AckAction, error) {
			processExecutions.Add(1)
			processStarted <- struct{}{}
			<-processDone

			return amqp.Ack, nil
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			opts := amqp.ConsumeOptions{Queue: queue, PrefetchCount: 1, PrefetchSize: 100}
			gotErr = amqp.ConsumeWithReconnect(ctx, ch, opts, workers, slowProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()

		for i := 0; i < workers; i++ {
			<-processStarted
		}

		synctest.Wait()

		cancel()
		close(d)

		synctest.Wait()

		close(processDone)

		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("ConsumeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}

		if got := processExecutions.Load(); got != int32(workers) {
			t.Fatalf("expected %d processed deliveries, got %d", workers, got)
		}

		select {
		case m, ok := <-d:
			if ok {
				t.Fatalf("expected deliveries to be drained, got %+v", m)
			}
		default:
		}
	})
}

func TestSubscribeWithReconnect_GivenSuccessfulSetup_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queue := "amq.gen-1"
		exchange := "test.exchange"
		workers := 2
		ch := mock_amqp.NewMockChannel(ctrl)
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		d := make(chan amqp091.Delivery)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)
		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
			PrefetchCount:  1,
			PrefetchSize:   100,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, ch, opts, workers, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("SubscribeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenBindQueueMissing_ShouldRetryDeclare(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		queue1 := "amq.gen-1"
		queue2 := "amq.gen-2"
		exchange := "test.exchange"
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue1}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp091.Error{Code: amqp091.NotFound, Reason: "NOT_FOUND - no queue '" + queue1 + "' in vhost '/'"})
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue2}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		d := make(chan amqp091.Delivery)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
			PrefetchCount:  1,
			PrefetchSize:   100,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("SubscribeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenConsumeQueueMissing_ShouldResetAndRetry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		queue1 := "amq.gen-1"
		queue2 := "amq.gen-2"
		exchange := "test.exchange"
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue1}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil, &amqp091.Error{Code: amqp091.NotFound, Reason: "NOT_FOUND - no queue '" + queue1 + "' in vhost '/'"})
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue2}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		d := make(chan amqp091.Delivery)
		ch.EXPECT().Qos(gomock.Any(), gomock.Eq(1), gomock.Eq(100), gomock.Eq(false))
		ch.EXPECT().
			ConsumeWithCloseNotify(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(d, make(chan *amqp091.Error), nil)

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
			PrefetchCount:  1,
			PrefetchSize:   100,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, ch, opts, 1, newNoopProcess, zerolog.Nop())
			close(done)
		}()

		synctest.Wait()
		cancel()
		close(d)
		synctest.Wait()

		select {
		case <-done:
		default:
			t.Fatal("SubscribeWithReconnect did not return after ctx cancel")
		}

		if gotErr != nil {
			t.Fatalf("expected nil, got %+v", gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenBindExchangeMissing_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp091.Error{Code: amqp091.NotFound, Reason: "NOT_FOUND - no exchange '" + exchange + "' in vhost '/'"})

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
			PrefetchCount:  1,
			PrefetchSize:   100,
		}

		gotErr := amqp.SubscribeWithReconnect(t.Context(), ch, opts, 1, newNoopProcess, zerolog.Nop())

		if aerr, ok := errors.AsType[*amqp091.Error](gotErr); !ok || aerr.Code != amqp091.NotFound {
			t.Fatalf("expected NOT_FOUND error, got %+v", gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenMaxRetriesExceeded_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ch := mock_amqp.NewMockChannel(ctrl)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue}, nil).
			Times(10)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp091.Error{Code: amqp091.NotFound, Reason: "NOT_FOUND - no queue '" + queue + "' in vhost '/'"}).
			Times(10)

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
			PrefetchCount:  1,
			PrefetchSize:   100,
		}

		gotErr := amqp.SubscribeWithReconnect(t.Context(), ch, opts, 1, newNoopProcess, zerolog.Nop())

		if gotErr == nil {
			t.Fatal("expected error after max retries, got nil")
		}

		if aerr, ok := errors.AsType[*amqp091.Error](gotErr); !ok || aerr.Code != amqp091.NotFound {
			t.Fatalf("expected NOT_FOUND error, got %+v", gotErr)
		}
	})
}

func newNoopProcess(context.Context, amqp091.Delivery) (amqp.AckAction, error) {
	return amqp.Ack, nil
}
