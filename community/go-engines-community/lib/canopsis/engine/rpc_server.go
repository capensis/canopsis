package engine

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
)

// NewRPCServer creates consumer.
func NewRPCServer(
	name, queue string,
	prefetchCount, prefetchSize int,
	workers int,
	publisher libamqp.Publisher,
	consumePool libamqp.ChannelPool,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &rpcServer{
		defaultConsumer: defaultConsumer{
			name:          name,
			queue:         queue,
			prefetchCount: prefetchCount,
			prefetchSize:  prefetchSize,
			publisher:     publisher,
			consumePool:   consumePool,
			processor:     processor,
			logger:        logger,
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
	return c.consume(ctx, c.workers)
}
