package action_test

import (
	"context"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_action "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/action"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestTaskManager_Run_GiveTask_ShouldSendResult(t *testing.T) {
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
	task := action.ExecuteScenariosTask{
		Triggers: []string{"create"},
		Entity:   types.Entity{ID: "test-entity"},
		Alarm:    types.Alarm{ID: "test-alarm"},
	}
	scenario := action.Scenario{
		ID:   "test-scenario",
		Name: "test-scenario-name",
		Actions: []action.Action{
			{
				Type: "snooze",
			},
		},
	}
	execution := action.ScenarioExecution{
		ScenarioID: scenario.ID,
		AlarmID:    task.Alarm.ID,
		Entity:     task.Entity,
		ActionExecutions: []action.Execution{
			{
				Action:   scenario.Actions[0],
				Executed: false,
			},
		},
	}
	rpcResultCh := make(chan action.RpcResult)
	defer close(rpcResultCh)
	inputCh := make(chan action.ExecuteScenariosTask)
	defer close(inputCh)
	taskResultCh := make(chan action.TaskResult)
	defer close(taskResultCh)
	mockWorkerPool := mock_action.NewMockWorkerPool(ctrl)
	mockWorkerPool.EXPECT().RunWorkers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, taskCh <-chan action.Task) {
			go func() {
				select {
				case <-ctx.Done():
				case task := <-taskCh:
					taskResultCh <- action.TaskResult{
						Source:            "test",
						Alarm:             task.Alarm,
						Step:              task.Step,
						ExecutionCacheKey: task.ExecutionCacheKey,
						Status:            action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().IncExecutedCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().IncExecutedWebhookCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(0)), gomock.Eq(true)).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().DelExecutingCount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().DelExecutedCount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().DelExecutedWebhookCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(true, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), gomock.Eq(execution.GetCacheKey())).Return(&execution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), gomock.Eq(execution.GetCacheKey())).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq(task.Triggers), gomock.Eq(task.Alarm)).
		Return(map[string][]action.Scenario{"create": {scenario}}, nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq(task.Triggers), gomock.Eq(task.Alarm), gomock.Eq(task.Entity), gomock.Eq(task.AdditionalData)).
		Return(nil)
	logger := zerolog.Logger{}
	manager := action.NewTaskManager(mockWorkerPool, mockExecutionStorage, mockScenarioStorage, logger)
	resultCh, err := manager.Run(ctx, rpcResultCh, inputCh)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	inputCh <- task

	time.Sleep(100 * time.Millisecond)

	select {
	case res := <-resultCh:
		if res.Err != nil {
			t.Errorf("expected no error but got %v", res.Err)
		}

		if res.Alarm.ID != task.Alarm.ID {
			t.Errorf("expected alarm but got %v", res.Alarm)
		}
	default:
		t.Errorf("expected result but got nothing")
	}
}

