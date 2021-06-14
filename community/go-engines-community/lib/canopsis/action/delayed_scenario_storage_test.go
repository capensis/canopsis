package action_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/action"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/rs/zerolog"
	"reflect"
	"testing"
	"time"
)

func TestRedisDelayedScenarioStorage_Add(t *testing.T) {
	session, err := redis.NewSession(redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.DelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	scenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: types.CpsTime{Time: time.Now()},
		Paused:        false,
		TimeLeft:      0,
	}

	_, err = storage.Add(scenario)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestRedisDelayedScenarioStorage_Get(t *testing.T) {
	session, err := redis.NewSession(redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.DelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: types.CpsTime{Time: time.Now()},
		Paused:        false,
		TimeLeft:      0,
	}

	id, _ := storage.Add(expectedScenario)
	scenario, err := storage.Get(id)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenario, expectedScenario) {
		t.Errorf("expected %v but got %v", expectedScenario, scenario)
	}
}

func TestRedisDelayedScenarioStorage_GetAll(t *testing.T) {
	session, err := redis.NewSession(redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.DelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenarios := []action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id-1",
			AlarmID:       "test-alarm-id-1",
			ExecutionTime: types.CpsTime{Time: time.Now()},
			Paused:        false,
			TimeLeft:      0,
		},
		{
			ScenarioID:    "test-scenario-id-2",
			AlarmID:       "test-alarm-id-2",
			ExecutionTime: types.CpsTime{Time: time.Now()},
			Paused:        false,
			TimeLeft:      0,
		},
	}

	for _, s := range expectedScenarios {
		_, _ = storage.Add(s)
	}

	scenarios, err := storage.GetAll()
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenarios, expectedScenarios) {
		t.Errorf("expected %v but got %v", expectedScenarios, scenarios)
	}

	scenarios, err = storage.GetAll()
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenarios, expectedScenarios) {
		t.Errorf("expected %v but got %v", expectedScenarios, scenarios)
	}
}

func TestRedisDelayedScenarioStorage_Delete(t *testing.T) {
	session, err := redis.NewSession(redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.DelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	scenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: types.CpsTime{Time: time.Now()},
		Paused:        false,
		TimeLeft:      0,
	}

	id, _ := storage.Add(scenario)
	ok, err := storage.Delete(id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !ok {
		t.Errorf("expected %v but got %v", true, ok)
	}

	res, err := storage.Get(id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if res != nil {
		t.Errorf("expected null but got %v", res)
	}
}

func TestRedisDelayedScenarioStorage_Update(t *testing.T) {
	session, err := redis.NewSession(redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.DelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: types.CpsTime{Time: time.Now()},
		Paused:        false,
		TimeLeft:      0,
	}

	id, _ := storage.Add(expectedScenario)
	expectedScenario.ID = id
	expectedScenario.Paused = true
	ok, err := storage.Update(expectedScenario)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !ok {
		t.Errorf("expected %v but got %v", true, ok)
	}

	scenario, err := storage.Get(id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenario, expectedScenario) {
		t.Errorf("expected %v but got %v", expectedScenario, scenario)
	}
}
