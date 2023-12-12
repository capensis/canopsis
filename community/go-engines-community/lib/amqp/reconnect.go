package amqp

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

const (
	maxReconnectTimeout             = time.Minute
	waitReconnectionTimeoutOverhead = 500 * time.Millisecond
)

// Dial accepts a string in the AMQP URI format and returns a new amqp connection.
// If connection is closed it tries to reconnect.
func Dial(url string, logger zerolog.Logger,
	reconnectCount int, minReconnectTimeout time.Duration) (Connection, error) {
	amqpConn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	if reconnectCount < -1 {
		return nil, fmt.Errorf("invalid reconnectCount %v, can be -1 or possitive int", reconnectCount)
	}

	isReconnectable := reconnectCount != 0 && minReconnectTimeout > 0

	var waitReconnectionTimeout time.Duration
	if isReconnectable && reconnectCount > 0 {
		for i := 0; i < reconnectCount; i++ {
			waitReconnectionTimeout += minReconnectTimeout << i
			if waitReconnectionTimeout > maxReconnectTimeout {
				waitReconnectionTimeout = maxReconnectTimeout
				break
			}
		}

		waitReconnectionTimeout = minTimeDuration(waitReconnectionTimeout+waitReconnectionTimeoutOverhead, maxReconnectTimeout)
	}

	conn := &baseConnection{
		amqpConn:                amqpConn,
		logger:                  logger,
		isReconnectable:         isReconnectable,
		reconnectCount:          reconnectCount,
		minReconnectTimeout:     minTimeDuration(minReconnectTimeout, maxReconnectTimeout),
		waitReconnectionTimeout: waitReconnectionTimeout,
	}

	if conn.isReconnectable {
		go reconnect(url, conn)
	}

	return conn, nil
}

// baseConnection is used to wrap amqp connection to reconnect it each time
// when connection is closed.
type baseConnection struct {
	// amqpConn contains real connection.
	amqpConn *amqp.Connection
	// listeners is used to notify opened channels that connection is reconnected.
	listenersM sync.Mutex
	listeners  []chan<- bool
	// closed is true if connection is closed normally or if reconnection failed after max retries.
	closedM                 sync.Mutex
	closed                  bool
	logger                  zerolog.Logger
	isReconnectable         bool
	reconnectCount          int
	minReconnectTimeout     time.Duration
	waitReconnectionTimeout time.Duration
}

// Channel opens new amqp channel.
// If connection or channel is closed it tries to reopen channel.
func (c *baseConnection) Channel() (res Channel, resErr error) {
	// Add listeners before check on isClosed to prevent race condition.
	listener := make(chan bool, 1)
	c.addListener(listener)

	defer func() {
		if resErr != nil {
			c.removeListener(listener)
			close(listener)
		}
	}()

	if c.IsClosed() {
		return nil, amqp.ErrClosed
	}

	amqpCh, err := c.amqpConn.Channel()
	if err != nil {
		if c.isReconnectable && errors.Is(err, amqp.ErrClosed) {
			if c.waitReconnectionTimeout == 0 {
				reconnected := <-listener
				if reconnected {
					amqpCh, err = c.amqpConn.Channel()
				}
			} else {
				select {
				case reconnected := <-listener:
					if reconnected {
						amqpCh, err = c.amqpConn.Channel()
					}
				case <-time.After(c.waitReconnectionTimeout):
					/* add timeout to prevent infinity waiting on bug */
				}
			}
		}

		if err != nil {
			return nil, err
		}
	}

	ch := &baseChannel{
		amqpCh:                  amqpCh,
		logger:                  c.logger,
		waitReconnectionTimeout: c.waitReconnectionTimeout,
		isReconnectable:         c.isReconnectable,
	}

	if ch.isReconnectable {
		go reconnectChannel(c, ch, listener)
	}

	return ch, nil
}

// IsClosed returns true if the connection is marked as closed.
func (c *baseConnection) IsClosed() bool {
	c.closedM.Lock()
	defer c.closedM.Unlock()

	return c.closed
}

// Close closes connection.
func (c *baseConnection) Close() error {
	c.closedM.Lock()
	defer c.closedM.Unlock()

	if c.closed {
		return amqp.ErrClosed
	}

	c.closed = true

	return c.amqpConn.Close()
}

// addListener adds listener.
func (c *baseConnection) addListener(listener chan<- bool) {
	c.listenersM.Lock()
	defer c.listenersM.Unlock()
	c.listeners = append(c.listeners, listener)
}

// removeListener removes listener.
func (c *baseConnection) removeListener(listener chan<- bool) {
	c.listenersM.Lock()
	defer c.listenersM.Unlock()
	removeIndex := -1
	for i := range c.listeners {
		if c.listeners[i] == listener {
			removeIndex = i
		}
	}

	if removeIndex >= 0 {
		c.listeners = append(c.listeners[:removeIndex], c.listeners[removeIndex+1:]...)
	}
}

