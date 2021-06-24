package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

//TODO: refactor all applicators with the new RedisGroupingStorage

type redisGroupingStorageNew struct{}

func (s *redisGroupingStorageNew) SetMany(ctx context.Context, tx *redis.Tx, timeInterval int64, alarmGroups ...TimeBasedAlarmGroup) error {
	pipe := tx.TxPipeline()

	for _, group := range alarmGroups {
		pipe.Del(ctx, group.GetKey())
		pipe.Set(ctx, group.GetKey(), group, time.Duration(timeInterval) * time.Second)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (s *redisGroupingStorageNew) Set(ctx context.Context, tx *redis.Tx, key string, alarmGroup TimeBasedAlarmGroup, timeInterval int64) error {
	pipe := tx.TxPipeline()

	pipe.Del(ctx, key)
	pipe.Set(ctx, key, alarmGroup, time.Duration(timeInterval) * time.Second)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *redisGroupingStorageNew) Clean(ctx context.Context, tx *redis.Tx, ruleID string) error {
	pipe := tx.TxPipeline()

	pipe.Del(ctx, ruleID)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *redisGroupingStorageNew) Get(ctx context.Context, tx *redis.Tx, key string) (TimeBasedAlarmGroup, error) {
	group := NewAlarmGroup(key)

	res := tx.Get(ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return group, nil
		}

		return group, err
	}

	err := json.Unmarshal([]byte(res.Val()), &group)
	return group, err
}

func NewRedisGroupingStorageNew() GroupingStorageNew {
	return &redisGroupingStorageNew{}
}
