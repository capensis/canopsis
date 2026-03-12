// Package pattern provides functionality for filtering and matching models.
package pattern

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrUnsupportedField = errors.New("unsupported field")
var ErrUnsupportedConditionType = errors.New("unsupported condition type")
var ErrWrongConditionValue = errors.New("wrong condition value")
var ErrOverwriteMatchedRegexp = errors.New("attempt to overwrite matched regex with groups")

const (
	ConditionEqual        = "eq"
	ConditionNotEqual     = "neq"
	ConditionGT           = "gt"
	ConditionLT           = "lt"
	ConditionRegexp       = "regexp"
	ConditionContain      = "contain"
	ConditionNotContain   = "not_contain"
	ConditionBeginWith    = "begin_with"
	ConditionNotBeginWith = "not_begin_with"
	ConditionEndWith      = "end_with"
	ConditionNotEndWith   = "not_end_with"
	ConditionHasEvery     = "has_every"
	ConditionHasOneOf     = "has_one_of"
	ConditionHasNot       = "has_not"
	ConditionIsEmpty      = "is_empty"
	ConditionExist        = "exist"
	ConditionTimeRelative = "relative_time"
	ConditionTimeAbsolute = "absolute_time"
	ConditionIsOneOf      = "is_one_of"
	ConditionIsNotOneOf   = "is_not_one_of"
	ConditionHasLabels    = "has_labels"
	ConditionHasNotLabels = "has_not_labels"
)

const (
	FieldTypeString      = "string"
	FieldTypeInt         = "int"
	FieldTypeBool        = "bool"
	FieldTypeStringArray = "string_array"
	FieldTypeTimestamp   = "timestamp"
)

const tagLabelSeparator = ":"

// FieldCondition represents a condition for a specific field.
type FieldCondition struct {
	Field string `json:"field" bson:"field"`
	// FieldType is only defined for custom fields, ex: infos.
	FieldType string    `json:"field_type,omitempty" bson:"field_type,omitempty" mapstructure:"field_type,omitempty"`
	Condition Condition `json:"cond" bson:"cond" mapstructure:"cond"`
	// Alias can only be defined for entity infos.
	Alias string `json:"alias,omitempty" bson:"alias,omitempty"`
}

// Condition represents an expression to decide if a value fits.
type Condition struct {
	Type  string `json:"type" bson:"type"`
	Value any    `json:"value" bson:"value" swaggertype:"object"`

	valueStr              *string
	valueRegexp           utils.RegexExpression
	valueInt              *int64
	valueBool             *bool
	valueStrArray         []string
	valueTimeIntervalFrom *int64
	valueTimeIntervalTo   *int64
	// valueDuration is a single duration for >, <  and "after relative_time" conditions
	// valueDurationTo is second duration for "before relative_time" condition
	valueDuration   *int64
	valueDurationTo *int64
}

// RegexMatches is a type that contains the values of the sub-expressions of a
// regular expression.
type RegexMatches map[string]string

func NewStringCondition(t, s string) Condition {
	return Condition{
		Type:     t,
		Value:    s,
		valueStr: &s,
	}
}

func NewRegexpCondition(t, s string) (Condition, error) {
	r, err := utils.NewRegexExpression(s)
	if err != nil {
		return Condition{}, err
	}

	return Condition{
		Type:        t,
		Value:       s,
		valueRegexp: r,
	}, nil
}

func NewIntCondition(t string, i int64) Condition {
	return Condition{
		Type:     t,
		Value:    i,
		valueInt: &i,
	}
}

func NewBoolCondition(t string, b bool) Condition {
	return Condition{
		Type:      t,
		Value:     b,
		valueBool: &b,
	}
}

func NewStringArrayCondition(t string, a []string) Condition {
	return Condition{
		Type:          t,
		Value:         a,
		valueStrArray: a,
	}
}

func NewTimeIntervalCondition(t string, from, to int64) Condition {
	return Condition{
		Type: t,
		Value: map[string]int64{
			"from": from,
			"to":   to,
		},
		valueTimeIntervalFrom: &from,
		valueTimeIntervalTo:   &to,
	}
}

func NewDurationCondition(t string, d ...datetime.DurationWithUnit) (Condition, error) {
	if len(d) == 0 || len(d) > 2 {
		return Condition{}, fmt.Errorf("wrong number of duration arguments: %d", len(d))
	}

	durationFrom, err := d[0].To(datetime.DurationUnitSecond)
	if err != nil {
		return Condition{}, err
	}

	c := Condition{
		Type:          t,
		Value:         d[0],
		valueDuration: &durationFrom.Value,
	}

	if len(d) == 1 {
		return c, nil
	}

	durationTo, err := d[1].To(datetime.DurationUnitSecond)
	if err != nil {
		return Condition{}, err
	}

	if durationTo.Value == 0 {
		return c, nil
	}

	c.valueDurationTo = &durationTo.Value
	if durationFrom.Value == 0 {
		c.Value = map[string]*datetime.DurationWithUnit{"to": &d[1]}
		c.valueDuration = nil
	} else {
		if durationFrom.Value < durationTo.Value {
			return Condition{}, fmt.Errorf("from (%d) < to (%d)", durationFrom.Value, durationTo.Value)
		}

		c.Value = map[string]*datetime.DurationWithUnit{"from": &d[0], "to": &d[1]}
	}

	return c, nil
}

