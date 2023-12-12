package action

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/redis/go-redis/v9"
)

type DelayedScenario struct {
	ID            string           `json:"id"`
	ScenarioID    string           `json:"scenario_id"`
	AlarmID       string           `json:"alarm_id"`
	ExecutionTime datetime.CpsTime `json:"execution_time"`
	Paused        bool             `json:"paused"`
	TimeLeft      time.Duration    `json:"time_left"`

	AdditionalData AdditionalData `json:"additional_data"`
}

type DelayedScenarioStorage interface {
	Add(ctx context.Context, scenario DelayedScenario) (string, error)
	GetAll(ctx context.Context) ([]DelayedScenario, error)
	Get(ctx context.Context, id string) (*DelayedScenario, error)
	Delete(ctx context.Context, id string) (bool, error)
	Update(ctx context.Context, scenario DelayedScenario) (bool, error)
}

func NewRedisDelayedScenarioStorage(
	key string,
	client redis.Cmdable,
	encoder encoding.Encoder,
	decoder encoding.Decoder,
) DelayedScenarioStorage {
	return &redisDelayedScenarioStorage{
		key:     key,
		client:  client,
		encoder: encoder,
		decoder: decoder,
	}
}

type redisDelayedScenarioStorage struct {
	key     string
	encoder encoding.Encoder
	decoder encoding.Decoder
	client  redis.Cmdable
}

func (s *redisDelayedScenarioStorage) GetAll(ctx context.Context) ([]DelayedScenario, error) {
	scenarios := make([]DelayedScenario, 0)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := s.client.Scan(ctx, cursor, s.key+"*", 50)
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
			resGet := s.client.MGet(ctx, unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return nil, err
			}

			for i, v := range resGet.Val() {
				if v == nil {
					continue
				}

				if s, ok := v.(string); ok {
					var scenario DelayedScenario
					err := json.Unmarshal([]byte(s), &scenario)
					if err != nil {
						return nil, err
					}

					scenarios = append(scenarios, scenario)
				} else {
					return nil, fmt.Errorf("unknown value type by key %q : expected string but got %+v", unprocessedKeys[i], v)
				}
			}
		}

		if cursor == 0 {
			break
		}
	}

	return scenarios, nil
}

func (s *redisDelayedScenarioStorage) Get(ctx context.Context, id string) (*DelayedScenario, error) {
	res := s.client.Get(ctx, s.getKey(id))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, err
	}

	var scenario DelayedScenario
	err := s.decoder.Decode([]byte(res.Val()), &scenario)
	if err != nil {
		return nil, err
	}

	return &scenario, nil
}

func (s *redisDelayedScenarioStorage) Add(ctx context.Context, scenario DelayedScenario) (string, error) {
	scenario.ID = utils.NewID()
	v, err := s.encoder.Encode(scenario)
	if err != nil {
		return "", err
	}

	res := s.client.SetNX(ctx, s.getKey(scenario.ID), v, 0)
	if err := res.Err(); err != nil {
		return "", err
	}

	return scenario.ID, nil
}

func (s *redisDelayedScenarioStorage) Delete(ctx context.Context, id string) (bool, error) {
	res := s.client.Del(ctx, s.getKey(id))
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, err

	}

	return res.Val() > 0, nil
}

func (s *redisDelayedScenarioStorage) Update(ctx context.Context, scenario DelayedScenario) (bool, error) {
	v, err := s.encoder.Encode(scenario)
	if err != nil {
		return false, err
	}

	res := s.client.SetXX(ctx, s.getKey(scenario.ID), v, 0)
	if err := res.Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *redisDelayedScenarioStorage) getKey(id string) string {
	return s.key + "-" + id
}
