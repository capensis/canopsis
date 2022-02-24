// Package pattern provides functionality for filtering and matching models.
package pattern

import (
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var ErrUnsupportedField = errors.New("unsupported field")
var ErrUnsupportedConditionType = errors.New("unsupported condition type")
var ErrWrongConditionValue = errors.New("wrong condition value")

const (
	ConditionEqual        = "eq"
	ConditionNotEqual     = "neq"
	ConditionGT           = "gt"
	ConditionLT           = "lt"
	ConditionRegexp       = "regexp"
	ConditionHasEvery     = "has_every"
	ConditionHasOneOf     = "has_one_of"
	ConditionHasNot       = "has_not"
	ConditionIsEmpty      = "is_empty"
	ConditionExist        = "exist"
	ConditionTimeRelative = "relative_time"
	ConditionTimeAbsolute = "absolute_time"
)

const (
	FieldTypeString      = "string"
	FieldTypeInt         = "int"
	FieldTypeBool        = "bool"
	FieldTypeStringArray = "string_array"
)

// FieldCondition represents a condition for a specific field.
type FieldCondition struct {
	Field string `json:"field" bson:"field"`
	// FieldType is only defined for custom fields, ex: infos.
	FieldType string    `json:"field_type,omitempty" bson:"field_type,omitempty"`
	Condition Condition `json:"cond" bson:"cond"`
}

// Condition represents an expression to decide if a value fits.
type Condition struct {
	Type  string      `json:"type" bson:"type"`
	Value interface{} `json:"value" bson:"value"`
}

// RegexMatches is a type that contains the values of the sub-expressions of a
// regular expression.
type RegexMatches map[string]string

func (c Condition) MatchString(value string) (bool, RegexMatches, error) {
	conditionValue, ok := c.Value.(string)
	if !ok {
		return false, nil, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionEqual:
		return value == conditionValue, nil, nil
	case ConditionNotEqual:
		return value != conditionValue, nil, nil
	case ConditionRegexp:
		r, err := utils.NewRegexExpression(conditionValue)
		if err != nil {
			return false, nil, err
		}

		regexMatches := utils.FindStringSubmatchMapWithRegexExpression(r, value)

		return regexMatches != nil, regexMatches, nil
	}

	return false, nil, ErrUnsupportedConditionType
}

func (c Condition) MatchInt(value int) (bool, error) {
	conditionValue, ok := c.Value.(int)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionEqual:
		return value == conditionValue, nil
	case ConditionNotEqual:
		return value != conditionValue, nil
	case ConditionGT:
		return value > conditionValue, nil
	case ConditionLT:
		return value < conditionValue, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) MatchBool(value bool) (bool, error) {
	conditionValue, ok := c.Value.(bool)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionEqual:
		return value == conditionValue, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) MatchRef(value interface{}) (bool, error) {
	conditionValue, ok := c.Value.(bool)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionExist:
		return conditionValue == (value != nil), nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) MatchStringArray(value []string) (bool, error) {
	if c.Type == ConditionIsEmpty {
		conditionValue, ok := c.Value.(bool)
		if !ok {
			return false, ErrWrongConditionValue
		}

		return conditionValue == (len(value) == 0), nil
	}

	conditionValue, ok := c.Value.([]string)
	if !ok {
		return false, ErrWrongConditionValue
	}

	valueMap := make(map[string]bool)
	for _, v := range value {
		valueMap[v] = true
	}

	switch c.Type {
	case ConditionEqual:
		for _, v := range conditionValue {
			if _, ok := valueMap[v]; !ok {
				return false, nil
			}

			delete(valueMap, v)
		}

		return len(valueMap) == 0, nil
	case ConditionHasEvery:
		for _, v := range conditionValue {
			_, exists := valueMap[v]
			if !exists {
				return false, nil
			}
		}

		return true, nil
	case ConditionHasOneOf:
		for _, v := range conditionValue {
			_, exists := valueMap[v]
			if exists {
				return true, nil
			}
		}

		return false, nil
	case ConditionHasNot:
		for _, v := range conditionValue {
			_, exists := valueMap[v]
			if exists {
				return false, nil
			}
		}

		return true, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) MatchTime(value time.Time) (bool, error) {
	switch c.Type {
	case ConditionTimeRelative:
		conditionValue, ok := c.Value.(int)
		if !ok {
			return false, ErrWrongConditionValue
		}

		return value.After(time.Now().Add(time.Duration(-conditionValue) * time.Second)), nil
	case ConditionTimeAbsolute:
		conditionValue, ok := c.Value.(map[string]int64)
		if !ok {
			return false, ErrWrongConditionValue
		}

		from, ok := conditionValue["from"]
		if !ok {
			return false, errors.New("condition value expected 'from' key")
		}

		to, ok := conditionValue["to"]
		if !ok {
			return false, errors.New("condition value expected 'to' key")
		}

		return value.After(time.Unix(from, 0)) && value.Before(time.Unix(to, 0)), nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) MatchDuration(value int64) (bool, error) {
	conditionValue, ok := c.Value.(map[string]interface{})
	if !ok {
		return false, ErrWrongConditionValue
	}

	rawVal, ok := conditionValue["value"]
	if !ok {
		return false, errors.New("condition value expected 'value' key")
	}

	val, ok := rawVal.(int)
	if !ok {
		return false, errors.New("value should be an int64")
	}

	rawUnit, ok := conditionValue["unit"]
	if !ok {
		return false, errors.New("condition value expected 'unit' key")
	}

	unit, ok := rawUnit.(string)
	if !ok {
		return false, errors.New("unit should be a string")
	}

	d := types.DurationWithUnit{
		Value: int64(val),
		Unit:  unit,
	}

	d, err := types.DurationWithUnit{
		Value: int64(val),
		Unit:  unit,
	}.To("s")

	if err != nil {
		return false, err
	}

	switch c.Type {
	case ConditionGT:
		return value > d.Value, nil
	case ConditionLT:
		return value < d.Value, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c Condition) ToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionEqual:
		return bson.M{f: bson.M{"$eq": c.Value}}, nil
	case ConditionNotEqual:
		return bson.M{f: bson.M{"$ne": c.Value}}, nil
	case ConditionGT:
		return bson.M{f: bson.M{"$gt": c.Value}}, nil
	case ConditionLT:
		return bson.M{f: bson.M{"$lt": c.Value}}, nil
	case ConditionRegexp:
		return bson.M{f: bson.M{"$regex": c.Value}}, nil
	case ConditionHasEvery:
		return bson.M{f: bson.M{"$all": c.Value}}, nil
	case ConditionHasOneOf:
		return bson.M{f: bson.M{"$elemMatch": bson.M{"$in": c.Value}}}, nil
	case ConditionHasNot:
		return bson.M{f: bson.M{"$elemMatch": bson.M{"$nin": c.Value}}}, nil
	case ConditionIsEmpty:
		conditionValue, ok := c.Value.(bool)
		if !ok {
			return nil, ErrWrongConditionValue
		}

		if conditionValue {
			return bson.M{f: bson.M{"$in": bson.A{nil, bson.A{}}}}, nil
		}

		return bson.M{f: bson.M{"$exists": true, "$type": "array", "$ne": bson.A{}}}, nil
	case ConditionExist:
		conditionValue, ok := c.Value.(bool)
		if !ok {
			return nil, ErrWrongConditionValue
		}

		if conditionValue {
			return bson.M{f: bson.M{"$exists": true, "$ne": nil}}, nil
		}

		return bson.M{"$or": []bson.M{{f: bson.M{"$exists": false}}, {f: bson.M{"$eq": nil}}}}, nil
	case ConditionTimeRelative:
		conditionValue, ok := c.Value.(int)
		if !ok {
			return nil, ErrWrongConditionValue
		}

		t := types.CpsTime{Time: time.Now().Add(time.Duration(-conditionValue) * time.Second)}

		return bson.M{f: bson.M{"$gt": t}}, nil
	case ConditionTimeAbsolute:
		conditionValue, ok := c.Value.(map[string]int64)
		if !ok {
			return nil, ErrWrongConditionValue
		}

		from, ok := conditionValue["from"]
		if !ok {
			return nil, errors.New("condition value expected 'from' key")
		}

		to, ok := conditionValue["to"]
		if !ok {
			return nil, errors.New("condition value expected 'to' key")
		}

		ft := types.NewCpsTime(from)
		tt := types.NewCpsTime(to)

		return bson.M{f: bson.M{"$gt": ft, "$lt": tt}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}
