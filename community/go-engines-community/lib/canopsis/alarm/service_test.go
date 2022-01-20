package alarm_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	mock_alarmstatus "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarmstatus"
	mock_resolverule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/resolverule"

	cps "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	"github.com/golang/mock/gomock"
)

func TestService_ResolveDone(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var dataSets = []struct {
		testName     string
		findAlarms   []types.Alarm
		findError    error
		expectedDone int
	}{
		{
			"given no alarms should return empty result",
			[]types.Alarm{},
			nil,
			0,
		},
		{
			"given done alarms with done time < DoneAutosolveDelay should return empty result",
			[]types.Alarm{
				newDoneAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now(),
				}),
			},
			nil,
			0,
		},
		{
			"given done alarms and done alarms with time > DoneAutosolveDelay should return count of alarms with time > DoneAutosolveDelay",
			[]types.Alarm{
				newDoneAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
			},
			nil,
			2,
		},
		{
			"given done alarms with valid time should return count of alarms",
			[]types.Alarm{
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
			},
			nil,
			3,
		},
		{
			"given find error should return error",
			[]types.Alarm{
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
				newDoneAlarm(types.CpsTime{
					Time: time.Now().Add(-cps.DoneAutosolveDelay * time.Second),
				}),
			},
			fmt.Errorf("not found"),
			0,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			mockResolveRuleAdapter := mock_resolverule.NewMockAdapter(ctrl)
			mockAlarmStatusService := mock_alarmstatus.NewMockService(ctrl)
			mockAlarmAdapter.
				EXPECT().
				GetAlarmsWithDoneMark(gomock.Any()).
				Return(dataset.findAlarms, dataset.findError)

			service := alarm.NewService(
				mockAlarmAdapter,
				mockResolveRuleAdapter,
				mockAlarmStatusService,
				log.NewLogger(true),
			)

			doneAlarms, err := service.ResolveDone(context.Background())
			if err != nil {
				expectedErr := fmt.Sprintf("done alarms error: %v", dataset.findError.Error())
				if errors.Is(err, errors.New(expectedErr)) {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(doneAlarms) != dataset.expectedDone {
				t.Errorf("expected %d done alarms but got %d", dataset.expectedDone, len(doneAlarms))
			}
		})
	}
}