// notifyListeners sends reconnection result to opened channels.
func (c *baseConnection) notifyListeners(reconnected bool) {
	c.listenersM.Lock()
	defer c.listenersM.Unlock()

	for _, listener := range c.listeners {
		listener <- reconnected
	}
}

// baseChannel is used to wrap amqp channel to reconnect it each time
// when connection or channel is closed.
type baseChannel struct {
	// amqpCh contains real channel.
	amqpCh *amqp.Channel
	// listeners is used to notify consumers and publishers that channel is reconnected.
	listenersM sync.Mutex
	listeners  []chan<- bool
	// closed is true if channel is closed normally or if reconnection failed.
	closedM                 sync.Mutex
	closed                  bool
	logger                  zerolog.Logger
	isReconnectable         bool
	waitReconnectionTimeout time.Duration
}

// Consume starts delivering queued messages.
// If connection or channel is closed it waits reconnection.
func (ch *baseChannel) Consume(
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (resCh <-chan amqp.Delivery, resErr error) {
	// Add listeners before check on isClosed to prevent race condition.
	listener := make(chan bool, 1)
	ch.addListener(listener)

	defer func() {
		if resErr != nil {
			ch.removeListener(listener)
			close(listener)
		}
	}()

	if ch.IsClosed() {
		return nil, amqp.ErrClosed
	}

	// Create custom delivery channel to manage its closing.
	res := make(chan amqp.Delivery)

	go func() {
		defer close(res)
		defer func() {
			ch.removeListener(listener)
			close(listener)
		}()

		for {
			msgs, err := ch.amqpCh.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				if errors.Is(err, amqp.ErrClosed) {
					if !ch.isReconnectable || ch.IsClosed() {
						return
					}

					ch.logger.Debug().Msgf("consume is waiting reconnection")

					if ch.waitReconnectionTimeout == 0 {
						reconnected := <-listener
						if !reconnected {
							return
						}
						msgs, err = ch.amqpCh.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
					} else {
						select {
						case reconnected := <-listener:
							if !reconnected {
								return
							}
							msgs, err = ch.amqpCh.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
						case <-time.After(ch.waitReconnectionTimeout):
							/* add timeout to prevent infinity waiting on bug */
						}
					}

					if err == nil {
						ch.logger.Debug().Msgf("consume restarted after reconnect")
					}
				}

				if err != nil {
					ch.logger.Err(err).Msgf("consume failed")
					break
				}
			}

			// Sends messages to custom channel.
			for msg := range msgs {
				res <- msg
			}
		}
	}()

	return res, nil
}

// PublishWithContext sends a message.
// If connection or channel is closed it waits reconnection.
func (ch *baseChannel) PublishWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	return ch.retry(func() error {
		return ch.amqpCh.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
	})
}

// Ack acknowledges a delivery
func (ch *baseChannel) Ack(tag uint64, multiple bool) error {
	return ch.retry(func() error {
		return ch.amqpCh.Ack(tag, multiple)
	})
}

func (ch *baseChannel) Nack(tag uint64, multiple bool, requeue bool) error {
	return ch.retry(func() error {
		return ch.amqpCh.Nack(tag, multiple, requeue)
	})
}

func (ch *baseChannel) Qos(prefetchCount, prefetchSize int, global bool) error {
	return ch.retry(func() error {
		return ch.amqpCh.Qos(prefetchCount, prefetchSize, global)
	})
}

func (ch *baseChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return ch.retry(func() error {
		return ch.amqpCh.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	})
}

func (ch *baseChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	var queue amqp.Queue
	err := ch.retry(func() error {
		var err error
		queue, err = ch.amqpCh.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
		return err
	})

	return queue, err
}

func (ch *baseChannel) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return ch.retry(func() error {
		return ch.amqpCh.QueueBind(name, key, exchange, noWait, args)
	})
}

func (ch *baseChannel) QueuePurge(name string, noWait bool) (int, error) {
	var res int
	err := ch.retry(func() error {
		var err error
		res, err = ch.amqpCh.QueuePurge(name, noWait)
		return err
	})

	return res, err
}

func (ch *baseChannel) QueueInspect(name string) (amqp.Queue, error) {
	var res amqp.Queue
	err := ch.retry(func() error {
		var err error
		res, err = ch.amqpCh.QueueDeclarePassive(name, false, false, false, false, nil)
		return err
	})

	return res, err
}

// IsClosed returns true if the channel is marked as closed.
func (ch *baseChannel) IsClosed() bool {
	ch.closedM.Lock()
	defer ch.closedM.Unlock()

	return ch.closed
}

// Close closes channel.
func (ch *baseChannel) Close() error {
	ch.closedM.Lock()
	defer ch.closedM.Unlock()

	if ch.closed {
		return amqp.ErrClosed
	}

	ch.closed = true

	return ch.amqpCh.Close()
}

