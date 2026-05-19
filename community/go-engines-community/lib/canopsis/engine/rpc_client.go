package engine

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// NewRPCClient creates new AMQP RPC client.
func NewRPCClient(
	name, serverExchangeName, serverRoutingKey, clientQueueName string,
	workers int,
	publishCh libamqp.Channel,
	consumeCh libamqp.Channel,
	processor RPCMessageProcessor,
	logger zerolog.Logger,
) RPCClient {
	return &rpcClient{
		defaultConsumer: defaultConsumer{
			name:      name,
			queue:     clientQueueName,
			processor: &rpcClientMessageProcessorWrapper{processor: processor},
			publishCh: publishCh,
			consumeCh: consumeCh,
			logger:    logger,
		},
		serverExchangeName: serverExchangeName,
		serverRoutingKey:   serverRoutingKey,
		publishCh:          publishCh,
		workers:            workers,
	}
}

func NewRPCClientWithoutReply(
	serverExchangeName string,
	serverRoutingKey string,
	publishCh libamqp.Channel,
) RPCClient {
	return &rpcClient{
		serverExchangeName: serverExchangeName,
		serverRoutingKey:   serverRoutingKey,
		publishCh:          publishCh,
	}
}

// rpcClient implements RPC client.
type rpcClient struct {
	defaultConsumer
	serverExchangeName string
	serverRoutingKey   string
	publishCh          libamqp.Channel
	// amount of workers which process events.
	workers int
}

func (c *rpcClient) Call(ctx context.Context, m RPCMessage) error {
	err := c.publishCh.PublishWithContext(
		ctx,
		c.serverExchangeName,
		c.serverRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: m.CorrelationID,
			ReplyTo:       c.queue,
			Body:          m.Body,
			DeliveryMode:  amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *rpcClient) Consume(ctx context.Context) (err error) {
	if c.consumeCh == nil {
		return errors.New("consume channel is nil")
	}

	return c.consume(ctx, c.workers)
}

type rpcClientMessageProcessorWrapper struct {
	processor RPCMessageProcessor
}

func (r *rpcClientMessageProcessorWrapper) Process(ctx context.Context, d amqp.Delivery) ([]byte, error) {
	err := r.processor.Process(ctx, RPCMessage{
		CorrelationID: d.CorrelationId,
		Body:          d.Body,
	})

	return nil, err
}
