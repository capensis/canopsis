package engine

import (
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
)

// IsConnectionError uses to check if stop engine or continue work.
func IsConnectionError(err error) bool {
	return mongo.IsConnectionError(err) || redis.IsConnectionError(err)
}