func (c *Condition) GetValueStr() *string {
	if c == nil {
		return nil
	}

	return c.valueStr
}

func (c *Condition) GetValueStrArray() []string {
	if c == nil {
		return nil
	}

	return c.valueStrArray
}

func (c *Condition) MatchString(value string) (bool, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueStr == nil {
			return false, ErrWrongConditionValue
		}

		return value == *c.valueStr, nil
	case ConditionNotEqual:
		if c.valueStr == nil {
			return false, ErrWrongConditionValue
		}

		return value != *c.valueStr, nil
	case ConditionIsOneOf:
		if len(c.valueStrArray) == 0 {
			return false, ErrWrongConditionValue
		}
		return slices.Contains(c.valueStrArray, value), nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return false, ErrWrongConditionValue
		}

		return !slices.Contains(c.valueStrArray, value), nil
	case ConditionRegexp,
		ConditionContain,
		ConditionNotContain,
		ConditionBeginWith,
		ConditionNotBeginWith,
		ConditionEndWith,
		ConditionNotEndWith:
		if c.valueRegexp == nil {
			return false, ErrWrongConditionValue
		}

		return utils.MatchWithRegexExpression(c.valueRegexp, value), nil
	case ConditionExist:
		if c.valueBool == nil {
			return false, ErrWrongConditionValue
		}

		return *c.valueBool == (value != ""), nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchStringWithRegexpMatches(value string) (bool, RegexMatches, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueStr == nil {
			return false, nil, ErrWrongConditionValue
		}

		return value == *c.valueStr, nil, nil
	case ConditionNotEqual:
		if c.valueStr == nil {
			return false, nil, ErrWrongConditionValue
		}

		return value != *c.valueStr, nil, nil
	case ConditionIsOneOf:
		if len(c.valueStrArray) == 0 {
			return false, nil, ErrWrongConditionValue
		}

		return slices.Contains(c.valueStrArray, value), nil, nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return false, nil, ErrWrongConditionValue
		}

		return !slices.Contains(c.valueStrArray, value), nil, nil
	case ConditionRegexp,
		ConditionContain,
		ConditionNotContain,
		ConditionBeginWith,
		ConditionNotBeginWith,
		ConditionEndWith,
		ConditionNotEndWith:
		if c.valueRegexp == nil {
			return false, nil, ErrWrongConditionValue
		}

		regexMatches := utils.FindStringSubmatchMapWithRegexExpression(c.valueRegexp, value)

		return regexMatches != nil, regexMatches, nil
	case ConditionExist:
		if c.valueBool == nil {
			return false, nil, ErrWrongConditionValue
		}

		return *c.valueBool == (value != ""), nil, nil
	}

	return false, nil, ErrUnsupportedConditionType
}

