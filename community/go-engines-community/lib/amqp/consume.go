package amqp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const maxSubscribeRetries = 10

// AckAction is the decision returned by a ProcessFunc:
// Ack to acknowledge the delivery, Nack to reject it (with requeue).
type AckAction int

const (
	Nack AckAction = iota
	Ack
)

// ProcessFunc handles a delivery and returns whether it should be Ack'd or Nack'd.
// A non-nil error is propagated to the caller and tears down the consume loop;
// the returned AckAction is still forwarded to the dispatcher before the worker exits.
type ProcessFunc func(ctx context.Context, msg amqp.Delivery) (AckAction, error)

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
// It opens a consumer on ch, spawns `workers` goroutines to process deliveries,
// and watches the channel's notifyClose. When the channel closes, it re-consumes
// against the freshly reopened channel and spawns a new batch of workers.
//
// Ack/Nack from in-flight workers on the previous channel is discarded after reconnect;
// RabbitMQ requeues those deliveries on the new consumer.
//
// Returns nil when ctx is canceled, or the first error returned by any spawned worker or
// by ConsumeWithContext on (re)consume.
func ConsumeWithReconnect(
	ctx context.Context,
	ch Channel,
	opts ConsumeOptions,
	workers int,
	process ProcessFunc,
	logger zerolog.Logger,
) error {
	setup := func(gctx context.Context) (<-chan amqp.Delivery, <-chan *amqp.Error, error) {
		d, notifyCh, err := ch.ConsumeWithCloseNotify(gctx, opts.Queue, opts.Consumer, opts.AutoAck, opts.Exclusive,
			opts.NoLocal, opts.NoWait, opts.Args)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot consume messages: %w", err)
		}

		return d, notifyCh, nil
	}

	return runWithReconnect(ctx, ch, workers, setup, process, logger)
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
// See ConsumeWithReconnect for reconnect and ack semantics.
func SubscribeWithReconnect(
	ctx context.Context,
	ch Channel,
	opts SubscribeOptions,
	workers int,
	process ProcessFunc,
	logger zerolog.Logger,
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

			d, notifyCh, err = ch.ConsumeWithCloseNotify(gctx, queue, opts.Consumer,
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

	return runWithReconnect(ctx, ch, workers, setup, process, logger)
}

// ackOp is a request from a worker to runWithReconnect's loop to Ack or Nack a delivery.
// generation identifies the channel incarnation the deliveryTag belongs to;
// requests from a stale generation are discarded after reconnect.
type ackOp struct {
	deliveryTag uint64
	action      AckAction
	generation  uint64
}

// runWithReconnect runs setup to obtain a delivery channel, spawns workers to process it,
// and re-runs setup whenever the channel notifies close.
//
// A generation counter is bumped on every reconnect and stamped on every ackOp; the loop
// drops acks whose generation no longer matches the current one, so stale delivery tags
// from a closed channel are never forwarded to the new one. RabbitMQ requeues those
// deliveries automatically when the previous channel closes.
//
// Returns when ctx is canceled, when setup returns an error, or when any worker returns an error.
func runWithReconnect(
	ctx context.Context,
	ch Channel,
	workers int,
	setup func(gctx context.Context) (<-chan amqp.Delivery, <-chan *amqp.Error, error),
	process ProcessFunc,
	logger zerolog.Logger,
) error {
	g, gctx := errgroup.WithContext(ctx)
	ackCh := make(chan ackOp, workers)

	spawnWorkers := func(d <-chan amqp.Delivery, generation uint64) {
		for i := 0; i < workers; i++ {
			g.Go(newWorker(gctx, d, ackCh, process, generation))
		}
	}

	d, notifyCh, err := setup(gctx)
	if err != nil {
		return err
	}

	generation := uint64(0)
	spawnWorkers(d, generation)

	ack := func(op ackOp) {
		if op.generation != generation {
			return
		}

		var err error
		switch op.action {
		case Ack:
			err = ch.Ack(op.deliveryTag, false)
		case Nack:
			err = ch.Nack(op.deliveryTag, false, true)
		default:
			err = fmt.Errorf("unknown ack action %d", op.action)
		}

		if err != nil {
			logger.Err(err).Msg("cannot acknowledge message")
		}
	}

	var setupErr error
loop:
	for {
		select {
		case <-gctx.Done():
			break loop
		case <-notifyCh:
			d, notifyCh, setupErr = setup(gctx)
			if setupErr != nil {
				break loop
			}

			generation++
			spawnWorkers(d, generation)
		case op := <-ackCh:
			ack(op)
		}
	}

	go func() {
		_ = g.Wait()
		close(ackCh)
	}()

	for op := range ackCh {
		ack(op)
	}

	return errors.Join(setupErr, g.Wait())
}

func newWorker(ctx context.Context, deliveries <-chan amqp.Delivery, ackCh chan ackOp, process ProcessFunc, generation uint64) func() error {
	return func() (resErr error) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("%v", r)
				}

				resErr = fmt.Errorf("consumer recovered from panic: %w", err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case d, ok := <-deliveries:
				if !ok {
					return nil
				}

				action, err := process(ctx, d)
				ackCh <- ackOp{deliveryTag: d.DeliveryTag, action: action, generation: generation}

				if err != nil {
					return fmt.Errorf("cannot process message: %w", err)
				}
			}
		}
	}
}
