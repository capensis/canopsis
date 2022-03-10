package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type ScenarioExecutionStorage interface {
	Get(ctx context.Context, executionID string) (*ScenarioExecution, error)
	GetAbandoned(ctx context.Context) ([]ScenarioExecution, error)
	Create(ctx context.Context, execution ScenarioExecution) (string, error)
	Update(ctx context.Context, execution ScenarioExecution) error
	Del(ctx context.Context, executionID string) error
	Inc(ctx context.Context, id string, inc int64, drop bool) (int64, error)
}

type redisScenarioExecutionStorage struct {
	key         string
	redisClient redis.Cmdable
	encoder     encoding.Encoder
	decoder     encoding.Decoder
	logger      zerolog.Logger

	abandonedInterval time.Duration
}

func NewRedisScenarioExecutionStorage(
	key string,
	redisClient redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	abandonedInterval time.Duration,
	logger zerolog.Logger,
) ScenarioExecutionStorage {
	return &redisScenarioExecutionStorage{
		key:         key,
		redisClient: redisClient,
		encoder:     encoder,
		decoder:     decoder,
		logger:      logger,

		abandonedInterval: abandonedInterval,
	}
}

func (s *redisScenarioExecutionStorage) Get(ctx context.Context, executionID string) (*ScenarioExecution, error) {
	res := s.redisClient.Get(ctx, s.getKey(executionID))
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var execution ScenarioExecution
	err := s.decoder.Decode([]byte(res.Val()), &execution)
	if err != nil {
		return nil, err
	}

	execution.ID = executionID
	execution.AlarmID, execution.ScenarioID, err = s.parseExecutionID(executionID)
	if err != nil {
		return nil, err
	}

	return &execution, nil
}

func (s *redisScenarioExecutionStorage) Create(
	ctx context.Context,
	execution ScenarioExecution,
) (string, error) {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return "", err
	}

	executionID := s.getExecutionID(execution)
	res := s.redisClient.SetNX(ctx, s.getKey(executionID), encoded, 0)
	if err := res.Err(); err != nil {
		return "", err
	}

	if !res.Val() {
		return "", nil
	}

	return executionID, nil
}

func (s *redisScenarioExecutionStorage) Update(ctx context.Context, execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(ctx, s.getKey(execution.ID), encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) updateWithoutPrefix(ctx context.Context, executionID string, execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(ctx, executionID, encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) Del(ctx context.Context, executionID string) error {
	return s.redisClient.Del(ctx, s.getKey(executionID)).Err()
}

func (s *redisScenarioExecutionStorage) delWithoutPrefix(ctx context.Context, executionID string) error {
	return s.redisClient.Del(ctx, executionID).Err()
}

func (s *redisScenarioExecutionStorage) Inc(ctx context.Context, id string, inc int64, drop bool) (int64, error) {
	key := s.getIncKey(id)
	if drop {
		res := s.redisClient.Del(ctx, key)
		if err := res.Err(); err != nil {
			return 0, err
		}
	}

	res := s.redisClient.IncrBy(ctx, key, inc)
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
		res := s.redisClient.Scan(ctx, cursor, s.getKey("*"), 50)
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
				key := unprocessedKeys[i]

				if se, ok := v.(string); ok {
					var execution ScenarioExecution
					err := json.Unmarshal([]byte(se), &execution)
					if err != nil {
						return nil, err
					}

					if execution.LastUpdate > 0 && time.Since(time.Unix(execution.LastUpdate, 0)) > s.abandonedInterval {
						execution.Tries++
						if execution.Tries > MaxRetries {
							err := s.delWithoutPrefix(ctx, key)
							if err != nil {
								s.logger.Warn().Err(err).Str("execution_id", key).Msg("Scenario execution storage: Failed to delete execution, since it has reached max number of retries.")
								continue
							}

							s.logger.Debug().Str("execution_id", key).Msg("Scenario execution storage: execution has been deleted, since it reached max number of retries.")

							continue
						}

						err := s.updateWithoutPrefix(ctx, key, execution)
						if err != nil {
							s.logger.Warn().Err(err).Msg("Scenario execution storage: Failed to update execution tries, abandoned execution will be skipped.")
						} else {
							executionID := s.getParseKey(key)
							execution.ID = executionID
							execution.AlarmID, execution.ScenarioID, err = s.parseExecutionID(executionID)
							if err != nil {
								s.logger.Warn().Err(err).Str("execution_id", key).Msgf("Scenario execution storage: execution will be removed")
								err = s.delWithoutPrefix(ctx, key)
								if err != nil {
									s.logger.Warn().Err(err).Str("execution_id", key).Msg("Scenario execution storage: Failed to delete execution.")
									continue
								}
							}

							executions = append(executions, execution)
						}
					}
				} else {
					return nil, fmt.Errorf("unknown value type")
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return executions, nil
}

func (s *redisScenarioExecutionStorage) getExecutionID(execution ScenarioExecution) string {
	return fmt.Sprintf("%s$$%s", execution.AlarmID, execution.ScenarioID)
}

func (s *redisScenarioExecutionStorage) parseExecutionID(executionID string) (string, string, error) {
	parts := strings.Split(executionID, "$$")
	if len(parts) < 2 {
		return "", "", errors.New("invalid execution id")

	}

	return parts[0], parts[1], nil
}

func (s *redisScenarioExecutionStorage) getKey(id string) string {
	return fmt.Sprintf("%s-execution-%s", s.key, id)
}

func (s *redisScenarioExecutionStorage) getParseKey(key string) string {
	return strings.ReplaceAll(key, fmt.Sprintf("%s-execution-", s.key), "")
}

func (s *redisScenarioExecutionStorage) getIncKey(id string) string {
	return fmt.Sprintf("%s-inc-%s", s.key, id)
}
