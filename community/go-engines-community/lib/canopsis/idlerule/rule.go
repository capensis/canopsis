/*
Package idlerule contains idle rule model and adapter.
*/
package idlerule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
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
	ID             string                    `bson:"_id,omitempty" json:"_id"`
	Name           string                    `bson:"name" json:"name"`
	Description    string                    `bson:"description" json:"description"`
	Author         string                    `bson:"author" json:"author"`
	Enabled        bool                      `bson:"enabled" json:"enabled"`
	Type           string                    `bson:"type" json:"type"`
	Priority       int64                     `bson:"priority" json:"priority"`
	Duration       types.DurationWithUnit    `bson:"duration" json:"duration"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	// DisableDuringPeriods is an option that allows to disable the rule
	// when entity is in listed periods due pbehavior schedule.
	DisableDuringPeriods []string      `bson:"disable_during_periods" json:"disable_during_periods"`
	Created              types.CpsTime `bson:"created" json:"created"`
	Updated              types.CpsTime `bson:"updated" json:"updated"`
	// Only for Alarm rules
	AlarmPatterns  pattern.AlarmPatternList `bson:"alarm_patterns,omitempty" json:"alarm_patterns,omitempty"`
	AlarmCondition string                   `bson:"alarm_condition,omitempty" json:"alarm_condition,omitempty"`
	Operation      *Operation               `bson:"operation,omitempty" json:"operation,omitempty"`
}

type Parameters struct {
	Output string `json:"output" bson:"output,omitempty" binding:"max=255"`

	// ChangeState
	State *types.CpsNumber `json:"state" bson:"state,omitempty"`
	// AssocTicket
	Ticket string `json:"ticket" bson:"ticket,omitempty" binding:"max=255" `
	// Snooze and Pbehavior
	Duration *types.DurationWithUnit `json:"duration" bson:"duration,omitempty"`
	// Pbehavior
	Name           string         `json:"name" bson:"name,omitempty" binding:"max=255"`
	Reason         string         `json:"reason" bson:"reason,omitempty"`
	Type           string         `json:"type" bson:"type,omitempty"`
	RRule          string         `json:"rrule" bson:"rrule,omitempty"`
	Tstart         *types.CpsTime `json:"tstart" bson:"tstart,omitempty"`
	Tstop          *types.CpsTime `json:"tstop" bson:"tstop,omitempty"`
	StartOnTrigger *bool          `json:"start_on_trigger" bson:"start_on_trigger,omitempty"`
}

// Matches returns true if alarm and entity match time condition and field patterns.
func (r *Rule) Matches(alarm *types.Alarm, entity *types.Entity, now types.CpsTime) bool {
	return r.matchesByAlarmLastEventDate(alarm, now) &&
		r.matchesByAlarmLastUpdateDate(alarm, now) &&
		r.matchesByEntityLastEventDate(entity, now) &&
		(alarm == nil || r.AlarmPatterns.Matches(alarm)) &&
		(entity == nil || r.EntityPatterns.Matches(entity))
}

func (r *Rule) matchesByAlarmLastEventDate(alarm *types.Alarm, now types.CpsTime) bool {
	before := r.Duration.SubFrom(now)

	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastEvent ||
		alarm.Value.LastEventDate.Before(before)
}

func (r *Rule) matchesByAlarmLastUpdateDate(alarm *types.Alarm, now types.CpsTime) bool {
	before := r.Duration.SubFrom(now)

	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastUpdate ||
		alarm.Value.LastUpdateDate.Before(before)
}

func (r *Rule) matchesByEntityLastEventDate(entity *types.Entity, now types.CpsTime) bool {
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

func bsonDtoMap(i interface{}) interface{} {
	if b, ok := i.(bson.D); ok {
		m := b.Map()
		for k := range m {
			if b, ok := m[k].(bson.D); ok {
				m[k] = bsonDtoMap(b)
			}
		}

		return m
	}

	return i
}