func (c *Condition) MatchInt(value int64) (bool, error) {
	if c.valueInt == nil {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionEqual:
		return value == *c.valueInt, nil
	case ConditionNotEqual:
		return value != *c.valueInt, nil
	case ConditionGT:
		return value > *c.valueInt, nil
	case ConditionLT:
		return value < *c.valueInt, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchBool(value bool) (bool, error) {
	if c.valueBool == nil {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionEqual:
		return value == *c.valueBool, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchRef(value any) (bool, error) {
	if c.valueBool == nil {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionExist:
		return *c.valueBool == (value != nil), nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchStringArray(value []string) (bool, error) {
	if c.Type == ConditionIsEmpty {
		if c.valueBool == nil {
			return false, ErrWrongConditionValue
		}

		return *c.valueBool == (len(value) == 0), nil
	}

	if len(c.valueStrArray) == 0 {
		return false, ErrWrongConditionValue
	}

	valueMap := make(map[string]bool)
	for _, v := range value {
		valueMap[v] = true
	}

	switch c.Type {
	case ConditionHasEvery:
		for _, v := range c.valueStrArray {
			_, exists := valueMap[v]
			if !exists {
				return false, nil
			}
		}

		return true, nil
	case ConditionHasOneOf:
		for _, v := range c.valueStrArray {
			_, exists := valueMap[v]
			if exists {
				return true, nil
			}
		}

		return false, nil
	case ConditionHasNot:
		for _, v := range c.valueStrArray {
			_, exists := valueMap[v]
			if exists {
				return false, nil
			}
		}

		return true, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchTags(tags []string) (bool, error) {
	switch c.Type {
	case ConditionIsEmpty, ConditionHasEvery, ConditionHasOneOf, ConditionHasNot:
		return c.MatchStringArray(tags)
	case ConditionHasLabels, ConditionHasNotLabels:
		hasLabels := c.Type == ConditionHasLabels

		if len(c.valueStrArray) == 0 {
			return false, ErrWrongConditionValue
		}

		labelsMap := make(map[string]bool)
		for _, v := range tags {
			label, _, _ := strings.Cut(v, tagLabelSeparator)

			labelsMap[label] = true
		}

		for _, v := range c.valueStrArray {
			if labelsMap[v] {
				return hasLabels, nil
			}
		}

		return !hasLabels, nil
	default:
		return false, ErrUnsupportedConditionType
	}
}

func (c *Condition) MatchTime(value, now time.Time) (bool, error) {
	switch c.Type {
	case ConditionTimeRelative:
		if c.valueDuration == nil && c.valueDurationTo == nil {
			return false, ErrWrongConditionValue
		}
		var t1, t2 time.Time
		if c.valueDurationTo != nil {
			t2 = now.Add(time.Duration(-*c.valueDurationTo) * time.Second)
		}
		if c.valueDuration != nil {
			t1 = now.Add(time.Duration(-*c.valueDuration) * time.Second)
			if c.valueDurationTo != nil {
				// when both durations are set, t1 must be not after t2
				if t1.After(t2) {
					return false, fmt.Errorf("duration from (%d) < to (%d)", *c.valueDuration, *c.valueDurationTo)
				}
				return value.After(t1) && value.Before(t2), nil
			}
		}

		return (!t1.IsZero() && value.After(t1)) || value.Before(t2), nil
	case ConditionTimeAbsolute:
		if c.valueTimeIntervalFrom == nil || c.valueTimeIntervalTo == nil {
			return false, ErrWrongConditionValue
		}

		return value.After(time.Unix(*c.valueTimeIntervalFrom, 0)) && value.Before(time.Unix(*c.valueTimeIntervalTo, 0)), nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) MatchDuration(value int64) (bool, error) {
	if c.valueDuration == nil {
		return false, ErrWrongConditionValue
	}

	switch c.Type {
	case ConditionGT:
		return value > *c.valueDuration, nil
	case ConditionLT:
		return value < *c.valueDuration, nil
	}

	return false, ErrUnsupportedConditionType
}

func (c *Condition) StringToMongoQuery(f string, checkExists bool) (bson.M, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueStr == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$eq": *c.valueStr}}, nil
	case ConditionNotEqual:
		if c.valueStr == nil {
			return nil, ErrWrongConditionValue
		}

		if checkExists {
			return bson.M{f: bson.M{"$nin": bson.A{*c.valueStr, nil}}}, nil
		}

		return bson.M{f: bson.M{"$ne": *c.valueStr}}, nil
	case ConditionIsOneOf:
		if len(c.valueStrArray) == 0 {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$in": c.valueStrArray}}, nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return nil, ErrWrongConditionValue
		}

		if checkExists {
			return bson.M{f: bson.M{
				"$nin": c.valueStrArray,
				"$ne":  nil,
			}}, nil
		}

		return bson.M{f: bson.M{"$nin": c.valueStrArray}}, nil
	case ConditionRegexp,
		ConditionContain,
		ConditionNotContain,
		ConditionBeginWith,
		ConditionNotBeginWith,
		ConditionEndWith,
		ConditionNotEndWith:
		if c.valueRegexp == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$regex": c.valueRegexp.String()}}, nil
	case ConditionExist:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			return bson.M{f: bson.M{
				"$exists": true,
				"$nin":    bson.A{nil, ""},
			}}, nil
		}

		return bson.M{"$or": []bson.M{
			{f: bson.M{"$exists": false}},
			{f: bson.M{"$eq": nil}},
			{f: bson.M{"$eq": ""}},
		}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) StringInArrayToMongoQuery(arrayField, arrayItemField string, checkExists bool) (bson.M, error) {
	switch c.Type {
	case ConditionExist:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			break
		}

		return bson.M{"$or": []bson.M{
			{arrayField: bson.M{"$in": bson.A{nil, bson.A{}}}},
			{arrayField: bson.M{"$not": bson.M{"$elemMatch": bson.M{
				arrayItemField: bson.M{
					"$exists": true,
					"$nin":    bson.A{nil, ""},
				},
			}}}},
		}}, nil
	}

	cond, err := c.StringToMongoQuery(arrayItemField, checkExists)
	if err != nil {
		return nil, err
	}

	return bson.M{arrayField: bson.M{"$elemMatch": cond}}, nil
}

func (c *Condition) IntToMongoQuery(f string, checkExists bool) (bson.M, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueInt == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$eq": *c.valueInt}}, nil
	case ConditionNotEqual:
		if c.valueInt == nil {
			return nil, ErrWrongConditionValue
		}

		if checkExists {
			return bson.M{f: bson.M{"$nin": bson.A{*c.valueInt, nil}}}, nil
		}

		return bson.M{f: bson.M{"$ne": *c.valueInt}}, nil
	case ConditionGT:
		if c.valueInt == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$gt": *c.valueInt}}, nil
	case ConditionLT:
		if c.valueInt == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$lt": *c.valueInt}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) IntInArrayToMongoQuery(arrayField, arrayItemField string, checkExists bool) (bson.M, error) {
	cond, err := c.IntToMongoQuery(arrayItemField, checkExists)
	if err != nil {
		return nil, err
	}

	return bson.M{arrayField: bson.M{"$elemMatch": cond}}, nil
}

func (c *Condition) BoolToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$eq": *c.valueBool}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) BoolInArrayToMongoQuery(arrayField, arrayItemField string) (bson.M, error) {
	cond, err := c.BoolToMongoQuery(arrayItemField)
	if err != nil {
		return nil, err
	}

	return bson.M{arrayField: bson.M{"$elemMatch": cond}}, nil
}

func (c *Condition) RefToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionExist:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			return bson.M{f: bson.M{
				"$exists": true,
				"$ne":     nil,
			}}, nil
		}

		return bson.M{"$or": []bson.M{
			{f: bson.M{"$exists": false}},
			{f: bson.M{"$eq": nil}},
		}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) RefInArrayToMongoQuery(arrayField, arrayItemField string) (bson.M, error) {
	switch c.Type {
	case ConditionExist:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			break
		}

		return bson.M{"$or": []bson.M{
			{arrayField: bson.M{"$in": bson.A{nil, bson.A{}}}},
			{arrayField: bson.M{"$not": bson.M{"$elemMatch": bson.M{
				arrayItemField: bson.M{
					"$exists": true,
					"$ne":     nil,
				},
			}}}},
		}}, nil
	}

	cond, err := c.RefToMongoQuery(arrayItemField)
	if err != nil {
		return nil, err
	}

	return bson.M{arrayField: bson.M{"$elemMatch": cond}}, nil
}

func (c *Condition) StringArrayToMongoQuery(f string, checkExists bool) (bson.M, error) {
	switch c.Type {
	case ConditionIsEmpty:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			return bson.M{f: bson.M{"$in": bson.A{nil, bson.A{}}}}, nil
		}

		return bson.M{f: bson.M{"$exists": true, "$type": "array", "$ne": bson.A{}}}, nil
	case ConditionHasEvery:
		if len(c.valueStrArray) == 0 {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$all": c.valueStrArray}}, nil
	case ConditionHasOneOf:
		if len(c.valueStrArray) == 0 {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$elemMatch": bson.M{"$in": c.valueStrArray}}}, nil
	case ConditionHasNot:
		if len(c.valueStrArray) == 0 {
			return nil, ErrWrongConditionValue
		}

		if checkExists {
			return bson.M{f: bson.M{
				"$nin": c.valueStrArray,
				"$ne":  nil,
			}}, nil
		}

		return bson.M{f: bson.M{"$nin": c.valueStrArray}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) TagsToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionIsEmpty, ConditionHasEvery, ConditionHasOneOf, ConditionHasNot:
		return c.StringArrayToMongoQuery(f, false)
	case ConditionHasLabels, ConditionHasNotLabels:
		expr := "^(" + strings.Join(c.valueStrArray, "|") + ")(?:\\:|$)"

		r, err := utils.NewRegexExpression(expr)
		if err != nil {
			return nil, fmt.Errorf("invalid regexp expression - %q in tags alarm pattern: %w", expr, err)
		}

		if c.Type == ConditionHasLabels {
			return bson.M{f: bson.M{"$regex": r.String()}}, nil
		} else {
			return bson.M{f: bson.M{"$not": bson.M{"$regex": r.String()}}}, nil
		}
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) StringArrayInArrayToMongoQuery(arrayField, arrayItemField string, checkExists bool) (bson.M, error) {
	switch c.Type {
	case ConditionIsEmpty:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if !*c.valueBool {
			break
		}

		return bson.M{"$or": []bson.M{
			{arrayField: bson.M{"$in": bson.A{nil, bson.A{}}}},
			{arrayField: bson.M{"$not": bson.M{"$elemMatch": bson.M{
				arrayItemField: bson.M{
					"$exists": true,
					"$type":   "array",
					"$ne":     bson.A{},
				},
			}}}},
		}}, nil
	}

	cond, err := c.StringArrayToMongoQuery(arrayItemField, checkExists)
	if err != nil {
		return nil, err
	}

	return bson.M{arrayField: bson.M{"$elemMatch": cond}}, nil
}

func (c *Condition) TimeToMongoQuery(f string, now time.Time) (bson.M, error) {
	switch c.Type {
	case ConditionTimeRelative:
		if c.valueDuration == nil && c.valueDurationTo == nil {
			return nil, ErrWrongConditionValue
		}
		var t1, t2 datetime.CpsTime
		if c.valueDurationTo != nil {
			t2 = datetime.NewCpsTime(now.Add(time.Duration(-*c.valueDurationTo) * time.Second).Unix())
		}
		if c.valueDuration != nil {
			t1 = datetime.NewCpsTime(now.Add(time.Duration(-*c.valueDuration) * time.Second).Unix())
			if c.valueDurationTo != nil {
				// when both durations are set, t1 must be not after t2
				if t1.Time.After(t2.Time) {
					return nil, fmt.Errorf("duration from (%d) < to (%d)", *c.valueDuration, *c.valueDurationTo)
				}
				return bson.M{f: bson.M{"$gt": t1, "$lt": t2}}, nil
			}
			return bson.M{f: bson.M{"$gt": t1}}, nil
		}
		return bson.M{f: bson.M{"$lt": t2}}, nil
	case ConditionTimeAbsolute:
		if c.valueTimeIntervalFrom == nil || c.valueTimeIntervalTo == nil {
			return nil, ErrWrongConditionValue
		}

		ft := datetime.NewCpsTime(*c.valueTimeIntervalFrom)
		tt := datetime.NewCpsTime(*c.valueTimeIntervalTo)

		return bson.M{f: bson.M{"$gt": ft, "$lt": tt}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) DurationToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionGT:
		if c.valueDuration == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$gt": c.valueDuration}}, nil
	case ConditionLT:
		if c.valueDuration == nil {
			return nil, ErrWrongConditionValue
		}

		return bson.M{f: bson.M{"$lt": c.valueDuration}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

func (c *Condition) StringToSql(f string) (string, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueStr == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s = %s", f, sqlQuoteString(*c.valueStr)), nil
	case ConditionNotEqual:
		if c.valueStr == nil {
			return "", ErrWrongConditionValue
		}

		// "IS NULL" is mandatory
		return fmt.Sprintf("(%[1]s IS NULL OR %[1]s != %s)", f, sqlQuoteString(*c.valueStr)), nil
	case ConditionIsOneOf:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("%s = ANY (ARRAY [%s])", f, strings.Join(values, ",")), nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		// "IS NULL" is mandatory
		return fmt.Sprintf("(%[1]s IS NULL OR NOT (%[1]s = ANY (ARRAY [%s])))", f, strings.Join(values, ",")), nil
	case ConditionRegexp,
		ConditionContain,
		ConditionNotContain,
		ConditionBeginWith,
		ConditionNotBeginWith,
		ConditionEndWith,
		ConditionNotEndWith:
		if c.valueRegexp == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s ~ %s", f, sqlQuoteString(c.valueRegexp.String())), nil
	case ConditionExist:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("(%[1]s IS NOT NULL AND %[1]s != '')", f), nil
		}

		return fmt.Sprintf("(%[1]s IS NULL OR %[1]s = '')", f), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) IntToSql(f string) (string, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s = %d", f, *c.valueInt), nil
	case ConditionNotEqual:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		// "IS NULL" is mandatory
		return fmt.Sprintf("(%[1]s IS NULL OR %[1]s != %d)", f, *c.valueInt), nil
	case ConditionGT:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s > %d", f, *c.valueInt), nil
	case ConditionLT:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s < %d", f, *c.valueInt), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) BoolToSql(f string) (string, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s = %t", f, *c.valueBool), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) DurationToSql(f string) (string, error) {
	switch c.Type {
	case ConditionGT:
		if c.valueDuration == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s > %d", f, *c.valueDuration), nil
	case ConditionLT:
		if c.valueDuration == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("%s < %d", f, *c.valueDuration), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) StringToSqlJson(field, key string) (string, error) {
	checkType := fmt.Sprintf("jsonb_typeof(%s->%s) = 'string'", field, sqlQuoteString(key))
	operand := fmt.Sprintf("%s->>%s", field, sqlQuoteString(key))
	cast := fmt.Sprintf("%s AND %s", checkType, operand)

	switch c.Type {
	case ConditionEqual:
		if c.valueStr == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s = %s)", cast, sqlQuoteString(*c.valueStr)), nil
	case ConditionNotEqual:
		if c.valueStr == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s != %s)", cast, sqlQuoteString(*c.valueStr)), nil
	case ConditionIsOneOf:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("(%s = ANY (ARRAY [%s]))", cast, strings.Join(values, ",")), nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("(%s AND NOT (%s = ANY (ARRAY [%s])))", checkType, operand, strings.Join(values, ",")), nil
	case ConditionRegexp,
		ConditionContain,
		ConditionNotContain,
		ConditionBeginWith,
		ConditionNotBeginWith,
		ConditionEndWith,
		ConditionNotEndWith:
		if c.valueRegexp == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s ~ %s)", cast, sqlQuoteString(c.valueRegexp.String())), nil
	case ConditionExist:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("(%s ? %s AND %s != '')", field, sqlQuoteString(key), cast), nil
		}

		return fmt.Sprintf("(NOT (%s ? %s) OR %s = '')", field, sqlQuoteString(key), cast), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) IntToSqlJson(field, key string) (string, error) {
	// "CASE" is mandatory to cast json value because Postgres "SELECT" with following condition returns an error
	// if there is a row with field -> key of another (not numeric) type:
	// jsonb_typeof(field -> key) = 'number' AND (field -> key)::numeric = 2
	checkType := fmt.Sprintf("jsonb_typeof(%s->%s) = 'number'", field, sqlQuoteString(key))
	operand := fmt.Sprintf("(CASE WHEN %s THEN (%s->%s)::numeric END)", checkType, field, sqlQuoteString(key))
	cast := fmt.Sprintf("%s AND %s", checkType, operand)

	switch c.Type {
	case ConditionEqual:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s = %d)", cast, *c.valueInt), nil
	case ConditionNotEqual:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s != %d)", cast, *c.valueInt), nil
	case ConditionGT:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s > %d)", cast, *c.valueInt), nil
	case ConditionLT:
		if c.valueInt == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s < %d)", cast, *c.valueInt), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) BoolToSqlJson(field, key string) (string, error) {
	checkType := fmt.Sprintf("jsonb_typeof(%s->%s) = 'boolean'", field, sqlQuoteString(key))
	operand := fmt.Sprintf("(CASE WHEN %s THEN (%s->%s)::boolean END)", checkType, field, sqlQuoteString(key))
	cast := fmt.Sprintf("%s AND %s", checkType, operand)

	switch c.Type {
	case ConditionEqual:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}
		return fmt.Sprintf("(%s = %t)", cast, *c.valueBool), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) RefToSqlJson(field, key string) (string, error) {
	switch c.Type {
	case ConditionExist:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("%s ? %s", field, sqlQuoteString(key)), nil
		}

		return fmt.Sprintf("NOT (%s ? %s)", field, sqlQuoteString(key)), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) StringArrayToSqlJson(field, key string) (string, error) {
	// "CASE" is mandatory to cast json value because Postgres "SELECT" with following condition returns an error
	// if there is a row with field -> key of another (not numeric) type:
	// jsonb_typeof(field -> key) = 'number' AND (field -> key)::numeric = 2
	checkType := fmt.Sprintf("jsonb_typeof(%s->%s) = 'array'", field, sqlQuoteString(key))
	operand := fmt.Sprintf("(CASE WHEN %s THEN %s->%s END)", checkType, field, sqlQuoteString(key))
	cast := fmt.Sprintf("%s AND %s->%s", checkType, field, sqlQuoteString(key))

	switch c.Type {
	case ConditionIsEmpty:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("(%s AND jsonb_array_length(%s) = 0)", checkType, operand), nil
		}

		return fmt.Sprintf("(%s AND jsonb_array_length(%s) > 0)", checkType, operand), nil
	case ConditionHasEvery:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("(%s ?& ARRAY [%s])", cast, strings.Join(values, ",")), nil
	case ConditionHasOneOf:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("(%s ?| ARRAY [%s])", cast, strings.Join(values, ",")), nil
	case ConditionHasNot:
		if len(c.valueStrArray) == 0 {
			return "", ErrWrongConditionValue
		}

		values := make([]string, len(c.valueStrArray))
		for i, s := range c.valueStrArray {
			values[i] = sqlQuoteString(s)
		}

		return fmt.Sprintf("(%s AND NOT (%s ?| ARRAY [%s]))", checkType, operand, strings.Join(values, ",")), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) TimeToSqlJson(field, key string) (string, error) {
	checkType := fmt.Sprintf("jsonb_typeof(%s->%s) = 'number'", field, sqlQuoteString(key))
	operand := fmt.Sprintf("(%s->%s)::bigint", field, sqlQuoteString(key))
	cast := fmt.Sprintf("%s AND %s", checkType, operand)

	switch c.Type {
	case ConditionTimeRelative:
		if c.valueDuration == nil && c.valueDurationTo == nil {
			return "", ErrWrongConditionValue
		}
		if c.valueDuration != nil && c.valueDurationTo != nil && *c.valueDuration < *c.valueDurationTo {
			return "", fmt.Errorf("duration from (%d) < to (%d)", *c.valueDuration, *c.valueDurationTo)
		}
		var conditions []string
		if c.valueDuration != nil {
			conditions = append(conditions, fmt.Sprintf("%s > EXTRACT(EPOCH FROM (NOW() - INTERVAL '%d seconds'))::bigint", operand, *c.valueDuration))
		}
		if c.valueDurationTo != nil {
			conditions = append(conditions, fmt.Sprintf("%s < EXTRACT(EPOCH FROM (NOW() - INTERVAL '%d seconds'))::bigint", operand, *c.valueDurationTo))
		}

		return fmt.Sprintf("(%s AND %s)", checkType, strings.Join(conditions, " AND ")), nil
	case ConditionTimeAbsolute:
		if c.valueTimeIntervalFrom == nil || c.valueTimeIntervalTo == nil {
			return "", ErrWrongConditionValue
		}

		return fmt.Sprintf("(%s > %d AND %s < %d)", cast, *c.valueTimeIntervalFrom, operand, *c.valueTimeIntervalTo), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

func (c *Condition) UnmarshalJSON(b []byte) error {
	type Alias Condition
	v := Alias{}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*c = Condition(v)
	c.parseValue()

	return nil
}

func (c *Condition) UnmarshalBSONValue(_ byte, b []byte) error {
	type Alias Condition
	v := Alias{}
	err := bson.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*c = Condition(v)
	c.parseValue()

	return nil
}

func (c *Condition) parseValue() {
	if s, err := GetStringValue(c.Value); err == nil {
		regexpStr := ""
		switch c.Type {
		case ConditionRegexp:
			regexpStr = s
		case ConditionContain:
			regexpStr = utils.EscapeRegex(s)
		case ConditionNotContain:
			regexpStr = "^((?!" + utils.EscapeRegex(s) + ").)*$"
		case ConditionBeginWith:
			regexpStr = "^" + utils.EscapeRegex(s)
		case ConditionNotBeginWith:
			regexpStr = "^(?!" + utils.EscapeRegex(s) + ")"
		case ConditionEndWith:
			regexpStr = utils.EscapeRegex(s) + "$"
		case ConditionNotEndWith:
			regexpStr = "(?<!" + utils.EscapeRegex(s) + ")$"
		}
		if regexpStr != "" {
			if r, err := utils.NewRegexExpression(regexpStr); err == nil {
				c.valueRegexp = r
			}

			return
		}

		c.Value = s
		c.valueStr = &s
		return
	}

	if i, err := GetIntValue(c.Value); err == nil {
		c.Value = i
		c.valueInt = &i
		return
	}

	if b, err := GetBoolValue(c.Value); err == nil {
		c.Value = b
		c.valueBool = &b
		return
	}

	if a, err := GetStringArrayValue(c.Value); err == nil {
		c.Value = a
		c.valueStrArray = a
		return
	}

	if from, to, err := getTimeIntervalValue(c.Value); err == nil {
		c.Value = map[string]int64{
			"from": from,
			"to":   to,
		}
		c.valueTimeIntervalFrom = &from
		c.valueTimeIntervalTo = &to
		return
	}

	if from, to, err := getDurationValues(c.Value); err == nil {
		if from == nil {
			if to == nil {
				return
			}

			dBySec, err := to.To(datetime.DurationUnitSecond)
			if err != nil {
				return
			}
			c.Value = map[string]*datetime.DurationWithUnit{
				"to": to,
			}
			c.valueDurationTo = &dBySec.Value
			return
		}

		fromBySec, err := from.To(datetime.DurationUnitSecond)
		if err != nil {
			return
		}

		c.valueDuration = &fromBySec.Value
		if to == nil {
			c.Value = *from
			return
		}

		c.Value = map[string]*datetime.DurationWithUnit{
			"from": from,
			"to":   to,
		}
		toBySec, err := to.To(datetime.DurationUnitSecond)
		if err != nil {
			return
		}

		c.valueDurationTo = &toBySec.Value
		return
	}
}

// ValidateInfoCondition is a helper function to validate FieldCondition when it's used to match various infos fields.
func (c *FieldCondition) ValidateInfoCondition() bool {
	var err error

	switch c.FieldType {
	case FieldTypeString:
		_, err = c.Condition.MatchString("")
	case FieldTypeInt:
		_, err = c.Condition.MatchInt(0)
	case FieldTypeBool:
		_, err = c.Condition.MatchBool(false)
	case FieldTypeStringArray:
		_, err = c.Condition.MatchStringArray([]string{})
	case "":
		_, err = c.Condition.MatchRef(nil)
	default:
		return false
	}

	return err == nil
}

// ValidateEntityInfoCondition is a helper function to validate FieldCondition when it's used to match various infos fields.
func (c *FieldCondition) ValidateEntityInfoCondition() bool {
	var err error

	switch c.FieldType {
	case FieldTypeString:
		_, err = c.Condition.MatchString("")
	case FieldTypeInt:
		_, err = c.Condition.MatchInt(0)
	case FieldTypeBool:
		_, err = c.Condition.MatchBool(false)
	case FieldTypeStringArray:
		_, err = c.Condition.MatchStringArray([]string{})
	case FieldTypeTimestamp:
		_, err = c.Condition.MatchTime(time.Time{}, time.Now())
	case "":
		_, err = c.Condition.MatchRef(nil)
	default:
		return false
	}

	return err == nil
}

// MatchAlarmInfoCondition is a helper function to match FieldCondition when it's used to match various alarm infos fields.
func (c *FieldCondition) MatchAlarmInfoCondition(infoVal any, infoExists bool) (bool, error) {
	var matched bool
	var err error

	if c.FieldType == "" {
		matched, err = c.Condition.MatchRef(infoVal)
	} else if infoExists {
		switch c.FieldType {
		case FieldTypeString:
			var s string
			if s, err = GetStringValue(infoVal); err == nil {
				matched, err = c.Condition.MatchString(s)
			}
		case FieldTypeInt:
			var i int64
			if i, err = GetIntValue(infoVal); err == nil {
				matched, err = c.Condition.MatchInt(i)
			}
		case FieldTypeBool:
			var b bool
			if b, err = GetBoolValue(infoVal); err == nil {
				matched, err = c.Condition.MatchBool(b)
			}
		case FieldTypeStringArray:
			var a []string
			if a, err = GetStringArrayValue(infoVal); err == nil {
				matched, err = c.Condition.MatchStringArray(a)
			}
		default:
			return false, fmt.Errorf("invalid field type for %q field: %s", c.Field, c.FieldType)
		}
	}

	return matched, err
}

// MatchEntityInfoCondition is a helper function to match FieldCondition when it's used to match various entity infos fields.
func (c *FieldCondition) MatchEntityInfoCondition(infoVal any, infoExists bool) (bool, error) {
	var matched bool
	var err error

	if c.FieldType == "" {
		matched, err = c.Condition.MatchRef(infoVal)
	} else if infoExists {
		switch c.FieldType {
		case FieldTypeString:
			var s string
			if s, err = GetStringValue(infoVal); err == nil {
				matched, err = c.Condition.MatchString(s)
			}
		case FieldTypeInt:
			var i int64
			if i, err = GetIntValue(infoVal); err == nil {
				matched, err = c.Condition.MatchInt(i)
			}
		case FieldTypeBool:
			var b bool
			if b, err = GetBoolValue(infoVal); err == nil {
				matched, err = c.Condition.MatchBool(b)
			}
		case FieldTypeStringArray:
			var a []string
			if a, err = GetStringArrayValue(infoVal); err == nil {
				matched, err = c.Condition.MatchStringArray(a)
			}
		case FieldTypeTimestamp:
			var t time.Time
			if t, err = GetTimeValue(infoVal); err == nil {
				matched, err = c.Condition.MatchTime(t, time.Now())
			}
		default:
			return false, fmt.Errorf("invalid field type for %q field: %s", c.Field, c.FieldType)
		}
	}

	return matched, err
}

func (c *Condition) GetRegexp() utils.RegexExpression {
	return c.valueRegexp
}

func GetStringValue(v any) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}

	return "", ErrWrongConditionValue
}

func GetIntValue(v any) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return i, nil
	case uint:
		return int64(i), nil
	case uint32:
		return int64(i), nil
	case uint64:
		return int64(i), nil
	case float32:
		a, b := math.Modf(float64(i))
		if b == 0 {
			return int64(a), nil
		}
	case float64:
		a, b := math.Modf(i)
		if b == 0 {
			return int64(a), nil
		}
	}

	return 0, ErrWrongConditionValue
}

func GetBoolValue(v any) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}

	return false, ErrWrongConditionValue
}

