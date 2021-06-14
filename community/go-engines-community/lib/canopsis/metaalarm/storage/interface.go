package storage

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/go-redis/redis/v8"
)

type GroupingStorage interface {
	Push(context.Context, *redis.Tx, metaalarm.Rule, types.Alarm, string) error
	CleanPush(context.Context, *redis.Tx, metaalarm.Rule, types.Alarm, string) error
	Clean(context.Context, *redis.Tx, string) error
	Get(context.Context, *redis.Tx, string) (AlarmGroup, error)
	GetGroupLen(context.Context, *redis.Tx, string) (int64, error)
}
