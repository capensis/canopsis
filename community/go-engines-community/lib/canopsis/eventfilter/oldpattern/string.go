package oldpattern

import (
	"errors"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// RegexMatches is a type that contains the values of the sub-expressions of a
// regular expression.
type RegexMatches map[string]string

// StringConditions is a struct representing a pattern that can be applied to
// the value of a string field of an event.
// Each field of a StringConditions represents a condition that is applied if
// the value of this field is not nil.
// The fields are not defined directly in the StringPattern struct to make the
// unmarshalling easier.
type StringConditions struct {
	// If Equal is set, the value of a field has to be equal to the value of
	// Equal to be matched by the pattern.
	Equal types.OptionalString

	// If RegexMatch is set, the value of a field has to be matched by this
	// regular expression to be matched by the pattern.
	RegexMatch types.OptionalRegexp `bson:"regex_match,omitempty"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	// UnexpectedFields should always be empty.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// AsMongoDriverQuery returns a mongodb filter from the StringConditions from mongo-driver
func (p StringConditions) AsMongoDriverQuery() bson.M {
	query := make(bson.M)

	if p.Equal.Set {
		query["$eq"] = p.Equal.Value
	}

	if p.RegexMatch.Set {
		query["$regex"] = p.RegexMatch.Value.String()
	}

	return query
}

func (p StringConditions) AsSqlQuery() string {
	if p.Equal.Set {
		return fmt.Sprintf("= '%s'", p.Equal.Value)
	}

	if p.RegexMatch.Set {
		return fmt.Sprintf("~ '%s'", p.RegexMatch.Value)
	}

	return ""
}

// Matches returns true if the value satisfies each of the conditions defined
// in the StringConditions. If the pattern contains a regular expression with
// sub-expressions, the values of the sub-expressions are written in the
// matches argument.
func (p StringConditions) Matches(value string, matches *RegexMatches) bool {
	if p.Equal.Set {
		if !(value == p.Equal.Value) {
			return false
		}
	}

	if p.RegexMatch.Set {

		submatches := utils.FindStringSubmatchMapWithRegexExpression(p.RegexMatch.Value, value)
		if submatches == nil {
			return false
		}
		*matches = submatches
	}

	return true
}

// Empty returns true if the none of the conditions have been set.
func (p StringConditions) Empty() bool {
	return !(p.Equal.Set || p.RegexMatch.Set)
}

// StringPattern is a type representing a pattern that can be applied to the
// value of a field of an event that contains a string.
type StringPattern struct {
	StringConditions
}

func (p StringPattern) IsSet() bool {
	return p.Equal.Set || p.RegexMatch.Set
}

func (p StringPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.Equal.Set {
		return bson.MarshalValue(p.Equal)
	}

	if p.RegexMatch.Set {
		bsonFieldName, err := GetFieldBsonName(p, "RegexMatch", "regexmatch")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		return bson.MarshalValue(bson.M{
			bsonFieldName: p.RegexMatch.Value.String(),
		})
	}

	return bsontype.Undefined, nil, nil
}

func (p *StringPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.String:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		p.Equal.Value = value
		p.Equal.Set = true
	case bsontype.EmbeddedDocument:
		err := bson.Unmarshal(b, &p.StringConditions)
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
	default:
		return fmt.Errorf("unable to parse string")
	}

	return nil
}

// StringRefPattern is a type representing a pattern that can be applied to the value
// of a field of an event that contains a reference to a string.
type StringRefPattern struct {
	// If EqualNil is true, the field should be nil or not be set to be matched
	// by the pattern.
	EqualNil bool

	StringPattern
}

// AsMongoDriverQuery returns a mongodb filter from the StringRefPattern for mongo-driver.
func (p StringRefPattern) AsMongoDriverQuery() bson.M {
	if p.EqualNil {
		return nil
	} else {
		return p.StringPattern.AsMongoDriverQuery()
	}
}

// Matches returns true if the value is matched by the pattern. If the pattern
// contains a regular expression with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p StringRefPattern) Matches(value *string, matches *RegexMatches) bool {
	if value == nil {
		return p.Empty()
	} else if p.EqualNil {
		return false
	} else {
		return p.StringPattern.Matches(*value, matches)
	}
}

func (p StringRefPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.EqualNil {
		return bsontype.Null, []byte{}, nil
	}

	return p.StringPattern.MarshalBSONValue()
}

func (p StringRefPattern) IsSet() bool {
	return p.StringPattern.IsSet() || p.EqualNil
}

func (p *StringRefPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		p.EqualNil = true
	default:
		return p.StringPattern.UnmarshalBSONValue(valueType, b)
	}

	return nil
}
