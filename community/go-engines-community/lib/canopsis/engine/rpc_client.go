package engine

import (
	"context"
	"errors"
	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
	"github.com/streadway/amqp"
)

// NewRPCClient creates new AMQP RPC client.
func NewRPCClient(
	name, serverQueueName, clientQueueName string,
	consumePrefetchCount, consumePrefetchSize int,
	processor RPCMessageProcessor,
	amqpChannel libamqp.Channel,
	logger zerolog.Logger,
) RPCClient {
	return &rpcClient{
		name:                 name,
		serverQueueName:      serverQueueName,
		clientQueueName:      clientQueueName,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		processor:            processor,
		amqpChannel:          amqpChannel,
		logger:               logger,
	}
}

// rpcClient implements RPC client.
type rpcClient struct {
	// name is consumer name.
	name string
	// serverQueueName is name of AMQP queue to where client sends RPC requests.
	serverQueueName string
	// clientQueueName is name of AMQP queue from where client receives RPC response.
	clientQueueName                           string
	consumePrefetchCount, consumePrefetchSize int
	// processor handles AMQP messages.
	processor RPCMessageProcessor
	// connection is AMQP connection.
	amqpChannel libamqp.Channel
	logger      zerolog.Logger
}

func (c *rpcClient) Call(m RPCMessage) error {
	err := c.amqpChannel.Publish(
		"",                // exchange
		c.serverQueueName, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: m.CorrelationID,
			ReplyTo:       c.clientQueueName,
			Body:          m.Body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *rpcClient) Consume(ctx context.Context) error {
	c.logger.Debug().Msg("consume started")
	defer c.logger.Debug().Msg("consume stopped")

	msgs, err := c.amqpChannel.Consume(
		c.clientQueueName, // queue
		"",                // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-msgs:
			if !ok {
				c.logger.Error().Msg("the rabbitmq channel has been closed")
				return errors.New("channel is closed")
			}

			c.logger.Debug().Str("msg", string(d.Body)).Msgf("received")
			err := c.processor.Process(ctx, RPCMessage{
				CorrelationID: d.CorrelationId,
				Body:          d.Body,
			})
			if err != nil {
				c.logger.Err(err).Msg("cannot process delivery")
				nackErr := c.amqpChannel.Nack(d.DeliveryTag, false, true)
				if nackErr != nil {
					c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
				}

				return err
			}

			err = c.amqpChannel.Ack(d.DeliveryTag, false)
			if err != nil {
				c.logger.Err(err).Msg("cannot ack amqp delivery")
			}
		}
	}
}
