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
// It opens a consumer on ch, spawns `workers` goroutines built by newWorker to process deliveries,
// and watches the channel's notifyClose.
// When the channel closes, it re-consumes against the freshly reopened channel and
// spawns a new batch of workers bound to the new delivery channel.
//
// newWorker is invoked once per worker per generation, with the errgroup's context and the current delivery channel.
// The returned func() error runs in an errgroup goroutine — returning a non-nil error cancels the whole consume loop
// and is propagated as the return value.
// Workers must select on gctx.Done() to exit cleanly on shutdown or sibling failure.
//
// Returns nil when ctx is canceled, or the first error returned by any spawned worker or
// by ConsumeWithContext on (re)consume.
func ConsumeWithReconnect(
	ctx context.Context,
	ch Channel,
	opts ConsumeOptions,
	workers int,
	newWorker func(gctx context.Context, d <-chan amqp.Delivery) func() error,
) error {
	setup := func(gctx context.Context) (<-chan amqp.Delivery, <-chan *amqp.Error, error) {
		d, notifyCh, err := ch.ConsumeWithContext(gctx, opts.Queue, opts.Consumer, opts.AutoAck, opts.Exclusive,
			opts.NoLocal, opts.NoWait, opts.Args)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot consume messages: %w", err)
		}

		return d, notifyCh, nil
	}

	return runWithReconnect(ctx, workers, setup, newWorker)
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

// SubscribeWithReconnect declares a queue, binds it to opts.Exchange with opts.RoutingKey, and consumes from it.
// On each channel reconnect it re-declares, re-binds, and re-consumes -
// so subscribers backed by exclusive or auto-delete queues survive connection drops without
// the caller managing queue lifecycle.
//
// On a NOT_FOUND from bind or consume whose reply text contains "no queue",
// the queue has died (typically the prior connection's exclusive queue lost on reconnect):
// re-declare and retry, up to maxSubscribeRetries times.
// Any other error - including a NOT_FOUND for a missing exchange - is returned to the caller.
//
// See ConsumeWithReconnect for newWorker semantics.
func SubscribeWithReconnect(
	ctx context.Context,
	ch Channel,
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
	setup := func(gctx context.Context) (d <-chan amqp.Delivery, notifyCh <-chan *amqp.Error, err error) {
		attempts := 0
		for {
			if attempts == maxSubscribeRetries {
				return nil, nil, fmt.Errorf("cannot resolve queue after %d attempts: %w", maxSubscribeRetries, err)
			}

			attempts++

			if queue == "" {
				// First time, or queue was lost on the previous attempt: declare and bind.
				var q amqp.Queue
				q, err = ch.QueueDeclare(gctx, opts.QueueName, opts.QueueDurable,
					opts.QueueAutoDelete, opts.QueueExclusive, false, opts.QueueArgs)
				if err != nil {
					return nil, nil, fmt.Errorf("cannot declare queue: %w", err)
				}

				err = ch.QueueBind(gctx, q.Name, opts.RoutingKey, opts.Exchange, false, nil)
				if err != nil {
					if queueMissing(err) {
						continue
					}

					return nil, nil, fmt.Errorf("cannot bind queue: %w", err)
				}

				queue = q.Name
			}

			d, notifyCh, err = ch.ConsumeWithContext(gctx, queue, opts.Consumer,
				opts.AutoAck, opts.QueueExclusive, opts.NoLocal, opts.NoWait, opts.Args)
			if err != nil {
				if queueMissing(err) {
					queue = ""

					continue
				}

				return nil, nil, fmt.Errorf("cannot consume messages: %w", err)
			}

			return d, notifyCh, nil
		}
	}

	return runWithReconnect(ctx, workers, setup, newWorker)
}

// runWithReconnect runs setup to obtain a delivery channel, spawns workers to process it,
// and re-runs setup whenever the channel notifies close.
// Returns when ctx is canceled, when setup returns an error, or when any worker returns an error.
func runWithReconnect(
	ctx context.Context,
	workers int,
	setup func(gctx context.Context) (<-chan amqp.Delivery, <-chan *amqp.Error, error),
	newWorker func(gctx context.Context, d <-chan amqp.Delivery) func() error,
) error {
	g, gctx := errgroup.WithContext(ctx)

	spawn := func(d <-chan amqp.Delivery) {
		for i := 0; i < workers; i++ {
			g.Go(newWorker(gctx, d))
		}
	}

	d, notifyCh, err := setup(gctx)
	if err != nil {
		return err
	}

	spawn(d)

	g.Go(func() error {
		for {
			select {
			case <-gctx.Done():
				return nil
			case <-notifyCh:
				d, notifyCh, err = setup(gctx)
				if err != nil {
					return err
				}

				spawn(d)
			}
		}
	})

	return g.Wait()
}
