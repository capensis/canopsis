package types_test

import (
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func TestGetTriggers(t *testing.T) {
	tests := []struct {
		name             string
		alarmChange      types.AlarmChange
		expectedTriggers []string
	}{
		{
			name: "none type, no events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeNone,
				EventsCount: 0,
			},
			expectedTriggers: []string{},
		},
		{
			name: "none type, > 0 events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeNone,
				EventsCount: 1,
			},
			expectedTriggers: []string{string(types.AlarmChangeEventsCount)},
		},
		{
			name: "none type, events >= threshold",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeNone,
				EventsCount: 3,
			},
			expectedTriggers: []string{
				string(types.AlarmChangeEventsCount) + "3",
				string(types.AlarmChangeEventsCount),
			},
		},
		{
			name: "state increase, no events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 0,
			},
			expectedTriggers: []string{string(types.AlarmChangeTypeStateIncrease)},
		},
		{
			name: "state increase, > 0 events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 1,
			},
			expectedTriggers: []string{
				string(types.AlarmChangeTypeStateIncrease),
				string(types.AlarmChangeEventsCount),
			},
		},
		{
			name: "state increase, events >= threshold",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 2,
			},
			expectedTriggers: []string{
				string(types.AlarmChangeTypeStateIncrease),
				string(types.AlarmChangeEventsCount) + "2",
				string(types.AlarmChangeEventsCount),
			},
		},
		{
			name:             "create and pbh enter",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeCreateAndPbhEnter},
			expectedTriggers: []string{string(types.AlarmChangeTypeCreate), string(types.AlarmChangeTypePbhEnter)},
		},
		{
			name:             "pbh leave and enter",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypePbhLeaveAndEnter},
			expectedTriggers: []string{string(types.AlarmChangeTypePbhEnter), string(types.AlarmChangeTypePbhLeave)},
		},
		{
			name:             "double ack maps to ack trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeDoubleAck},
			expectedTriggers: []string{string(types.AlarmChangeTypeAck)},
		},
		{
			name:             "webhook start has no triggers",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeWebhookStart},
			expectedTriggers: []string{},
		},
		{
			name:             "declare ticket webhook trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeDeclareTicketWebhook},
			expectedTriggers: []string{string(types.AlarmChangeTypeDeclareTicketWebhook)},
		},
		{
			name:             "snooze trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeSnooze},
			expectedTriggers: []string{string(types.AlarmChangeTypeSnooze)},
		},
		{
			name:        "output change trigger",
			alarmChange: types.AlarmChange{ChangedOutput: true},
			expectedTriggers: []string{
				string(types.AlarmChangeTypeChangedOutput),
			},
		},
		{
			name:        "long output change trigger",
			alarmChange: types.AlarmChange{ChangedLongOutput: true},
			expectedTriggers: []string{
				string(types.AlarmChangeTypeChangedLongOutput),
			},
		},
		{
			name:        "output and long output change triggers",
			alarmChange: types.AlarmChange{ChangedOutput: true, ChangedLongOutput: true},
			expectedTriggers: []string{
				string(types.AlarmChangeTypeChangedOutput),
				string(types.AlarmChangeTypeChangedLongOutput),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.alarmChange.GetTriggers()
			if slices.Compare(result, tt.expectedTriggers) != 0 {
				t.Errorf("expected = %v, want %v", tt.expectedTriggers, result)
			}
		})
	}
}