func TestTaskManager_Run_GiveTaskWithEmitTrigger_ShouldSendResult(t *testing.T) {
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
	task := action.ExecuteScenariosTask{
		Triggers: []string{"create"},
		Entity:   types.Entity{ID: "test-entity"},
		Alarm:    types.Alarm{ID: "test-alarm"},
	}
	firstScenario := action.Scenario{
		ID:   "test-scenario-1",
		Name: "test-scenario-1-name",
		Actions: []action.Action{
			{
				Type:        "snooze",
				EmitTrigger: true,
			},
		},
	}
	secondScenario := action.Scenario{
		ID:   "test-scenario-2",
		Name: "test-scenario-2-name",
		Actions: []action.Action{
			{
				Type: "ack",
			},
		},
	}
	firstExecution := action.ScenarioExecution{
		ScenarioID: firstScenario.ID,
		AlarmID:    task.Alarm.ID,
		Entity:     task.Entity,
		ActionExecutions: []action.Execution{
			{
				Action:   firstScenario.Actions[0],
				Executed: false,
			},
		},
	}
	secondExecution := action.ScenarioExecution{
		ScenarioID: secondScenario.ID,
		AlarmID:    task.Alarm.ID,
		Entity:     task.Entity,
		ActionExecutions: []action.Execution{
			{
				Action:   secondScenario.Actions[0],
				Executed: false,
			},
		},
	}
	rpcResultCh := make(chan action.RpcResult)
	defer close(rpcResultCh)
	inputCh := make(chan action.ExecuteScenariosTask)
	defer close(inputCh)
	taskResultCh := make(chan action.TaskResult)
	defer close(taskResultCh)
	mockWorkerPool := mock_action.NewMockWorkerPool(ctrl)
	mockWorkerPool.EXPECT().RunWorkers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, taskCh <-chan action.Task) {
			go func() {
				for i := 0; i < 2; i++ {
					select {
					case <-ctx.Done():
						return
					case task := <-taskCh:
						taskResultCh <- action.TaskResult{
							Source:            "test",
							Alarm:             task.Alarm,
							Step:              task.Step,
							ExecutionCacheKey: task.ExecutionCacheKey,
							AlarmChangeType:   types.AlarmChangeType(task.Action.Type),
							Status:            action.TaskNotMatched,
						}
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).
		Return(int64(1), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutedCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).
		Return(int64(1), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutedWebhookCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(0)), gomock.Eq(true)).
		Return(int64(0), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(false)).
		Return(int64(2), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutedCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(false)).
		Return(int64(1), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).
		Return(int64(1), nil).Times(1)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).
		Return(int64(0), nil).Times(1)
	mockExecutionStorage.EXPECT().DelExecutingCount(gomock.Any(), gomock.Any()).Return(int64(2), nil)
	mockExecutionStorage.EXPECT().DelExecutedCount(gomock.Any(), gomock.Any()).Return(int64(2), nil)
	mockExecutionStorage.EXPECT().DelExecutedWebhookCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	createCall := 0
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ action.ScenarioExecution) (bool, error) {
		createCall++
		return true, nil
	}).Times(2)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), firstExecution.GetCacheKey()).Return(&firstExecution, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), secondExecution.GetCacheKey()).Return(&secondExecution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), firstExecution.GetCacheKey()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), secondExecution.GetCacheKey()).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq(task.Triggers), gomock.Eq(task.Alarm)).
		Return(map[string][]action.Scenario{"create": {firstScenario}}, nil)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq([]string{firstScenario.Actions[0].Type}), gomock.Eq(task.Alarm)).
		Return(map[string][]action.Scenario{"create": {secondScenario}}, nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq(task.Triggers), gomock.Eq(task.Alarm), gomock.Eq(task.Entity), gomock.Eq(task.AdditionalData)).
		Return(nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq([]string{firstScenario.Actions[0].Type}), gomock.Eq(task.Alarm), gomock.Eq(task.Entity), gomock.Any()).
		Return(nil)
	logger := zerolog.Logger{}
	manager := action.NewTaskManager(mockWorkerPool, mockExecutionStorage, mockScenarioStorage, logger)
	resultCh, err := manager.Run(ctx, rpcResultCh, inputCh)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	inputCh <- task

	time.Sleep(100 * time.Millisecond)

	select {
	case res := <-resultCh:
		if res.Err != nil {
			t.Errorf("expected no error but got %v", res.Err)
		}

		if res.Alarm.ID != task.Alarm.ID {
			t.Errorf("expected alarm but got %v", res.Alarm)
		}
	default:
		t.Errorf("expected result but got nothing")
	}
}

