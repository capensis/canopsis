package action_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	mock_action "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/action"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
)

func TestDelayedScenarioManager_AddDelayedScenario_GivenNotDelayedScenario_ShouldReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second

	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{}
	scenario := action.Scenario{}

	err := manager.AddDelayedScenario(context.Background(), alarm, scenario, action.AdditionalData{})
	if err == nil {
		t.Errorf("expected error but nothing")
	}
}

func TestDelayedScenarioManager_AddDelayedScenario_GivenMatchedDelayedScenario_ShouldReturnTrue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second

	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{
		ID: "test-alarm-id",
		Value: types.AlarmValue{
			Resource: "test-resource",
		},
	}
	scenario := action.Scenario{
		ID: "test-scenario-id",
		Delay: &datetime.DurationWithUnit{
			Value: 10,
			Unit:  "s",
		},
		Actions: []action.Action{
			{
				AlarmPatternFields: newTestMatchResourceAlarmPattern("test-resource"),
			},
		},
	}
	mockStorage.EXPECT().Add(gomock.Any(), gomock.Any()).Do(func(_ context.Context, delayedScenario action.DelayedScenario) {
		if delayedScenario.ScenarioID != scenario.ID {
			t.Errorf("expected scenario id %v but got %v", scenario.ID, delayedScenario.ScenarioID)
		}
		if delayedScenario.AlarmID != alarm.ID {
			t.Errorf("expected alarm id %v but got %v", alarm.ID, delayedScenario.AlarmID)
		}
		expectedExecutionTime := time.Now().Add(time.Second * 10)
		if delayedScenario.ExecutionTime.Unix() != expectedExecutionTime.Unix() {
			t.Errorf("expected execution time %v but got %v", expectedExecutionTime, delayedScenario.ExecutionTime)
		}
	}).Return("test-id", nil)

	err := manager.AddDelayedScenario(context.Background(), alarm, scenario, action.AdditionalData{})
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestDelayedScenarioManager_PauseDelayedScenarios_GivenNotPausedScenario_ShouldUpdateStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{ID: "test-alarm-id"}

	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        false,
			TimeLeft:      0,
		},
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id-2",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        false,
			TimeLeft:      0,
		},
	}, nil)
	mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Do(func(_ context.Context, scenario action.DelayedScenario) {
		if !scenario.Paused {
			t.Errorf("expected paued %v but got %v", true, scenario.Paused)
		}

		if scenario.TimeLeft > time.Second*10 || time.Second*10-scenario.TimeLeft > time.Second {
			t.Errorf("expected time left %v but got %v", time.Second*10, scenario.TimeLeft)
		}
	})

	err := manager.PauseDelayedScenarios(context.Background(), alarm)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestDelayedScenarioManager_PauseDelayedScenarios_GivenPausedScenario_ShouldNotUpdateStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{ID: "test-alarm-id"}

	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        true,
			TimeLeft:      time.Second * 5,
		},
	}, nil)
	mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	err := manager.PauseDelayedScenarios(context.Background(), alarm)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestDelayedScenarioManager_ResumeDelayedScenarios_GivenPausedScenario_ShouldUpdateStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{ID: "test-alarm-id"}

	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        true,
			TimeLeft:      time.Second * 5,
		},
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id-2",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        true,
			TimeLeft:      time.Second * 5,
		},
	}, nil)
	mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Do(func(_ context.Context, scenario action.DelayedScenario) {
		if scenario.Paused {
			t.Errorf("expected paused %v but got %v", false, scenario.Paused)
		}

		if scenario.TimeLeft != 0 {
			t.Errorf("expected time left %v but got %v", 0, scenario.TimeLeft)
		}

		expectedExecutionTime := time.Now().Add(time.Second * 5)
		if scenario.ExecutionTime.Unix() != expectedExecutionTime.Unix() {
			t.Errorf("expected execution time %v but got %v", expectedExecutionTime, scenario.ExecutionTime)
		}
	})

	err := manager.ResumeDelayedScenarios(context.Background(), alarm)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestDelayedScenarioManager_ResumeDelayedScenarios_GivenNotPausedScenario_ShouldNotUpdateStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})
	alarm := types.Alarm{ID: "test-alarm-id"}

	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		{
			ScenarioID:    "test-scenario-id",
			AlarmID:       "test-alarm-id",
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * time.Second)},
			Paused:        false,
			TimeLeft:      0,
		},
	}, nil)
	mockStorage.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)

	err := manager.ResumeDelayedScenarios(context.Background(), alarm)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}
}

