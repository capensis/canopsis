package engine

import (
	"context"
	"fmt"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

func NewConcurrentConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	workers int,
	exclusive bool,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &concurrentConsumer{
		defaultConsumer: defaultConsumer{
			name:                 name,
			queue:                queue,
			consumePrefetchCount: consumePrefetchCount,
			consumePrefetchSize:  consumePrefetchSize,
			purgeQueue:           purgeQueue,
			nextExchange:         nextExchange,
			nextQueue:            nextQueue,
			fifoExchange:         fifoExchange,
			fifoQueue:            fifoQueue,
			exclusive:            exclusive,
			processor:            processor,
			connection:           connection,
			logger:               logger,
		},
		workers: workers,
	}
}

type concurrentConsumer struct {
	defaultConsumer
	// amount of workers which process events.
	workers int
}

func (c *concurrentConsumer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := c.getConsumeChannel()
	if err != nil {
		return err
	}

	var publishCh libamqp.Channel
	if c.nextQueue != "" || c.fifoQueue != "" {
		publishCh, err = c.connection.Channel()
		if err != nil {
			return fmt.Errorf("unable to open channel: %w", err)
		}
	}

	defer func() {
		_ = consumeCh.Close()
		if publishCh != nil {
			_ = publishCh.Close()
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < c.workers; i++ {
		g.Go(c.getWorkerFunc(ctx, msgs, consumeCh, publishCh))
	}

	return g.Wait()
}
