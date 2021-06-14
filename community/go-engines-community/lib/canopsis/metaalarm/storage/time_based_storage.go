package storage

import (
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"

	"github.com/go-redis/redis/v7"
)

type RedisGroupingStorage struct {}

func (s *RedisGroupingStorage) Push(tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, key string) error {
	return s.set(tx, rule, newAlarm, false, key)
}

func (s *RedisGroupingStorage) CleanPush(tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, key string) error {
	return s.set(tx, rule, newAlarm, true, key)
}

func (s *RedisGroupingStorage) set(tx *redis.Tx, rule metaalarm.Rule, newAlarm types.Alarm, clean bool, key string) error {
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
		alarmGroup, err = s.Get(tx, key)
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

	pipe.Del(key)
	pipe.RPush(key, EncodeGroup(alarmGroup))
	pipe.Expire(key, time.Duration(rule.Config.TimeInterval) * time.Second)

	_, err = pipe.Exec()
	return err
}

func (s *RedisGroupingStorage) Clean(tx *redis.Tx, ruleID string) error {
	pipe := tx.TxPipeline()

	pipe.Del(ruleID)

	_, err := pipe.Exec()
	return err
}

func (s *RedisGroupingStorage) Get(tx *redis.Tx, ruleID string) (AlarmGroup, error) {
	var err error

	getGroupResult := tx.LRange(ruleID, 0, -1)
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

func (s *RedisGroupingStorage) GetGroupLen(tx *redis.Tx, ruleID string) (int64, error) {
	getLenResult := tx.LLen(ruleID)
	if err := getLenResult.Err(); err != nil {
		return 0, err
	}

	return getLenResult.Val(), nil
}

func NewRedisGroupingStorage() *RedisGroupingStorage {
	var storage RedisGroupingStorage

	return &storage
}
