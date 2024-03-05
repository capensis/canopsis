package action_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	redislib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog"
)

func TestRedisScenarioExecutionStorage_GetAbandoned_GivenTooLongNotUpdatedExecutions_ShouldReturnThem(t *testing.T) {
	ctx := context.Background()

	timestamp := time.Now().Unix()
	storage := createTestStorage()
	_, err := storage.Create(ctx, action.ScenarioExecution{
		AlarmID:    "1",
		ScenarioID: "1",
		LastUpdate: timestamp - 10,
	})
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}
	_, err = storage.Create(ctx, action.ScenarioExecution{
		AlarmID:    "2",
		ScenarioID: "2",
		LastUpdate: timestamp - 30,
	})
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}
	_, err = storage.Create(ctx, action.ScenarioExecution{
		AlarmID:    "3",
		ScenarioID: "3",
		LastUpdate: timestamp - 40,
	})
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}
	zeroTime := datetime.CpsTime{}
	firstExecution := action.ScenarioExecution{
		AlarmID:    "4",
		ScenarioID: "4",
		LastUpdate: timestamp - 70,
		Entity:     types.Entity{Created: zeroTime},
	}
	_, err = storage.Create(ctx, firstExecution)
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}
	firstExecution.Tries = 1
	secondExecution := action.ScenarioExecution{
		AlarmID:    "5",
		ScenarioID: "5",
		LastUpdate: timestamp - 90,
		Entity:     types.Entity{Created: zeroTime},
	}
	_, err = storage.Create(ctx, secondExecution)
	secondExecution.Tries = 1
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}

	abandonedExecutions, err := storage.GetAbandoned(ctx)
	if err != nil {
		t.Errorf("Error %s is not expected", err.Error())
	}

	if len(abandonedExecutions) != 2 {
		t.Errorf("Expected 2 abandoned executions, got %d", len(abandonedExecutions))
	}

	for _, exec := range abandonedExecutions {
		exec.Entity.Created = zeroTime
		exec.FifoAckEvent = types.Event{}
		diff1 := pretty.Compare(exec, firstExecution)
		diff2 := pretty.Compare(exec, secondExecution)

		if diff1 != "" && diff2 != "" {
			t.Errorf("GetAbandoned should return\n%s\nor\b%s", diff1, diff2)
		}

		if exec.Tries != 1 {
			t.Errorf("GetAbandoned should increase tries")
		}
	}
}

func TestRedisScenarioExecutionStorage_GetAbandoned_GivenExecutionWithMaxRetries_ShouldDeleteIt(t *testing.T) {
	ctx := context.Background()

	timestamp := time.Now().Unix()
	storage := createTestStorage()
	execution := action.ScenarioExecution{
		AlarmID:    "6",
		ScenarioID: "6",
		LastUpdate: timestamp - 90,
		Tries:      action.MaxRetries,
	}
	_, err := storage.Create(ctx, execution)
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}

	_, _ = storage.GetAbandoned(ctx)

	exec, err := storage.Get(ctx, execution.GetCacheKey())
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}

	if exec != nil {
		t.Errorf("execution should be deleted")
	}
}

func createTestStorage() action.ScenarioExecutionStorage {
	ctx := context.Background()

	session, err := redislib.NewSession(ctx, redislib.ActionScenarioStorage, log.NewLogger(true), 0, 0)
	if err != nil {
		panic(err)
	}

	cmdRes := session.FlushDB(ctx)
	if cmdRes.Err() != nil {
		panic(err)
	}

	key := "test-scenario-execution-key"

	return action.NewRedisScenarioExecutionStorage(key, session, json.NewEncoder(),
		json.NewDecoder(), time.Minute, zerolog.Nop())
}