func TestTaskManager_Run_GiveDelayedTask_ShouldSendResult(t *testing.T) {
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
	scenario := action.Scenario{
		ID:   "test-scenario",
		Name: "test-scenario-name",
		Actions: []action.Action{
			{
				Type: "snooze",
			},
		},
	}
	task := action.ExecuteScenariosTask{
		DelayedScenarioID: scenario.ID,
		Entity:            types.Entity{ID: "test-entity"},
		Alarm:             types.Alarm{ID: "test-alarm"},
	}
	execution := action.ScenarioExecution{
		ScenarioID: scenario.ID,
		AlarmID:    task.Alarm.ID,
		Entity:     task.Entity,
		ActionExecutions: []action.Execution{
			{
				Action:   scenario.Actions[0],
				Executed: false,
			},
		},
	}
	rpcResultCh := make(chan action.RpcResult)
	defer close(rpcResultCh)
	inputCh := make(chan action.ExecuteScenariosTask)
	defer close(inputCh)
	taskResultCh := make(chan action.TaskResult)
	defer close(taskResultCh)
	mockWorkerPool := mock_action.NewMockWorkerPool(ctrl)
	mockWorkerPool.EXPECT().RunWorkers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, taskCh <-chan action.Task) {
			go func() {
				select {
				case <-ctx.Done():
				case task := <-taskCh:
					taskResultCh <- action.TaskResult{
						Source:            "test",
						Alarm:             task.Alarm,
						Step:              task.Step,
						ExecutionCacheKey: task.ExecutionCacheKey,
						Status:            action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().DelExecutingCount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().DelExecutedCount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().DelExecutedWebhookCount(gomock.Any(), gomock.Any()).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(true, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), execution.GetCacheKey()).Return(&execution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), execution.GetCacheKey()).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	mockScenarioStorage.EXPECT().
		GetScenario(gomock.Eq(task.DelayedScenarioID)).
		Return(&scenario)
	logger := zerolog.Logger{}
	manager := action.NewTaskManager(mockWorkerPool, mockExecutionStorage, mockScenarioStorage, logger)
	resultCh, err := manager.Run(ctx, rpcResultCh, inputCh)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	inputCh <- task

	time.Sleep(100 * time.Millisecond)

	select {
	case res := <-resultCh:
		if res.Err != nil {
			t.Errorf("expected no error but got %v", res.Err)
		}

		if res.Alarm.ID != task.Alarm.ID {
			t.Errorf("expected alarm but got %v", res.Alarm)
		}
	default:
		t.Errorf("expected result but got nothing")
	}
}

func TestTaskManager_Run_GiveAbandonedTask_ShouldSendResult(t *testing.T) {
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
	executionCacheKey := "test-alarm$$test-scenario"
	task := action.ExecuteScenariosTask{
		AbandonedExecutionCacheKey: executionCacheKey,
		Entity:                     types.Entity{ID: "test-entity"},
		Alarm:                      types.Alarm{ID: "test-alarm"},
	}
	scenario := action.Scenario{
		ID:   "test-scenario",
		Name: "test-scenario-name",
		Actions: []action.Action{
			{
				Type: "snooze",
			},
		},
	}
	execution := action.ScenarioExecution{
		ScenarioID: scenario.ID,
		AlarmID:    task.Alarm.ID,
		Entity:     task.Entity,
		ActionExecutions: []action.Execution{
			{
				Action:   scenario.Actions[0],
				Executed: false,
			},
		},
		Tries: 1,
	}
	rpcResultCh := make(chan action.RpcResult)
	defer close(rpcResultCh)
	inputCh := make(chan action.ExecuteScenariosTask)
	defer close(inputCh)
	taskResultCh := make(chan action.TaskResult)
	mockWorkerPool := mock_action.NewMockWorkerPool(ctrl)
	mockWorkerPool.EXPECT().RunWorkers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, taskCh <-chan action.Task) {
			go func() {
				defer close(taskResultCh)

				select {
				case <-ctx.Done():
				case task := <-taskCh:
					taskResultCh <- action.TaskResult{
						Source:            "test",
						Alarm:             task.Alarm,
						Step:              task.Step,
						ExecutionCacheKey: task.ExecutionCacheKey,
						Status:            action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().IncExecutingCount(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), gomock.Eq(execution.GetCacheKey())).Return(&execution, nil).Times(2)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), gomock.Eq(execution.GetCacheKey())).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	logger := zerolog.Logger{}
	manager := action.NewTaskManager(mockWorkerPool, mockExecutionStorage, mockScenarioStorage, logger)
	resultCh, err := manager.Run(ctx, rpcResultCh, inputCh)
	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

	inputCh <- task

	time.Sleep(100 * time.Millisecond)

	select {
	case res := <-resultCh:

		t.Errorf("expected not result but got %+v", res)
	default:
	}
}
