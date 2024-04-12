package alarm_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	mock_alarm "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarm"
	mock_alarmstatus "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/alarmstatus"
	mock_resolverule "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/resolverule"
	"github.com/golang/mock/gomock"
)

func TestService_ResolveCancels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var dataSets = []struct {
		testName       string
		findAlarms     []types.AlarmWithEntity
		findError      error
		expectedCancel int
	}{
		{
			"given no alarms should return empty result",
			[]types.AlarmWithEntity{},
			nil,
			0,
		},
		{
			"given canceled alarms with cancel time < CancelAutosolveDelay should return empty result",
			[]types.AlarmWithEntity{
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now(),
				}),
			},
			nil,
			0,
		},
		{
			"given canceled alarms and canceled alarms with cancel time > CancelAutosolveDelay should return count of alarms with time > CancelAutosolveDelay",
			[]types.AlarmWithEntity{
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now(),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			nil,
			2,
		},
		{
			"given canceled alarms with valid time should return count of alarms",
			[]types.AlarmWithEntity{
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			nil,
			3,
		},
		{
			"given find error should return error",
			[]types.AlarmWithEntity{
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
				newCancelAlarm(datetime.CpsTime{
					Time: time.Now().Add(-config.AlarmCancelAutosolveDelay),
				}),
			},
			errors.New("not found"),
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
				event.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector),
				log.NewLogger(true),
			)

			events, err := service.ResolveCancels(context.Background(), config.AlarmConfig{
				CancelAutosolveDelay: time.Minute * 60,
			})
			if err != nil {
				expectedErr := fmt.Sprintf("cancel alarms error: %v", dataset.findError.Error())
				if errors.Is(err, errors.New(expectedErr)) {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(events) != dataset.expectedCancel {
				t.Errorf("expected %d cancel alarms but got %d", dataset.expectedCancel, len(events))
			}
		})
	}
}

func TestService_ResolveSnoozes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var dataSets = []struct {
		testName          string
		findAlarms        []types.AlarmWithEntity
		findError         error
		expectedUnsnoozes int
	}{
		{
			"given no alarms should return empty result",
			[]types.AlarmWithEntity{},
			nil,
			0,
		},
		{
			"given snoozed alarms and none unsnoozed alarms should return empty result",
			[]types.AlarmWithEntity{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*3)),
			},
			nil,
			0,
		},
		{
			"given snoozed alarms and snoozed alarms unsnoozed time <= now should return count of unsnoozed alarms",
			[]types.AlarmWithEntity{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			nil,
			2,
		},
		{
			"given snoozed alarms with unsnoozed time <= now should return count of alarms",
			[]types.AlarmWithEntity{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			nil,
			3,
		},
		{
			"given snoozed alarms and snoozed alarms with active pbehavior should return count of active alarms",
			[]types.AlarmWithEntity{
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
			[]types.AlarmWithEntity{
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
			[]types.AlarmWithEntity{
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-1)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-2)),
				newSnoozedAlarm(time.Now(), time.Now().Add(time.Minute*-3)),
			},
			errors.New("not found"),
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
				event.NewGenerator(canopsis.AxeConnector, canopsis.AxeConnector),
				log.NewLogger(true),
			)

			events, err := service.ResolveSnoozes(context.Background(), config.AlarmConfig{})
			if err != nil {
				expectedErr := fmt.Sprintf("snooze alarms error: %v", dataset.findError.Error())
				if errors.Is(err, errors.New(expectedErr)) {
					t.Errorf("expected err %v but got %v", expectedErr, err)
				}
			}

			if len(events) != dataset.expectedUnsnoozes {
				t.Errorf("expected %d unsnoozed alarms but got %d", dataset.expectedUnsnoozes, len(events))
			}
		})
	}
}

func newCancelAlarm(time datetime.CpsTime) types.AlarmWithEntity {
	return types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				Canceled: &types.AlarmStep{
					Type:      types.AlarmStepStatusIncrease,
					Timestamp: time,
				},
			},
		},
		Entity: types.Entity{
			Type: types.EntityTypeResource,
		},
	}
}

func newSnoozedAlarm(snoozeStart time.Time, snoozeEnd time.Time) types.AlarmWithEntity {
	return types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				Snooze: &types.AlarmStep{
					Type:      types.AlarmStepSnooze,
					Timestamp: datetime.NewCpsTime(snoozeStart.Unix()),
					Author:    "",
					Message:   "",
					Value:     types.CpsNumber(snoozeEnd.Unix()),
				},
			},
		},
		Entity: types.Entity{
			Type: types.EntityTypeResource,
		},
	}
}

func newSnoozedAlarmWithActivePbh(snoozeStart time.Time, snoozeEnd time.Time) types.AlarmWithEntity {
	return types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeActive,
				},
				Snooze: &types.AlarmStep{
					Type:      types.AlarmStepSnooze,
					Timestamp: datetime.NewCpsTime(snoozeStart.Unix()),
					Author:    "",
					Message:   "",
					Value:     types.CpsNumber(snoozeEnd.Unix()),
				},
			},
		},
		Entity: types.Entity{
			Type: types.EntityTypeResource,
		},
	}
}

func newSnoozedAlarmWithMaintenancePbh(snoozeStart time.Time, snoozeEnd time.Time) types.AlarmWithEntity {
	return types.AlarmWithEntity{
		Alarm: types.Alarm{
			Value: types.AlarmValue{
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
				Snooze: &types.AlarmStep{
					Type:      types.AlarmStepSnooze,
					Timestamp: datetime.NewCpsTime(snoozeStart.Unix()),
					Author:    "",
					Message:   "",
					Value:     types.CpsNumber(snoozeEnd.Unix()),
				},
			},
		},
		Entity: types.Entity{
			Type: types.EntityTypeResource,
		},
	}
}
