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
			expectedTriggers: []string{"eventscount"},
		},
		{
			name: "none type, events >= threshold",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeNone,
				EventsCount: 3,
			},
			expectedTriggers: []string{"eventscount3", "eventscount"},
		},
		{
			name: "state increase, no events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 0,
			},
			expectedTriggers: []string{"stateinc"},
		},
		{
			name: "state increase, > 0 events",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 1,
			},
			expectedTriggers: []string{"stateinc", "eventscount"},
		},
		{
			name: "state increase, events >= threshold",
			alarmChange: types.AlarmChange{
				Type:        types.AlarmChangeTypeStateIncrease,
				EventsCount: 2,
			},
			expectedTriggers: []string{"stateinc", "eventscount2", "eventscount"},
		},
		{
			name:             "create and pbh enter",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeCreateAndPbhEnter},
			expectedTriggers: []string{"create", "pbhenter"},
		},
		{
			name:             "pbh leave and enter",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypePbhLeaveAndEnter},
			expectedTriggers: []string{"pbhenter", "pbhleave"},
		},
		{
			name:             "double ack maps to ack trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeDoubleAck},
			expectedTriggers: []string{"ack"},
		},
		{
			name:             "webhook start has no triggers",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeWebhookStart},
			expectedTriggers: []string{},
		},
		{
			name:             "declare ticket webhook trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeDeclareTicketWebhook},
			expectedTriggers: []string{"declareticketwebhook"},
		},
		{
			name:             "snooze trigger",
			alarmChange:      types.AlarmChange{Type: types.AlarmChangeTypeSnooze},
			expectedTriggers: []string{"snooze"},
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
