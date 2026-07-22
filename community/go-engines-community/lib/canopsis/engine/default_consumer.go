package engine

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type defaultConsumer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue         string
	prefetchCount int
	prefetchSize  int
	purgeQueue    bool
	exclusive     bool
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

func (c *defaultConsumer) processMessage(ctx context.Context, d amqp.Delivery) (libamqp.AckAction, error) {
	c.logger.Debug().
		Str("consumer", c.name).Str("queue", c.queue).
		Str("msg", string(d.Body)).
		Msgf("received")
	msgToNext, err := c.processor.Process(ctx, d)
	if err != nil {
		return libamqp.Nack, fmt.Errorf("cannot process message: %w", err)
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
			return libamqp.Ack, fmt.Errorf("cannot sent message to next queue: %w", err)
		}

		return libamqp.Ack, nil
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
			return libamqp.Ack, fmt.Errorf("cannot sent message result back to sender: %w", err)
		}

		return libamqp.Ack, nil
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
			return libamqp.Ack, fmt.Errorf("cannot sent message to fifo queue: %w", err)
		}
	}

	return libamqp.Ack, nil
}

func (c *defaultConsumer) consume(ctx context.Context, workers int) error {
	ch, err := c.consumePool.Get(ctx)
	if err != nil {
		return err
	}

	defer c.consumePool.Put(ch)

	if c.purgeQueue {
		_, err = ch.QueuePurge(ctx, c.queue, false)
		if err != nil {
			return fmt.Errorf("error while purging queue: %w", err)
		}
	}

	opts := libamqp.ConsumeOptions{
		Queue:         c.queue,
		Consumer:      c.name,
		Exclusive:     c.exclusive,
		PrefetchSize:  c.prefetchSize,
		PrefetchCount: c.prefetchCount,
	}

	return libamqp.ConsumeWithReconnect(ctx, ch, opts, workers, c.processMessage, c.logger)
}
