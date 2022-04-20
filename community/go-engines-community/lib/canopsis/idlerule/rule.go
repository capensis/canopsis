/*
Package idlerule contains idle rule model and adapter.
*/
package idlerule

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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
	Created              types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated              types.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
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
