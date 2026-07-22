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
	prefetchCount, prefetchSize int,
	workers int,
	publisher libamqp.Publisher,
	consumePool libamqp.ChannelPool,
	processor RPCMessageProcessor,
	logger zerolog.Logger,
) RPCClient {
	return &rpcClient{
		defaultConsumer: defaultConsumer{
			name:          name,
			queue:         clientQueueName,
			prefetchCount: prefetchCount,
			prefetchSize:  prefetchSize,
			processor:     &rpcClientMessageProcessorWrapper{processor: processor},
			publisher:     publisher,
			consumePool:   consumePool,
			logger:        logger,
		},
		serverExchangeName: serverExchangeName,
		serverRoutingKey:   serverRoutingKey,
		publisher:          publisher,
		workers:            workers,
	}
}

func NewRPCClientWithoutReply(
	serverExchangeName string,
	serverRoutingKey string,
	publisher libamqp.Publisher,
) RPCClient {
	return &rpcClient{
		serverExchangeName: serverExchangeName,
		serverRoutingKey:   serverRoutingKey,
		publisher:          publisher,
	}
}

// rpcClient implements RPC client.
type rpcClient struct {
	defaultConsumer
	serverExchangeName string
	serverRoutingKey   string
	publisher          libamqp.Publisher
	// amount of workers which process events.
	workers int
}

func (c *rpcClient) Call(ctx context.Context, m RPCMessage) error {
	return c.publisher.PublishWithContext(
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
}

func (c *rpcClient) Consume(ctx context.Context) (err error) {
	if c.consumePool == nil {
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
