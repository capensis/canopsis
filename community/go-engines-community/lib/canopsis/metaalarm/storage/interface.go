package storage

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-redis/redis/v8"
)

type TimeBasedAlarmGroup interface {
	GetKey() string
	GetAlarmIds() []string
	GetTimes() []int64
	GetGroupLength() int
	GetOpenTime() int64
	Push(newAlarm types.Alarm, ruleTimeInterval int64)
	RemoveBefore(timestamp int64)
}

type GroupingStorageNew interface {
	SetMany(ctx context.Context, tx *redis.Tx, timeInterval int64, alarmGroups ...TimeBasedAlarmGroup) error
	Set(ctx context.Context, tx *redis.Tx, alarmGroup TimeBasedAlarmGroup, timeInterval int64) error
	Clean(ctx context.Context, tx *redis.Tx, ruleID string) error
	Get(ctx context.Context, tx *redis.Tx, key string) (TimeBasedAlarmGroup, error)
}

type GroupingStorage interface {
	Push(context.Context, *redis.Tx, metaalarm.Rule, types.Alarm, string) error
	CleanPush(context.Context, *redis.Tx, metaalarm.Rule, types.Alarm, string) error
	Clean(context.Context, *redis.Tx, string) error
	Get(context.Context, *redis.Tx, string) (AlarmGroup, error)
	GetGroupLen(context.Context, *redis.Tx, string) (int64, error)
}
