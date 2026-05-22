package amqp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const maxSubscribeRetries = 10

// ConsumeOptions bundles the parameters forwarded to Channel.ConsumeWithContext.
type ConsumeOptions struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// ConsumeWithReconnect drives a consume loop that survives AMQP channel reconnects.
// It spawns `workers` goroutines, each borrowing a Channel from pool, calling
// ConsumeWithContext against opts.Queue, and running a worker function built by newWorker
// to process the resulting deliveries.
//
// Channel reconnects are handled transparently by the pool: when the delivery channel ends,
// the worker function returns and the goroutine re-runs ConsumeWithContext to bind to the
// new delivery channel.
//
// newWorker is invoked once per delivery-channel generation, with the errgroup's context and
// the current delivery channel. The returned func() error runs synchronously in the worker
// goroutine — returning a non-nil error cancels the whole consume loop and is propagated as
// the return value. A nil return means the delivery channel ended and the goroutine is ready
// to re-consume.
// Workers must select on gctx.Done() to exit cleanly on shutdown or sibling failure.
//
// Returns nil when ctx is canceled, or the first error returned by any spawned worker or
// by ConsumeWithContext on (re)consume.
func ConsumeWithReconnect(
	ctx context.Context,
	pool ChannelPool,
	opts ConsumeOptions,
	workers int,
	newWorker func(gctx context.Context, d <-chan amqp.Delivery) func() error,
) error {
	setup := func(gctx context.Context, ch Channel) (<-chan amqp.Delivery, error) {
		d, err := ch.ConsumeWithContext(gctx, opts.Queue, opts.Consumer, opts.AutoAck, opts.Exclusive,
			opts.NoLocal, opts.NoWait, opts.Args)
		if err != nil {
			return nil, fmt.Errorf("cannot consume messages: %w", err)
		}

		return d, nil
	}

	return runWithReconnect(ctx, pool, workers, setup, newWorker)
}

// SubscribeOptions bundles the parameters for declaring a transient queue,
// binding it to an exchange, and consuming from it.
type SubscribeOptions struct {
	Exchange   string
	RoutingKey string

	// Queue declare params. QueueName="" lets the broker generate a name.
	QueueName       string
	QueueDurable    bool
	QueueAutoDelete bool
	QueueExclusive  bool
	QueueArgs       amqp.Table

	// Consume params.
	Consumer string
	AutoAck  bool
	NoLocal  bool
	NoWait   bool
	Args     amqp.Table
}

// SubscribeWithReconnect declares a queue, binds it to opts.Exchange with opts.RoutingKey,
// and consumes from it. On every (re)consume — including after a channel reconnect — it
// reuses the previously declared queue and just re-runs ConsumeWithContext; if the broker
// reports the queue as missing (typically the prior connection's exclusive queue lost on
// reconnect), it re-declares, re-binds, and retries, up to maxSubscribeRetries times.
// This lets subscribers backed by exclusive or auto-delete queues survive connection drops
// without the caller managing queue lifecycle.
//
// A NOT_FOUND whose reply text contains "no queue" is treated as queue-missing and retried.
// Any other error - including a NOT_FOUND for a missing exchange - is returned to the caller.
//
// See ConsumeWithReconnect for newWorker semantics.
func SubscribeWithReconnect(
	ctx context.Context,
	pool ChannelPool,
	opts SubscribeOptions,
	workers int,
	newWorker func(gctx context.Context, d <-chan amqp.Delivery) func() error,
) error {
	queueMissing := func(err error) bool {
		if aerr, ok := errors.AsType[*amqp.Error](err); ok {
			return aerr.Code == amqp.NotFound && strings.Contains(aerr.Reason, "no queue")
		}

		return false
	}

	queue := ""
	setup := func(gctx context.Context, ch Channel) (d <-chan amqp.Delivery, err error) {
		attempts := 0
		for {
			if attempts == maxSubscribeRetries {
				return nil, fmt.Errorf("cannot resolve queue after %d attempts: %w", maxSubscribeRetries, err)
			}

			attempts++

			if queue == "" {
				// First time, or queue was lost on the previous attempt: declare and bind.
				var q amqp.Queue
				q, err = ch.QueueDeclare(gctx, opts.QueueName, opts.QueueDurable,
					opts.QueueAutoDelete, opts.QueueExclusive, false, opts.QueueArgs)
				if err != nil {
					return nil, fmt.Errorf("cannot declare queue: %w", err)
				}

				err = ch.QueueBind(gctx, q.Name, opts.RoutingKey, opts.Exchange, false, nil)
				if err != nil {
					if queueMissing(err) {
						continue
					}

					return nil, fmt.Errorf("cannot bind queue: %w", err)
				}

				queue = q.Name
			}

			d, err = ch.ConsumeWithContext(gctx, queue, opts.Consumer,
				opts.AutoAck, opts.QueueExclusive, opts.NoLocal, opts.NoWait, opts.Args)
			if err != nil {
				if queueMissing(err) {
					queue = ""

					continue
				}

				return nil, fmt.Errorf("cannot consume messages: %w", err)
			}

			return d, nil
		}
	}

	return runWithReconnect(ctx, pool, workers, setup, newWorker)
}

// runWithReconnect spawns `workers` goroutines, each borrowing a Channel from pool.
// Each goroutine loops: call setup to obtain a delivery channel, then run newWorker against it.
// When newWorker returns nil (typically because the delivery channel ended on a reconnect)
// the loop re-runs setup against the same Channel — the pool handles the underlying AMQP
// reconnect transparently.
// Returns when ctx is canceled, when setup returns an error, or when any worker returns a non-nil error.
func runWithReconnect(
	ctx context.Context,
	pool ChannelPool,
	workers int,
	setup func(gctx context.Context, ch Channel) (<-chan amqp.Delivery, error),
	newWorker func(gctx context.Context, d <-chan amqp.Delivery) func() error,
) error {
	g, gctx := errgroup.WithContext(ctx)

	for i := 0; i < workers; i++ {
		g.Go(func() error {
			ch, err := pool.Get(gctx)
			if err != nil {
				return err
			}

			defer pool.Put(ch)

			for {
				select {
				case <-gctx.Done():
					return nil
				default:
				}

				d, err := setup(gctx, ch)
				if err != nil {
					return err
				}

				err = newWorker(gctx, d)()
				if err != nil {
					return err
				}
			}
		})
	}

	return g.Wait()
}
