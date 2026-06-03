package amqp

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

var ErrChannelClosed = errors.New("channel is closed")

type channel struct {
	conn   *connection
	logger zerolog.Logger

	mx          sync.RWMutex
	amqpChannel amqp091Channel
	isClosed    bool

	closeOnce sync.Once
	doneOnce  sync.Once
	done      chan struct{}

	reconnNotify chan struct{}
}

func newChannel(c *connection, amqpConn amqp091Conn) *channel {
	ch := &channel{
		conn:         c,
		logger:       c.logger,
		done:         make(chan struct{}),
		reconnNotify: make(chan struct{}),
	}

	go ch.handleReconnect(amqpConn)

	return ch
}

func (c *channel) ConsumeWithContext(
	ctx context.Context,
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, error) {
	var d <-chan amqp.Delivery
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		d, err = amqpCh.ConsumeWithContext(ctx, queue, consumer, autoAck, exclusive, noLocal, noWait, args)
		if err != nil {
			return err
		}

		return nil
	})

	return d, err
}

func (c *channel) ConsumeWithCloseNotify(
	ctx context.Context,
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args amqp.Table,
) (<-chan amqp.Delivery, <-chan *amqp.Error, error) {
	var d <-chan amqp.Delivery
	var notifyCh chan *amqp.Error
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		d, err = amqpCh.ConsumeWithContext(ctx, queue, consumer, autoAck, exclusive, noLocal, noWait, args)
		if err != nil {
			return err
		}

		notifyCh = make(chan *amqp.Error, 1)
		amqpCh.NotifyClose(notifyCh)

		return nil
	})

	return d, notifyCh, err
}

func (c *channel) PublishWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg amqp.Publishing,
) error {
	return c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		return amqpCh.PublishWithContext(ctx, exchange, key, mandatory, immediate, msg)
	})
}

func (c *channel) Ack(tag uint64, multiple bool) error {
	c.mx.RLock()
	ch := c.amqpChannel
	c.mx.RUnlock()

	if ch == nil {
		return ErrChannelClosed
	}

	return ch.Ack(tag, multiple)
}

func (c *channel) Nack(tag uint64, multiple, requeue bool) error {
	c.mx.RLock()
	ch := c.amqpChannel
	c.mx.RUnlock()

	if ch == nil {
		return ErrChannelClosed
	}

	return ch.Nack(tag, multiple, requeue)
}

func (c *channel) Cancel(consumer string, noWait bool) error {
	c.mx.RLock()
	ch := c.amqpChannel
	c.mx.RUnlock()

	if ch == nil {
		return ErrChannelClosed
	}

	return ch.Cancel(consumer, noWait)
}

func (c *channel) Qos(ctx context.Context, prefetchCount, prefetchSize int, global bool) error {
	return c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		return amqpCh.Qos(prefetchCount, prefetchSize, global)
	})
}

func (c *channel) ExchangeDeclare(ctx context.Context, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		return amqpCh.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
	})
}

func (c *channel) QueueDeclare(ctx context.Context, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	var q amqp.Queue
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		q, err = amqpCh.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)

		return err
	})

	return q, err
}

func (c *channel) QueueDeclarePassive(ctx context.Context, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	var q amqp.Queue
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		q, err = amqpCh.QueueDeclarePassive(name, durable, autoDelete, exclusive, noWait, args)

		return err
	})

	return q, err
}

func (c *channel) QueueBind(ctx context.Context, name, key, exchange string, noWait bool, args amqp.Table) error {
	return c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		return amqpCh.QueueBind(name, key, exchange, noWait, args)
	})
}

func (c *channel) QueuePurge(ctx context.Context, name string, noWait bool) (int, error) {
	var count int
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		count, err = amqpCh.QueuePurge(name, noWait)

		return err
	})

	return count, err
}

func (c *channel) QueueDelete(ctx context.Context, name string, ifUnused, ifEmpty, noWait bool) (int, error) {
	var count int
	err := c.withChannel(ctx, func(amqpCh amqp091Channel) error {
		var err error
		count, err = amqpCh.QueueDelete(name, ifUnused, ifEmpty, noWait)

		return err
	})

	return count, err
}

func (c *channel) Close() error {
	c.logger.Debug().CallerSkipFrame(1).Msg("closing channel")
	var err error
	c.closeOnce.Do(func() {
		c.mx.Lock()

		amqpCh := c.amqpChannel
		c.amqpChannel = nil
		c.isClosed = true
		c.doneOnce.Do(func() { close(c.done) })

		c.mx.Unlock()

		if amqpCh != nil {
			err = amqpCh.Close()
		}
	})

	return err
}

