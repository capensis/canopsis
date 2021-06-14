package pattern

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
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
	Equal utils.OptionalString

	// If RegexMatch is set, the value of a field has to be matched by this
	// regular expression to be matched by the pattern.
	RegexMatch utils.OptionalRegexp `bson:"regex_match,omitempty"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	// UnexpectedFields should always be empty.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// AsMongoQuery returns a mongodb filter from the StringConditions for mgo-driver
func (p StringConditions) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)

	if p.Equal.Set {
		query["$eq"] = p.Equal.Value
	}

	if p.RegexMatch.Set {
		query["$regex"] = p.RegexMatch.Value.String()
	}

	return query
}

// AsMongoDriverQuery returns a mongodb filter from the StringConditions from mongo-driver
func (p StringConditions) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)

	if p.Equal.Set {
		query["$eq"] = p.Equal.Value
	}

	if p.RegexMatch.Set {
		query["$regex"] = p.RegexMatch.Value.String()
	}

	return query
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

// SetBSON unmarshals a BSON value into a StringPattern.
func (p *StringPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementString:
		// The BSON value is a string. The field should be equal to this
		// string.
		var s string
		err := raw.Unmarshal(&s)
		if err != nil {
			return err
		}
		p.Equal.Set = true
		p.Equal.Value = s
		return nil

	case mgobson.ElementDocument:
		// The BSON value is a document. Parse this document as a
		// StringConditions.
		err := raw.Unmarshal(&p.StringConditions)
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
		return nil

	default:
		return fmt.Errorf("A pattern on a string should be a string or an object with the following optional keys: regex_match")
	}
}

func (p StringPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.Equal.Set {
		return mongobson.MarshalValue(p.Equal.Value)
	}

	if p.RegexMatch.Set {
		bsonFieldName, err := GetFieldBsonName(p, "RegexMatch", "regexmatch")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		return mongobson.MarshalValue(mongobson.M{
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
		err := mongobson.Unmarshal(b, &p.StringConditions)
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

// AsMongoQuery returns a mongodb filter from the StringRefPattern for mgo-driver
func (p StringRefPattern) AsMongoQuery() mgobson.M {
	if p.EqualNil {
		return nil
	} else {
		return p.StringPattern.AsMongoQuery()
	}
}

// AsMongoQuery returns a mongodb filter from the StringRefPattern for mongo-driver
func (p StringRefPattern) AsMongoDriverQuery() mongobson.M {
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

// SetBSON unmarshals a BSON value into a StringRefPattern.
func (p *StringRefPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		// The BSON value is null. The field should be nil or not be set.
		p.EqualNil = true
		return nil

	default:
		return raw.Unmarshal(&p.StringPattern)
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
