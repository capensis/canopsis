package amqp_test

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/mock/gomock"
)

func TestConsumeWithReconnect_GivenSuccessfulConsume_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := mock_amqp.NewMockChannelPool(ctrl)

		queue := "amq.queue"
		workers := 3
		for i := 0; i < workers; i++ {
			ch := mock_amqp.NewMockChannel(ctrl)
			pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
			pool.EXPECT().Put(gomock.Any())
			ch.EXPECT().
				ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(make(chan amqp091.Delivery), nil)
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.ConsumeWithReconnect(ctx, pool, amqp.ConsumeOptions{Queue: queue}, workers, newNoopWorker())
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

		ch := mock_amqp.NewMockChannel(ctrl)
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

		cerr := errors.New("test")
		queue := "amq.queue"
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, cerr)

		gotErr := amqp.ConsumeWithReconnect(t.Context(), pool, amqp.ConsumeOptions{Queue: queue}, 1, newNoopWorker())

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
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

		var msgCh1, msgCh2 chan amqp091.Delivery
		queue := "amq.queue"
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(
				ctx context.Context,
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args amqp091.Table,
			) (<-chan amqp091.Delivery, error) {
				msgCh1 = make(chan amqp091.Delivery)

				return msgCh1, nil
			})
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(
				ctx context.Context,
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args amqp091.Table,
			) (<-chan amqp091.Delivery, error) {
				msgCh2 = make(chan amqp091.Delivery)

				return msgCh2, nil
			})

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.ConsumeWithReconnect(ctx, pool, amqp.ConsumeOptions{Queue: queue}, 1, newNoopWorker())
			close(done)
		}()

		synctest.Wait()
		close(msgCh1)
		synctest.Wait()

		if msgCh2 == nil {
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

		ch := mock_amqp.NewMockChannel(ctrl)
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

		queue := "amq.queue"
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp091.Delivery), nil)

		werr := errors.New("worker error")
		failFast := func(gctx context.Context, d <-chan amqp091.Delivery) func() error {
			return func() error {
				return werr
			}
		}

		gotErr := amqp.ConsumeWithReconnect(t.Context(), pool, amqp.ConsumeOptions{Queue: queue}, 1, failFast)
		if !errors.Is(gotErr, werr) {
			t.Fatalf("expected %+v, got %+v", werr, gotErr)
		}
	})
}

func TestSubscribeWithReconnect_GivenSuccessfulSetup_ShouldRunUntilCtxCancel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pool := mock_amqp.NewMockChannelPool(ctrl)

		queue := "amq.gen-1"
		exchange := "test.exchange"
		workers := 2
		for i := 0; i < workers; i++ {
			ch := mock_amqp.NewMockChannel(ctrl)
			pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
			pool.EXPECT().Put(gomock.Any())
			if i == 0 {
				ch.EXPECT().
					QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(amqp091.Queue{Name: queue}, nil)
				ch.EXPECT().
					QueueBind(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
					Return(nil)
			}

			ch.EXPECT().
				ConsumeWithContext(gomock.Any(), gomock.Eq(queue), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(make(chan amqp091.Delivery), nil)
		}

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, pool, opts, workers, newNoopWorker())
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

		ch := mock_amqp.NewMockChannel(ctrl)
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

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
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp091.Delivery), nil)

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, pool, opts, 1, newNoopWorker())
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

		ch := mock_amqp.NewMockChannel(ctrl)
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

		queue1 := "amq.gen-1"
		queue2 := "amq.gen-2"
		exchange := "test.exchange"
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue1}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue1), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, &amqp091.Error{Code: amqp091.NotFound, Reason: "NOT_FOUND - no queue '" + queue1 + "' in vhost '/'"})
		ch.EXPECT().
			QueueDeclare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(amqp091.Queue{Name: queue2}, nil)
		ch.EXPECT().
			QueueBind(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Eq(exchange), gomock.Any(), gomock.Any()).
			Return(nil)
		ch.EXPECT().
			ConsumeWithContext(gomock.Any(), gomock.Eq(queue2), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make(chan amqp091.Delivery), nil)

		opts := amqp.SubscribeOptions{
			Exchange:       exchange,
			QueueExclusive: true,
		}

		ctx, cancel := context.WithCancel(t.Context())
		var gotErr error
		done := make(chan struct{})
		go func() {
			gotErr = amqp.SubscribeWithReconnect(ctx, pool, opts, 1, newNoopWorker())
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

		ch := mock_amqp.NewMockChannel(ctrl)
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

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
		}

		gotErr := amqp.SubscribeWithReconnect(t.Context(), pool, opts, 1, newNoopWorker())

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
		pool := mock_amqp.NewMockChannelPool(ctrl)
		pool.EXPECT().Get(gomock.Any()).Return(ch, nil)
		pool.EXPECT().Put(gomock.Any())

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
		}

		gotErr := amqp.SubscribeWithReconnect(t.Context(), pool, opts, 1, newNoopWorker())

		if gotErr == nil {
			t.Fatal("expected error after max retries, got nil")
		}

		if aerr, ok := errors.AsType[*amqp091.Error](gotErr); !ok || aerr.Code != amqp091.NotFound {
			t.Fatalf("expected NOT_FOUND error, got %+v", gotErr)
		}
	})
}

func newNoopWorker() func(context.Context, <-chan amqp091.Delivery) func() error {
	return func(gctx context.Context, d <-chan amqp091.Delivery) func() error {
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