func TestDelayedScenarioManager_Run_GivenExpiredScenario_ShouldReturnItByTick(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter,
		mockStorage, periodicalTimeout, zerolog.Logger{})

	expectedAlarm := types.Alarm{ID: "test-alarm-id"}
	expectedScenario := action.Scenario{ID: "test-scenario-id"}
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		{
			ID:            "test-delayed-id",
			ScenarioID:    expectedScenario.ID,
			AlarmID:       expectedAlarm.ID,
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(periodicalTimeout + time.Millisecond)},
		},
		{
			ID:            "test-delayed-id-2",
			ScenarioID:    expectedScenario.ID,
			AlarmID:       expectedAlarm.ID,
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * periodicalTimeout)},
		},
	}, nil).Times(2)
	mockStorage.EXPECT().Delete(gomock.Any(), gomock.Eq("test-delayed-id")).Return(true, nil)
	mockActionAdapter.EXPECT().GetEnabledByIDs(gomock.Any(), gomock.Eq([]string{expectedScenario.ID})).Return([]action.Scenario{expectedScenario}, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsWithEntityByAlarmIDs(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, ids []string, alarms *[]types.AlarmWithEntity) {
			if !reflect.DeepEqual(ids, []string{expectedAlarm.ID}) {
				t.Errorf("expected %v but got %v", []string{expectedAlarm.ID}, ids)
			}

			*alarms = []types.AlarmWithEntity{{Alarm: expectedAlarm}}
		}).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := manager.Run(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	time.Sleep(periodicalTimeout + 10*time.Millisecond)

	select {
	case task := <-ch:
		if !reflect.DeepEqual(task.Scenario, expectedScenario) {
			t.Errorf("expected scenario %v but got %v", expectedScenario, task.Scenario)
		}
		if !reflect.DeepEqual(task.Alarm, expectedAlarm) {
			t.Errorf("expected alarm %v but got %v", expectedAlarm, task.Alarm)
		}
	default:
		t.Errorf("expected task but go nothing")
	}
}

func TestDelayedScenarioManager_Run_GivenExpiredScenario_ShouldReturnItByWaitingGoroutine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockActionAdapter := mock_action.NewMockAdapter(ctrl)
	mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
	mockStorage := mock_action.NewMockDelayedScenarioStorage(ctrl)
	periodicalTimeout := 2 * time.Second
	manager := action.NewDelayedScenarioManager(mockActionAdapter, mockAlarmAdapter, mockStorage, periodicalTimeout, zerolog.Logger{})

	expectedAlarm := types.Alarm{ID: "test-alarm-id"}
	expectedScenario := action.Scenario{ID: "test-scenario-id"}
	delayedScenario := action.DelayedScenario{
		ID:            "test-delayed-id",
		ScenarioID:    expectedScenario.ID,
		AlarmID:       expectedAlarm.ID,
		ExecutionTime: datetime.CpsTime{Time: time.Now().Add(3500 * time.Millisecond)},
	}
	mockStorage.EXPECT().GetAll(gomock.Any()).Return([]action.DelayedScenario{
		delayedScenario,
		{
			ID:            "test-delayed-id-2",
			ScenarioID:    expectedScenario.ID,
			AlarmID:       expectedAlarm.ID,
			ExecutionTime: datetime.CpsTime{Time: time.Now().Add(10 * periodicalTimeout)},
		},
	}, nil).Times(2)
	mockStorage.EXPECT().Get(gomock.Any(), gomock.Eq(delayedScenario.ID)).Return(&delayedScenario, nil)
	mockStorage.EXPECT().Delete(gomock.Any(), gomock.Eq("test-delayed-id")).Return(true, nil)
	mockActionAdapter.EXPECT().GetEnabledByIDs(gomock.Any(), gomock.Eq([]string{expectedScenario.ID})).Return([]action.Scenario{expectedScenario}, nil)
	mockAlarmAdapter.EXPECT().GetOpenedAlarmsWithEntityByAlarmIDs(gomock.Any(), gomock.Any(), gomock.Any()).Do(func(_ context.Context, ids []string, alarms *[]types.AlarmWithEntity) {
		if !reflect.DeepEqual(ids, []string{expectedAlarm.ID}) {
			t.Errorf("expected %v but got %v", []string{expectedAlarm.ID}, ids)
		}

		*alarms = []types.AlarmWithEntity{{Alarm: expectedAlarm}}
	}).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch, err := manager.Run(ctx)
	if err != nil {
		t.Errorf("expected not error but got %v", err)
	}

	time.Sleep(3600 * time.Millisecond)

	select {
	case task := <-ch:
		if !reflect.DeepEqual(task.Scenario, expectedScenario) {
			t.Errorf("expected scenario %v but got %v", expectedScenario, task.Scenario)
		}
		if !reflect.DeepEqual(task.Alarm, expectedAlarm) {
			t.Errorf("expected alarm %v but got %v", expectedAlarm, task.Alarm)
		}
	default:
		t.Errorf("expected task but go nothing")
	}
}

func newTestMatchResourceAlarmPattern(resource string) savedpattern.AlarmPatternFields {
	return savedpattern.AlarmPatternFields{
		AlarmPattern: pattern.Alarm{
			{
				{
					Field:     "v.resource",
					Condition: pattern.NewStringCondition("eq", resource),
				},
			},
		},
	}
}
