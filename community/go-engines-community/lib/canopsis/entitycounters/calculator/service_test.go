package calculator_test

import (
	"context"
	"maps"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entitycounters/calculator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_config "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/config"
	mock_entitycounters "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/canopsis/entitycounters"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type serviceDataset struct {
	name                      string
	alarm                     *types.Alarm
	entity                    types.Entity
	alarmChange               types.AlarmChange
	counters                  entitycounters.EntityCounters
	expectedDiff              map[string]int
	expectedStateServicesInfo map[string]entitycounters.UpdatedServicesInfo
	countersCallDoNotExpected bool
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeNone(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity without an alarm and alarm change is none, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is none, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"acked":          1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity without an alarm in inactive pbh and alarm change is none, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm in inactive pbh and alarm change is none, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given counted entity without an alarm and alarm change is none, shouldn't check counters",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters:                  entitycounters.EntityCounters{},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm and alarm change is none, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange:               types.AlarmChange{Type: types.AlarmChangeTypeNone},
			counters:                  entitycounters.EntityCounters{},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters without an alarm and alarm change is none, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is none, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:           "service-1",
				Depends:      1,
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given an entity to be removed from counters without an alarm in inactive pbh and alarm change is none, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm in inactive pbh and alarm change is none, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm in inactive pbh and alarm change is none, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID:                   "service-1",
				Depends:              1,
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is none, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-2",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is none, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"inherited_state.major": 1,
				"state.major":           1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is none, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeNone,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "2"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeState(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"acked":          1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.minor":           1,
				"inherited_state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is changestate, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"unacked":        1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is changestate, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given counted entity with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major":           1,
				"inherited_state.major": 1,
				"state.minor":           -1,
				"inherited_state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.ok":    1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.minor":           1,
				"inherited_state.minor": 1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is changestate, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.minor":    -1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given counted entity with a ko alarm in inactive pbh and alarm change is stateinc, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with a ko alarm in inactive pbh and alarm change is statedec, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is changestate, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters entity with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity matched to inherited rule with a ko alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.minor":           -1,
				"inherited_state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters entity with an acked alarm and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters entity with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters entity matched to inherited rule with a ko alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters entity with an acked alarm and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is changestate, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given an entity to be removed from counters entity with a ko alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateOK,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity with an acked alarm in inactive pbh and alarm change is stateinc, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity with a ko alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateMajor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters entity with an acked alarm in inactive pbh and alarm change is statedec, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm in inactive pbh and alarm change is changestate, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given counted entity with a major alarm and alarm change is stateinc, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
					Minor:    1,
				},
			},
			expectedDiff: map[string]int{
				"state.critical": 1,
				"state.minor":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a major alarm and alarm change is stateinc, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"state.major": 1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity with a major alarm and alarm change is statedec, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"state.major":    1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a major alarm and alarm change is statedec when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateDecrease,
				PreviousState: types.AlarmStateCritical,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       8,
					Major:    1,
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  19,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.major":    1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is changestate, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeChangeState,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"state.critical": 1,
				"state.minor":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is changestate, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"inherited_state.minor":    -1,
				"state.minor":              -1,
				"inherited_state.critical": 1,
				"state.critical":           1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		{
			name:  "given counted entity with a major alarm and alarm change is stateinc, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor:    1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"state.major": 1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a major alarm and alarm change is stateinc, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:          types.AlarmChangeTypeStateIncrease,
				PreviousState: types.AlarmStateMinor,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
					Minor: 1,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"state.major": 1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "2"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeCreate(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is create, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is create, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// when already counted
		{
			name:  "given counted entity with a ko alarm and alarm change is create, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.minor": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is create, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		// when entity should be removed from counters
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is create, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is create, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":            -1,
				"state.ok":           -1,
				"inherited_state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given given counted entity with a major alarm and alarm change is create, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
					Major: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  65,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is create, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is create, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"inherited_state.major": 1,
				"state.major":           1,
				"inherited_state.ok":    -1,
				"state.ok":              -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreate,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "2"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeCreateAndPbhEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is create with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is create with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is create with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is create with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given counted entity with a ko alarm and alarm change is create with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.minor": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is create with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is create with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is create with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
				Depends: 1,
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when entity should be removed from counters
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is create with active pbh, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is create with active pbh, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":            -1,
				"state.ok":           -1,
				"inherited_state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is create with inactive pbh, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is create with inactive pbh, should change service counters and shouldn't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":            -1,
				"state.ok":           -1,
				"inherited_state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
					Major: 3,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						Major: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleMajor,
							Value:  65,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is create with active pbh, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is create with active pbh, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"inherited_state.major": 1,
				"state.major":           1,
				"inherited_state.ok":    -1,
				"state.ok":              -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is create with active pbh, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeCreateAndPbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 1,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
				"state.ok":    -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "2"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangePbhEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhenter with active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhenter with active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		// active
		{
			name:  "given counted entity without alarm and alarm change is pbhenter with active pbh, shouldn't check counters",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm and alarm change is pbhenter with active pbh, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// inactive
		{
			name:  "given counted entity without alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an acked ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"acked_under_pbh":       1,
				"acked":                 -1,
				"state.ok":              1,
				"state.minor":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                   -1,
				"state.ok":                 1,
				"state.critical":           -1,
				"inherited_state.ok":       1,
				"inherited_state.critical": -1,
				"under_pbh":                1,
				"pbehavior.maintenance":    1,
				"unacked":                  -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// when entity should be removed from counters
		{
			name:  "given entity to be removed from counters without alarm and alarm change is pbhenter with active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an ok alarm and alarm change is pbhenter with active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an acked ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters with a ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhenter with active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  -1,
				"all":                      -1,
				"active":                   -1,
				"unacked":                  -1,
				"state.critical":           -1,
				"inherited_state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given entity to be removed from counters without alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an ok alarm and alarm change is pbhenter with inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an acked ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhenter with inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  -1,
				"all":                      -1,
				"active":                   -1,
				"unacked":                  -1,
				"state.critical":           -1,
				"inherited_state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// state and output changes
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
							Value:  1,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 3,
					Major: 6,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Value:  15,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhenter with inactive pbh, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 2,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhenter with inactive pbh, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"inherited_state.ok":    1,
				"inherited_state.major": -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhenter with inactive pbh, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypePbhEnter,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 2,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "2",
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "1"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangePbhLeave(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleave from active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleave from active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// when already counted
		// active
		{
			name:  "given counted entity without alarm and alarm change is pbhleave from active pbh, shouldn't check counters",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm and alarm change is pbhleave from active pbh, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// inactive
		{
			name:  "given counted entity without alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an acked ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"acked_under_pbh":       -1,
				"acked":                 1,
				"state.ok":              -1,
				"state.minor":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                   1,
				"state.ok":                 -1,
				"state.critical":           1,
				"inherited_state.ok":       -1,
				"inherited_state.critical": 1,
				"under_pbh":                -1,
				"pbehavior.maintenance":    -1,
				"unacked":                  1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// when entity should be removed from counters
		{
			name:  "given entity to be removed from counters without alarm and alarm change is pbhleave from active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an ok alarm and alarm change is pbhleave from active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an acked ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters with a ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleave from active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  -1,
				"all":                      -1,
				"active":                   -1,
				"unacked":                  -1,
				"state.critical":           -1,
				"inherited_state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given entity to be removed from counters without alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an ok alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with an acked ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleave from inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
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
							Value:  1,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    2,
					Minor: 3,
					Major: 5,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Value:  15,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleave from inactive pbh, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleave from inactive pbh, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"inherited_state.ok":    -1,
				"inherited_state.major": 1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive pbh, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeave,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 2,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "2",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "3"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangePbhLeaveAndEnter(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active to active
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// active to inactive
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive to active
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  1,
				"all":                      1,
				"active":                   1,
				"unacked":                  1,
				"state.critical":           1,
				"inherited_state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive to inactive
		{
			name:  "given new not yet counted entity without alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		// active to active
		{
			name:  "given counted entity without alarm and alarm change is pbhleaveandenter from active to active pbh, shouldn't do anything",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhleaveandenter from active to active pbh, shouldn't do anything",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// active to inactive
		{
			name:  "given counted entity without alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an acked ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"acked_under_pbh":       1,
				"acked":                 -1,
				"state.ok":              1,
				"state.minor":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                   -1,
				"state.ok":                 1,
				"state.critical":           -1,
				"inherited_state.ok":       1,
				"inherited_state.critical": -1,
				"under_pbh":                1,
				"pbehavior.maintenance":    1,
				"unacked":                  -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive to active
		{
			name:  "given counted entity without alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an acked ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"acked_under_pbh":       -1,
				"acked":                 1,
				"state.ok":              -1,
				"state.minor":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                   1,
				"state.ok":                 -1,
				"state.critical":           1,
				"inherited_state.ok":       -1,
				"inherited_state.critical": 1,
				"under_pbh":                -1,
				"pbehavior.maintenance":    -1,
				"unacked":                  1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive to inactive
		{
			name:  "given counted entity without alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should update counters",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"pbehavior.prev-maintenance": -1,
				"pbehavior.maintenance":      1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an ok alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should update counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"pbehavior.prev-maintenance": -1,
				"pbehavior.maintenance":      1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when entity should be removed from counters
		// active to active
		{
			name:  "given an entity to be removed from counters without alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from active to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  -1,
				"all":                      -1,
				"active":                   -1,
				"unacked":                  -1,
				"state.critical":           -1,
				"inherited_state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// active to inactive
		{
			name:  "given an entity to be removed from counters without alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"acked":       -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Critical: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                  -1,
				"all":                      -1,
				"active":                   -1,
				"unacked":                  -1,
				"state.critical":           -1,
				"inherited_state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive to active
		{
			name:  "given an entity to be removed from counters without alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive to inactive
		{
			name:  "given an entity to be removed from counters without alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":                    -1,
				"state.ok":                   -1,
				"under_pbh":                  -1,
				"pbehavior.prev-maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":                    -1,
				"all":                        -1,
				"state.ok":                   -1,
				"under_pbh":                  -1,
				"pbehavior.prev-maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":                    -1,
				"all":                        -1,
				"acked_under_pbh":            -1,
				"state.ok":                   -1,
				"under_pbh":                  -1,
				"pbehavior.prev-maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":                    -1,
				"all":                        -1,
				"state.ok":                   -1,
				"under_pbh":                  -1,
				"pbehavior.prev-maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is pbhleaveandenter from inactive to inactive pbh, should change service counters and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":                    -1,
				"all":                        -1,
				"state.ok":                   -1,
				"inherited_state.ok":         -1,
				"under_pbh":                  -1,
				"pbehavior.prev-maintenance": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from active to inactive pbh, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 2,
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
					Major: 1,
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
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
							Value:  1,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    2,
					Minor: 3,
					Major: 5,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Value:  15,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from active to inactive pbh, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 2,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodNumber,
							State:  types.AlarmStateTitleOK,
							Value:  0,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from active to inactive pbh when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 3,
					Major: 5,
				},
				Rule: &statesetting.StateSetting{
					Method: statesetting.MethodDependencies,
					StateThresholds: &statesetting.StateThresholds{
						OK: &statesetting.StateThreshold{
							Cond:   statesetting.CalculationCondGT,
							Method: statesetting.CalculationMethodShare,
							State:  types.AlarmStateTitleOK,
							Value:  15,
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.major":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleaveandenter from inactive to active pbh, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleave from inactive to active pbh, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					CanonicalType: types.PbhCanonicalTypeActive,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"inherited_state.ok":    -1,
				"inherited_state.major": 1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleaveandenter from active to inactive pbh, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 2,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.minor":           -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a ko alarm and alarm change is pbhleave from active to inactive pbh, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: types.PbhCanonicalTypeActive,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Minor: 2,
				},
				InheritedState: entitycounters.StateCounters{
					Ok:    1,
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"active":                -1,
				"state.ok":              1,
				"state.minor":           -1,
				"inherited_state.ok":    1,
				"inherited_state.minor": -1,
				"under_pbh":             1,
				"pbehavior.maintenance": 1,
				"unacked":               -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleaveandenter from inactive to active pbh, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorTypeID:         "maintenance",
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:       1,
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"active":                1,
				"state.ok":              -1,
				"state.major":           1,
				"under_pbh":             -1,
				"pbehavior.maintenance": -1,
				"unacked":               1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given given counted entity with a major alarm and alarm change is pbhleave from inactive to inactive pbh, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type:                            types.AlarmChangeTypePbhLeaveAndEnter,
				PreviousPbehaviorCannonicalType: pbehavior.TypeMaintenance,
				PreviousPbehaviorTypeID:         "prev-maintenance",
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Ok:    1,
					Major: 2,
				},
				PbehaviorCounters: map[string]int{
					"prev-maintenance": 1,
				},
				OutputTemplate: "{{ .PbehaviorCounters }}",
				Output:         "map[prev-maintenance:1]",
			},
			expectedDiff: map[string]int{
				"pbehavior.prev-maintenance": -1,
				"pbehavior.maintenance":      1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "map[maintenance:1 prev-maintenance:0]"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeResolve(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":            1,
				"state.ok":           1,
				"inherited_state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive
		{
			name:  "given new not yet counted entity with an ok alarm in inactive pbh and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given counted entity with an ok alarm and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"all":     -1,
				"active":  -1,
				"unacked": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"state.ok":              1,
				"inherited_state.major": -1,
				"inherited_state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given counted entity with an acked alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
				"state.ok":       1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given counted entity with an ok alarm in inactive pbh and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"all": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"all": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given counted entity with an acked alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"all":             -1,
				"acked_under_pbh": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:           "service-1",
				Depends:      1,
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given an entity to be removed from counters with an ok alarm in inactive pbh and alarm change is resolve, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm in inactive pbh and alarm change is resolve, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID:                   "service-1",
				Depends:              1,
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major:    1,
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
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
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
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
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
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
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is resolve, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-2",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
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
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is resolve, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Name:     "test-resource-1",
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
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
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"inherited_state.major": -1,
				"state.major":           -1,
				"inherited_state.ok":    1,
				"state.ok":              1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
					Major:    1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is resolve, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeResolve,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 2,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "2",
			},
			expectedDiff: map[string]int{
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.major": -1,
				"state.ok":    1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "1"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeAck(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// active
		{
			name:  "given new not yet counted entity matched to inherited rule with an acked alarm and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"acked":                 1,
				"state.major":           1,
				"inherited_state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"acked":          1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
				"acked_under_pbh":       1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"acked":   1,
				"unacked": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"acked_under_pbh": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given new not yet counted entity matched to inherited rule with an acked alarm and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"unacked":        -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is ack, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is ack, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-2",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is ack, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"acked":                 1,
				"inherited_state.major": 1,
				"state.major":           1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"acked":       1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ack, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAck,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				NotAcknowledged: 1,
				OutputTemplate:  "{{.NotAcknowledged}} {{.Acknowledged}}",
				Output:          "1 0",
			},
			expectedDiff: map[string]int{
				"acked":   1,
				"unacked": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "0 1"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeAckRemove(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// active
		{
			name:  "given new not yet counted entity matched to inherited rule with an alarm and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with an alarm and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"unacked":        1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an alarm in inactive pbh and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given new not yet counted entity with an alarm and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"unacked": 1,
				"acked":   -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// inactive
		{
			name:  "given new not yet counted entity with an alarm in inactive pbh and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"acked_under_pbh": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given new not yet counted entity matched to inherited rule with an alarm and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"acked":                 -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given new not yet counted entity with an alarm and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
				"acked_under_pbh":       -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an alarm in inactive pbh and alarm change is ackremove, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
				"acked_under_pbh":       -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is ackremove, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-2",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is ackremove, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"inherited_state.major": 1,
				"state.major":           1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is ackremove, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeAckremove,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				Acknowledged:   1,
				OutputTemplate: "{{.NotAcknowledged}} {{.Acknowledged}}",
				Output:         "0 1",
			},
			expectedDiff: map[string]int{
				"acked":   -1,
				"unacked": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "1 0"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func TestEntityServiceService_ProcessCounters_GivenAlarmChangeEnabled(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	componentService, entityHelper, countersCollection := prepareEntityServiceTest(ctrl)

	dataSets := []serviceDataset{
		// when entity should be added to counters
		// active
		{
			name:  "given new not yet counted entity without an alarm and alarm change is enabled, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm and alarm change is enabled, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":  1,
				"all":      1,
				"active":   1,
				"unacked":  1,
				"state.ok": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.minor": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMinor},
			},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"state.major":           1,
				"inherited_state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with an acked alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":        1,
				"all":            1,
				"active":         1,
				"acked":          1,
				"state.critical": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateCritical},
			},
		},
		// inactive
		{
			name:  "given new not yet counted entity without an alarm in inactive pbh and alarm change is enabled, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an ok alarm in inactive pbh and alarm change is enabled, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity matched to inherited rule with a ko alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"state.ok":              1,
				"inherited_state.ok":    1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with an acked alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"acked_under_pbh":       1,
				"state.ok":              1,
				"pbehavior.maintenance": 1,
				"under_pbh":             1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// when already counted
		{
			name:  "given counted entity without an alarm and alarm change is enabled, shouldn't check counters",
			alarm: nil,
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters:                  entitycounters.EntityCounters{},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		{
			name:  "given counted entity with an alarm and alarm change is enabled, shouldn't check counters",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:  true,
				Type:     types.EntityTypeResource,
				Services: []string{"service-1"},
			},
			alarmChange:               types.AlarmChange{Type: types.AlarmChangeTypeEnabled},
			counters:                  entitycounters.EntityCounters{},
			expectedDiff:              map[string]int{},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
			countersCallDoNotExpected: true,
		},
		// when entity should be removed from counters
		// active
		{
			name:  "given an entity to be removed from counters without an alarm and alarm change is enabled, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm and alarm change is enabled, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":  -1,
				"all":      -1,
				"active":   -1,
				"unacked":  -1,
				"state.ok": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     -1,
				"all":         -1,
				"active":      -1,
				"unacked":     -1,
				"state.minor": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:              "service-1",
				Depends:         1,
				All:             1,
				Active:          1,
				NotAcknowledged: 1,
				State: entitycounters.StateCounters{
					Major: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Major: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"active":                -1,
				"unacked":               -1,
				"state.major":           -1,
				"inherited_state.major": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:           "service-1",
				Depends:      1,
				All:          1,
				Active:       1,
				Acknowledged: 1,
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":        -1,
				"all":            -1,
				"active":         -1,
				"acked":          -1,
				"state.critical": -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateOK},
			},
		},
		// inactive
		{
			name:  "given an entity to be removed from counters without an alarm in inactive pbh and alarm change is enabled, should change service counters and don't update service",
			alarm: nil,
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an ok alarm in inactive pbh and alarm change is enabled, should change service counter and don't update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateOK}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with a ko alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMinor}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters matched to inherited rule with a ko alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:          true,
				Name:             "test-resource-1",
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:      "service-1",
				Depends: 1,
				All:     1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"state.ok":              -1,
				"inherited_state.ok":    -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given an entity to be removed from counters with an acked alarm in inactive pbh and alarm change is enabled, should change service counters and update service",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{ACK: &types.AlarmStep{}, State: &types.AlarmStep{Value: types.AlarmStateCritical}}},
			entity: types.Entity{
				Enabled:          true,
				Type:             types.EntityTypeResource,
				ServicesToRemove: []string{"service-1"},
				PbehaviorInfo: types.PbehaviorInfo{
					TypeID:        "maintenance",
					CanonicalType: pbehavior.TypeMaintenance,
				},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID:                   "service-1",
				Depends:              1,
				All:                  1,
				AcknowledgedUnderPbh: 1,
				State: entitycounters.StateCounters{
					Ok: 1,
				},
				PbehaviorCounters: map[string]int{"maintenance": 1},
				UnderPbehavior:    1,
			},
			expectedDiff: map[string]int{
				"depends":               -1,
				"all":                   -1,
				"acked_under_pbh":       -1,
				"state.ok":              -1,
				"pbehavior.maintenance": -1,
				"under_pbh":             -1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		// state and output changes
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled, when rule type is dependencies, shouldn't do any updates as worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled when rule type is dependencies, should update service, because the worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
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
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is enabled, when rule type is inherited, shouldn't update service, because the inherited worst state didn't changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-2",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a ko alarm and alarm change is enabled, when rule type is inherited, should update service, because the inherited worst state changed",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Name:          "test-resource-1",
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Minor: 1,
				},
				InheritedState: entitycounters.StateCounters{
					Minor: 1,
				},
				Rule: &statesetting.StateSetting{
					Type:   statesetting.RuleTypeService,
					Method: statesetting.MethodInherited,
					InheritedEntityPattern: &pattern.Entity{
						{
							{
								Field:     "name",
								Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-resource-1"),
							},
						},
					},
				},
			},
			expectedDiff: map[string]int{
				"depends":               1,
				"all":                   1,
				"active":                1,
				"unacked":               1,
				"inherited_state.major": 1,
				"state.major":           1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{
				"service-1": {State: types.AlarmStateMajor},
			},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled, shouldn't update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Critical: 1,
				},
				OutputTemplate: "{{.State.Critical}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{},
		},
		{
			name:  "given new not yet counted entity with a major alarm and alarm change is enabled, should update output",
			alarm: &types.Alarm{ID: "test-alarm", Value: types.AlarmValue{State: &types.AlarmStep{Value: types.AlarmStateMajor}}},
			entity: types.Entity{
				Enabled:       true,
				Type:          types.EntityTypeResource,
				Services:      []string{"service-1"},
				ServicesToAdd: []string{"service-1"},
			},
			alarmChange: types.AlarmChange{
				Type: types.AlarmChangeTypeEnabled,
			},
			counters: entitycounters.EntityCounters{
				ID: "service-1",
				State: entitycounters.StateCounters{
					Major: 1,
				},
				OutputTemplate: "{{.State.Major}}",
				Output:         "1",
			},
			expectedDiff: map[string]int{
				"depends":     1,
				"all":         1,
				"active":      1,
				"unacked":     1,
				"state.major": 1,
			},
			expectedStateServicesInfo: map[string]entitycounters.UpdatedServicesInfo{"service-1": {State: types.AlarmStateMajor, Output: "2"}},
		},
	}

	for _, dSet := range dataSets {
		t.Run(dSet.name, func(t *testing.T) {
			runServicesDataset(ctx, ctrl, t, entityHelper, countersCollection, componentService, dSet)
		})
	}
}

func runServicesDataset(
	ctx context.Context,
	ctrl *gomock.Controller,
	t *testing.T,
	entityHelper *mock_mongo.MockSingleResultHelper,
	countersCollection *mock_mongo.MockDbCollection,
	componentService calculator.EntityServiceCountersCalculator,
	dSet serviceDataset) {
	entityHelper.EXPECT().Decode(gomock.Any()).Do(func(v *entitycounters.ServicesInfo) {
		*v = entitycounters.ServicesInfo{
			Services:         dSet.entity.Services,
			ServicesToAdd:    dSet.entity.ServicesToAdd,
			ServicesToRemove: dSet.entity.ServicesToRemove,
		}
	}).Return(nil)

	if !dSet.countersCallDoNotExpected {
		mockCursor := mock_mongo.NewMockCursor(ctrl)
		mockCursor.EXPECT().Next(gomock.Any()).Return(true)
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
		mockCursor.EXPECT().Decode(gomock.Any()).Do(func(v *entitycounters.EntityCounters) {
			*v = dSet.counters
		}).Return(nil)
		mockCursor.EXPECT().Close(gomock.Any()).Return(nil)

		countersCollection.EXPECT().Find(gomock.Any(), gomock.Any()).Return(mockCursor, nil)
	}

	if len(dSet.expectedDiff) > 0 {
		countersCollection.EXPECT().BulkWrite(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, models []mongo.WriteModel, _ ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
				for _, model := range models {
					updateOneModel, ok := model.(*mongo.UpdateOneModel)
					if !ok {
						t.Fatal("write models should be mongo.UpdateOneModel models")
					}

					m, ok := updateOneModel.Update.(bson.M)
					if !ok {
						t.Fatal("mongo.UpdateOneModel Update struct should be a bson.M")
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
				}

				return nil, nil
			},
		)
	}

	resultStateServicesInfo, err := componentService.CalculateCounters(ctx, &dSet.entity, dSet.alarm, dSet.alarmChange)
	if err != nil {
		t.Fatalf("expected no err, but got %v", err)
	}

	if !maps.Equal(resultStateServicesInfo, dSet.expectedStateServicesInfo) {
		t.Fatalf("expected state services info = %v, but got = %v", dSet.expectedStateServicesInfo, resultStateServicesInfo)
	}
}

func prepareEntityServiceTest(ctrl *gomock.Controller) (
	calculator.EntityServiceCountersCalculator,
	*mock_mongo.MockSingleResultHelper,
	*mock_mongo.MockDbCollection,
) {
	entityHelper := mock_mongo.NewMockSingleResultHelper(ctrl)

	entityCollection := mock_mongo.NewMockDbCollection(ctrl)
	entityCollection.EXPECT().FindOne(gomock.Any(), gomock.Any(), gomock.Any()).Return(entityHelper).AnyTimes()
	entityCollection.EXPECT().UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	countersCollection := mock_mongo.NewMockDbCollection(ctrl)

	dbClient := mock_mongo.NewMockDbClient(ctrl)
	dbClient.EXPECT().Collection(libmongo.EntityMongoCollection).Return(entityCollection).AnyTimes()
	dbClient.EXPECT().Collection(libmongo.EntityCountersCollection).Return(countersCollection).AnyTimes()

	templateConfigProvider := mock_config.NewMockTemplateConfigProvider(ctrl)
	templateConfigProvider.EXPECT().Get().Return(config.SectionTemplate{}).AnyTimes()
	timezoneConfigProvider := mock_config.NewMockTimezoneConfigProvider(ctrl)
	timezoneConfigProvider.EXPECT().Get().Return(config.TimezoneConfig{}).AnyTimes()

	templateExecutor := template.NewExecutor(templateConfigProvider, timezoneConfigProvider)

	eventSender := mock_entitycounters.NewMockEventsSender(ctrl)

	return calculator.NewEntityServiceCountersCalculator(dbClient, templateExecutor, eventSender), entityHelper, countersCollection
}
