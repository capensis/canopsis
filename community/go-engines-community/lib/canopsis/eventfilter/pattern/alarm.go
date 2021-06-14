package pattern

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
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

// NewAlarmRegexMatches creates an AlarmValueRegexMatches, with the Value.Extra
// field initialized to an empty map.
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

// AsMongoQuery returns a mongodb filter from the AlarmFields
func (p AlarmFields) AsMongoQuery() mgobson.M {
	query := mgobson.M{}
	if !p.ID.Empty() {
		query["_id"] = p.ID.AsMongoQuery()
	}
	if !p.Time.Empty() {
		query["t"] = p.Time.AsMongoQuery()
	}
	if !p.EntityID.Empty() {
		query["d"] = p.EntityID.AsMongoQuery()
	}
	p.Value.AsMongoQuery("v", query)

	return query
}

func (p AlarmFields) AsMongoDriverQuery() mongobson.M {
	query := mongobson.M{}
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

// AsMongoQuery returns a mongodb filter from the AlarmPattern
func (p AlarmPattern) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)

	if p.ShouldBeNil {
		return nil
	} else if p.ShouldNotBeNil {
		return p.AlarmFields.AsMongoQuery()
	}
	return query
}

func (p AlarmPattern) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)

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

// SetBSON unmarshals a BSON value into an AlarmPattern.
func (p *AlarmPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		// The BSON value is null. The field should not be set.
		p.ShouldBeNil = true
		p.ShouldNotBeNil = false
		return nil

	default:
		err := raw.Unmarshal(&p.AlarmFields)
		if err != nil {
			return err
		}

		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}
			return fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", "))
		}

		p.ShouldBeNil = false
		p.ShouldNotBeNil = true
		return nil
	}
}

func (p AlarmPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.ShouldBeNil {
		return bsontype.Null, []byte{}, nil
	}

	resultBson := mongobson.M{}

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
		return mongobson.MarshalValue(resultBson)
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
		err := mongobson.Unmarshal(b, &p.AlarmFields)
		if err != nil {
			return err
		}

		if len(p.UnexpectedFields) != 0 {
			unexpectedFieldNames := make([]string, 0, len(p.UnexpectedFields))
			for key := range p.UnexpectedFields {
				unexpectedFieldNames = append(unexpectedFieldNames, key)
			}

			return UnexpectedFieldsError{
				Err: fmt.Errorf("Unexpected pattern fields: %s", strings.Join(unexpectedFieldNames, ", ")),
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
type AlarmPatternList struct {
	Patterns []AlarmPattern

	// isSet is a boolean indicating whether the AlarmPatternList has been set
	// explicitly or not.
	Set bool

	// isValid is a boolean indicating whether the event patterns or valid or
	// not.
	// isValid is also false if the AlarmPatternList has not been set.
	Valid bool
}

func (l *AlarmPatternList) UnmarshalJSON(b []byte) error {
	var jsonPatterns interface{}
	err := json.Unmarshal(b, &jsonPatterns)
	if err != nil {
		return err
	}

	marshalled, err := mongobson.Marshal(mongobson.M{
		"list": jsonPatterns,
	})
	if err != nil {
		return err
	}

	var patternWrapper struct {
		PatternList AlarmPatternList `bson:"list"`
	}

	err = mongobson.Unmarshal(marshalled, &patternWrapper)
	if err != nil {
		return err
	}

	*l = patternWrapper.PatternList
	return nil
}

func (l AlarmPatternList) MarshalJSON() ([]byte, error) {
	bsonType, bson, err := mongobson.MarshalValue(l)
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
	raw := mongobson.RawValue{
		Type:  bsontype.Array,
		Value: bson,
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

	return mongobson.MarshalValue(l.Patterns)
}

func NewAlarmPatternList(p []AlarmPattern) AlarmPatternList {
	return AlarmPatternList{
		Patterns: p,
		Set:      true,
		Valid:    true,
	}
}

// AsMongoQuery returns a mongodb filter from the AlarmPatternList
func (l AlarmPatternList) AsMongoQuery() mgobson.M {
	var patternFilters []mgobson.M
	if !l.Set {
		return mgobson.M{}
	}
	for _, alarmPattern := range l.Patterns {
		patternFilters = append(patternFilters, alarmPattern.AsMongoQuery())
	}
	return mgobson.M{"$or": patternFilters}
}

func (l AlarmPatternList) AsMongoDriverQuery() mongobson.M {
	var patternFilters []mongobson.M
	if !l.Set {
		return mongobson.M{}
	}
	for _, alarmPattern := range l.Patterns {
		patternFilters = append(patternFilters, alarmPattern.AsMongoDriverQuery())
	}
	return mongobson.M{"$or": patternFilters}
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

// SetBSON unmarshals a BSON value into a AlarmPatternList.
// If it cannot be unmarshalled, it is marked as invalid.
func (l *AlarmPatternList) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		l.Set = false
		l.Valid = false
	default:
		l.Set = true

		err := raw.Unmarshal(&l.Patterns)
		if err != nil {
			log.Printf("unable to parse alarm pattern list: %v", err)
			l.Valid = false
			return nil
		}

		l.Valid = true
	}
	return nil
}

func (l *AlarmPatternList) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		l.Set = false
		l.Valid = false
	case bsontype.Array:
		l.Set = true
		l.Valid = false

		var raw mongobson.Raw
		err := mongobson.Unmarshal(b, &raw)
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

			err = mongobson.Unmarshal(document, &pattern)
			if err != nil {
				if _, ok = err.(UnexpectedFieldsError); ok {
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
