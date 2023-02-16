package entityservice_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestGetAlarmCountersFromEvent(t *testing.T) {
	dataSets := getGetAlarmCountersFromEventDataSets()
	for testName, dataSet := range dataSets {
		t.Run(testName, func(t *testing.T) {
			oldCounters, newCounters, isChanged := entityservice.GetAlarmCountersFromEvent(dataSet.event)
			if diff := pretty.Compare(oldCounters, dataSet.expectedOldCounters); diff != "" {
				t.Errorf("unexpected old counters %s", diff)
			}
			if diff := pretty.Compare(newCounters, dataSet.expectedNewCounters); diff != "" {
				t.Errorf("unexpected new counters %s", diff)
			}
			if isChanged != dataSet.expectedIsAlarmChanged {
				t.Errorf("expected %+v but got %+v", dataSet.expectedIsAlarmChanged, isChanged)
			}
		})
	}
}

type getAlarmCountersFromEventDataSet struct {
	event                  types.Event
	expectedOldCounters    *entityservice.AlarmCounters
	expectedNewCounters    *entityservice.AlarmCounters
	expectedIsAlarmChanged bool
}

func getGetAlarmCountersFromEventDataSets() map[string]getAlarmCountersFromEventDataSet {
	maintenancePbhInfo := types.PbehaviorInfo{
		ID:            "test-pbh-1",
		Name:          "test-pbh-1-name",
		ReasonID:      "test-pbh-reason-1",
		TypeID:        "test-pbh-type-1",
		TypeName:      "test-pbh-type-1-name",
		CanonicalType: pbehavior.TypeMaintenance,
	}
	pausePbhInfo := types.PbehaviorInfo{
		ID:            "test-pbh-2",
		Name:          "test-pbh-2-name",
		ReasonID:      "test-pbh-reason-2",
		TypeID:        "test-pbh-type-2",
		TypeName:      "test-pbh-type-2-name",
		CanonicalType: pbehavior.TypePause,
	}
	activePbhInfo := types.PbehaviorInfo{
		ID:            "test-pbh-3",
		Name:          "test-pbh-3-name",
		ReasonID:      "test-pbh-reason-3",
		TypeID:        "test-pbh-type-3",
		TypeName:      "test-pbh-type-3-name",
		CanonicalType: pbehavior.TypeActive,
	}
	return map[string]getAlarmCountersFromEventDataSet{
		"given event without alarm change should return isChanged false": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			},
			expectedIsAlarmChanged: false,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given ack event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeAck,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State: &types.AlarmStep{Value: types.AlarmStateCritical},
					ACK:   &types.AlarmStep{},
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeAck,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State:        entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given ackremove event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeAckremove,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeAckremove,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State:        entityservice.StateCounters{Critical: 1},
			},
		},
		"given create event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeCreate,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given create and pbhenter event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypeCreateAndPbhEnter,
					PreviousPbehaviorTypeID:         maintenancePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: maintenancePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedOldCounters: &entityservice.AlarmCounters{
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given pbhenter event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypePbhEnter,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given pbhenter event to active pbh should return isChanged false": {
			event: types.Event{
				EventType: types.EventTypePbhEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: activePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypePbhEnter,
				},
			},
			expectedIsAlarmChanged: false,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given pbhleave event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhLeave,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State: &types.AlarmStep{Value: types.AlarmStateCritical},
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeave,
					PreviousPbehaviorTypeID:         maintenancePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: maintenancePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given pbhleave from active pbh event should return isChanged false": {
			event: types.Event{
				EventType: types.EventTypePbhLeave,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State: &types.AlarmStep{Value: types.AlarmStateCritical},
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeave,
					PreviousPbehaviorTypeID:         activePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: activePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: false,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given pbhleaveandenter event from maintenance to pause pbh should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhLeaveAndEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: pausePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
					PreviousPbehaviorTypeID:         maintenancePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: maintenancePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{pausePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given pbhleaveandenter event from maintenance to active pbh should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhLeaveAndEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: activePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
					PreviousPbehaviorTypeID:         maintenancePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: maintenancePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given pbhleaveandenter event from active to maintenance pbh should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhLeaveAndEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
					PreviousPbehaviorTypeID:         activePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: activePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given pbhleaveandenter event from active to active pbh should return isChanged false": {
			event: types.Event{
				EventType: types.EventTypePbhLeaveAndEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: activePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
					PreviousPbehaviorTypeID:         activePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: activePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: false,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given resolve event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeResolveCancel,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
				Entity:    &types.Entity{PbehaviorInfo: maintenancePbhInfo},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeResolve,
				},
			},
			expectedIsAlarmChanged: true,
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedNewCounters: &entityservice.AlarmCounters{
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given stateinc event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
				AlarmChange: &types.AlarmChange{
					Type:          types.AlarmChangeTypeStateIncrease,
					PreviousState: types.AlarmStateMajor,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Major: 1},
			},
		},
		"given statedec event should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm:     &types.Alarm{Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
				AlarmChange: &types.AlarmChange{
					Type:          types.AlarmChangeTypeStateDecrease,
					PreviousState: types.AlarmStateCritical,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Major: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State:           entityservice.StateCounters{Critical: 1},
			},
		},
		"given stateinc event with pbehavior should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeCheck,
				Alarm: &types.Alarm{
					Value: types.AlarmValue{
						State:         &types.AlarmStep{Value: types.AlarmStateCritical},
						PbehaviorInfo: maintenancePbhInfo,
					},
				},
				AlarmChange: &types.AlarmChange{
					Type:          types.AlarmChangeTypeStateIncrease,
					PreviousState: types.AlarmStateMajor,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given ack event under maintenance pbehavior should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeAck,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					ACK:           &types.AlarmStep{},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeAck,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State:                entityservice.StateCounters{Ok: 1},
				PbehaviorCounters:    map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given ackremove event under maintenance pbehavior should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypeAckremove,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypeAckremove,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:               1,
				State:             entityservice.StateCounters{Ok: 1},
				PbehaviorCounters: map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State:                entityservice.StateCounters{Ok: 1},
				PbehaviorCounters:    map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
		"given pbhenter event under ack should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhEnter,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State:         &types.AlarmStep{Value: types.AlarmStateCritical},
					ACK:           &types.AlarmStep{},
					PbehaviorInfo: maintenancePbhInfo,
				}},
				AlarmChange: &types.AlarmChange{
					Type: types.AlarmChangeTypePbhEnter,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State:                entityservice.StateCounters{Ok: 1},
				PbehaviorCounters:    map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State:        entityservice.StateCounters{Critical: 1},
			},
		},
		"given pbhleave event under ack should return isChanged true": {
			event: types.Event{
				EventType: types.EventTypePbhLeave,
				Alarm: &types.Alarm{Value: types.AlarmValue{
					State: &types.AlarmStep{Value: types.AlarmStateCritical},
					ACK:   &types.AlarmStep{},
				}},
				AlarmChange: &types.AlarmChange{
					Type:                            types.AlarmChangeTypePbhLeave,
					PreviousPbehaviorTypeID:         maintenancePbhInfo.TypeID,
					PreviousPbehaviorCannonicalType: maintenancePbhInfo.CanonicalType,
				},
			},
			expectedIsAlarmChanged: true,
			expectedNewCounters: &entityservice.AlarmCounters{
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State:        entityservice.StateCounters{Critical: 1},
			},
			expectedOldCounters: &entityservice.AlarmCounters{
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State:                entityservice.StateCounters{Ok: 1},
				PbehaviorCounters:    map[string]int64{maintenancePbhInfo.TypeID: 1},
			},
		},
	}
}
