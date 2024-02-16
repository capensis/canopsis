package engine

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

func NewConcurrentConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	workers int,
	connection libamqp.Connection,
	processor MessageProcessor,
	logger zerolog.Logger,
) Consumer {
	return &concurrentConsumer{
		name:                 name,
		queue:                queue,
		consumePrefetchCount: consumePrefetchCount,
		consumePrefetchSize:  consumePrefetchSize,
		purgeQueue:           purgeQueue,
		nextExchange:         nextExchange,
		nextQueue:            nextQueue,
		fifoExchange:         fifoExchange,
		fifoQueue:            fifoQueue,
		workers:              workers,
		connection:           connection,
		processor:            processor,
		logger:               logger,
	}
}

type concurrentConsumer struct {
	// name is consumer name.
	name string
	// queue is name of AMQP queue from where consumer receives messages.
	queue                                     string
	consumePrefetchCount, consumePrefetchSize int
	purgeQueue                                bool
	// processor handles AMQP messages.
	processor MessageProcessor
	// nextQueue is name of AMQP queue to where consumer sends message after succeeded processing.
	nextQueue    string
	nextExchange string
	// fifoQueue is name of AMQP queue to where consumer sends message after failed processing
	// or if nextQueue is not defined.
	fifoQueue    string
	fifoExchange string
	// amount of workers which process events
	workers int
	// connection is AMQP connection.
	connection libamqp.Connection
	logger     zerolog.Logger
}

func (c *concurrentConsumer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := getConsumeChannel(c.connection, c.name, c.queue,
		c.consumePrefetchCount, c.consumePrefetchSize, c.purgeQueue)
	if err != nil {
		return err
	}

	var publishCh libamqp.Channel
	if c.nextQueue != "" || c.fifoQueue != "" {
		publishCh, err = c.connection.Channel()
		if err != nil {
			return err
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
		g.Go(func() (resErr error) {
			defer func() {
				if r := recover(); r != nil {
					var err error
					var ok bool
					if err, ok = r.(error); !ok {
						err = fmt.Errorf("%v", r)
					}

					c.logger.Err(err).Msgf("consumer recovered from panic\n%s\n", debug.Stack())
					resErr = fmt.Errorf("consumer recovered from panic: %w", err)
				}
			}()

			for {
				select {
				case <-ctx.Done():
					return nil
				case d, ok := <-msgs:
					if !ok {
						return errors.New("the rabbitmq channel has been closed")
					}

					err := c.processMessage(ctx, d, consumeCh, publishCh)
					if err != nil {
						return err
					}
				}
			}
		})
	}

	return g.Wait()
}

func (c *concurrentConsumer) processMessage(ctx context.Context, d amqp.Delivery, consumeCh, publishCh libamqp.Channel) error {
	c.logger.Debug().
		Str("consumer", c.name).Str("queue", c.queue).
		Str("msg", string(d.Body)).
		Msgf("received")
	msgToNext, err := c.processor.Process(ctx, d)

	if err != nil {
		nackErr := consumeCh.Nack(d.DeliveryTag, false, true)
		if nackErr != nil {
			c.logger.Err(nackErr).Msg("cannot nack amqp delivery")
		}

		return fmt.Errorf("cannot process message: %w", err)
	}

	err = consumeCh.Ack(d.DeliveryTag, false)
	if err != nil {
		c.logger.Err(err).Msg("cannot ack amqp delivery")
	}

	if c.nextQueue != "" && msgToNext != nil {
		err := publishToChannel(ctx, publishCh, c.nextExchange, c.nextQueue, msgToNext)
		if err != nil {
			return fmt.Errorf("cannot sent message to next queue: %w", err)
		}
	} else if c.fifoQueue != "" {
		err := publishToChannel(ctx, publishCh, c.fifoExchange, c.fifoQueue, d.Body)
		if err != nil {
			return fmt.Errorf("cannot sent message to fifo queue: %w", err)
		}
	}

	return nil
}
