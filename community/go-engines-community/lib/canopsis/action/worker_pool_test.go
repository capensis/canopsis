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
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_engine "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/engine"
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
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)

	taskChannel := make(chan action.Task)
	defer close(taskChannel)

	mockTimezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	mockTimezoneConfigProvider.EXPECT().Get().Return(config.TimezoneConfig{}).AnyTimes()

	pool := action.NewWorkerPool(5, axeRpcMock, webhookRpcMock, alarmAdapterMock, json.NewEncoder(), zerolog.Nop(), mockTimezoneConfigProvider)
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
				ExecutionID: "execution_1",
				Step:        2,
			},
			expectedOutput: "output 1",
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
				ExecutionID: "execution_2",
				Step:        1,
			},
			expectedNotMatched: true,
		},
		{
			task: action.Task{
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
				ExecutionID: "execution_3",
				Step:        4,
			},
			expectedOutput: "output 3",
			expectedAuthor: "author 3",
		},
		{
			testName: "should render templates",
			task: action.Task{
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
				ExecutionID: "execution_1",
				Step:        2,
			},
			expectedOutput: "rendered output: test",
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

						expectedCorrelationID := fmt.Sprintf("%s&&%d", dataset.task.ExecutionID, dataset.task.Step)
						if expectedCorrelationID != correlationID {
							t.Errorf("expected correlation_id = %s, got %s", expectedCorrelationID, correlationID)
						}

						var event types.RPCAxeEvent
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
						t.Errorf("Task for action executionID=%s should be matched", result.ExecutionID)
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
	alarmAdapterMock := mock_alarm.NewMockAdapter(ctrl)

	taskChannel := make(chan action.Task)

	poolSize := 5
	pool := action.NewWorkerPool(poolSize, axeRpcMock, webhookRpcMock, alarmAdapterMock, json.NewEncoder(), zerolog.Nop(), nil)
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
