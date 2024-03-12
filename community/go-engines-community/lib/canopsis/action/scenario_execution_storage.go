package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type ScenarioExecutionStorage interface {
	Get(ctx context.Context, key string) (*ScenarioExecution, error)
	GetAbandoned(ctx context.Context) ([]ScenarioExecution, error)
	Create(ctx context.Context, execution ScenarioExecution) (bool, error)
	Update(ctx context.Context, execution ScenarioExecution) error
	Del(ctx context.Context, key string) error
	Inc(ctx context.Context, key string, inc int64, drop bool) (int64, error)
}

type redisScenarioExecutionStorage struct {
	redisKeyPrefix string
	redisClient    redis.Cmdable
	encoder        encoding.Encoder
	decoder        encoding.Decoder
	logger         zerolog.Logger

	lastRetryInterval time.Duration
}

func NewRedisScenarioExecutionStorage(
	redisKeyPrefix string,
	redisClient redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	lastRetryInterval time.Duration,
	logger zerolog.Logger,
) ScenarioExecutionStorage {
	return &redisScenarioExecutionStorage{
		redisKeyPrefix: redisKeyPrefix,
		redisClient:    redisClient,
		encoder:        encoder,
		decoder:        decoder,
		logger:         logger,

		lastRetryInterval: lastRetryInterval,
	}
}

func (s *redisScenarioExecutionStorage) Get(ctx context.Context, key string) (*ScenarioExecution, error) {
	res := s.redisClient.Get(ctx, s.getRedisKey(key))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	var execution ScenarioExecution
	err := s.decoder.Decode([]byte(res.Val()), &execution)
	if err != nil {
		return nil, err
	}

	return &execution, nil
}

func (s *redisScenarioExecutionStorage) Create(
	ctx context.Context,
	execution ScenarioExecution,
) (bool, error) {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return false, err
	}

	key := execution.GetCacheKey()
	res := s.redisClient.SetNX(ctx, s.getRedisKey(key), encoded, 0)
	if err := res.Err(); err != nil {
		return false, err
	}

	return res.Val(), nil
}

func (s *redisScenarioExecutionStorage) Update(ctx context.Context, execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(ctx, s.getRedisKey(execution.GetCacheKey()), encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) updateWithoutPrefix(ctx context.Context, redisKey string, execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(ctx, redisKey, encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) Del(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, s.getRedisKey(key)).Err()
}

func (s *redisScenarioExecutionStorage) delWithoutPrefix(ctx context.Context, redisKey string) error {
	return s.redisClient.Del(ctx, redisKey).Err()
}

func (s *redisScenarioExecutionStorage) Inc(ctx context.Context, key string, inc int64, drop bool) (int64, error) {
	incRedisKey := s.getRedisIncKey(key)
	if drop {
		res := s.redisClient.Del(ctx, incRedisKey)
		if err := res.Err(); err != nil {
			return 0, err
		}
	}

	res := s.redisClient.IncrBy(ctx, incRedisKey, inc)
	if err := res.Err(); err != nil {
		return 0, err
	}

	return res.Val(), nil
}

func (s *redisScenarioExecutionStorage) GetAbandoned(ctx context.Context) ([]ScenarioExecution, error) {
	executions := make([]ScenarioExecution, 0)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := s.redisClient.Scan(ctx, cursor, s.getRedisKey("*"), 50)
		if err := res.Err(); err != nil {
			return nil, err
		}

		var keys []string
		keys, cursor = res.Val()
		unprocessedKeys := make([]string, 0)
		for _, key := range keys {
			if !processedKeys[key] {
				unprocessedKeys = append(unprocessedKeys, key)
				processedKeys[key] = true
			}
		}

		if len(unprocessedKeys) > 0 {
			resGet := s.redisClient.MGet(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return nil, err
			}

			for i, v := range resGet.Val() {
				redisKey := unprocessedKeys[i]
				if v == nil {
					continue
				}

				if se, ok := v.(string); ok {
					var execution ScenarioExecution
					err := json.Unmarshal([]byte(se), &execution)
					if err != nil {
						return nil, err
					}

					if execution.LastUpdate > 0 && time.Since(time.Unix(execution.LastUpdate, 0)) > s.lastRetryInterval {
						key := s.parseRedisKey(redisKey)
						execution.Tries++
						if execution.Tries > MaxRetries {
							err := s.delWithoutPrefix(ctx, redisKey)
							if err != nil {
								s.logger.Warn().Err(err).Str("execution", key).Msg("Scenario execution storage: Failed to delete execution, since it has reached max number of retries.")
								continue
							}

							s.logger.Debug().Str("execution", key).Msg("Scenario execution storage: execution has been deleted, since it reached max number of retries.")

							continue
						}

						err := s.updateWithoutPrefix(ctx, redisKey, execution)
						if err != nil {
							s.logger.Warn().Err(err).Msg("Scenario execution storage: Failed to update execution tries, abandoned execution will be skipped.")
						} else {
							executions = append(executions, execution)
						}
					}
				} else {
					return nil, fmt.Errorf("unknown value type by key %q : expected string but got %+v", redisKey, v)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return executions, nil
}

func (s *redisScenarioExecutionStorage) getRedisKey(key string) string {
	return s.redisKeyPrefix + "-execution-" + key
}

func (s *redisScenarioExecutionStorage) parseRedisKey(key string) string {
	return strings.ReplaceAll(key, s.redisKeyPrefix+"-execution-", "")
}

func (s *redisScenarioExecutionStorage) getRedisIncKey(key string) string {
	return s.redisKeyPrefix + "-inc-" + key
}
