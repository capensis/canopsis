package pattern

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"time"
)

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
	ConditionTimeRelative = "relative_time"
	ConditionTimeAbsolute = "absolute_time"
)

type Condition struct {
	Type  string      `json:"cond" bson:"cond"`
	Value interface{} `json:"value" bson:"value"`
}

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
	case ConditionIsEmpty:
		return conditionValue == (value == nil), nil
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
