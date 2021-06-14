package storage

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/go-redis/redis/v7"
)

type GroupingStorage interface {
	Push(*redis.Tx, metaalarm.Rule, types.Alarm, string) error
	CleanPush(*redis.Tx, metaalarm.Rule, types.Alarm, string) error
	Clean(*redis.Tx, string) error
	Get(*redis.Tx, string) (AlarmGroup, error)
	GetGroupLen(*redis.Tx, string) (int64, error)
}
