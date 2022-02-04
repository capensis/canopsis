package pattern

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

var ErrUnsupportedConditionType = errors.New("unsupported condition type")
var ErrWrongConditionValue = errors.New("wrong condition value")

const (
	ConditionEqual    = "eq"
	ConditionNotEqual = "neq"
	ConditionGT       = "gt"
	ConditionLT       = "lt"
	ConditionRegexp   = "regexp"
	ConditionHasEvery = "has_every"
	ConditionHasOneOf = "has_one_of"
	ConditionHasNot   = "has_not"
	ConditionIsEmpty  = "is_empty"
	ConditionTime     = "time"
)

type Condition struct {
	Type  string      `json:"cond" bson:"cond"`
	Value interface{} `json:"value" bson:"value"`
}

func MatchString(conf Condition, value string) (bool, RegexMatches, error) {
	conditionValue, ok := conf.Value.(string)
	if !ok {
		return false, nil, ErrWrongConditionValue
	}

	switch conf.Type {
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

func MatchInt(conf Condition, value int) (bool, error) {
	conditionValue, ok := conf.Value.(int)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch conf.Type {
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

func MatchBool(conf Condition, value bool) (bool, error) {
	conditionValue, ok := conf.Value.(bool)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch conf.Type {
	case ConditionEqual:
		return value == conditionValue, nil
	}

	return false, ErrUnsupportedConditionType
}

func MatchRef(conf Condition, value interface{}) (bool, error) {
	conditionValue, ok := conf.Value.(bool)
	if !ok {
		return false, ErrWrongConditionValue
	}

	switch conf.Type {
	case ConditionIsEmpty:
		return conditionValue == (value == nil), nil
	}

	return false, ErrUnsupportedConditionType
}

func MatchStringArray(conf Condition, value []string) (bool, error) {
	if conf.Type == ConditionIsEmpty {
		conditionValue, ok := conf.Value.(bool)
		if !ok {
			return false, ErrWrongConditionValue
		}

		return conditionValue == (len(value) == 0), nil
	}

	conditionValue, ok := conf.Value.([]string)
	if !ok {
		return false, ErrWrongConditionValue
	}

	valueMap := make(map[string]bool)
	for _, v := range value {
		valueMap[v] = true
	}

	switch conf.Type {
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
