package engine

import (
	"context"
	"errors"
	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

// NewDefaultConsumer creates consumer.
func NewDefaultConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	nextQueue, fifoQueue string,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &defaultConsumer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		nextQueue:            nextQueue,
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
	// processor handles AMQP messages.
	processor MessageProcessor
	// nextQueue is name of AMQP queue to where consumer sends message after succeeded processing.
	nextQueue string
	// fifoQueue is name of AMQP queue to where consumer sends message after failed processing
	// or if nextQueue is not defined.
	fifoQueue string
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}

func (c *defaultConsumer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := getConsumeChannel(c.connection, c.name, c.queue,
		c.consumePrefetchCount, c.consumePrefetchSize)
	if err != nil {
		return err
	}

	publishCh, err := c.connection.Channel()
	if err != nil {
		return err
	}

	defer func() {
		_ = consumeCh.Close()
		_ = publishCh.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-msgs:
			if !ok {
				c.logger.Error().Msg("the rabbitmq channel has been closed")
				return errors.New("channel is closed")
			}

			err := c.processMessage(d, consumeCh, publishCh)
			if err != nil {
				return err
			}
		}
	}
}

func (c *defaultConsumer) processMessage(d amqp.Delivery, consumeCh, publishCh libamqp.Channel) error {
	c.logger.Debug().Str("msg", string(d.Body)).Msgf("received")
	msgToNext, err := c.processor.Process(d)

	if err != nil {
		c.logger.Err(err).Msg("cannot process delivery")
		nackErr := consumeCh.Nack(d.DeliveryTag, false, true)
		if nackErr != nil {
			c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
		}

		return err
	}

	err = consumeCh.Ack(d.DeliveryTag, false)
	if err != nil {
		c.logger.Err(err).Msg("cannot ack amqp delivery")
	}

	if c.nextQueue != "" && msgToNext != nil {
		err := publishToChannel(publishCh, c.nextQueue, msgToNext)
		if err != nil {
			c.logger.Err(err).Msg("cannot sent message to next queue")
			return err
		}
	} else if c.fifoQueue != "" {
		err := publishToChannel(publishCh, c.fifoQueue, d.Body)
		if err != nil {
			c.logger.Err(err).Msg("cannot sent message to fifo queue")
			return err
		}
	}

	return nil
}

func getConsumeChannel(
	connection libamqp.Connection,
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
) (libamqp.Channel, <-chan amqp.Delivery, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = channel.Qos(consumePrefetchCount, consumePrefetchSize, false)
	if err != nil {
		return nil, nil, err
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

func publishToChannel(channel libamqp.Channel, queue string, body []byte) error {
	return channel.Publish(
		"",
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
