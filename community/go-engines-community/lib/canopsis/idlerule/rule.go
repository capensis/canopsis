// Package idlerule contains idle rule model and adapter.
package idlerule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	RuleTypeAlarm                = "alarm"
	RuleTypeEntity               = "entity"
	RuleAlarmConditionLastEvent  = "last_event"  // alarms which don't receive events after "now - duration"
	RuleAlarmConditionLastUpdate = "last_update" // alarms which don't change after "now - duration"
)

type Operation struct {
	Type       string     `bson:"type" json:"type"`
	Parameters Parameters `bson:"parameters,omitempty" json:"parameters"`
}

// Rule represents alarm modification condition and operation.
type Rule struct {
	ID          string                    `bson:"_id,omitempty" json:"_id"`
	Name        string                    `bson:"name" json:"name"`
	Description string                    `bson:"description" json:"description"`
	Author      string                    `bson:"author" json:"author"`
	Enabled     bool                      `bson:"enabled" json:"enabled"`
	Type        string                    `bson:"type" json:"type"`
	Priority    int64                     `bson:"priority" json:"priority"`
	Duration    datetime.DurationWithUnit `bson:"duration" json:"duration"`
	Comment     string                    `bson:"comment" json:"comment"`
	// DisableDuringPeriods is an option that allows to disable the rule
	// when entity is in listed periods due pbehavior schedule.
	DisableDuringPeriods []string         `bson:"disable_during_periods" json:"disable_during_periods"`
	Created              datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty"`
	Updated              datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty"`
	// Only for Alarm rules
	AlarmCondition string     `bson:"alarm_condition,omitempty" json:"alarm_condition,omitempty"`
	Operation      *Operation `bson:"operation,omitempty" json:"operation,omitempty"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

type Parameters struct {
	Output string `json:"output" bson:"output,omitempty" binding:"max=255"`

	// ChangeState
	State *types.CpsNumber `json:"state" bson:"state,omitempty"`
	// AssocTicket
	// Ticket is used in assocticket action.
	Ticket string `json:"ticket,omitempty" binding:"max=255" bson:"ticket,omitempty"`
	// TicketURL is used in assocticket action.
	TicketURL string `json:"ticket_url,omitempty" binding:"max=255" bson:"ticket_url,omitempty"`
	// TicketSystemName is used in assocticket action.
	TicketSystemName string `json:"ticket_system_name,omitempty" binding:"max=255" bson:"ticket_system_name,omitempty"`
	// TicketData is used in assocticket action.
	TicketData map[string]string `json:"ticket_data,omitempty" bson:"ticket_data,omitempty"`
	// Snooze and Pbehavior
	Duration *datetime.DurationWithUnit `json:"duration" bson:"duration,omitempty"`
	// Pbehavior
	Name           string            `json:"name" bson:"name,omitempty" binding:"max=255"`
	Reason         string            `json:"reason" bson:"reason,omitempty"`
	Type           string            `json:"type" bson:"type,omitempty"`
	RRule          string            `json:"rrule" bson:"rrule,omitempty"`
	Tstart         *datetime.CpsTime `json:"tstart" bson:"tstart,omitempty"`
	Tstop          *datetime.CpsTime `json:"tstop" bson:"tstop,omitempty"`
	StartOnTrigger *bool             `json:"start_on_trigger" bson:"start_on_trigger,omitempty"`
}

// Matches returns true if alarm and entity match time condition and field patterns.
func (r *Rule) Matches(alarmWithEntity types.AlarmWithEntity, now datetime.CpsTime) (bool, error) {
	alarm := alarmWithEntity.Alarm
	entity := alarmWithEntity.Entity

	matched, err := match.Match(&alarmWithEntity.Entity, &alarmWithEntity.Alarm, r.EntityPattern, r.AlarmPattern)
	if err != nil {
		return false, err
	}

	return matched &&
		r.matchesByAlarmLastEventDate(alarm, now) &&
		r.matchesByAlarmLastUpdateDate(alarm, now) &&
		r.matchesByEntityLastEventDate(entity, now), nil
}

func (r *Rule) matchesByAlarmLastEventDate(alarm types.Alarm, now datetime.CpsTime) bool {
	before := r.Duration.SubFrom(now)

	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastEvent ||
		alarm.Value.LastEventDate.Before(before)
}

func (r *Rule) matchesByAlarmLastUpdateDate(alarm types.Alarm, now datetime.CpsTime) bool {
	before := r.Duration.SubFrom(now)

	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastUpdate ||
		alarm.Value.LastUpdateDate.Before(before)
}

func (r *Rule) matchesByEntityLastEventDate(entity types.Entity, now datetime.CpsTime) bool {
	if r.Type != RuleTypeEntity {
		return true
	}

	before := r.Duration.SubFrom(now)

	if entity.LastEventDate != nil {
		return entity.LastEventDate.Before(before)
	}

	if !entity.Created.IsZero() {
		return entity.Created.Before(before)
	}

	return true
}
