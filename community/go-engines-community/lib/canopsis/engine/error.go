package engine

import (
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	amqp "github.com/rabbitmq/amqp091-go"
)

// IsConnectionError uses to check if stop engine or continue work.
func IsConnectionError(err error) bool {
	return mongo.IsConnectionError(err) ||
		redis.IsConnectionError(err) ||
		errors.Is(err, amqp.ErrClosed)
}
