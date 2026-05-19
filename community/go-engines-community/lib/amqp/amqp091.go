package amqp

import (
	"context"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

// The interfaces below are narrow internal seams over the amqp091 library types.
// They exist so the reconnect wrappers (connection, channel) can be unit-tested
// against generated mocks instead of a real broker -
// production code uses the adapters at the bottom of this file, which forward to *amqp.Connection and *amqp.Channel.
// Each interface lists exactly the subset of methods the wrappers actually call; do not add methods here speculatively.

//go:generate go tool go.uber.org/mock/mockgen -package=amqp -destination=./amqp091_mock.go git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/amqp amqpDialer,amqp091Conn,amqp091Channel

// amqpDialer opens a fresh AMQP connection.
// In production this is *defaultDialer, which calls amqp.DialConfig with the URL read from EnvURL at construction.
type amqpDialer interface {
	Dial() (amqp091Conn, error)
}

// amqp091Conn is the subset of *amqp.Connection used by the connection wrapper.
type amqp091Conn interface {
	Channel() (amqp091Channel, error)
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	Close() error
}

// amqp091Channel is the subset of *amqp.Channel used by the channel wrapper
// (and forwarded through the public Channel interface to callers).
type amqp091Channel interface {
	ConsumeWithContext(ctx context.Context, queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple, requeue bool) error
	Reject(tag uint64, requeue bool) error
	PublishWithContext(ctx context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Qos(prefetchCount, prefetchSize int, global bool) error
	Close() error
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueDeclarePassive(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	QueuePurge(name string, noWait bool) (int, error)
	QueueDelete(name string, ifUnused, ifEmpty, noWait bool) (int, error)
}

// defaultDialer is the production amqpDialer.
type defaultDialer struct {
	url    string
	config amqp.Config
}

func newDefaultDialer(config amqp.Config) (amqpDialer, error) {
	url := os.Getenv(EnvURL)
	if url == "" {
		return nil, fmt.Errorf("environment variable %s empty", EnvURL)
	}

	return &defaultDialer{
		url:    url,
		config: config,
	}, nil
}

func (d *defaultDialer) Dial() (amqp091Conn, error) {
	amqpConn, err := amqp.DialConfig(d.url, d.config)
	if err != nil {
		return nil, err
	}

	return &amqp091ConnAdapter{Connection: amqpConn}, nil
}

// amqp091ConnAdapter adapts *amqp.Connection to amqp091Conn -
// the only purpose is to change the return type of Channel() from *amqp.Channel to amqp091Channel
// so the wrapper code can stay interface-typed.
type amqp091ConnAdapter struct {
	*amqp.Connection
}

func (c *amqp091ConnAdapter) Channel() (amqp091Channel, error) {
	return c.Connection.Channel()
}
