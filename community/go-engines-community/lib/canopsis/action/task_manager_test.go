package action_test

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/action"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mock_action "git.canopsis.net/canopsis/go-engines/mocks/lib/canopsis/action"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"testing"
	"time"
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
		Triggers:     []string{"create"},
		Entity:       types.Entity{ID: "test-entity"},
		Alarm:        types.Alarm{ID: "test-alarm"},
		AckResources: false,
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
	executionID := "test-alarm&&test-scenario"
	execution := action.ScenarioExecution{
		ID:         executionID,
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
						Source:      "test",
						Alarm:       task.Alarm,
						Step:        task.Step,
						ExecutionID: task.ExecutionID,
						Status:      action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(executionID, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), executionID).Return(&execution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), executionID).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq(task.Triggers), gomock.Eq(task.Alarm)).
		Return([]action.Scenario{scenario}, nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq(task.Triggers), gomock.Eq(task.Alarm), gomock.Eq(task.Entity)).
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
		Triggers:     []string{"create"},
		Entity:       types.Entity{ID: "test-entity"},
		Alarm:        types.Alarm{ID: "test-alarm"},
		AckResources: false,
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
	firstExecutionID := "test-alarm&&test-scenario-1"
	firstExecution := action.ScenarioExecution{
		ID:         firstExecutionID,
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
	secondExecutionID := "test-alarm&&test-scenario-2"
	secondExecution := action.ScenarioExecution{
		ID:         secondExecutionID,
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
							Source:          "test",
							Alarm:           task.Alarm,
							Step:            task.Step,
							ExecutionID:     task.ExecutionID,
							AlarmChangeType: types.AlarmChangeType(task.Action.Type),
							Status:          action.TaskNotMatched,
						}
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).
		Return(int64(1), nil).Times(1)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(false)).
		Return(int64(1), nil).Times(1)
	decCall := int64(2)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).
		DoAndReturn(func(_ context.Context, _ string, _ int64, _ bool) (int64, error) {
			decCall--
			return decCall, nil
		}).Times(2)
	createCall := 0
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, _ action.ScenarioExecution) (string, error) {
		createCall++
		return fmt.Sprintf("test-alarm&&test-scenario-%d", createCall), nil
	}).Times(2)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), firstExecutionID).Return(&firstExecution, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), secondExecutionID).Return(&secondExecution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), firstExecutionID).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), secondExecutionID).Return(nil)
	mockScenarioStorage := mock_action.NewMockScenarioStorage(ctrl)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq(task.Triggers), gomock.Eq(task.Alarm)).
		Return([]action.Scenario{firstScenario}, nil)
	mockScenarioStorage.EXPECT().
		GetTriggeredScenarios(gomock.Eq([]string{firstScenario.Actions[0].Type}), gomock.Eq(task.Alarm)).
		Return([]action.Scenario{secondScenario}, nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq(task.Triggers), gomock.Eq(task.Alarm), gomock.Eq(task.Entity)).
		Return(nil)
	mockScenarioStorage.EXPECT().
		RunDelayedScenarios(gomock.Any(), gomock.Eq([]string{firstScenario.Actions[0].Type}), gomock.Eq(task.Alarm), gomock.Eq(task.Entity)).
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
		AckResources:      false,
	}
	executionID := "test-alarm&&test-scenario"
	execution := action.ScenarioExecution{
		ID:         executionID,
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
						Source:      "test",
						Alarm:       task.Alarm,
						Step:        task.Step,
						ExecutionID: task.ExecutionID,
						Status:      action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(1)), gomock.Eq(true)).Return(int64(1), nil)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Eq(int64(-1)), gomock.Eq(false)).Return(int64(0), nil)
	mockExecutionStorage.EXPECT().Create(gomock.Any(), gomock.Any()).Return(executionID, nil)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), executionID).Return(&execution, nil)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), executionID).Return(nil)
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
	executionID := "test-alarm&&test-scenario"
	task := action.ExecuteScenariosTask{
		AbandonedExecutionID: executionID,
		Entity:               types.Entity{ID: "test-entity"},
		Alarm:                types.Alarm{ID: "test-alarm"},
		AckResources:         false,
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
		ID:         executionID,
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
	defer close(taskResultCh)
	mockWorkerPool := mock_action.NewMockWorkerPool(ctrl)
	mockWorkerPool.EXPECT().RunWorkers(gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, taskCh <-chan action.Task) {
			go func() {
				select {
				case <-ctx.Done():
				case task := <-taskCh:
					taskResultCh <- action.TaskResult{
						Source:      "test",
						Alarm:       task.Alarm,
						Step:        task.Step,
						ExecutionID: task.ExecutionID,
						Status:      action.TaskNotMatched,
					}
				}
			}()
		}).
		Return(taskResultCh, nil)
	mockExecutionStorage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	mockExecutionStorage.EXPECT().Inc(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	mockExecutionStorage.EXPECT().Get(gomock.Any(), executionID).Return(&execution, nil).Times(2)
	mockExecutionStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	mockExecutionStorage.EXPECT().Del(gomock.Any(), executionID).Return(nil)
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
