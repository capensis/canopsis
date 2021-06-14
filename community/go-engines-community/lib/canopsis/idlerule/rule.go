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

type Operation struct {
	Type       string      `bson:"type"`
	Parameters interface{} `bson:"parameters,omitempty"`
}

// Rule represents alarm modification condition and operation.
// Condition "type=last_event" applies alarms which don't receive events after "now - duration".
// Condition "type=last_update" applies alarms which don't change after "now - duration".
type Rule struct {
	ID             string                    `bson:"_id"`
	Type           RuleType                  `bson:"type"`
	Duration       types.CpsDuration         `bson:"duration"`
	AlarmPatterns  pattern.AlarmPatternList  `bson:"alarm_patterns"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns"`
	Operation      Operation                 `bson:"operation"`
	// DisableDuringPeriods is an option that allows to disable the rule
	// when entity is in listed periods due pbehavior schedule.
	DisableDuringPeriods []string `bson:"disable_during_periods" json:"disable_during_periods"`
}

func (r *Rule) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	type Alias Rule
	var tmp Alias

	err := bson.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	*r = Rule(tmp)
	r.Operation.Parameters = bsonDtoMap(r.Operation.Parameters)

	switch r.Operation.Type {
	case types.ActionTypeAssocTicket:
		var params types.OperationAssocTicketParameters
		err := mapstructure.Decode(r.Operation.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Operation.Parameters = params
	case types.ActionTypeChangeState:
		var params types.OperationChangeStateParameters
		err := mapstructure.Decode(r.Operation.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Operation.Parameters = params
	case types.ActionTypeSnooze:
		var params types.OperationSnoozeParameters
		err := mapstructure.Decode(r.Operation.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Operation.Parameters = params
	case types.ActionTypePbehavior:
		var params types.ActionPBehaviorParameters
		err := mapstructure.Decode(r.Operation.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Operation.Parameters = params
	default:
		var params types.OperationParameters
		err := mapstructure.Decode(r.Operation.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Operation.Parameters = params
	}

	return nil
}

// Valid values for type field.
type RuleType string

const (
	RuleTypeLastEvent  RuleType = "last_event"
	RuleTypeLastUpdate RuleType = "last_update"
)

// Matches returns true if alarm and entity match time condition and field patterns.
func (r *Rule) Matches(alarm *types.Alarm, entity *types.Entity) bool {
	return r.matchesByLastEventDate(alarm) &&
		r.matchesByLastUpdateDate(alarm) &&
		r.AlarmPatterns.Matches(alarm) &&
		r.EntityPatterns.Matches(entity)
}

func (r *Rule) matchesByLastEventDate(alarm *types.Alarm) bool {
	if r.Type != RuleTypeLastEvent {
		return true
	}

	now := time.Now()
	t := alarm.Value.LastEventDate.Add(r.Duration.Duration())

	return t.Before(now)
}

func (r *Rule) matchesByLastUpdateDate(alarm *types.Alarm) bool {
	if r.Type != RuleTypeLastUpdate {
		return true
	}

	now := time.Now()
	t := alarm.Value.LastUpdateDate.Add(r.Duration.Duration())

	return t.Before(now)
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
