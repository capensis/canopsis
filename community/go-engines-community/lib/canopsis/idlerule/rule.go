/*
Package idlerule contains idle rule model and adapter.
*/
package idlerule

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"time"
)

const (
	RuleTypeAlarm                = "alarm"
	RuleTypeEntity               = "entity"
	RuleAlarmConditionLastEvent  = "last_event"  // alarms which don't receive events after "now - duration"
	RuleAlarmConditionLastUpdate = "last_update" // alarms which don't change after "now - duration"
)

type Operation struct {
	Type       string      `bson:"type" json:"type"`
	Parameters interface{} `bson:"parameters,omitempty" json:"parameters"`
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

func (o *Operation) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	type Alias Operation
	var tmp Alias

	err := bson.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	*o = Operation(tmp)
	o.Parameters = bsonDtoMap(o.Parameters)

	switch o.Type {
	case types.ActionTypeAssocTicket:
		var params types.OperationAssocTicketParameters
		err := mapstructure.Decode(o.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		o.Parameters = params
	case types.ActionTypeChangeState:
		var params types.OperationChangeStateParameters
		err := mapstructure.Decode(o.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		o.Parameters = params
	case types.ActionTypeSnooze:
		var params types.OperationSnoozeParameters
		err := mapstructure.Decode(o.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		o.Parameters = params
	case types.ActionTypePbehavior:
		var params types.ActionPBehaviorParameters
		err := mapstructure.Decode(o.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		o.Parameters = params
	default:
		var params types.OperationParameters
		err := mapstructure.Decode(o.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		o.Parameters = params
	}

	return nil
}

// Matches returns true if alarm and entity match time condition and field patterns.
func (r *Rule) Matches(alarm *types.Alarm, entity *types.Entity) bool {
	return r.matchesByAlarmLastEventDate(alarm) &&
		r.matchesByAlarmLastUpdateDate(alarm) &&
		r.matchesByEntityLastEventDate(entity) &&
		(alarm == nil || r.AlarmPatterns.Matches(alarm)) &&
		(entity == nil || r.EntityPatterns.Matches(entity))
}

func (r *Rule) matchesByAlarmLastEventDate(alarm *types.Alarm) bool {
	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastEvent ||
		time.Since(alarm.Value.LastEventDate.Time) >= r.Duration.Duration()
}

func (r *Rule) matchesByAlarmLastUpdateDate(alarm *types.Alarm) bool {
	return r.Type != RuleTypeAlarm || r.AlarmCondition != RuleAlarmConditionLastUpdate ||
		time.Since(alarm.Value.LastUpdateDate.Time) >= r.Duration.Duration()
}

func (r *Rule) matchesByEntityLastEventDate(entity *types.Entity) bool {
	if r.Type != RuleTypeEntity {
		return true
	}

	if entity.LastEventDate != nil {
		return time.Since(entity.LastEventDate.Time) >= r.Duration.Duration()
	}

	if !entity.Created.IsZero() {
		return time.Since(entity.Created.Time) >= r.Duration.Duration()
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
