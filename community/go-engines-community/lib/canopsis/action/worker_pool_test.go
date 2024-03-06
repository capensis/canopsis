package action_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestPool_RunWorkers_GivenMatchedTask_ShouldDoRpcCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timerCtx, timerCancel := context.WithCancel(context.Background())
	defer timerCancel()

	go func(ctx context.Context) {
		deadlockTimer := time.NewTimer(5 * time.Second)

		select {
		case <-ctx.Done():
			return
		case <-deadlockTimer.C:
			panic("workers or test are deadlocked")
		}
	}(timerCtx)

	cond, err := pattern.NewRegexpCondition("regexp", "abc-.*-def")
	if err != nil {
		t.Fatal("error shouldn't be raised")
	}

	p := pattern.Entity{{
		{
			Field:     "type",
			Condition: cond,
		},
		{
			Field:     "infos.output",
			FieldType: "string",
			Condition: pattern.NewStringCondition("eq", "debian9"),
		},
	}}

	axeRpcMock := mock_engine.NewMockRPCClient(ctrl)
	webhookRpcMock := mock_engine.NewMockRPCClient(ctrl)
	mongoClientMock := mock_mongo.NewMockDbClient(ctrl)
	alarmCollectionMock := mock_mongo.NewMockDbCollection(ctrl)
	webhookHistoryCollectionMock := mock_mongo.NewMockDbCollection(ctrl)

	mongoClientMock.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.AlarmMongoCollection:
			return alarmCollectionMock
		case mongo.WebhookHistoryMongoCollection:
			return webhookHistoryCollectionMock
		}

		return nil
	}).AnyTimes()

	taskChannel := make(chan action.Task)
	defer close(taskChannel)

	tplExecutor := template.NewExecutor(config.NewTemplateConfigProvider(config.CanopsisConf{}, zerolog.Nop()), config.NewTimezoneConfigProvider(config.CanopsisConf{}, zerolog.Nop()))

	pool := action.NewWorkerPool(5, mongoClientMock, axeRpcMock, webhookRpcMock, json.NewEncoder(), zerolog.Nop(), tplExecutor)
	resultChannel, err := pool.RunWorkers(ctx, taskChannel)
	if err != nil {
		t.Fatal("error shouldn't be raised")
	}

	var dataSets = []struct {
		testName           string
		expectedNotMatched bool
		task               action.Task
		expectedOutput     string
		expectedAuthor     string
	}{
		{
			task: action.Task{
				ScenarioName: "test-scenario-name",
				Action: action.Action{
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: p,
					},
					Type: "action_1",
					Parameters: action.Parameters{
						Output: "output 1",
						Author: "author 1",
					},
				},
				Alarm: types.Alarm{
					ID: "9",
				},
				Entity: types.Entity{
					Type:    "abc-ghi-def",
					Enabled: true,
					Infos: map[string]types.Info{
						"output": {
							Value: "debian9",
						},
					},
				},
				ExecutionCacheKey: "execution_1",
				Step:              2,
			},
			expectedOutput: "Scenario: test-scenario-name. Comment: output 1.",
			expectedAuthor: "author 1",
		},
		{
			task: action.Task{
				Action: action.Action{
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: p,
					},
					Type: "action_2",
					Parameters: action.Parameters{
						Output: "output 2",
						Author: "author 2",
					},
				},
				Alarm: types.Alarm{
					ID: "8",
				},
				ExecutionCacheKey: "execution_2",
				Step:              1,
			},
			expectedNotMatched: true,
		},
		{
			task: action.Task{
				ScenarioName: "test-scenario-name",
				Action: action.Action{
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: p,
					},
					Type: "action_3",
					Parameters: action.Parameters{
						Output: "output 3",
						Author: "author 3",
					},
				},
				Alarm: types.Alarm{
					ID: "7",
				},
				Entity: types.Entity{
					Type:    "abc-ghi-def",
					Enabled: true,
					Infos: map[string]types.Info{
						"output": {
							Value: "debian9",
						},
					},
				},
				ExecutionCacheKey: "execution_3",
				Step:              4,
			},
			expectedOutput: "Scenario: test-scenario-name. Comment: output 3.",
			expectedAuthor: "author 3",
		},
		{
			testName: "should render templates",
			task: action.Task{
				ScenarioName: "test-scenario-name",
				Action: action.Action{
					EntityPatternFields: savedpattern.EntityPatternFields{
						EntityPattern: p,
					},
					Type: "action_1",
					Parameters: action.Parameters{
						Output: "rendered output: {{.Entity.ID}}",
						Author: "rendered author: {{.Alarm.ID}}",
					},
				},
				Alarm: types.Alarm{
					ID: "9",
				},
				Entity: types.Entity{
					ID:      "test",
					Type:    "abc-ghi-def",
					Enabled: true,
					Infos: map[string]types.Info{
						"output": {
							Value: "debian9",
						},
					},
				},
				ExecutionCacheKey: "execution_1",
				Step:              2,
			},
			expectedOutput: "Scenario: test-scenario-name. Comment: rendered output: test.",
			expectedAuthor: "rendered author: 9",
		},
	}

	rpcWait := make(chan int)
	defer close(rpcWait)

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			if !dataset.expectedNotMatched {
				axeRpcMock.
					EXPECT().
					Call(gomock.Any(), gomock.Any()).
					Times(1).
					Do(func(_ context.Context, val1 interface{}) {
						decoder := json.NewDecoder()
						message := val1.(engine.RPCMessage)
						correlationID := message.CorrelationID

						expectedCorrelationID := fmt.Sprintf("%s&&%d", dataset.task.ExecutionCacheKey, dataset.task.Step)
						if expectedCorrelationID != correlationID {
							t.Errorf("expected correlation_id = %s, got %s", expectedCorrelationID, correlationID)
						}

						var event rpc.AxeEvent
						err := decoder.Decode(message.Body, &event)
						if err != nil {
							t.Error("failed to decode rpc axe event")
						} else {
							if event.Parameters.Output != dataset.expectedOutput {
								t.Errorf("expected output = %s, got %s", dataset.expectedOutput, event.Parameters.Output)
							}

							if event.Parameters.Author != dataset.expectedAuthor {
								t.Errorf("expected author = %s, got %s", dataset.expectedAuthor, event.Parameters.Author)
							}
						}

						rpcWait <- 1
					})
			}

			taskChannel <- dataset.task
			if dataset.expectedNotMatched {
				result := <-resultChannel
				if result.Status == action.TaskNotMatched {
					if !dataset.expectedNotMatched {
						t.Errorf("Task for action executionCacheKey=%s should be matched", result.ExecutionCacheKey)
					}
				}
			} else {
				<-rpcWait
			}
		})
	}
}

