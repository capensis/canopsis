package storage

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-redis/redis/v8"
)

type TimeBasedAlarmGroup interface {
	GetKey() string
	GetAlarmIds() []string
	GetTimes() []int64
	GetGroupLength() int64
	GetOpenTime() int64
	Push(newAlarm types.Alarm, ruleTimeInterval int64)
	RemoveBefore(timestamp int64)
}

type GroupingStorage interface {
	SetMany(ctx context.Context, tx *redis.Tx, timeInterval int64, alarmGroups ...TimeBasedAlarmGroup) error
	Set(ctx context.Context, tx *redis.Tx, alarmGroup TimeBasedAlarmGroup, timeInterval int64) error
	Clean(ctx context.Context, tx *redis.Tx, ruleID string) error
	Get(ctx context.Context, tx *redis.Tx, key string) (TimeBasedAlarmGroup, error)
}
