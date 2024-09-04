package engine

import (
	"context"
	"errors"
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
	queue                                     string
	consumePrefetchCount, consumePrefetchSize int
	purgeQueue                                bool
	// processor handles AMQP messages.
	processor MessageProcessor
	// nextQueue is name of AMQP queue to where consumer sends message after succeeded processing.
	nextQueue    string
	nextExchange string
	// fifoQueue is name of AMQP queue to where consumer sends message after failed processing
	// or if nextQueue is not defined.
	fifoQueue    string
	fifoExchange string
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}

func (c *defaultConsumer) getConsumeChannel() (libamqp.Channel, <-chan amqp.Delivery, error) {
	channel, err := c.connection.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open channel: %w", err)
	}

	err = channel.Qos(c.consumePrefetchCount, c.consumePrefetchSize, false)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot set QoS: %w", err)
	}

	if c.purgeQueue {
		_, err := channel.QueuePurge(c.queue, false)
		if err != nil {
			return nil, nil, fmt.Errorf("error while purging queue: %w", err)
		}
	}

	msgs, err := channel.Consume(
		c.queue, // queue
		c.name,  // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot consume messages: %w", err)
	}

	return channel, msgs, nil
}

func (c *defaultConsumer) processMessage(ctx context.Context, d amqp.Delivery, consumeCh, publishCh libamqp.Channel) error {
	c.logger.Debug().
		Str("consumer", c.name).Str("queue", c.queue).
		Str("msg", string(d.Body)).
		Msgf("received")
	msgToNext, err := c.processor.Process(ctx, d)
	if err != nil {
		nackErr := consumeCh.Nack(d.DeliveryTag, false, true)
		if nackErr != nil {
			c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
		}

		return fmt.Errorf("cannot process message: %w", err)
	}

	err = consumeCh.Ack(d.DeliveryTag, false)
	if err != nil {
		c.logger.Err(err).Msg("cannot ack amqp delivery")
	}

	if c.nextQueue != "" && msgToNext != nil {
		err = publishCh.PublishWithContext(
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
		err = publishCh.PublishWithContext(
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
		err = publishCh.PublishWithContext(
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

func (c *defaultConsumer) getWorkerFunc(
	ctx context.Context,
	ch <-chan amqp.Delivery,
	consumeCh, publishCh libamqp.Channel,
) func() error {
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
					return errors.New("the rabbitmq channel has been closed")
				}

				err := c.processMessage(ctx, d, consumeCh, publishCh)
				if err != nil {
					return err
				}
			}
		}
	}
}
