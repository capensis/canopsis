package ratelimit

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Adapter interface {
	DeleteBefore(context.Context, types.CpsTime) (int64, error)
}
