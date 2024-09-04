package engine

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"

	libamqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/valyala/fastjson"
	"golang.org/x/sync/errgroup"
)

func NewDivergingConsumer(
	name, queue string,
	consumePrefetchCount, consumePrefetchSize int,
	purgeQueue bool,
	nextExchange, nextQueue, fifoExchange, fifoQueue string,
	workers int,
	connection libamqp.Connection,
	processor MessageProcessor,
	routingField string,
	routingKeys []string,
	logger zerolog.Logger,
) Consumer {
	return &divergingConsumer{
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
			processor:            processor,
			connection:           connection,
			logger:               logger,
		},
		workers:      workers,
		routingField: routingField,
		routingKeys:  routingKeys,
	}
}

type divergingConsumer struct {
	defaultConsumer
	// amount of workers which process events.
	workers int
	// routingField is used to route a message.
	routingField string
	routingKeys  []string
}

func (c *divergingConsumer) Consume(ctx context.Context) error {
	consumeCh, msgs, err := c.getConsumeChannel()
	if err != nil {
		return err
	}

	var publishCh libamqp.Channel
	if c.nextQueue != "" || c.fifoQueue != "" {
		publishCh, err = c.connection.Channel()
		if err != nil {
			return fmt.Errorf("cannot open channel: %w", err)
		}
	}

	defer func() {
		_ = consumeCh.Close()
		if publishCh != nil {
			_ = publishCh.Close()
		}
	}()

	g, ctx := errgroup.WithContext(ctx)
	workersPerRoute := c.workers / (len(c.routingKeys) + 1)
	if workersPerRoute == 0 {
		workersPerRoute = 1
	}

	chansByRoute := make(map[string]chan amqp.Delivery, len(c.routingKeys))
	fallbackCh := make(chan amqp.Delivery)
	for _, routingKey := range c.routingKeys {
		ch := make(chan amqp.Delivery)
		chansByRoute[routingKey] = ch
		for i := 0; i < workersPerRoute; i++ {
			g.Go(c.getWorkerFunc(ctx, ch, consumeCh, publishCh))
		}
	}

	for i := 0; i < workersPerRoute; i++ {
		g.Go(c.getWorkerFunc(ctx, fallbackCh, consumeCh, publishCh))
	}

	defer func() {
		close(fallbackCh)
		for _, ch := range chansByRoute {
			close(ch)
		}
	}()

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

					msg, err := fastjson.ParseBytes(d.Body)
					if err != nil {
						c.logger.Error().Msg("message is not JSON")
						err = consumeCh.Ack(d.DeliveryTag, false)
						if err != nil {
							c.logger.Err(err).Msg("cannot ack amqp delivery")
						}

						continue
					}

					routingKey := msg.GetStringBytes(c.routingField)
					ch, ok := chansByRoute[string(routingKey)]
					if !ok {
						ch = fallbackCh
					}

					select {
					case <-ctx.Done():
						return nil
					case ch <- d:
					}
				}
			}
		})
	}

	return g.Wait()
}
