package ratelimit

import "context"

type Adapter interface {
	DeleteBefore(context.Context, int64) (int64, error)
}
