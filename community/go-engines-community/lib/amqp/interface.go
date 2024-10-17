package amqp

//go:generate mockgen -destination=../../mocks/lib/amqp/amqp.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp Connection,Channel,Publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connection is used to implement amqp connection.
type Connection interface {
	Channel() (Channel, error)
	IsClosed() bool
	Close() error
}

// Channel is used to implement amqp channel.
type Channel interface {
	Consume(
		queue, consumer string,
		autoAck, exclusive, noLocal, noWait bool,
		args amqp.Table,
	) (<-chan amqp.Delivery, error)
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Reject(tag uint64, requeue bool) error
	PublishWithContext(
		ctx context.Context,
		exchange, key string,
		mandatory, immediate bool,
		msg amqp.Publishing,
	) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueInspect(name string) (amqp.Queue, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
}