// withChannel runs op against the currently live amqp channel,
// transparently reacquiring it (via waitReconnect) when the previous attempt failed because the channel was closed.
// Errors other than amqp.ErrClosed are returned immediately - they're application errors, not transport ones.
func (c *channel) withChannel(ctx context.Context, op func(amqp091Channel) error) error {
	var stale amqp091Channel
	for {
		amqpCh := c.waitReconnect(ctx, stale)
		if amqpCh == nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}

			return ErrChannelClosed
		}

		err := op(amqpCh)
		if err != nil {
			if errors.Is(err, amqp.ErrClosed) {
				stale = amqpCh

				continue
			}

			return err
		}

		return nil
	}
}

// handleReconnect mirrors connection.handleReconnect at the channel layer.
// Two extra wrinkles:
//  1. amqp.ErrClosed when opening a channel means the underlying connection died -
//     wait for the connection layer to reconnect, then retry on the fresh conn.
//  2. We watch *both* notifyChClose and notifyConnClose so
//     a dropped connection forces us to re-acquire it before opening a channel.
func (c *channel) handleReconnect(amqpConn amqp091Conn) {
	c.logger.Debug().Msg("starting channel reconnect loop")
	defer c.logger.Debug().Msg("channel reconnect loop stopped")
	attempt := 0
	reconnectTimeout := c.conn.reconnectTimeout
	for {
		select {
		case <-c.done:
			return
		default:
		}

		ch, err := amqpConn.Channel()
		if err != nil {
			// Underlying connection died between us being handed it and trying to open a channel on it.
			// Don't burn a retry - wait for the connection layer to surface a fresh one.
			if errors.Is(err, amqp.ErrClosed) {
				amqpConn = c.conn.waitReconnect(c.done, amqpConn)
				if amqpConn == nil {
					c.doneOnce.Do(func() { close(c.done) })

					return
				}

				attempt = 0
				reconnectTimeout = c.conn.reconnectTimeout

				continue
			}

			attempt++

			le := c.logger.Err(err)
			if c.conn.reconnectCount > 0 {
				le = le.Str("retries", strconv.Itoa(attempt)+"/"+strconv.Itoa(c.conn.reconnectCount))
			}

			le.Msg("cannot create amqp channel")

			if c.conn.reconnectCount == 0 || attempt == c.conn.reconnectCount {
				c.doneOnce.Do(func() { close(c.done) })

				return
			}

			t := time.NewTimer(reconnectTimeout)
			select {
			case <-c.done:
				t.Stop()

				return
			case <-t.C:
				reconnectTimeout = min(2*reconnectTimeout, maxReconnectTimeout)

				continue
			}
		}

		attempt = 0
		reconnectTimeout = c.conn.reconnectTimeout

		var prevReconnNotify chan struct{}

		c.mx.Lock()

		isClosed := c.isClosed
		if !isClosed {
			c.amqpChannel = ch
			prevReconnNotify = c.reconnNotify
			c.reconnNotify = make(chan struct{})
		}

		c.mx.Unlock()

		if isClosed {
			_ = ch.Close()

			return
		}

		close(prevReconnNotify)

		notifyConnClose := make(chan *amqp.Error, 1)
		amqpConn.NotifyClose(notifyConnClose)
		notifyChClose := make(chan *amqp.Error, 1)
		ch.NotifyClose(notifyChClose)

		// notifyChClose: only the channel died, the connection is still good.
		// notifyConnClose: the whole connection went;
		//   connDied tells the post-block code to re-acquire a connection before reopening.

		connDied := false
		select {
		case <-c.done:
			return
		case <-notifyChClose:
		case <-notifyConnClose:
			connDied = true
		}

		c.mx.Lock()
		c.amqpChannel = nil
		c.mx.Unlock()

		if c.conn.reconnectCount == 0 {
			c.doneOnce.Do(func() { close(c.done) })

			return
		}

		if connDied {
			amqpConn = c.conn.waitReconnect(c.done, amqpConn)
			if amqpConn == nil {
				c.doneOnce.Do(func() { close(c.done) })

				return
			}
		}
	}
}

func (c *channel) waitReconnect(ctx context.Context, stale amqp091Channel) amqp091Channel {
	for {
		c.mx.RLock()
		amqpChannel := c.amqpChannel
		notify := c.reconnNotify
		c.mx.RUnlock()

		if amqpChannel != nil && amqpChannel != stale {
			return amqpChannel
		}

		select {
		case <-ctx.Done():
			return nil
		case <-c.done:
			return nil
		case <-notify:
		}
	}
}
