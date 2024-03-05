package action_test

import (
	"context"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/rs/zerolog"
)

func TestRedisDelayedScenarioStorage_Add(t *testing.T) {
	session, err := redis.NewSession(context.Background(), redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	scenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: datetime.NewCpsTime(),
		Paused:        false,
		TimeLeft:      0,
	}

	_, err = storage.Add(context.Background(), scenario)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestRedisDelayedScenarioStorage_Get(t *testing.T) {
	session, err := redis.NewSession(context.Background(), redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: datetime.NewCpsTime(),
		Paused:        false,
		TimeLeft:      0,
	}

	ctx := context.Background()

	id, _ := storage.Add(ctx, expectedScenario)
	scenario, err := storage.Get(ctx, id)

	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenario, expectedScenario) {
		t.Errorf("expected %v but got %v", expectedScenario, scenario)
	}
}

func TestRedisDelayedScenarioStorage_GetAll(t *testing.T) {
	ctx := context.Background()

	session, err := redis.NewSession(ctx, redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenarios := []action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id-1",
			AlarmID:       "test-alarm-id-1",
			ExecutionTime: datetime.NewCpsTime(),
			Paused:        false,
			TimeLeft:      0,
		},
		{
			ScenarioID:    "test-scenario-id-2",
			AlarmID:       "test-alarm-id-2",
			ExecutionTime: datetime.NewCpsTime(),
			Paused:        false,
			TimeLeft:      0,
		},
	}

	for _, s := range expectedScenarios {
		_, _ = storage.Add(ctx, s)
	}

	scenarios, err := storage.GetAll(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenarios, expectedScenarios) {
		t.Errorf("expected %v but got %v", expectedScenarios, scenarios)
	}

	scenarios, err = storage.GetAll(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenarios, expectedScenarios) {
		t.Errorf("expected %v but got %v", expectedScenarios, scenarios)
	}
}

func TestRedisDelayedScenarioStorage_Delete(t *testing.T) {
	ctx := context.Background()

	session, err := redis.NewSession(ctx, redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	scenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: datetime.NewCpsTime(),
		Paused:        false,
		TimeLeft:      0,
	}

	id, _ := storage.Add(ctx, scenario)
	ok, err := storage.Delete(ctx, id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !ok {
		t.Errorf("expected %v but got %v", true, ok)
	}

	res, err := storage.Get(ctx, id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if res != nil {
		t.Errorf("expected null but got %v", res)
	}
}

func TestRedisDelayedScenarioStorage_Update(t *testing.T) {
	ctx := context.Background()

	session, err := redis.NewSession(ctx, redis.ActionScenarioStorage, zerolog.Logger{}, 0, 0)
	if err != nil {
		panic(err)
	}

	storage := action.NewRedisDelayedScenarioStorage(redis.ActionDelayedScenarioKey, session, json.NewEncoder(), json.NewDecoder())
	expectedScenario := action.DelayedScenario{
		ScenarioID:    "test-scenario-id",
		AlarmID:       "test-alarm-id",
		ExecutionTime: datetime.NewCpsTime(),
		Paused:        false,
		TimeLeft:      0,
	}

	id, _ := storage.Add(ctx, expectedScenario)
	expectedScenario.ID = id
	expectedScenario.Paused = true
	ok, err := storage.Update(ctx, expectedScenario)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if !ok {
		t.Errorf("expected %v but got %v", true, ok)
	}

	scenario, err := storage.Get(ctx, id)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	if reflect.DeepEqual(scenario, expectedScenario) {
		t.Errorf("expected %v but got %v", expectedScenario, scenario)
	}
}
