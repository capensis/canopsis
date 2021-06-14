package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"github.com/go-redis/redis/v7"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type ScenarioExecutionStorage interface {
	Get(executionID string) (*ScenarioExecution, error)
	GetAbandoned() ([]ScenarioExecution, error)
	Create(execution ScenarioExecution) (string, error)
	Update(execution ScenarioExecution) error
	Del(executionID string) error
	Inc(id string, inc int64, drop bool) (int64, error)
}

type redisScenarioExecutionStorage struct {
	key         string
	redisClient redis.Cmdable
	encoder     encoding.Encoder
	decoder     encoding.Decoder
	logger      zerolog.Logger
}

func NewRedisScenarioExecutionStorage(
	key string,
	redisClient redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
	logger zerolog.Logger,
) ScenarioExecutionStorage {
	return &redisScenarioExecutionStorage{
		key:         key,
		redisClient: redisClient,
		encoder:     encoder,
		decoder:     decoder,
		logger:      logger,
	}
}

func (s *redisScenarioExecutionStorage) Get(executionID string) (*ScenarioExecution, error) {
	res := s.redisClient.Get(s.getKey(executionID))
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
	execution ScenarioExecution,
) (string, error) {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return "", err
	}

	executionID := s.getExecutionID(execution)
	res := s.redisClient.SetNX(s.getKey(executionID), encoded, 0)
	if err := res.Err(); err != nil {
		return "", err
	}

	if !res.Val() {
		return "", nil
	}

	return executionID, nil
}

func (s *redisScenarioExecutionStorage) Update(execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(s.getKey(execution.ID), encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) updateWithoutPrefix(executionID string, execution ScenarioExecution) error {
	encoded, err := s.encoder.Encode(execution)
	if err != nil {
		return err
	}

	res := s.redisClient.SetXX(executionID, encoded, 0)
	if err := res.Err(); err != nil {
		return err
	}

	if !res.Val() {
		return errors.New("key not found")
	}

	return nil
}

func (s *redisScenarioExecutionStorage) Del(executionID string) error {
	return s.redisClient.Del(s.getKey(executionID)).Err()
}

func (s *redisScenarioExecutionStorage) delWithoutPrefix(executionID string) error {
	return s.redisClient.Del(executionID).Err()
}

func (s *redisScenarioExecutionStorage) Inc(id string, inc int64, drop bool) (int64, error) {
	key := s.getIncKey(id)
	if drop {
		res := s.redisClient.Del(key)
		if err := res.Err(); err != nil {
			return 0, err
		}
	}

	res := s.redisClient.IncrBy(key, inc)
	if err := res.Err(); err != nil {
		return 0, err
	}

	return res.Val(), nil
}

func (s *redisScenarioExecutionStorage) GetAbandoned() ([]ScenarioExecution, error) {
	executions := make([]ScenarioExecution, 0)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := s.redisClient.Scan(cursor, s.getKey("*"), 50)
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
			resGet := s.redisClient.MGet(unprocessedKeys...)
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

					if execution.LastUpdate > 0 && time.Now().Unix()-execution.LastUpdate > AbandonedDuration {
						execution.Tries++
						if execution.Tries > MaxRetries {
							err := s.delWithoutPrefix(key)
							if err != nil {
								s.logger.Warn().Err(err).Str("execution_id", key).Msg("Scenario execution storage: Failed to delete execution, since it has reached max number of retries.")
								continue
							}

							s.logger.Debug().Str("execution_id", key).Msg("Scenario execution storage: execution has been deleted, since it reached max number of retries.")

							continue
						}

						err := s.updateWithoutPrefix(key, execution)
						if err != nil {
							s.logger.Warn().Err(err).Msg("Scenario execution storage: Failed to update execution tries, abandoned execution will be skipped.")
						} else {
							executionID := s.getParseKey(key)
							execution.ID = executionID
							execution.AlarmID, execution.ScenarioID, err = s.parseExecutionID(executionID)
							if err != nil {
								s.logger.Warn().Err(err).Str("execution_id", key).Msgf("Scenario execution storage: execution will be removed")
								err = s.delWithoutPrefix(key)
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
