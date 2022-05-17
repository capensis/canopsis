package amqp

import amqp "github.com/rabbitmq/amqp091-go"

// Publisher is an interface that represents a non-consumable AMQP
// channel. This interface is implemented by amqp.Channel. It should be used
// in services that only publish to amqp, in order to be able to test them
// easily by mocking this interface.
type Publisher interface {
	// Publish sends an amqp.Publishing from the client to an exchange on the server.
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}
