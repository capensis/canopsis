package statistics

import (
	"context"
)

type StatsListener interface {
	Listen(ctx context.Context, channel <-chan Message)
}
