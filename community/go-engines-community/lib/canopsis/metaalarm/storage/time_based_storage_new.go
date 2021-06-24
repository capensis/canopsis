package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

//TODO: refactor all applicators with the new RedisGroupingStorage

type RedisGroupingStorageNew struct{}

func (s *RedisGroupingStorageNew) Pipe(ctx context.Context, tx *redis.Tx) redis.Pipeliner {
	return tx.TxPipeline()
}

func (s *RedisGroupingStorageNew) Set(ctx context.Context, tx *redis.Tx, key string, alarmGroup TimeBasedAlarmGroup, timeInterval int64) error {
	var err error

	pipe := tx.TxPipeline()

	pipe.Del(ctx, key)
	pipe.Set(ctx, key, alarmGroup, time.Duration(timeInterval) * time.Second)

	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisGroupingStorageNew) Clean(ctx context.Context, tx *redis.Tx, ruleID string) error {
	pipe := tx.TxPipeline()

	pipe.Del(ctx, ruleID)

	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisGroupingStorageNew) Get(ctx context.Context, tx *redis.Tx, key string) (TimeBasedAlarmGroup, error) {
	var group TimeBasedAlarmGroup

	res := tx.Get(ctx, key)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return TimeBasedAlarmGroup{}, nil
		}

		return TimeBasedAlarmGroup{}, err
	}

	err := json.Unmarshal([]byte(res.Val()), &group)
	return group, err
}

func NewRedisGroupingStorageNew() *RedisGroupingStorageNew {
	var storage RedisGroupingStorageNew

	return &storage
}
