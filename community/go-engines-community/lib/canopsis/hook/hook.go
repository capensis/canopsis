package hook

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// Hook represents the conditions when the hook must be triggered
type Hook struct {
	Triggers       []string                  `bson:"triggers" json:"triggers"`
	AlarmPatterns  pattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EventPatterns  pattern.EventPatternList  `bson:"event_patterns" json:"event_patterns"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
}

// IsTriggeredByEvent checks if a hook must be activated according an event
func (h *Hook) IsTriggeredByEvent(event types.Event) bool {
	return h.AlarmPatterns.Matches(event.Alarm) && h.EventPatterns.Matches(event) && h.EntityPatterns.Matches(event.Entity) && h.IsTriggeredByAlarmChange(*event.AlarmChange)
}

// IsTriggeredByAlarmChange checks if a hook must be activated according the triggers list defined be the user
func (h *Hook) IsTriggeredByAlarmChange(alarmChange types.AlarmChange) bool {
	for _, trigger := range h.Triggers {
		for _, acTrigger := range alarmChange.GetTriggers() {
			if acTrigger == trigger {
				return true
			}
		}
	}
	return false
}

// IsTriggeredByAlarmChangeType checks if a hook must be activated accoding the triggers list defined be the user
func (h *Hook) IsTriggeredByAlarmChangeType(change types.AlarmChangeType) bool {
	alarmChange := types.AlarmChange{Type: change}
	for _, trigger := range h.Triggers {
		for _, acTrigger := range alarmChange.GetTriggers() {
			if acTrigger == trigger {
				return true
			}
		}
	}
	return false
}
