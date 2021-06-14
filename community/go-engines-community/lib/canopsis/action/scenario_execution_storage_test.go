package action_test

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	redislib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	"github.com/influxdata/influxdb/pkg/deep"
	"testing"
	"time"
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
	zeroTime := types.CpsTime{}
	firstExecution := action.ScenarioExecution{
		AlarmID:    "4",
		ScenarioID: "4",
		LastUpdate: timestamp - action.AbandonedDuration - 10,
		Entity:     types.Entity{Created: zeroTime},
	}
	firstExecutionID, err := storage.Create(ctx, firstExecution)
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}
	firstExecution.ID = firstExecutionID
	firstExecution.Tries = 1
	secondExecution := action.ScenarioExecution{
		AlarmID:    "5",
		ScenarioID: "5",
		LastUpdate: timestamp - action.AbandonedDuration - 30,
		Entity:     types.Entity{Created: zeroTime},
	}
	secondExecutionID, err := storage.Create(ctx, secondExecution)
	secondExecution.ID = secondExecutionID
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
		if !deep.Equal(exec, firstExecution) && !deep.Equal(exec, secondExecution) {
			t.Errorf("GetAbandoned should return %+v or %+v but got %v",
				firstExecution, secondExecution, exec)
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
	executionID, err := storage.Create(ctx, action.ScenarioExecution{
		AlarmID:    "6",
		ScenarioID: "6",
		LastUpdate: timestamp - action.AbandonedDuration - 30,
		Tries:      action.MaxRetries,
	})
	if err != nil {
		t.Fatalf("Error %s is not expected", err.Error())
	}

	_, _ = storage.GetAbandoned(ctx)

	exec, err := storage.Get(ctx, executionID)
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
		json.NewDecoder(), log.NewLogger(true))
}
