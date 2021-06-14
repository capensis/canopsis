package storage

import (
	"context"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"

	"github.com/go-redis/redis/v8"
)

type RedisGroupingStorage struct {}

func (s *RedisGroupingStorage) Push(ctx context.Context, tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, key string) error {
	return s.set(ctx, tx, rule, newAlarm, false, key)
}

func (s *RedisGroupingStorage) CleanPush(ctx context.Context, tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, key string) error {
	return s.set(ctx, tx, rule, newAlarm, true, key)
}

func (s *RedisGroupingStorage) set(ctx context.Context, tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, clean bool, key string) error {
	var alarmGroup AlarmGroup
	var err error

	if key == "" {
		key = rule.ID
	}

	if clean {
		alarmGroup = AlarmGroup{
			OpenTime: types.CpsTime{},
			Group:    make(map[string]types.CpsTime),
		}
	} else {
		alarmGroup, err = s.Get(ctx, tx, key)
		if err != nil {
			return err
		}
	}

	if rule.Config.TimeInterval == 0 {
		rule.Config.TimeInterval = 86400
	}

	ruleTimeInterval := time.Duration(rule.Config.TimeInterval) * time.Second
	newAlarmTime := newAlarm.Value.LastUpdateDate

	if len(alarmGroup.Group) != 0 {
		openTime := alarmGroup.OpenTime.Time

		//if alarm is late
		if newAlarmTime.Before(openTime) {
			newIntervalEnd := newAlarmTime.Add(ruleTimeInterval)
			//check if interval can be shifted
			for _, alarmTime := range alarmGroup.Group {
				//if any alarm in the Group will be lost => then we cannot shift time
				if alarmTime.After(newIntervalEnd) {
					return nil
				}
			}
		} else if newAlarmTime.After(openTime.Add(ruleTimeInterval)) {
			newOpenTime := newAlarmTime
			newIntervalStart := newAlarmTime.Add(-ruleTimeInterval)
			for alarmID, alarmTime := range alarmGroup.Group {
				if newIntervalStart.After(alarmTime.Time) {
					delete(alarmGroup.Group, alarmID)
				} else if newOpenTime.After(alarmTime.Time) {
					newOpenTime = alarmTime
				}
			}
		}
	}

	alarmGroup.Group[newAlarm.ID] = newAlarmTime

	pipe := tx.TxPipeline()

	pipe.Del(ctx, key)
	pipe.RPush(ctx, key, EncodeGroup(alarmGroup))
	pipe.Expire(ctx, key, time.Duration(rule.Config.TimeInterval) * time.Second)

	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisGroupingStorage) Clean(ctx context.Context, tx *redis.Tx, ruleID string) error {
	pipe := tx.TxPipeline()

	pipe.Del(ctx, ruleID)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisGroupingStorage) Get(ctx context.Context, tx *redis.Tx, ruleID string) (AlarmGroup, error) {
	var err error

	getGroupResult := tx.LRange(ctx, ruleID, 0, -1)
	if err = getGroupResult.Err(); err != nil {
		return AlarmGroup{}, err
	}

	var group []string
	err = getGroupResult.ScanSlice(&group)
	if err != nil {
		return AlarmGroup{}, err
	}

	decodedGroup, err := DecodeGroup(group)
	if err != nil {
		return AlarmGroup{}, err
	}

	return decodedGroup, nil
}

func (s *RedisGroupingStorage) GetGroupLen(ctx context.Context, tx *redis.Tx, ruleID string) (int64, error) {
	getLenResult := tx.LLen(ctx, ruleID)
	if err := getLenResult.Err(); err != nil {
		return 0, err
	}

	return getLenResult.Val(), nil
}

func NewRedisGroupingStorage() *RedisGroupingStorage {
	var storage RedisGroupingStorage

	return &storage
}