func TestService_ResolveCancels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var dataSets = []struct {
		testName       string
		findAlarms     []types.Alarm
		findError      error
		expectedCancel int
	}{
		{
			"given no alarms should return empty result",
			[]types.Alarm{},
			nil,
			0,
		},
		{
			"given canceled alarms with cancel time < CancelAutosolveDelay should return empty result",
			[]types.Alarm{
				newCancelAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now(),
				}),
			},
			nil,
			0,
		},
		{
			"given canceled alarms and canceled alarms with cancel time > CancelAutosolveDelayshould return count of alarms with time > CancelAutosolveDelayshould",
			[]types.Alarm{
				newCancelAlarm(types.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			nil,
			2,
		},
		{
			"given canceled alarms with valid time should return count of alarms",
			[]types.Alarm{
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			nil,
			3,
		},
		{
			"given find error should return error",
			[]types.Alarm{
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(types.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			fmt.Errorf("not found"),
			0,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			mockResolveRuleAdapter := mock_resolverule.NewMockAdapter(ctrl)
			mockAlarmStatusService := mock_alarmstatus.NewMockService(ctrl)
			mockAlarmAdapter.
				EXPECT().
				GetAlarmsWithCancelMark(gomock.Any()).
				Return(dataset.findAlarms, dataset.findError)

			service := alarm.NewService(
				mockAlarmAdapter,
				mockResolveRuleAdapter,
				mockAlarmStatusService,
				log.NewLogger(true),
			)

			cancelAlarms, err := service.ResolveCancels(context.Background(), config.AlarmConfig{
				CancelAutosolveDelay: time.Minute * 60,
			})
			if err != nil {
				expectedErr := fmt.Sprintf("cancel alarms error: %v", dataset.findError.Error())
				if errors.Is(err, errors.New(expectedErr)) {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(cancelAlarms) != dataset.expectedCancel {
				t.Errorf("expected %d done alarms but got %d", dataset.expectedCancel, len(cancelAlarms))
			}
		})
	}
}

func TestService_ResolveSnoozes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var dataSets = []struct {
		testName          string
		findAlarms        []types.Alarm
		findError         error
		expectedUnsnoozes int
	}{
		{
			"given no alarms should return empty result",
			[]types.Alarm{},
			nil,
			0,
		},
		{
			"given snoozed alarms and none unsnoozed alarms should return empty result",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*3)),
			},
			nil,
			0,
		},
		{
			"given snoozed alarms and snoozed alarms unsnoozed time <= now should return count of unsnoozed alarms",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			nil,
			2,
		},
		{
			"given snoozed alarms with unsnoozed time <= now should return count of alarms",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			nil,
			3,
		},
		{
			"given snoozed alarms and snoozed alarms with active pbehavior should return count of active alarms",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarmWithActivePbh(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarmWithActivePbh(time.Now(), time.Now().Add(time.Minute*-3)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-4)),
			},
			nil,
			4,
		},
		{
			"given snoozed alarms and snoozed alarms with maintenance pbehavior should return count of active alarms",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarmWithMaintenancePbh(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarmWithMaintenancePbh(time.Now(), time.Now().Add(time.Minute*-3)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-4)),
			},
			nil,
			2,
		},
		{
			"given find error should return error",
			[]types.Alarm{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			fmt.Errorf("not found"),
			0,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			mockAlarmAdapter := mock_alarm.NewMockAdapter(ctrl)
			mockResolveRuleAdapter := mock_resolverule.NewMockAdapter(ctrl)
			mockAlarmStatusService := mock_alarmstatus.NewMockService(ctrl)
			mockAlarmAdapter.
				EXPECT().
				GetAlarmsWithSnoozeMark(gomock.Any()).
				Return(dataset.findAlarms, dataset.findError)

			service := alarm.NewService(
				mockAlarmAdapter,
				mockResolveRuleAdapter,
				mockAlarmStatusService,
				log.NewLogger(true),
			)

			unsnoozedAlarms, err := service.ResolveSnoozes(context.Background(), config.AlarmConfig{})
			if err != nil {
				expectedErr := fmt.Sprintf("snooze alarms error: %v", dataset.findError.Error())
				if errors.Is(err, errors.New(expectedErr)) {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(unsnoozedAlarms) != dataset.expectedUnsnoozes {
				t.Errorf("expected %d unsnoozed alarms but got %d", dataset.expectedUnsnoozes, len(unsnoozedAlarms))
			}
		})
	}
}

func newDoneAlarm(time types.CpsTime) types.Alarm {
	return types.Alarm{
		Value: types.AlarmValue{
			Done: &types.AlarmStep{
				Type:      types.AlarmStepDone,
				Timestamp: time,
			},
		},
	}
}

func newCancelAlarm(time types.CpsTime) types.Alarm {
	return types.Alarm{
		Value: types.AlarmValue{
			Canceled: &types.AlarmStep{
				Type:      types.AlarmStepCancel,
				Timestamp: time,
			},
		},
	}
}

func newSnoozedAlarm(snoozeStart time.Time, snoozeEnd time.Time) types.Alarm {
	return types.Alarm{
		Value: types.AlarmValue{
			Snooze: &types.AlarmStep{
				Type:      types.AlarmStepSnooze,
				Timestamp: types.NewCpsTime(snoozeStart.Unix()),
				Author:    "",
				Message:   "",
				Value:     types.CpsNumber(snoozeEnd.Unix()),
			},
		},
	}
}

func newSnoozedAlarmWithActivePbh(snoozeStart time.Time, snoozeEnd time.Time) types.Alarm {
	return types.Alarm{
		Value: types.AlarmValue{
			PbehaviorInfo: types.PbehaviorInfo{
				CanonicalType: pbehavior.TypeActive,
			},
			Snooze: &types.AlarmStep{
				Type:      types.AlarmStepSnooze,
				Timestamp: types.NewCpsTime(snoozeStart.Unix()),
				Author:    "",
				Message:   "",
				Value:     types.CpsNumber(snoozeEnd.Unix()),
			},
		},
	}
}

func newSnoozedAlarmWithMaintenancePbh(snoozeStart time.Time, snoozeEnd time.Time) types.Alarm {
	return types.Alarm{
		Value: types.AlarmValue{
			PbehaviorInfo: types.PbehaviorInfo{
				CanonicalType: pbehavior.TypeMaintenance,
			},
			Snooze: &types.AlarmStep{
				Type:      types.AlarmStepSnooze,
				Timestamp: types.NewCpsTime(snoozeStart.Unix()),
				Author:    "",
				Message:   "",
				Value:     types.CpsNumber(snoozeEnd.Unix()),
			},
		},
	}
}
