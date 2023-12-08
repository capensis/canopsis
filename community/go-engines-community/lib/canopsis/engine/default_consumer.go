package engine

import (
	"context"
	"errors"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// NewDefaultConsumer creates consumer.
func NewDefaultConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &defaultConsumer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		purgeQueue:           purgeQueue,
		nextExchange:         nextExchange,
		nextQueue:            nextQueue,
		fifoExchange:         fifoExchange,
		fifoQueue:            fifoQueue,
		connection:           connection,
		processor:            processor,
		logger:               logger,
	}
}

// defaultConsumer implements AMQP consumer.
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

func (c *defaultConsumer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := getConsumeChannel(c.connection, c.name, c.queue,
		c.consumePrefetchCount, c.consumePrefetchSize, c.purgeQueue)
	if err != nil {
		return err
	}

	var publishCh libamqp.Channel
	if c.nextQueue != "" || c.fifoQueue != "" {
		publishCh, err = c.connection.Channel()
		if err != nil {
			return err
		}
	}

	defer func() {
		_ = consumeCh.Close()
		if publishCh != nil {
			_ = publishCh.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-msgs:
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
		err := publishToChannel(ctx, publishCh, c.nextExchange, c.nextQueue, msgToNext)
		if err != nil {
			return fmt.Errorf("cannot sent message to next queue: %w", err)
		}
	} else if c.fifoQueue != "" {
		err := publishToChannel(ctx, publishCh, c.fifoExchange, c.fifoQueue, d.Body)
		if err != nil {
			return fmt.Errorf("cannot sent message to fifo queue: %w", err)
		}
	}

	return nil
}

func getConsumeChannel(
	connection libamqp.Connection,
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
) (libamqp.Channel, <-chan amqp.Delivery, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = channel.Qos(consumePrefetchCount, consumePrefetchSize, false)
	if err != nil {
		return nil, nil, err
	}

	if purgeQueue {
		_, err := channel.QueuePurge(queue, false)
		if err != nil {
			return nil, nil, fmt.Errorf("error while purging queue: %w", err)
		}
	}

	msgs, err := channel.Consume(
		queue, // queue
		name,  // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, nil, err
	}

	return channel, msgs, nil
}

func publishToChannel(ctx context.Context, channel libamqp.Channel, exchange, queue string, body []byte) error {
	return channel.PublishWithContext(
		ctx,
		exchange,
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
}
