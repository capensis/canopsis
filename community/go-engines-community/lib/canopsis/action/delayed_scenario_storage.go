package action

import (
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"github.com/go-redis/redis/v7"
	"sync"
	"time"
)

type DelayedScenario struct {
	ID            string        `json:"id"`
	ScenarioID    string        `json:"scenario_id"`
	AlarmID       string        `json:"alarm_id"`
	ExecutionTime types.CpsTime `json:"execution_time"`
	Paused        bool          `json:"paused"`
	TimeLeft      time.Duration `json:"time_left"`
}

type DelayedScenarioStorage interface {
	Add(DelayedScenario) (string, error)
	GetAll() ([]DelayedScenario, error)
	Get(id string) (*DelayedScenario, error)
	Delete(id string) (bool, error)
	Update(DelayedScenario) (bool, error)
}

func NewInMemoryDelayedScenarioStorage(init map[string]DelayedScenario) DelayedScenarioStorage {
	if init == nil {
		init = make(map[string]DelayedScenario)
	}

	return &inMemoryDelayedScenarioStorage{
		delayedScenarios: init,
	}
}

type inMemoryDelayedScenarioStorage struct {
	mx               sync.Mutex
	delayedScenarios map[string]DelayedScenario
}

func (s *inMemoryDelayedScenarioStorage) GetAll() ([]DelayedScenario, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	delayedScenarios := make([]DelayedScenario, 0)
	for _, scenario := range s.delayedScenarios {
		delayedScenarios = append(delayedScenarios, scenario)
	}

	return delayedScenarios, nil
}

func (s *inMemoryDelayedScenarioStorage) Get(id string) (*DelayedScenario, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	scenario, ok := s.delayedScenarios[id]
	if !ok {
		return nil, nil
	}

	return &scenario, nil
}

func (s *inMemoryDelayedScenarioStorage) Add(scenario DelayedScenario) (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	scenario.ID = utils.NewID()
	s.delayedScenarios[scenario.ID] = scenario

	return scenario.ID, nil
}

func (s *inMemoryDelayedScenarioStorage) Delete(id string) (bool, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.delayedScenarios[id]; !ok {
		return false, nil
	}

	delete(s.delayedScenarios, id)

	return true, nil
}

func (s *inMemoryDelayedScenarioStorage) Update(scenario DelayedScenario) (bool, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if _, ok := s.delayedScenarios[scenario.ID]; !ok {
		return false, nil
	}

	s.delayedScenarios[scenario.ID] = scenario

	return true, nil
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

func (s *redisDelayedScenarioStorage) GetAll() ([]DelayedScenario, error) {
	scenarios := make([]DelayedScenario, 0)
	var cursor uint64
	processedKeys := make(map[string]bool)

	for {
		res := s.client.Scan(cursor, fmt.Sprintf("%s*", s.key), 50)
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
			resGet := s.client.MGet(unprocessedKeys...)
			if err := resGet.Err(); err != nil {
				return nil, err
			}

			for i, v := range resGet.Val() {
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

func (s *redisDelayedScenarioStorage) Get(id string) (*DelayedScenario, error) {
	res := s.client.Get(s.getKey(id))
	if err := res.Err(); err != nil {
		if err == redis.Nil {
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

func (s *redisDelayedScenarioStorage) Add(scenario DelayedScenario) (string, error) {
	scenario.ID = utils.NewID()
	v, err := s.encoder.Encode(scenario)
	if err != nil {
		return "", err
	}

	res := s.client.SetNX(s.getKey(scenario.ID), v, 0)
	if err := res.Err(); err != nil {
		return "", err
	}

	return scenario.ID, nil
}

func (s *redisDelayedScenarioStorage) Delete(id string) (bool, error) {
	res := s.client.Del(s.getKey(id))
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err

	}

	return res.Val() > 0, nil
}

func (s *redisDelayedScenarioStorage) Update(scenario DelayedScenario) (bool, error) {
	v, err := s.encoder.Encode(scenario)
	if err != nil {
		return false, err
	}

	res := s.client.SetXX(s.getKey(scenario.ID), v, 0)
	if err := res.Err(); err != nil {
		if err == redis.Nil {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (s *redisDelayedScenarioStorage) getKey(id string) string {
	return fmt.Sprintf("%s-%s", s.key, id)
}
