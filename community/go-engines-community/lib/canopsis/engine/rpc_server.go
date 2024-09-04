package engine

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// NewRPCServer creates consumer.
func NewRPCServer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	workers int,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &rpcServer{
		defaultConsumer: defaultConsumer{
			name:                 name,
			queue:                queue,
			consumePrefetchCount: consumePrefetchCount,
			consumePrefetchSize:  consumePrefetchSize,
			connection:           connection,
			processor:            processor,
			logger:               logger,
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
	consumeCh, msgs, err := c.getConsumeChannel()
	if err != nil {
		return err
	}

	publishCh, err := c.connection.Channel()
	if err != nil {
		return fmt.Errorf("cannot open channel: %w", err)
	}

	defer func() {
		_ = consumeCh.Close()
		_ = publishCh.Close()
	}()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.workers; i++ {
		g.Go(c.getWorkerFunc(ctx, msgs, consumeCh, publishCh))
	}

	return g.Wait()
}
