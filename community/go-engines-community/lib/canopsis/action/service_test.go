package action_test

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_amqp "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/amqp"
	mock_action "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/action"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_encoding "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/encoding"
	"github.com/golang/mock/gomock"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestService_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

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

	logger := zerolog.Nop()
	scenarioExecChan := make(chan action.ExecuteScenariosTask)
	defer close(scenarioExecChan)

	amqpChannelMock := mock_amqp.NewMockChannel(ctrl)
	delayedScenarioManager := mock_action.NewMockDelayedScenarioManager(ctrl)
	alarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	storage := mock_action.NewMockScenarioExecutionStorage(ctrl)
	activationService := mock_alarm.NewMockActivationService(ctrl)
	actionService := action.NewService(alarmAdapter, scenarioExecChan, delayedScenarioManager,
		storage, json.NewEncoder(), json.NewDecoder(), amqpChannelMock, canopsis.FIFOAckExchangeName,
		canopsis.FIFOAckQueueName, activationService, logger)

	var dataSets = []struct {
		testName string
		event    *types.Event
		triggers []string
	}{
		{
			testName: "given create event should exec scenario",
			event: &types.Event{
				Alarm: &types.Alarm{
					ID: "alarm-1",
				},
				Entity: &types.Entity{
					ID: "entity-1",
				},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeCreate,
				},
			},
			triggers: []string{"create"},
		},
		{
			testName: "given snooze event should exec scenario",
			event: &types.Event{
				Alarm: &types.Alarm{
					ID: "alarm-2",
				},
				Entity: &types.Entity{
					ID: "entity-2",
				},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeSnooze,
				},
			},
			triggers: []string{"snooze"},
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			go func() {
				err := actionService.Process(ctx, dataset.event)
				if err != nil {
					t.Errorf("expected not error but got %v", err)
				}
			}()

			time.Sleep(time.Millisecond * 150)

			select {
			case scenarioExec := <-scenarioExecChan:
				if dataset.event.Alarm.ID != scenarioExec.Alarm.ID {
					t.Errorf("ScenarioExec message should have an alarm with id = %s, got %s", dataset.event.Alarm.ID, scenarioExec.Alarm.ID)
				}
				if dataset.event.Entity.ID != scenarioExec.Entity.ID {
					t.Errorf("ScenarioExec message should have an entity with id = %s, got %s", dataset.event.Entity.ID, scenarioExec.Entity.ID)
				}
				if !reflect.DeepEqual(dataset.triggers, scenarioExec.Triggers) {
					t.Errorf("ScenarioExec message should have an triggers %+v, got %+v", dataset.triggers, scenarioExec.Triggers)
				}
			default:
				t.Error("ScenarioExec message is expected")
			}
		})
	}
}

