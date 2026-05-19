package amqp

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/mock/gomock"
)

func TestConsumeWithReconnect_GivenSuccessfulConsume_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue := "amq.queue"
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any())

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = ConsumeWithReconnect(ctx, ch, ConsumeOptions{Queue: queue}, 3, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		cancel()
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

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		cerr := errors.New("test")
		queue := "amq.queue"
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, cerr)

		gotErr := ConsumeWithReconnect(t.Context(), ch, ConsumeOptions{Queue: queue}, 1, newNoopWorker())

		if !errors.Is(gotErr, cerr) {
			t.Fatalf("expected %+v, got %+v", cerr, gotErr)
		}
	})
}

func TestConsumeWithReconnect_GivenNotifyClose_ShouldReConsume(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		var notifyCh1, notifyCh2 chan *amqp.Error
		queue := "amq.queue"
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any()).DoAndReturn(func(ch chan *amqp.Error) chan *amqp.Error {
			notifyCh1 = ch

			return ch
		})
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any()).DoAndReturn(func(ch chan *amqp.Error) chan *amqp.Error {
			notifyCh2 = ch

			return ch
		})

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = ConsumeWithReconnect(ctx, ch, ConsumeOptions{Queue: queue}, 1, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		close(notifyCh1)
		synctest.Wait()

		if notifyCh2 == nil {
			t.Fatal("re-consume did not happen after notifyCh fired")
		}

		cancel()
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

func TestConsumeWithReconnect_GivenWorkerError_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue := "amq.queue"
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any())

		werr := errors.New("worker error")
		failFast := func(gctx context.Context, d <-chan amqp.Delivery) func() error {
			return func() error {
				return werr
			}
		}

		gotErr := ConsumeWithReconnect(t.Context(), ch, ConsumeOptions{Queue: queue}, 1, failFast)
		if !errors.Is(gotErr, werr) {
			t.Fatalf("expected %+v, got %+v", werr, gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenSuccessfulSetup_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any())

		opts := SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = SubscribeWithReconnect(ctx, ch, opts, 2, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		cancel()
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

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue1 := "amq.gen-1"
		queue2 := "amq.gen-2"
		exchange := "test.exchange"
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue1}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue1), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp.Error{Code: amqp.NotFound, Reason: "NOT_FOUND - no queue '" + queue1 + "' in vhost '/'"})
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue2}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue2), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any())

		opts := SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = SubscribeWithReconnect(ctx, ch, opts, 1, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		cancel()
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

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue1 := "amq.gen-1"
		queue2 := "amq.gen-2"
		exchange := "test.exchange"
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue1}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue1), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &amqp.Error{Code: amqp.NotFound, Reason: "NOT_FOUND - no queue '" + queue1 + "' in vhost '/'"})
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue2}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue2), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		mockCh.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp.Delivery), nil)
		mockCh.EXPECT().NotifyClose(gomock.Any())

		opts := SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = SubscribeWithReconnect(ctx, ch, opts, 1, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		cancel()
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

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue}, nil)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp.Error{Code: amqp.NotFound, Reason: "NOT_FOUND - no exchange '" + exchange + "' in vhost '/'"})

		opts := SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		gotErr := SubscribeWithReconnect(t.Context(), ch, opts, 1, newNoopWorker())

		if aerr, ok := errors.AsType[*amqp.Error](gotErr); !ok || aerr.Code != amqp.NotFound {
			t.Fatalf("expected NOT_FOUND error, got %+v", gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenMaxRetriesExceeded_ShouldReturnErr(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCh := NewMockamqp091Channel(ctrl)
		ch := newTestChannel(nil, mockCh)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		mockCh.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp.Queue{Name: queue}, nil).
			Times(maxSubscribeRetries)
		mockCh.EXPECT().
			QueueBind(gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(&amqp.Error{Code: amqp.NotFound, Reason: "NOT_FOUND - no queue '" + queue + "' in vhost '/'"}).
			Times(maxSubscribeRetries)

		opts := SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		gotErr := SubscribeWithReconnect(t.Context(), ch, opts, 1, newNoopWorker())

		if gotErr == nil {
			t.Fatal("expected error after max retries, got nil")
		}

		if aerr, ok := errors.AsType[*amqp.Error](gotErr); !ok || aerr.Code != amqp.NotFound {
			t.Fatalf("expected NOT_FOUND error, got %+v", gotErr)
		}
	})
}

func newNoopWorker() func(context.Context, <-chan amqp.Delivery) func() error {
	return func(gctx context.Context, d <-chan amqp.Delivery) func() error {
		return func() error {
			for {
				select {
				case <-gctx.Done():
					return nil
				case _, ok := <-d:
					if !ok {
						return nil
					}
				}
			}
		}
	}
}
