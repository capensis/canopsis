package engine

import (
	"context"
	"errors"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// NewRPCClient creates new AMQP RPC client.
func NewRPCClient(
	name, serverExchangeName, serverRoutingKey, clientQueueName string,
	consumePrefetchCount, consumePrefetchSize int,
	workers int,
	connection libamqp.Connection,
	publishCh libamqp.Channel,
	processor RPCMessageProcessor,
	logger zerolog.Logger,
) RPCClient {
	return &rpcClient{
		defaultConsumer: defaultConsumer{
			name:                 name,
			queue:                clientQueueName,
			consumePrefetchCount: consumePrefetchCount,
			consumePrefetchSize:  consumePrefetchSize,
			processor:            &rpcClientMessageProcessorWrapper{processor: processor},
			connection:           connection,
			logger:               logger,
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
	if c.connection == nil {
		return errors.New("connection is nil")
	}

	consumeCh, msgs, err := c.getConsumeChannel()
	if err != nil {
		return err
	}

	defer func() {
		closeError := consumeCh.Close()
		if err == nil {
			err = closeError
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.workers; i++ {
		g.Go(c.getWorkerFunc(ctx, msgs, consumeCh, nil))
	}

	return g.Wait()
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