func TestService_ListenScenarioFinish(t *testing.T) {
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

	alarm1 := types.Alarm{
		ID: "alarm-1",
		Value: types.AlarmValue{
			Component: "component-1",
			Resource:  "resource-1",
		},
	}
	alarm2 := types.Alarm{
		ID: "alarm-2",
		Value: types.AlarmValue{
			Component: "component-2",
			Resource:  "resource-2",
		},
	}
	logger := zerolog.Nop()

	var dataSets = []struct {
		testName      string
		scenarioInfos []action.ScenarioResult
	}{
		{
			testName: "given succes message should publish fifo ack event",
			scenarioInfos: []action.ScenarioResult{
				{Alarm: alarm1},
			},
		},
		{
			testName: "given failed message should do nothing",
			scenarioInfos: []action.ScenarioResult{
				{
					Alarm: alarm1,
					Err:   errors.New("test error"),
				},
			},
		},
		{
			testName: "given failed message and then success should continue work",
			scenarioInfos: []action.ScenarioResult{
				{
					Alarm: alarm1,
					Err:   errors.New("test error"),
				},
				{
					Alarm: alarm2,
				},
			},
		},
	}

	scenarioInfoChannel := make(chan action.ScenarioResult)

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx, cancel := context.WithCancel(timerCtx)
			defer cancel()

			scenarioExecChan := make(chan action.ExecuteScenariosTask)
			defer close(scenarioExecChan)
			amqpChannelMock := mock_amqp.NewMockChannel(ctrl)
			encoderMock := mock_encoding.NewMockEncoder(ctrl)
			decoderMock := mock_encoding.NewMockDecoder(ctrl)
			delayedScenarioManager := mock_action.NewMockDelayedScenarioManager(ctrl)
			alarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			storage := mock_action.NewMockScenarioExecutionStorage(ctrl)
			activationService := mock_alarm.NewMockActivationService(ctrl)
			actionService := action.NewService(alarmAdapter, scenarioExecChan, delayedScenarioManager,
				storage, encoderMock, decoderMock, amqpChannelMock, canopsis.FIFOAckExchangeName,
				canopsis.FIFOAckQueueName, activationService, logger)

			actionService.ListenScenarioFinish(ctx, scenarioInfoChannel)

			getInOrder := make([]*gomock.Call, 0)
			processInOrder := make([]*gomock.Call, 0)
			encodeInOrder := make([]*gomock.Call, 0)
			publishInOrder := make([]*gomock.Call, 0)

			for _, v := range dataset.scenarioInfos {
				info := v
				get := alarmAdapter.EXPECT().GetAlarmByAlarmId(gomock.Any(), gomock.Eq(info.Alarm.ID)).
					Return(info.Alarm, nil)
				getInOrder = append(getInOrder, get)

				if info.Err == nil {
					process := activationService.EXPECT().Process(gomock.Any(), gomock.Any(), gomock.Any()).
						Do(func(_ context.Context, alarm types.Alarm, _ types.Event) {
							if alarm.ID != info.Alarm.ID {
								t.Errorf("expected alarm %s but got %s", info.Alarm.ID, alarm.ID)
							}
						}).
						Return(false, nil)
					processInOrder = append(processInOrder, process)
				}

				body := []byte("body " + info.Alarm.ID)
				encode := encoderMock.
					EXPECT().
					Encode(gomock.Any()).
					Do(func(event types.Event) {
						if event.Alarm.ID != info.Alarm.ID {
							t.Errorf("expected alarm %s but got %s", info.Alarm.ID, event.Alarm.ID)
						}
					}).
					Return(body, nil)
				encodeInOrder = append(encodeInOrder, encode)
				publish := amqpChannelMock.
					EXPECT().
					PublishWithContext(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(_ context.Context, exchange, key string, mandatory, immediate bool, msg amqp.Publishing) {
						if !reflect.DeepEqual(msg.Body, body) {
							t.Errorf("expected event %s but got %s", body, msg.Body)
						}
					}).
					Return(nil)
				publishInOrder = append(publishInOrder, publish)
			}

			gomock.InOrder(getInOrder...)
			gomock.InOrder(processInOrder...)
			gomock.InOrder(encodeInOrder...)
			gomock.InOrder(publishInOrder...)

			go func() {
				for _, info := range dataset.scenarioInfos {
					select {
					case <-ctx.Done():
					case scenarioInfoChannel <- info:
					}
				}
			}()

			time.Sleep(time.Millisecond * 150)
		})
	}
}