func TestPool_RunWorkers_GivenCancelContext_ShouldCancelTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	timerCtx, timerCancel := context.WithCancel(context.Background())
	defer timerCancel()

	go func(ctx context.Context) {
		deadlockTimer := time.NewTimer(5 * time.Second)

		select {
		case <-ctx.Done():
			return
		case <-deadlockTimer.C:
			panic("workers or test are deadlocked")
		}
	}(timerCtx)

	axeRpcMock := mock_engine.NewMockRPCClient(ctrl)
	webhookRpcMock := mock_engine.NewMockRPCClient(ctrl)
	mongoClientMock := mock_mongo.NewMockDbClient(ctrl)
	alarmCollectionMock := mock_mongo.NewMockDbCollection(ctrl)
	webhookHistoryCollectionMock := mock_mongo.NewMockDbCollection(ctrl)

	mongoClientMock.EXPECT().Collection(gomock.Any()).DoAndReturn(func(name string) mongo.DbCollection {
		switch name {
		case mongo.AlarmMongoCollection:
			return alarmCollectionMock
		case mongo.WebhookHistoryMongoCollection:
			return webhookHistoryCollectionMock
		}

		return nil
	}).AnyTimes()

	taskChannel := make(chan action.Task)

	poolSize := 5
	pool := action.NewWorkerPool(poolSize, mongoClientMock, axeRpcMock, webhookRpcMock, json.NewEncoder(), zerolog.Nop(), nil)
	resultChannel, err := pool.RunWorkers(ctx, taskChannel)
	if err != nil {
		t.Fatal("error shouldn't be raised")
	}

	cancels := 0

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for result := range resultChannel {
			if result.Status == action.TaskCancelled {
				cancels++
			}
		}

		wg.Done()
	}()

	cancel()
	wg.Wait()

	if cancels != poolSize {
		t.Fatal("Not all workers cancelled properly")
	}
}
