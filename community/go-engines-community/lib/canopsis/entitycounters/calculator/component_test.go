package calculator_test

import (
	"context"
	"maps"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_entitycounters "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entitycounters"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type componentDataset struct {
	name                              string
	alarm                             *types.Alarm
	entity                            types.Entity
	alarmChange                       types.AlarmChange
	counters                          entitycounters.EntityCounters
	expectedDiff                      map[string]int
	expectedState                     int
	expectedUpdated                   bool
	countersCallDoNotExpected         bool
	cleanAddRemoveFieldsDoNotExpected bool
}

func TestComponentService_ProcessCounters_GivenAlarmChangeNone(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity without an alarm and alarm change is none, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is none, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is none, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		// inactive
		{
			name:  "given new not yet counted entity without an alarm with inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an ok alarm with inactive pbh, and alarm change is none, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a ko alarm with inactive pbh, and alarm change is none, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// when already counted
		// active
		{
			name:  "given counted entity without an alarm, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is none, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is none, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// inactive
		{
			name:  "given counted entity without an alarm with inactive pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is none with inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is none with inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters without an alarm, should update counters and shouldn't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is none, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is none, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// inactive
		{
			name:  "given an entity to be removed from counters without an alarm with inactive pbh, should update counters and shouldn't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is none with inactive pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is none with inactive pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given new not yet counted entity with a minor alarm and alarm change is none, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a critical alarm and alarm change is none, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.critical": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateCritical,
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMajor,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangeState(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity with a state change from ok to ko alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		{
			name:  "given new not yet counted entity with a state change from ko to ok alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a state change action, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.critical": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateCritical,
		},
		// inactive
		{
			name:  "given new not yet counted entity with a state change from ok to ko alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a state change from ko to ok alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a state change action with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// when already counted
		// active
		{
			name:  "given counted entity with a state change from ok to ko alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a state change from ko to ok alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateOK,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a state change action, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor":    -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// inactive
		{
			name:  "given counted entity with a state change from ok to ko alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a state change from ko to ok alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a state change action with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters entity with a state change from ok to ko alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with a state change from ko to ok alarm, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		{
			name:  "given an entity to be removed from counters entity with a state change action, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// inactive
		{
			name:  "given an entity to be removed from counters entity with a state change from ok to ko alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with a state change from ko to ok alarm with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with a state change action with inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmStepChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is stateinc, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is stateinc, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 2,
				},
			},
			expectedDiff: map[string]int{
				"state.major":    -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is statedec, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.critical": -1,
				"state.major":    1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is changestate, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangeCreate(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity with a new alarm in minor state, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		{
			name:  "given new not yet counted entity with a new alarm in major state, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMajor,
		},
		{
			name:  "given new not yet counted entity with a new alarm in critical state, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.critical": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateCritical,
		},
		// when already counted
		{
			name:  "given counted entity with an alarm in minor state, should increment minor and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in major state, should increment major and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in critical state, should increment critical and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		{
			name:  "given an entity to be removed from counters with a new alarm in minor state, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in major state, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in critical state, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is create, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is create, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 2,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is create, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is create, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangeCreateAndPbhEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity with a new alarm in minor state with inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a new alarm in major state with inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a new alarm in critical state with inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with a new alarm in minor state with active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		{
			name:  "given new not yet counted entity with a new alarm in major state with active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMajor,
		},
		{
			name:  "given new not yet counted entity with a new alarm in critical state with active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.critical": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateCritical,
		},
		// when already counted
		{
			name:  "given counted entity with a new alarm in minor state with inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a new alarm in major state with inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a new alarm in critical state with inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a new alarm in minor state with active pbh, should increment minor and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a new alarm in major state with active pbh, should increment major and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a new alarm alarm in critical state with active pbh, should increment critical and decrement ok counter and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		{
			name:  "given an entity to be removed from counters with a new alarm in minor state with inactive pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in major state with inactive pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in critical state with inactive pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in minor state with active pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in major state with active pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters with a new alarm in critical state with active pbh, should update counters and shouldn't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is createandpbhenter, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is createandpbhenter, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 2,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is createandpbhenter, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is createandpbhenter, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangePbhEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity without alarm when enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMajor,
		},
		// when already counted
		{
			name:  "given counted entity without alarm when enters inactive pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when enters inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when enters inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateOK,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when enters active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when enters active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
			countersCallDoNotExpected:         true,
		},
		// when entity should be removed from counters
		{
			name:  "given entity to be removed from counters without alarm when enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ok state when enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ko state when enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ok state when enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ko state when enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is pbhenter, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is pbhenter, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       1,
				"state.critical": -1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhenter, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    1,
				"state.major": -1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhenter, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    1,
				"state.major": -1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangePbhLeave(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity without alarm when leaves inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.major": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMajor,
		},
		// when already counted
		{
			name:  "given counted entity without alarm when leaves inactive pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": 1,
				"state.ok":    -1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
			countersCallDoNotExpected:         true,
		},
		// when entity should be removed from counters
		{
			name:  "given entity to be removed from counters without alarm when leaves inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ok state when leaves inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ko state when leaves inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ok state when leaves active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given entity to be removed from counters with an alarm in ko state when leaves active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is pbhleave, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": 1,
				"state.ok":    -1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is pbhleave, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major: 1,
					Ok:    1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       -1,
				"state.critical": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateCritical,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhleave, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhleave, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangePbhLeaveAndEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	dataSets := []componentDataset{
		// when entity should be added to counters
		// active to active
		{
			name:  "given new not yet counted entity without alarm when leaves active pbh and enters active pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves active pbh and enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves active pbh and enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		// active to inactive
		{
			name:  "given new not yet counted entity without alarm when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// inactive to active
		{
			name:  "given new not yet counted entity without alarm when leaves inactive pbh and enters active pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves inactive pbh and enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves inactive pbh and enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.minor": 1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateMinor,
		},
		// inactive to inactive
		{
			name:  "given new not yet counted entity without alarm when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// when already counted
		// active to active
		{
			name:  "given counted entity without alarm when leaves active pbh and enters active pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves active pbh and enters active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves active pbh and enters active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// active to inactive
		{
			name:  "given counted entity without alarm when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves active pbh and enters inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateOK,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// inactive to active
		{
			name:  "given counted entity without alarm when leaves inactive pbh and enters active pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves inactive pbh and enters active pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves inactive pbh and enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.minor": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// inactive to inactive
		{
			name:  "given counted entity without alarm when leaves inactive pbh and enters inactive pbh, shouldn't do any updates",
			alarm: nil,
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when leaves inactive pbh and enters inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when leaves inactive pbh and enters inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active to active
		{
			name:  "given an entity to be removed from counters entity without alarm when leaves active pbh and enters active pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ok state when leaves active pbh and enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ko state when leaves active pbh and enters active pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// active to inactive
		{
			name:  "given an entity to be removed from counters entity without alarm when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ok state when leaves active pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ko state when leaves active pbh and enters inactive pbh, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		// inactive to active
		{
			name:  "given an entity to be removed from counters entity without alarm when leaves inactive pbh and enters active pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ok state when leaves inactive pbh and enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ko state when leaves inactive pbh and enters active pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// inactive to inactive
		{
			name:  "given an entity to be removed from counters entity without alarm when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: nil,
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ok state when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given an entity to be removed from counters entity with an alarm in ko state when leaves inactive pbh and enters inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given counted entity with a minor alarm and alarm change is pbhleaveandenter, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a critical alarm and alarm change is pbhleaveandenter, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":       1,
				"state.critical": -1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhleaveandenter, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with a major alarm and alarm change is pbhleaveandenter, when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  60,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    -1,
				"state.major": 1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMajor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func TestComponentService_ProcessCounters_GivenAlarmChangeResolve(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersHelper, entityCollection, countersCollection := prepareComponentTest(ctrl)

	datasets := []componentDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity with an alarm in ok state when resolves, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when resolves, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ok state when resolves in inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given new not yet counted entity with an alarm in ko state when resolves in inactive, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                     true,
				Name:                        "test-resource-1",
				Type:                        types.EntityTypeResource,
				ComponentStateSettings:      true,
				ComponentStateSettingsToAdd: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{},
			expectedDiff: map[string]int{
				"state.ok": 1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// when already counted
		{
			name:  "given counted entity with an alarm in ok state when resolves, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when resolves, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateOK,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ok state when resolves in inactive pbh, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm in ko state when resolves in inactive, shouldn't do any updates",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                true,
				Name:                   "test-resource-1",
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters:                          entitycounters.EntityCounters{},
			expectedDiff:                      map[string]int{},
			expectedUpdated:                   false,
			expectedState:                     0,
			countersCallDoNotExpected:         true,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		// when entity should be removed from counters
		{
			name:  "given counted entity with an alarm in ok state when resolves, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given counted entity with an alarm in ko state when resolves, should update counters and update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": -1,
			},
			expectedUpdated: true,
			expectedState:   types.AlarmStateOK,
		},
		{
			name:  "given counted entity with an alarm in ok state when resolves in inactive pbh, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		{
			name:  "given counted entity with an alarm in ko state when resolves in inactive, should update counters and don't update component state",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:                        true,
				Name:                           "test-resource-1",
				Type:                           types.EntityTypeResource,
				ComponentStateSettingsToRemove: true,
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok": -1,
			},
			expectedUpdated: false,
			expectedState:   0,
		},
		// state changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, should update component state, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateMinor,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Critical: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleMajor,
							Value:  3,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   false,
			expectedState:                     0,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve when rule type is dependencies, should update component state, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:                true,
				Type:                   types.EntityTypeResource,
				ComponentStateSettings: true,
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondLT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  40,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major": -1,
				"state.ok":    1,
			},
			expectedUpdated:                   true,
			expectedState:                     types.AlarmStateOK,
			cleanAddRemoveFieldsDoNotExpected: true,
		},
	}

	for _, dSet := range datasets {
		t.Run(dSet.name, func(t *testing.T) {
			runComponentsDataset(ctx, t, entityHelper, countersHelper, entityCollection, countersCollection, componentService, dSet)
		})
	}
}

func prepareComponentTest(ctrl *gomock.Controller) (
	calculator.ComponentCountersCalculator,
	*mock_mongo.MockSingleResultHelper,
	*mock_mongo.MockSingleResultHelper,
	*mock_mongo.MockDbCollection,
	*mock_mongo.MockDbCollection,
) {
	entityHelper := mock_mongo.NewMockSingleResultHelper(ctrl)
	countersHelper := mock_mongo.NewMockSingleResultHelper(ctrl)

	entityCollection := mock_mongo.NewMockDbCollection(ctrl)
	entityCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(entityHelper).AnyTimes()

	countersCollection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(libmongo.EntityMongoCollection).Return(entityCollection).AnyTimes()
	dbClient.EXPECT().Collection(libmongo.EntityCountersCollection).Return(countersCollection).AnyTimes()

	eventSender := mock_entitycounters.NewMockEventsSender(ctrl)

	return calculator.NewComponentCountersCalculator(dbClient, eventSender), entityHelper, countersHelper, entityCollection, countersCollection
}

func runComponentsDataset(
	ctx context.Context,
	t *testing.T,
	entityHelper,
	countersHelper *mock_mongo.MockSingleResultHelper,
	entityCollection *mock_mongo.MockDbCollection,
	countersCollection *mock_mongo.MockDbCollection,
	componentService calculator.ComponentCountersCalculator,
	dSet componentDataset) {
	entityHelper.EXPECT().Decode(gomock.Any()).Do(func(v *entitycounters.StateSettingsInfo) {
		*v = entitycounters.StateSettingsInfo{
			ComponentStateSettings:         dSet.entity.ComponentStateSettings,
			ComponentStateSettingsToAdd:    dSet.entity.ComponentStateSettingsToAdd,
			ComponentStateSettingsToRemove: dSet.entity.ComponentStateSettingsToRemove,
		}
	}).Return(nil)

	if !dSet.countersCallDoNotExpected {
		countersCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(countersHelper)
		countersHelper.EXPECT().Decode(gomock.Any()).Do(func(v *entitycounters.EntityCounters) {
			*v = dSet.counters
		}).Return(nil)
	}

	if len(dSet.expectedDiff) > 0 {
		countersCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, _ any, update any, _ ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
				m, ok := update.(bson.M)
				if !ok {
					t.Fatal("update argument should be a bson.M")
				}

				rawInc, ok := m["$inc"]
				if !ok {
					t.Fatal("update document doesn't contain $inc")
				}

				inc, ok := rawInc.(map[string]int)
				if !ok {
					t.Fatal("wrong")
				}

				if !maps.Equal(inc, dSet.expectedDiff) {
					t.Errorf("expected counters diff = %v, but got %v", dSet.expectedDiff, inc)
				}

				return nil, nil
			},
		)
	}

	if !dSet.cleanAddRemoveFieldsDoNotExpected {
		entityCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)
	}

	_, resUpdated, resState, err := componentService.CalculateCounters(ctx, &dSet.entity, dSet.alarm, dSet.alarmChange)
	if err != nil {
		t.Fatalf("expected no err, but got %v", err)
	}

	if resUpdated != dSet.expectedUpdated {
		t.Errorf("expected updated = %v, but got %v", dSet.expectedUpdated, resUpdated)
	}

	if resState != dSet.expectedState {
		t.Errorf("expected state = %d, but got %d", dSet.expectedState, resState)
	}
}
