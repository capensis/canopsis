// Package pattern.
// Deprecated: use git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern instead.
package oldpattern

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// AlarmRegexMatches is a type that contains the values of the sub-expressions
// of regular expressions for each of the fields of an Alarm that contain
// strings.
type AlarmRegexMatches struct {
	ID       RegexMatches
	EntityID RegexMatches
	Value    AlarmValueRegexMatches
}

// NewAlarmRegexMatches creates an AlarmValueRegexMatches.
func NewAlarmRegexMatches() AlarmRegexMatches {
	return AlarmRegexMatches{
		Value: NewAlarmValueRegexMatches(),
	}
}

// AlarmFields is a type representing a pattern that can be applied to an
// alarm
type AlarmFields struct {
	ID       StringPattern     `bson:"_id"`
	Time     TimePattern       `bson:"t"`
	EntityID StringPattern     `bson:"d"`
	Value    AlarmValuePattern `bson:"v"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

func (p AlarmFields) AsMongoDriverQuery() bson.M {
	query := bson.M{}
	if !p.ID.Empty() {
		query["_id"] = p.ID.AsMongoDriverQuery()
	}
	if !p.Time.Empty() {
		query["t"] = p.Time.AsMongoDriverQuery()
	}
	if !p.EntityID.Empty() {
		query["d"] = p.EntityID.AsMongoDriverQuery()
	}
	p.Value.AsMongoDriverQuery("v", query)

	return query
}

// AlarmPattern is a type representing a pattern that can be applied to an
// alarm
type AlarmPattern struct {
	// ShouldNotBeNil is a boolean indicating that the alarm should not be nil,
	// and ShouldBeNil is a boolean indicating that the alarm should be nil.
	// The two booleans are needed to be able to make the difference between
	// the case where no alarm pattern was defined (in which case the alarm can
	// be nil or not), the case where a nil alarm pattern was defined (in which
	// case the alarm should be nil), and the case where a non-nil alarm
	// pattern was defined (in which case the alarm should not be nil).
	ShouldNotBeNil bool
	ShouldBeNil    bool

	AlarmFields
}

func (p AlarmPattern) AsMongoDriverQuery() bson.M {
	query := make(bson.M)

	if p.ShouldBeNil {
		return nil
	} else if p.ShouldNotBeNil {
		return p.AlarmFields.AsMongoDriverQuery()
	}
	return query
}

// Matches returns true if an alarm is matched by a pattern. If the pattern
// contains regular expressions with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p AlarmPattern) Matches(alarm *types.Alarm, matches *AlarmRegexMatches) bool {
	if p.ShouldBeNil {
		return alarm == nil
	}

	if p.ShouldNotBeNil {
		return alarm != nil &&
			p.ID.Matches(alarm.ID, &matches.ID) &&
			p.Time.Matches(alarm.Time) &&
			p.EntityID.Matches(alarm.EntityID, &matches.EntityID) &&
			p.Value.Matches(alarm.Value, &matches.Value)
	}

	return true
}

func (p AlarmPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := bson.M{}

	if p.ID.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "ID", "id")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.ID
	}

	if p.Time.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Time", "time")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Time
	}

	if p.EntityID.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "EntityID", "entityid")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.EntityID
	}

	if p.Value.IsSet() {
		bsonFieldName, err := GetFieldBsonName(p, "Value", "value")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Value
	}

	if len(resultBson) > 0 {
		return bson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p *AlarmPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		// The BSON value is null. The field should not be set.
		p.ShouldBeNil = true
		p.ShouldNotBeNil = false
		return nil
	default:
		err := bson.Unmarshal(b, &p.AlarmFields)
		if err != nil {
			return err
		}

		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}

			return UnexpectedFieldsError{
				Err: fmt.Errorf("unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
			}
		}
	}

	p.ShouldBeNil = false
	p.ShouldNotBeNil = true
	return nil
}

// AlarmPatternList is a type representing a list of alarm patterns.
// An alarm is matched by an AlarmPatternList if it is matched by one of its
// AlarmPatterns.
// The zero value of an AlarmPatternList (i.e. an AlarmPatternList that has
// not been set) is considered valid, and matches all alarms.
// Deprecated : community/go-engines-community/lib/canopsis/pattern/Alarm
type AlarmPatternList struct {
	Patterns []AlarmPattern `swaggerignore:"true"`

	// isSet is a boolean indicating whether the AlarmPatternList has been set
	// explicitly or not.
	Set bool `swaggerignore:"true"`

	// isValid is a boolean indicating whether the event patterns or valid or
	// not.
	// isValid is also false if the AlarmPatternList has not been set.
	Valid bool `swaggerignore:"true"`
}

func (l *AlarmPatternList) UnmarshalJSON(b []byte) error {
	var jsonPatterns interface{}
	err := json.Unmarshal(b, &jsonPatterns)
	if err != nil {
		return err
	}

	marshalled, err := bson.Marshal(bson.M{
		"list": jsonPatterns,
	})
	if err != nil {
		return err
	}

	var patternWrapper struct {
		PatternList AlarmPatternList `bson:"list"`
	}

	err = bson.Unmarshal(marshalled, &patternWrapper)
	if err != nil {
		return err
	}

	*l = patternWrapper.PatternList
	return nil
}

func (l AlarmPatternList) MarshalJSON() ([]byte, error) {
	bsonType, b, err := bson.MarshalValue(l)
	if err != nil {
		return nil, err
	}

	if bsonType == bsontype.Null {
		res, err := json.Marshal(nil)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	var unmarshalledBson []map[string]interface{}
	raw := bson.RawValue{
		Type:  bsontype.Array,
		Value: b,
	}
	err = raw.Unmarshal(&unmarshalledBson)
	if err != nil {
		return nil, err
	}

	res, err := json.Marshal(unmarshalledBson)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (l AlarmPatternList) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if !l.Set {
		return bsontype.Null, nil, nil
	}

	return bson.MarshalValue(l.Patterns)
}

func NewAlarmPatternList(p []AlarmPattern) AlarmPatternList {
	return AlarmPatternList{
		Patterns: p,
		Set:      true,
		Valid:    true,
	}
}

func (l AlarmPatternList) AsMongoDriverQuery() bson.M {
	var patternFilters []bson.M
	if !l.Set {
		return bson.M{}
	}
	for _, alarmPattern := range l.Patterns {
		patternFilters = append(patternFilters, alarmPattern.AsMongoDriverQuery())
	}
	return bson.M{"$or": patternFilters}
}

// Matches returns true if the alarm is matched by the AlarmPatternList.
func (l AlarmPatternList) Matches(alarm *types.Alarm) bool {
	if !l.Set {
		return true
	}

	regexMatches := NewAlarmRegexMatches()
	for _, pattern := range l.Patterns {
		if pattern.Matches(alarm, &regexMatches) {
			return true
		}
	}

	return false
}

// IsSet returns true if the AlarmPatternList has been set explicitly.
func (l AlarmPatternList) IsSet() bool {
	return l.Set
}

// IsValid returns true if the AlarmPatternList is valid.
func (l AlarmPatternList) IsValid() bool {
	return !l.Set || l.Valid
}

func (l AlarmPatternList) IsZero() bool {
	return !l.Set
}

func (l *AlarmPatternList) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		l.Set = false
		l.Valid = false
	case bsontype.Array:
		l.Set = true
		l.Valid = false

		var raw bson.Raw
		err := bson.Unmarshal(b, &raw)
		if err != nil {
			return err
		}

		array, err := raw.Values()
		if err != nil {
			return err
		}

		for _, v := range array {
			document, ok := v.DocumentOK()
			if !ok {
				return fmt.Errorf("unable to parse alarm pattern list element")
			}

			var pattern AlarmPattern

			err = bson.Unmarshal(document, &pattern)
			if err != nil {
				if errors.As(err, &UnexpectedFieldsError{}) {
					return nil
				}

				return err
			}

			l.Patterns = append(l.Patterns, pattern)
		}

		l.Valid = true
	default:
		return fmt.Errorf("alarm pattern list should be an array or nil")
	}

	return nil
}

func (l *AlarmPatternList) AsInterface() (interface{}, error) {
	if l == nil {
		return nil, nil
	}
	b, err := l.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var m interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
