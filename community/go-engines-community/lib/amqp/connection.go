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

var ErrConnectionClosed = errors.New("connection is closed")

type connection struct {
	dialer amqpDialer

	mx       sync.RWMutex
	amqpConn amqp091Conn
	isClosed bool

	closeOnce sync.Once
	doneOnce  sync.Once
	done      chan struct{}

	reconnectCount   int
	reconnectTimeout time.Duration

	reconnNotify chan struct{}

	logger zerolog.Logger
}

func New(
	reconnectCount int,
	reconnectTimeout time.Duration,
	logger zerolog.Logger,
) (Connection, error) {
	d, err := newDefaultDialer(amqp.Config{})
	if err != nil {
		return nil, err
	}

	return newConn(d, reconnectCount, reconnectTimeout, logger), nil
}

func NewConfig(
	reconnectCount int,
	reconnectTimeout time.Duration,
	config amqp.Config,
	logger zerolog.Logger,
) (Connection, error) {
	d, err := newDefaultDialer(config)
	if err != nil {
		return nil, err
	}

	return newConn(d, reconnectCount, reconnectTimeout, logger), nil
}

func newConn(
	dialer amqpDialer,
	reconnectCount int,
	reconnectTimeout time.Duration,
	logger zerolog.Logger,
) *connection {
	conn := &connection{
		dialer:           dialer,
		done:             make(chan struct{}),
		reconnectCount:   reconnectCount,
		reconnectTimeout: reconnectTimeout,
		reconnNotify:     make(chan struct{}),
		logger:           logger,
	}

	go conn.handleReconnect()

	return conn
}

func (c *connection) Channel(ctx context.Context) (Channel, error) {
	amqpConn := c.waitReconnect(ctx.Done(), nil)
	if amqpConn == nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		return nil, ErrConnectionClosed
	}

	return newChannel(c, amqpConn), nil
}

func (c *connection) Close() error {
	c.logger.Debug().CallerSkipFrame(1).Msg("closing connection")
	err := ErrConnectionClosed
	c.closeOnce.Do(func() {
		err = nil

		c.mx.Lock()

		amqpConn := c.amqpConn
		c.amqpConn = nil
		c.isClosed = true
		c.doneOnce.Do(func() { close(c.done) })

		c.mx.Unlock()

		if amqpConn != nil {
			err = amqpConn.Close()
		}
	})

	return err
}

// handleReconnect runs in its own goroutine for the lifetime of the connection.
// It loops: dial → publish the new conn to waiters via reconnNotify → block until the conn dies (notifyClose) or
// the connection is shut down (c.done) → loop.
// Bounded by reconnectCount; reconnectCount==0 means "no retries: give up after the first drop".
func (c *connection) handleReconnect() {
	c.logger.Debug().Msg("starting connection reconnect loop")
	defer c.logger.Debug().Msg("connection reconnect loop stopped")
	attempt := 0
	for {
		select {
		case <-c.done:
			return
		default:
		}

		conn, err := c.dialer.Dial()
		if err != nil {
			attempt++

			c.logger.
				Err(err).
				Str("retries", strconv.Itoa(attempt)+"/"+strconv.Itoa(c.reconnectCount)).
				Msg("cannot connect to amqp server")

			if c.reconnectCount == 0 || attempt == c.reconnectCount {
				c.doneOnce.Do(func() { close(c.done) })

				return
			}

			t := time.NewTimer(c.reconnectTimeout)
			select {
			case <-c.done:
				t.Stop()

				return
			case <-t.C:
				continue
			}
		}

		// Publish the new conn under the lock, then close the previous reconnNotify *outside* the lock.
		// Any goroutine parked in waitReconnect on the old notify wakes up, re-reads c.amqpConn under the lock,
		// and either returns it or loops again (e.g. if it had flagged this conn as stale).

		attempt = 0
		var prevReconnNotify chan struct{}

		c.mx.Lock()

		isClosed := c.isClosed
		if !isClosed {
			c.amqpConn = conn
			prevReconnNotify = c.reconnNotify
			c.reconnNotify = make(chan struct{})
		}

		c.mx.Unlock()

		// Close was called between Dial returning and us taking the lock.
		// Discard the just-dialed conn; the consumer is shutting down.
		if isClosed {
			_ = conn.Close()

			return
		}

		close(prevReconnNotify)

		notifyClose := make(chan *amqp.Error, 1)
		conn.NotifyClose(notifyClose)

		select {
		case <-c.done:
			return
		case <-notifyClose:
		}

		c.mx.Lock()
		c.amqpConn = nil
		c.mx.Unlock()

		if c.reconnectCount == 0 {
			c.doneOnce.Do(func() { close(c.done) })

			return
		}
	}
}

// waitReconnect blocks until a non-stale amqpConn is available, or until the connection is closed, or until abort fires.
// Pass stale to skip a known-bad conn (typical when the previous operation returned amqp.ErrClosed against it).
func (c *connection) waitReconnect(callerDone <-chan struct{}, stale amqp091Conn) amqp091Conn {
	for {
		c.mx.RLock()
		conn := c.amqpConn
		notify := c.reconnNotify
		c.mx.RUnlock()

		if conn != nil && conn != stale {
			return conn
		}

		select {
		case <-c.done:
			return nil
		case <-callerDone:
			return nil
		case <-notify:
		}
	}
}
