package amqp

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

//go:generate go tool go.uber.org/mock/mockgen -destination=../../mocks/lib/amqp/amqp.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp Connection,Channel,Publisher,ChannelPool

// Connection is a long-lived AMQP connection that transparently reconnects on failure. Built by New / NewConfig.
//
// Channel blocks until a live underlying connection is available, then returns a Channel wrapper.
// The ctx bounds that wait: a canceled or deadlined ctx surfaces as ctx.Err();
// a closed Connection surfaces as ErrConnectionClosed.
// Channels returned by Channel outlive ctx - they're managed by their own
// reconnect goroutine until Close is called on them.
//
// Close is idempotent and safe to call from any goroutine. It tears down the underlying connection;
// in-flight Channel callers unblock with ErrConnectionClosed.
type Connection interface {
	Channel(ctx context.Context) (Channel, error)
	Close() error
}

// Channel is a long-lived AMQP channel that transparently reconnects when the
// underlying AMQP channel (or its connection) is closed by the broker.
//
// All methods retry once on amqp.ErrClosed against a freshly reopened channel,
// then surface other errors as-is. ctx cancels the wait for a reconnected channel:
// a canceled ctx returns ctx.Err(); a Close()d wrapper returns ErrChannelClosed.
// Application-level AMQP errors (e.g. NOT_FOUND on a bind against a missing exchange) are returned to the caller -
// they are not retried, since reopening the channel would not fix them.
//
// ConsumeWithContext returns a delivery channel and a notify channel;
// the notify channel fires when the underlying AMQP channel closes,
// signalling that the delivery channel will close imminently and the consumer should
// re-call ConsumeWithContext to obtain a fresh one
// (or use ConsumeWithReconnect / SubscribeWithReconnect, which handle this).
//
// Close is idempotent and safe to call from any goroutine.
type Channel interface {
	ConsumeWithContext(
		ctx context.Context,
		queue, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, error)
	ConsumeWithCloseNotify(
		ctx context.Context,
		queue, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, <-chan *amqp.Error, error)
	PublishWithContext(
		ctx context.Context,
		exchange, key string,
		mandatory, immediate bool,
		msg amqp.Publishing,
	) error
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple, requeue bool) error
	Cancel(consumer string, noWait bool) error
	Qos(ctx context.Context, prefetchCount, prefetchSize int, global bool) error
	Close() error
	ExchangeDeclare(ctx context.Context, name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(ctx context.Context, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(ctx context.Context, name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(ctx context.Context, name, key, exchange string, noWait bool, args amqp.Table) error
	QueuePurge(ctx context.Context, name string, noWait bool) (int, error)
	QueueDelete(ctx context.Context, name string, ifUnused, ifEmpty, noWait bool) (int, error)
}
