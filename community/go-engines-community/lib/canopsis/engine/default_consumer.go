package engine

import (
	"context"
	"fmt"
	"runtime/debug"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type defaultConsumer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue      string
	purgeQueue bool
	exclusive  bool
	// processor handles AMQP messages.
	processor MessageProcessor
	// nextQueue is name of AMQP queue to where consumer sends message after succeeded processing.
	nextQueue    string
	nextExchange string
	// fifoQueue is name of AMQP queue to where consumer sends message after failed processing
	// or if nextQueue is not defined.
	fifoQueue    string
	fifoExchange string

	publisher   libamqp.Publisher
	consumePool libamqp.ChannelPool

	logger zerolog.Logger
}

func (c *defaultConsumer) queuePurge(ctx context.Context) error {
	if c.purgeQueue {
		ch, err := c.consumePool.Get(ctx)
		if err != nil {
			return err
		}

		defer c.consumePool.Put(ch)

		_, err = ch.QueuePurge(ctx, c.queue, false)
		if err != nil {
			return fmt.Errorf("error while purging queue: %w", err)
		}
	}

	return nil
}

func (c *defaultConsumer) processMessage(ctx context.Context, d amqp.Delivery) error {
	c.logger.Debug().
		Str("consumer", c.name).Str("queue", c.queue).
		Str("msg", string(d.Body)).
		Msgf("received")
	msgToNext, err := c.processor.Process(ctx, d)
	if err != nil {
		nackErr := d.Nack(false, true)
		if nackErr != nil {
			c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
		}

		return fmt.Errorf("cannot process message: %w", err)
	}

	err = d.Ack(false)
	if err != nil {
		c.logger.Err(err).Msg("cannot ack amqp delivery")
	}

	if c.nextQueue != "" && msgToNext != nil {
		err = c.publisher.PublishWithContext(
			ctx,
			c.nextExchange,
			c.nextQueue,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         msgToNext,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot sent message to next queue: %w", err)
		}

		return nil
	}

	if d.ReplyTo != "" && msgToNext != nil {
		err = c.publisher.PublishWithContext(
			ctx,
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          msgToNext,
				DeliveryMode:  amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot sent message result back to sender: %w", err)
		}

		return nil
	}

	if c.fifoQueue != "" {
		err = c.publisher.PublishWithContext(
			ctx,
			c.fifoExchange,
			c.fifoQueue,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         d.Body,
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("cannot sent message to fifo queue: %w", err)
		}
	}

	return nil
}

func (c *defaultConsumer) consume(ctx context.Context, workers int) error {
	opts := libamqp.ConsumeOptions{
		Queue:     c.queue,
		Consumer:  c.name,
		Exclusive: c.exclusive,
	}

	return libamqp.ConsumeWithReconnect(ctx, c.consumePool, opts, workers, c.newWorkerFunc)
}

func (c *defaultConsumer) newWorkerFunc(ctx context.Context, ch <-chan amqp.Delivery) func() error {
	return func() (resErr error) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				var ok bool
				if err, ok = r.(error); !ok {
					err = fmt.Errorf("%v", r)
				}

				c.logger.Err(err).Msgf("consumer recovered from panic\n%s\n", debug.Stack())
				resErr = fmt.Errorf("consumer recovered from panic: %w", err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return nil
			case d, ok := <-ch:
				if !ok {
					return nil
				}

				err := c.processMessage(ctx, d)
				if err != nil {
					return err
				}
			}
		}
	}
}