func (ch *baseChannel) retry(f func() error) error {
	if !ch.isReconnectable {
		return f()
	}

	// Add listeners before check on isClosed to prevent race condition.
	listener := make(chan bool, 1)
	ch.addListener(listener)

	defer func() {
		ch.removeListener(listener)
		close(listener)
	}()

	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	err := f()
	if errors.Is(err, amqp.ErrClosed) {
		if ch.waitReconnectionTimeout == 0 {
			reconnected := <-listener
			if reconnected {
				err = f()
			}
		} else {
			select {
			case reconnected := <-listener:
				if reconnected {
					err = f()
				}
			case <-time.After(ch.waitReconnectionTimeout):
				/* add timeout to prevent infinity waiting on bug */
			}
		}
	}

	return err
}

// addListener adds listener.
func (ch *baseChannel) addListener(listener chan<- bool) {
	ch.listenersM.Lock()
	defer ch.listenersM.Unlock()
	ch.listeners = append(ch.listeners, listener)
}

// removeListener removes listener.
func (ch *baseChannel) removeListener(listener chan<- bool) {
	ch.listenersM.Lock()
	defer ch.listenersM.Unlock()
	removeIndex := -1
	for i := range ch.listeners {
		if ch.listeners[i] == listener {
			removeIndex = i
		}
	}

	if removeIndex >= 0 {
		ch.listeners = append(ch.listeners[:removeIndex], ch.listeners[removeIndex+1:]...)
	}
}

// notifyListeners sends reconnection result to consumers and publishers.
func (ch *baseChannel) notifyListeners(reconnected bool) {
	ch.listenersM.Lock()
	defer ch.listenersM.Unlock()

	for _, listener := range ch.listeners {
		listener <- reconnected
	}
}

// reconnect listens connection.NotifyClose and tries to restore connection.
// It notifies opened channels about reconnection result.
func reconnect(url string, c *baseConnection) {
	defer func() {
		// Close channel before notify listeners to prevent race condition.
		_ = c.Close()
		c.notifyListeners(false)
	}()

	for {
		closeErr, ok := <-c.amqpConn.NotifyClose(make(chan *amqp.Error))
		if !ok {
			c.logger.Debug().Msgf("connection normally closed")
			break
		}

		c.logger.Err(closeErr).Msgf("connection closed, try to reconnect")
		timeout := c.minReconnectTimeout
		var amqpConn *amqp.Connection
		var err error

		for try := 0; c.reconnectCount == -1 || try < c.reconnectCount; try++ {
			time.Sleep(timeout)
			timeout = minTimeDuration(timeout<<1, maxReconnectTimeout)

			amqpConn, err = amqp.Dial(url)
			if err == nil {
				break
			}

			c.logger.Debug().Err(err).Msgf("%d try to reconnect failed", try+1)
		}

		if err != nil {
			c.logger.Debug().Err(err).Msgf("connection reconnect failed")
			break
		}

		c.logger.Info().Msgf("connection reconnected")
		c.amqpConn = amqpConn
		c.notifyListeners(true)
	}
}

// reconnectChannel listens channel.NotifyClose and tries to restore channel.
// If connection is closed it waits restore connection before try to reopen channel.
// It notifies consumers and publishers about reconnection result.
func reconnectChannel(conn *baseConnection, ch *baseChannel, reconnectListener chan bool) {
	defer func() {
		// Close channel before notify listeners to prevent race condition.
		_ = ch.Close()
		ch.notifyListeners(false)
		conn.removeListener(reconnectListener)
		close(reconnectListener)
	}()

	for {
		closeErr, ok := <-ch.amqpCh.NotifyClose(make(chan *amqp.Error))
		if !ok {
			conn.logger.Debug().Msgf("channel normally closed")
			break
		}

		conn.logger.Err(closeErr).Msgf("channel closed, try to reconnect")
		// Check real connection
		if conn.amqpConn.IsClosed() {
			conn.logger.Debug().Msgf("channel is waiting reconnection")

			if conn.waitReconnectionTimeout == 0 {
				reconnected := <-reconnectListener
				if !reconnected {
					return
				}
			} else {
				select {
				case reconnected := <-reconnectListener:
					if !reconnected {
						return
					}
				case <-time.After(conn.waitReconnectionTimeout):
					/* add timeout to prevent infinity waiting on bug */
					return
				}
			}
		}

		amqpCh, err := conn.amqpConn.Channel()
		if err != nil {
			ch.logger.Debug().Err(err).Msgf("reconnect channel failed")
			break
		}

		conn.logger.Info().Msgf("channel reconnected")
		ch.amqpCh = amqpCh
		ch.notifyListeners(true)
	}
}

func minTimeDuration(l, r time.Duration) time.Duration {
	if l < r {
		return l
	}

	return r
}
