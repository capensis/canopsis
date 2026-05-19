package engine

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
)

// NewRPCServer creates consumer.
func NewRPCServer(
	name, queue string,
	workers int,
	publishCh libamqp.Channel,
	consumeCh libamqp.Channel,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &rpcServer{
		defaultConsumer: defaultConsumer{
			name:      name,
			queue:     queue,
			publishCh: publishCh,
			consumeCh: consumeCh,
			processor: processor,
			logger:    logger,
		},
		workers: workers,
	}
}

// rpcServer implements AMQP consumer of RPC requests.
type rpcServer struct {
	defaultConsumer
	// amount of workers which process events.
	workers int
}

func (c *rpcServer) Consume(ctx context.Context) error {
	err := c.queuePurge(ctx)
	if err != nil {
		return err
	}

	return c.consume(ctx, c.workers)
}
