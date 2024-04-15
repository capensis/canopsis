package pbehavior_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	libevent "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/event"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/golang/mock/gomock"
)

func TestGetEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()

	canonicalActiveInfo := types.PbehaviorInfo{}
	activeInfo := types.PbehaviorInfo{
		ID:            "test-pbh-active",
		Name:          "Pbh active",
		ReasonID:      "Reason active",
		ReasonName:    "Reason active name",
		TypeID:        "test-type-active",
		TypeName:      "Active type",
		CanonicalType: pbehavior.TypeActive,
		Author:        canopsis.DefaultEventAuthor,
	}
	anotherActiveInfo := types.PbehaviorInfo{
		ID:            "test-another-pbh-active",
		Name:          "Another pbh active",
		ReasonID:      "Another reason active",
		ReasonName:    "Another reason active name",
		TypeID:        "test-another-type-active",
		TypeName:      "Another active type",
		CanonicalType: pbehavior.TypeActive,
		Author:        canopsis.DefaultEventAuthor,
	}
	maintenanceInfo := types.PbehaviorInfo{
		ID:            "test-pbh-maintenance",
		Name:          "Pbh maintenance",
		ReasonID:      "Reason maintenance",
		ReasonName:    "Reason maintenance name",
		TypeID:        "test-type-maintenance",
		TypeName:      "Maintenance type",
		CanonicalType: pbehavior.TypeMaintenance,
		Author:        canopsis.DefaultEventAuthor,
	}
	anotherMaintenanceInfo := types.PbehaviorInfo{
		ID:            "test-another-pbh-maintenance",
		Name:          "Another pbh maintenance",
		ReasonID:      "Another reason maintenance",
		ReasonName:    "Another reason maintenance name",
		TypeID:        "test-another-type-maintenance",
		TypeName:      "Another maintenance type",
		CanonicalType: pbehavior.TypeMaintenance,
		Author:        canopsis.DefaultEventAuthor,
	}

	resolvedCanonicalActive := pbehavior.ResolveResult{}
	resolvedActive := pbehavior.ResolveResult{
		Type: pbehavior.Type{
			ID:   "test-type-active",
			Name: "Active type",
			Type: pbehavior.TypeActive,
		},
		ID:         "test-pbh-active",
		Name:       "Pbh active",
		ReasonID:   "Reason active",
		ReasonName: "Reason active name",
	}
	resolvedAnotherActive := pbehavior.ResolveResult{
		Type: pbehavior.Type{
			ID:   "test-another-type-active",
			Name: "Another active type",
			Type: pbehavior.TypeActive,
		},
		ID:         "test-another-pbh-active",
		Name:       "Another pbh active",
		ReasonID:   "Another reason active",
		ReasonName: "Another reason active name",
	}
	resolvedMaintenance := pbehavior.ResolveResult{
		Type: pbehavior.Type{
			ID:   "test-type-maintenance",
			Name: "Maintenance type",
			Type: pbehavior.TypeMaintenance,
		},
		ID:         "test-pbh-maintenance",
		Name:       "Pbh maintenance",
		ReasonID:   "Reason maintenance",
		ReasonName: "Reason maintenance name",
	}
	resolvedAnotherMaintenance := pbehavior.ResolveResult{
		Type: pbehavior.Type{
			ID:   "test-another-type-maintenance",
			Name: "Another maintenance type",
			Type: pbehavior.TypeMaintenance,
		},
		ID:         "test-another-pbh-maintenance",
		Name:       "Another pbh maintenance",
		ReasonID:   "Another reason maintenance",
		ReasonName: "Another reason maintenance name",
	}

	var dataSets = []struct {
		testName             string
		entity               types.Entity
		resolveResult        pbehavior.ResolveResult
		expectedEventType    string
		expectedAlarmPbhInfo types.PbehaviorInfo
	}{
		{
			"An alarm has no behaviors and resolved type is canonical active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: canonicalActiveInfo,
			},
			resolvedCanonicalActive,
			"",
			types.PbehaviorInfo{},
		},
		{
			"An alarm has no behaviors and resolved type is active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: canonicalActiveInfo,
			},
			resolvedActive,
			types.EventTypePbhEnter,
			activeInfo,
		},
		{
			"An alarm has no behaviors and resolved type is maintenance behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: canonicalActiveInfo,
			},
			resolvedMaintenance,
			types.EventTypePbhEnter,
			maintenanceInfo,
		},
		{
			"An alarm has an active behavior and resolved type is canonical active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: activeInfo,
			},
			resolvedCanonicalActive,
			types.EventTypePbhLeave,
			types.PbehaviorInfo{},
		},
		{
			"An alarm has an active behavior and resolved type is maintenance behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: activeInfo,
			},
			resolvedMaintenance,
			types.EventTypePbhLeaveAndEnter,
			maintenanceInfo,
		},
		{
			"An alarm has an active behavior and resolved type is the same behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: activeInfo,
			},
			resolvedActive,
			"",
			types.PbehaviorInfo{},
		},
		{
			"An alarm has an active behavior and resolved type is another active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: activeInfo,
			},
			resolvedAnotherActive,
			types.EventTypePbhLeaveAndEnter,
			anotherActiveInfo,
		},
		{
			"An alarm has a maintenance behavior and resolved type is canonical active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: maintenanceInfo,
			},
			resolvedCanonicalActive,
			types.EventTypePbhLeave,
			types.PbehaviorInfo{},
		},
		{
			"An alarm has a maintenance behavior and resolved type is the same behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: maintenanceInfo,
			},
			resolvedMaintenance,
			"",
			types.PbehaviorInfo{},
		},
		{
			"An alarm has a maintenance behavior and resolved type is active behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: maintenanceInfo,
			},
			resolvedActive,
			types.EventTypePbhLeaveAndEnter,
			activeInfo,
		},
		{
			"An alarm has a maintenance behavior and resolved type is another maintenance behavior",
			types.Entity{
				Type:          types.EntityTypeResource,
				PbehaviorInfo: maintenanceInfo,
			},
			resolvedAnotherMaintenance,
			types.EventTypePbhLeaveAndEnter,
			anotherMaintenanceInfo,
		},
	}

	for _, dataset := range dataSets {
		t.Run(dataset.testName, func(t *testing.T) {
			manager := pbehavior.NewEventManager(libevent.NewGenerator("", ""))

			event, err := manager.GetEvent(dataset.resolveResult, dataset.entity, datetime.NewCpsTime())
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if event.EventType != dataset.expectedEventType {
				t.Errorf("expected %s event type, got %s", dataset.expectedEventType, event.EventType)
			}

			event.PbehaviorInfo.Timestamp = nil
			if event.PbehaviorInfo != dataset.expectedAlarmPbhInfo {
				t.Errorf("expected events's pbehavior info = %v, got %v", dataset.expectedAlarmPbhInfo, event.PbehaviorInfo)
			}
		})
	}
}
