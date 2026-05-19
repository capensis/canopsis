package engine

import (
	"context"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	"github.com/rs/zerolog"
)

func NewConcurrentConsumer(
	name, queue string,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	workers int,
	exclusive bool,
	publishCh libamqp.Channel,
	consumeCh libamqp.Channel,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &concurrentConsumer{
		defaultConsumer: defaultConsumer{
			name:         name,
			queue:        queue,
			purgeQueue:   purgeQueue,
			nextExchange: nextExchange,
			nextQueue:    nextQueue,
			fifoExchange: fifoExchange,
			fifoQueue:    fifoQueue,
			exclusive:    exclusive,
			processor:    processor,
			publishCh:    publishCh,
			consumeCh:    consumeCh,
			logger:       logger,
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
	err := c.queuePurge(ctx)
	if err != nil {
		return err
	}

	return c.consume(ctx, c.workers)
}