func GetStringArrayValue(v any) ([]string, error) {
	var interfaceArr []any

	switch a := v.(type) {
	case []string:
		return a, nil
	case []any:
		interfaceArr = a
	case bson.A:
		interfaceArr = a
	default:
		return nil, ErrWrongConditionValue
	}

	l := len(interfaceArr)
	strArr := make([]string, l)
	for i := 0; i < l; i++ {
		if s, ok := interfaceArr[i].(string); ok {
			strArr[i] = s
		} else {
			return nil, ErrWrongConditionValue
		}
	}

	return strArr, nil
}

func GetTimeValue(v any) (time.Time, error) {
	intVal, err := GetIntValue(v)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(intVal, 0), nil
}

func getTimeIntervalValue(v any) (int64, int64, error) {
	mapVal, err := parseMapVal(v)
	if err != nil {
		return 0, 0, err
	}

	rawFrom, ok := mapVal["from"]
	if !ok {
		return 0, 0, errors.New("condition value expected 'from' key")
	}

	from, err := GetIntValue(rawFrom)
	if err != nil {
		return 0, 0, err
	}

	rawTo, ok := mapVal["to"]
	if !ok {
		return 0, 0, errors.New("condition value expected 'to' key")
	}

	to, err := GetIntValue(rawTo)
	if err != nil {
		return 0, 0, err
	}

	return from, to, nil
}

