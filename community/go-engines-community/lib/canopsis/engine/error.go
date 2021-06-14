package engine

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
)

// IsConnectionError uses to check if stop engine or continue work.
func IsConnectionError(err error) bool {
	return mongo.IsConnectionError(err) || redis.IsConnectionError(err)
}
