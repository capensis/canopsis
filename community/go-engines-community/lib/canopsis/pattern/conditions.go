// Package pattern provides functionality for filtering and matching models.
package pattern

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
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
	ConditionIsOneOf      = "is_one_of"
	ConditionIsNotOneOf   = "is_not_one_of"
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

	valueStr              *string
	valueRegexp           utils.RegexExpression
	valueInt              *int64
	valueBool             *bool
	valueStrArray         []string
	valueTimeIntervalFrom *int64
	valueTimeIntervalTo   *int64
	valueDuration         *int64
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

func NewDurationCondition(t string, d types.DurationWithUnit) (Condition, error) {
	var err error
	d, err = d.To("s")
	if err != nil {
		return Condition{}, err
	}

	return Condition{
		Type:          t,
		Value:         d,
		valueDuration: &d.Value,
	}, nil
}

func (c *Condition) MatchString(value string) (bool, RegexMatches, error) {
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
		for _, item := range c.valueStrArray {
			if item == value {
				return true, nil, nil
			}
		}

		return false, nil, nil
	case ConditionIsNotOneOf:
		if len(c.valueStrArray) == 0 {
			return false, nil, ErrWrongConditionValue
		}

		for _, item := range c.valueStrArray {
			if item == value {
				return false, nil, nil
			}
		}

		return true, nil, nil
	case ConditionRegexp:
		if c.valueRegexp == nil {
			return false, nil, ErrWrongConditionValue
		}

		regexMatches := utils.FindStringSubmatchMapWithRegexExpression(c.valueRegexp, value)

		return regexMatches != nil, regexMatches, nil
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

func (c *Condition) MatchRef(value interface{}) (bool, error) {
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

func (c *Condition) MatchTime(value time.Time) (bool, error) {
	switch c.Type {
	case ConditionTimeRelative:
		if c.valueDuration == nil {
			return false, ErrWrongConditionValue
		}

		return value.After(time.Now().Add(time.Duration(-*c.valueDuration) * time.Second)), nil
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

func (c *Condition) ToMongoQuery(f string) (bson.M, error) {
	switch c.Type {
	case ConditionEqual:
		return bson.M{f: bson.M{"$eq": c.Value}}, nil
	case ConditionNotEqual:
		return bson.M{f: bson.M{"$ne": c.Value}}, nil
	case ConditionGT:
		if c.valueDuration != nil {
			return bson.M{f: bson.M{"$gt": c.valueDuration}}, nil
		}

		return bson.M{f: bson.M{"$gt": c.Value}}, nil
	case ConditionLT:
		if c.valueDuration != nil {
			return bson.M{f: bson.M{"$lt": c.valueDuration}}, nil
		}

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
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			return bson.M{f: bson.M{"$in": bson.A{nil, bson.A{}}}}, nil
		}

		return bson.M{f: bson.M{"$exists": true, "$type": "array", "$ne": bson.A{}}}, nil
	case ConditionExist:
		if c.valueBool == nil {
			return nil, ErrWrongConditionValue
		}

		if *c.valueBool {
			return bson.M{f: bson.M{"$exists": true, "$ne": nil}}, nil
		}

		return bson.M{"$or": []bson.M{{f: bson.M{"$exists": false}}, {f: bson.M{"$eq": nil}}}}, nil
	case ConditionTimeRelative:
		if c.valueDuration == nil {
			return nil, ErrWrongConditionValue
		}

		t := types.CpsTime{Time: time.Now().Add(time.Duration(-*c.valueDuration) * time.Second)}

		return bson.M{f: bson.M{"$gt": t}}, nil
	case ConditionTimeAbsolute:
		if c.valueTimeIntervalFrom == nil || c.valueTimeIntervalTo == nil {
			return nil, ErrWrongConditionValue
		}

		ft := types.NewCpsTime(*c.valueTimeIntervalFrom)
		tt := types.NewCpsTime(*c.valueTimeIntervalTo)

		return bson.M{f: bson.M{"$gt": ft, "$lt": tt}}, nil
	case ConditionIsOneOf:
		return bson.M{f: bson.M{"$in": c.Value}}, nil
	case ConditionIsNotOneOf:
		return bson.M{f: bson.M{"$nin": c.Value}}, nil
	default:
		return nil, ErrUnsupportedConditionType
	}
}

// ToSql doesn't support all conditions. Add on demand.
func (c *Condition) ToSql(f string) (string, error) {
	switch c.Type {
	case ConditionEqual:
		if c.valueStr != nil {
			return fmt.Sprintf("%s = %s", f, sqlQuoteString(*c.valueStr)), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("%s = %d", f, *c.valueInt), nil
		}
		if c.valueBool != nil {
			return fmt.Sprintf("%s = %t", f, *c.valueBool), nil
		}

		return "", ErrWrongConditionValue
	case ConditionNotEqual:
		// "IS NULL" is mandatory
		if c.valueStr != nil {
			return fmt.Sprintf("(%[1]s IS NULL OR %[1]s != %s)", f, sqlQuoteString(*c.valueStr)), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("(%[1]s IS NULL OR %[1]s != %d)", f, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionGT:
		if c.valueDuration != nil {
			return fmt.Sprintf("%s > %d", f, *c.valueDuration), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("%s > %d", f, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionLT:
		if c.valueDuration != nil {
			return fmt.Sprintf("%s < %d", f, *c.valueDuration), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("%s < %d", f, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionRegexp:
		if c.valueRegexp != nil {
			return fmt.Sprintf("%s ~ %s", f, sqlQuoteString(c.valueRegexp.String())), nil
		}

		return "", ErrWrongConditionValue
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
		return fmt.Sprintf("(%[1]s IS NULL OR NOT (%[1]s = ANY (ARRAY [%s]))", f, strings.Join(values, ",")), nil
	default:
		return "", ErrUnsupportedConditionType
	}
}

// ToSqlJson doesn't support all conditions. Add on demand.
func (c *Condition) ToSqlJson(field, key, keyType string) (string, error) {
	cast := ""
	checkType := ""
	operand := fmt.Sprintf("%s->%s", field, sqlQuoteString(key))
	// "CASE" is mandatory to cast json value because Postgres "SELECT" with following condition returns an error
	// if there is a row with field -> key of another (not numeric) type:
	// jsonb_typeof(field -> key) = 'number' AND (field -> key)::numeric = 2
	switch keyType {
	case FieldTypeString:
		checkType = fmt.Sprintf("jsonb_typeof(%s) = 'string'", operand)
		operand = fmt.Sprintf("%s->>%s", field, sqlQuoteString(key))
	case FieldTypeInt:
		checkType = fmt.Sprintf("jsonb_typeof(%s) = 'number'", operand)
		operand = fmt.Sprintf("(CASE WHEN %s THEN (%s)::numeric END)", checkType, operand)
	case FieldTypeBool:
		checkType = fmt.Sprintf("jsonb_typeof(%s) = 'boolean'", operand)
		operand = fmt.Sprintf("(CASE WHEN %s THEN (%s)::bool END)", checkType, operand)
	case FieldTypeStringArray:
		checkType = fmt.Sprintf("jsonb_typeof(%s) = 'array'", operand)
		operand = fmt.Sprintf("(CASE WHEN %s THEN %s END)", checkType, operand)
	case "":
		/*do nothing*/
	default:
		return "", ErrUnsupportedField
	}

	if checkType == "" {
		cast = operand
	} else {
		cast = fmt.Sprintf("%s AND %s", checkType, operand)
	}

	switch c.Type {
	case ConditionEqual:
		if c.valueStr != nil {
			return fmt.Sprintf("(%s = %s)", cast, sqlQuoteString(*c.valueStr)), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("(%s = %d)", cast, *c.valueInt), nil
		}
		if c.valueBool != nil {
			return fmt.Sprintf("(%s = %t)", cast, *c.valueBool), nil
		}

		return "", ErrWrongConditionValue
	case ConditionNotEqual:
		if c.valueStr != nil {
			return fmt.Sprintf("(%s != %s)", cast, sqlQuoteString(*c.valueStr)), nil
		}
		if c.valueInt != nil {
			return fmt.Sprintf("(%s != %d)", cast, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionGT:
		if c.valueInt != nil {
			return fmt.Sprintf("(%s > %d)", cast, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionLT:
		if c.valueInt != nil {
			return fmt.Sprintf("(%s < %d)", cast, *c.valueInt), nil
		}

		return "", ErrWrongConditionValue
	case ConditionRegexp:
		if c.valueRegexp != nil {
			return fmt.Sprintf("(%s ~ %s)", cast, sqlQuoteString(c.valueRegexp.String())), nil
		}
		return "", ErrWrongConditionValue
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
	case ConditionIsEmpty:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("(%s AND jsonb_array_length(%s) = 0)", checkType, operand), nil
		}

		return fmt.Sprintf("(%s AND jsonb_array_length(%s) > 0)", checkType, operand), nil
	case ConditionExist:
		if c.valueBool == nil {
			return "", ErrWrongConditionValue
		}

		if *c.valueBool {
			return fmt.Sprintf("%s ? %s", field, sqlQuoteString(key)), nil
		}

		return fmt.Sprintf("NOT (%s ? %s)", field, sqlQuoteString(key)), nil
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

func (c *Condition) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
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
	if s, err := getStringValue(c.Value); err == nil {
		if c.Type == ConditionRegexp {
			if r, err := utils.NewRegexExpression(s); err == nil {
				c.valueRegexp = r
			}

			return
		}

		c.Value = s
		c.valueStr = &s
		return
	}

	if i, err := getIntValue(c.Value); err == nil {
		c.Value = i
		c.valueInt = &i
		return
	}

	if b, err := getBoolValue(c.Value); err == nil {
		c.Value = b
		c.valueBool = &b
		return
	}

	if a, err := getStringArrayValue(c.Value); err == nil {
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

	if d, err := getDurationValue(c.Value); err == nil {
		dBySec, err := d.To("s")
		if err == nil {
			c.Value = d
			c.valueDuration = &dBySec.Value
		}
		return
	}
}

func getStringValue(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}

	return "", ErrWrongConditionValue
}

func getIntValue(v interface{}) (int64, error) {
	switch i := v.(type) {
	case int:
		return int64(i), nil
	case int32:
		return int64(i), nil
	case int64:
		return i, nil
	case float32:
		a, b := math.Modf(float64(i))
		if b == 0 {
			return int64(a), nil
		}

		return 0, ErrWrongConditionValue
	case float64:
		a, b := math.Modf(i)
		if b == 0 {
			return int64(a), nil
		}

		return 0, ErrWrongConditionValue
	default:
		return 0, ErrWrongConditionValue
	}
}

func getBoolValue(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}

	return false, ErrWrongConditionValue
}

func getStringArrayValue(v interface{}) ([]string, error) {
	var interfaceArr []interface{}

	switch a := v.(type) {
	case []string:
		return a, nil
	case []interface{}:
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

func getTimeIntervalValue(v interface{}) (int64, int64, error) {
	var mapVal map[string]interface{}
	if m, ok := v.(map[string]interface{}); ok {
		mapVal = m
	} else if m, ok := v.(bson.D); ok {
		mapVal = m.Map()
	} else if m, ok := v.(bson.M); ok {
		mapVal = m
	} else {
		return 0, 0, ErrWrongConditionValue
	}

	rawFrom, ok := mapVal["from"]
	if !ok {
		return 0, 0, errors.New("condition value expected 'from' key")
	}

	from, err := getIntValue(rawFrom)
	if err != nil {
		return 0, 0, err
	}

	rawTo, ok := mapVal["to"]
	if !ok {
		return 0, 0, errors.New("condition value expected 'to' key")
	}

	to, err := getIntValue(rawTo)
	if err != nil {
		return 0, 0, err
	}

	return from, to, nil
}

func getDurationValue(v interface{}) (types.DurationWithUnit, error) {
	var mapVal map[string]interface{}
	if m, ok := v.(map[string]interface{}); ok {
		mapVal = m
	} else if m, ok := v.(bson.D); ok {
		mapVal = m.Map()
	} else if m, ok := v.(bson.M); ok {
		mapVal = m
	} else {
		return types.DurationWithUnit{}, ErrWrongConditionValue
	}

	rawVal, ok := mapVal["value"]
	if !ok {
		return types.DurationWithUnit{}, errors.New("condition value expected 'value' key")
	}

	val, err := getIntValue(rawVal)
	if err != nil {
		return types.DurationWithUnit{}, err
	}

	rawUnit, ok := mapVal["unit"]
	if !ok {
		return types.DurationWithUnit{}, errors.New("condition value expected 'unit' key")
	}

	unit, err := getStringValue(rawUnit)
	if err != nil {
		return types.DurationWithUnit{}, err
	}

	return types.DurationWithUnit{
		Value: val,
		Unit:  unit,
	}, nil
}

func sqlQuoteString(str string) string {
	return "'" + strings.Replace(str, "'", "''", -1) + "'"
}
