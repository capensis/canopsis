package engine

import (
	"context"
	"errors"
	libamqp "git.canopsis.net/canopsis/go-engines/lib/amqp"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

// NewRPCServer creates consumer.
func NewRPCServer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &rpcServer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		connection:           connection,
		processor:            processor,
		logger:               logger,
	}
}

// rpcServer implements AMQP consumer of RPC requests.
type rpcServer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue                                     string
	consumePrefetchCount, consumePrefetchSize int
	// processor handles AMQP messages.
	processor MessageProcessor
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}

func (c *rpcServer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := getConsumeChannel(c.connection, c.name, c.queue,
		c.consumePrefetchCount, c.consumePrefetchSize, false)
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

			err := c.processMessage(ctx, d, consumeCh, publishCh)
			if err != nil {
				return err
			}
		}
	}
}

func (c *rpcServer) processMessage(ctx context.Context, d amqp.Delivery, consumeCh, publishCh libamqp.Channel) error {
	c.logger.Debug().Str("msg", string(d.Body)).Msgf("received")
	body, err := c.processor.Process(ctx, d)

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

	if body != nil {
		err = publishCh.Publish(
			"",        // exchange
			d.ReplyTo, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType:   "application/json",
				CorrelationId: d.CorrelationId,
				Body:          body,
				DeliveryMode:  amqp.Persistent,
			},
		)
		if err != nil {
			c.logger.Err(err).Msg("cannot sent message result back to sender")
			return err
		}
	}

	return nil
}