func getDurationValues(v any) (from, to *datetime.DurationWithUnit, err error) {
	mapVal, err := parseMapVal(v)
	if err != nil {
		return nil, nil, err
	}

	rawTo, ok := mapVal["to"]
	if !ok {
		from, err = parseDurationWithUnit(mapVal)
		if err != nil {
			return nil, nil, err
		}

		return from, nil, nil
	}

	rawVal, err := parseMapVal(rawTo)
	if err != nil {
		return nil, nil, err
	}

	to, err = parseDurationWithUnit(rawVal)
	if err != nil {
		return nil, nil, err
	}

	rawVal, err = parseMapVal(mapVal["from"])
	if err == nil {
		from, err = parseDurationWithUnit(rawVal)
		if err != nil {
			return nil, nil, err
		}
	}

	return from, to, nil
}

func parseMapVal(v any) (map[string]any, error) {
	var mapVal map[string]any
	switch m := v.(type) {
	case map[string]any:
		mapVal = m
	case bson.D:
		mapVal = make(map[string]any, len(m))
		for _, e := range m {
			mapVal[e.Key] = e.Value
		}
	case bson.M:
		mapVal = m
	default:
		return nil, ErrWrongConditionValue
	}
	return mapVal, nil
}

func parseDurationWithUnit(m map[string]any) (*datetime.DurationWithUnit, error) {
	rawVal, ok := m["value"]
	if !ok {
		return nil, errors.New("condition value expected 'value' key")
	}

	val, err := GetIntValue(rawVal)
	if err != nil {
		return nil, err
	}

	rawUnit, ok := m["unit"]
	if !ok {
		return nil, errors.New("condition value expected 'unit' key")
	}

	unit, err := GetStringValue(rawUnit)
	if err != nil {
		return nil, err
	}

	return &datetime.DurationWithUnit{
		Value: val,
		Unit:  unit,
	}, nil
}

func sqlQuoteString(str string) string {
	return "'" + strings.ReplaceAll(str, "'", "''") + "'"
}