func TestService_ProcessAbandonedExecutions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	timerCtx, timerCancel := context.WithCancel(context.Background())
	defer timerCancel()

	go func(ctx context.Context) {
		deadlockTimer := time.NewTimer(5000 * time.Second)

		select {
		case <-ctx.Done():
			return
		case <-deadlockTimer.C:
			panic("workers or test are deadlocked")
		}
	}(timerCtx)

	logger := zerolog.Nop()
	amqpChannelMock := mock_amqp.NewMockChannel(ctrl)
	encoderMock := mock_encoding.NewMockEncoder(ctrl)
	decoderMock := mock_encoding.NewMockDecoder(ctrl)
	delayedScenarioManager := mock_action.NewMockDelayedScenarioManager(ctrl)
	activationService := mock_alarm.NewMockActivationService(ctrl)
	alarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	alarmAdapter.
		EXPECT().
		GetOpenedAlarmByAlarmId(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(_ context.Context, id string) (types.Alarm, error) {
			var err error
			if id != "test-alarm" {
				err = mongo.ErrNoDocuments
			}

			return types.Alarm{}, err
		})

	var dataSets = []struct {
		testName            string
		abandonedExecutions []action.ScenarioExecution
		executionKey        string
		expectDelete        bool
		expectExecute       bool
	}{
		{
			testName: "given completely executed abandoned execution should be deleted",
			abandonedExecutions: []action.ScenarioExecution{
				{
					AlarmID:    "test-alarm",
					ScenarioID: "test-scenario",
					Entity:     types.Entity{},
					ActionExecutions: []action.Execution{
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: true,
						},
					},
				},
			},
			executionKey:  "test-alarm$$test-scenario",
			expectDelete:  true,
			expectExecute: false,
		},
		{
			testName: "given execution with resolved/deleted alarm should be deleted",
			abandonedExecutions: []action.ScenarioExecution{
				{
					AlarmID:    "test-alarm-not-exist",
					ScenarioID: "test-scenario",
					Entity:     types.Entity{},
					ActionExecutions: []action.Execution{
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
					},
				},
			},
			executionKey:  "test-alarm-not-exist$$test-scenario",
			expectDelete:  true,
			expectExecute: false,
		},
		{
			testName: "given execution should send exec task",
			abandonedExecutions: []action.ScenarioExecution{
				{
					AlarmID:    "test-alarm",
					ScenarioID: "test-scenario",
					Entity:     types.Entity{},
					ActionExecutions: []action.Execution{
						{
							Action:   action.Action{},
							Executed: false,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
					},
				},
			},
			executionKey:  "test-alarm$$test-scenario",
			expectDelete:  false,
			expectExecute: true,
		},
		{
			testName: "given execution should send exec task",
			abandonedExecutions: []action.ScenarioExecution{
				{
					AlarmID:    "test-alarm",
					ScenarioID: "test-scenario",
					Entity:     types.Entity{},
					ActionExecutions: []action.Execution{
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
					},
				},
			},
			executionKey:  "test-alarm$$test-scenario",
			expectDelete:  false,
			expectExecute: true,
		},
		{
			testName: "given execution should send exec task",
			abandonedExecutions: []action.ScenarioExecution{
				{
					AlarmID:    "test-alarm",
					ScenarioID: "test-scenario",
					Entity:     types.Entity{},
					ActionExecutions: []action.Execution{
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: true,
						},
						{
							Action:   action.Action{},
							Executed: false,
						},
					},
				},
			},
			executionKey:  "test-alarm$$test-scenario",
			expectDelete:  false,
			expectExecute: true,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			scenarioExecChan := make(chan action.ExecuteScenariosTask)
			defer close(scenarioExecChan)

			storage := mock_action.NewMockScenarioExecutionStorage(ctrl)
			storage.
				EXPECT().
				GetAbandoned(gomock.Any()).
				Times(1).
				Return(dataset.abandonedExecutions, nil)

			if dataset.expectDelete {
				storage.
					EXPECT().
					Del(gomock.Any(), gomock.Eq(dataset.executionKey)).
					Times(1).
					Return(nil)
			}

			actionService := action.NewService(alarmAdapter, scenarioExecChan, delayedScenarioManager,
				storage, encoderMock, decoderMock, amqpChannelMock, canopsis.FIFOAckExchangeName,
				canopsis.FIFOAckQueueName, activationService, logger)

			var wg sync.WaitGroup
			if dataset.expectExecute {
				wg.Add(1)

				go func() {
					defer wg.Done()
					scenarioExec := <-scenarioExecChan
					if scenarioExec.AbandonedExecutionCacheKey != dataset.executionKey {
						t.Errorf("Scenario exec task should be marked as 'abandoned' but got %+v", scenarioExec)
					}
				}()
			}

			err := actionService.ProcessAbandonedExecutions(ctx)
			if err != nil {
				t.Errorf("error %s is not expected.", err.Error())
			}

			wg.Wait()
		})
	}
}
