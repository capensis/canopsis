package pattern

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

// IntegerConditions is a struct representing a pattern that can be applied to
// the value of a integer field of an event.
// Each field of a IntegerConditions represents a condition that is applied if
// the value of this field is not nil.
// The fields are not defined directly in the IntegerPattern struct to make the
// unmarshalling easier.
type IntegerConditions struct {
	// If Equal is set, the value of a field has to be equal to the value of
	// Equal to be matched by the pattern.
	Equal utils.OptionalInt64

	// If Gt is set, the value of a field has to be greater than the value
	// of Gt to be matched by the pattern.
	Gt utils.OptionalInt64 `bson:">,omitempty"`

	// If Gte is set, the value of a field has to be greater than the value
	// of Gte to be matched by the pattern.
	Gte utils.OptionalInt64 `bson:">=,omitempty"`

	// If Lt is set, the value of a field has to be greater than the value
	// of Lt to be matched by the pattern.
	Lt utils.OptionalInt64 `bson:"<,omitempty"`

	// If Lte is set, the value of a field has to be greater than the value
	// of Lte to be matched by the pattern.
	Lte utils.OptionalInt64 `bson:"<=,omitempty"`

	// When unmarshalling a BSON document, the fields of this document that are
	// not defined in this struct are added to UnexpectedFields.
	// UnexpectedFields should always be empty.
	UnexpectedFields map[string]interface{} `bson:",inline"`
}

// AsMongoQuery returns a mongodb filter from the IntegerConditions for mgo-driver
func (p IntegerConditions) AsMongoQuery() mgobson.M {
	var query mgobson.M
	query = make(mgobson.M)
	if p.Equal.Set {
		query["$eq"] = p.Equal.Value
	}

	if p.Gt.Set {
		query["$gt"] = p.Gt.Value
	}

	if p.Gte.Set {
		query["$gte"] = p.Gte.Value
	}

	if p.Lt.Set {
		query["$lt"] = p.Lt.Value
	}

	if p.Lte.Set {
		query["$lte"] = p.Lte.Value
	}

	return query
}

// AsMongoDriverQuery returns a mongodb filter from the IntegerConditions for mongo-driver
func (p IntegerConditions) AsMongoDriverQuery() mongobson.M {
	var query mongobson.M
	query = make(mongobson.M)
	if p.Equal.Set {
		query["$eq"] = p.Equal.Value
	}

	if p.Gt.Set {
		query["$gt"] = p.Gt.Value
	}

	if p.Gte.Set {
		query["$gte"] = p.Gte.Value
	}

	if p.Lt.Set {
		query["$lt"] = p.Lt.Value
	}

	if p.Lte.Set {
		query["$lte"] = p.Lte.Value
	}

	return query
}

// Matches returns true if the value satisfies each of the conditions defined
// in the IntegerConditions.
func (p IntegerConditions) Matches(value types.CpsNumber) bool {
	if p.Equal.Set {
		if !(int64(value) == p.Equal.Value) {
			return false
		}
	}

	if p.Gt.Set {
		if !(int64(value) > p.Gt.Value) {
			return false
		}
	}

	if p.Gte.Set {
		if !(int64(value) >= p.Gte.Value) {
			return false
		}
	}

	if p.Lt.Set {
		if !(int64(value) < p.Lt.Value) {
			return false
		}
	}

	if p.Lte.Set {
		if !(int64(value) <= p.Lte.Value) {
			return false
		}
	}

	return true
}

// Empty returns true if the none of the conditions have been set.
func (p IntegerConditions) Empty() bool {
	return !(p.Equal.Set || p.Lt.Set || p.Lte.Set || p.Gt.Set || p.Gte.Set)
}

// IntegerPattern is a type representing a pattern that can be applied to the
// value of a field of an event that contains an integer.
type IntegerPattern struct {
	IntegerConditions
}

// SetBSON unmarshals a BSON value into an IntegerPattern.
func (p *IntegerPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementInt32, mgobson.ElementInt64, mgobson.ElementFloat64:
		// The BSON value is an integer. The field should be equal to this
		// integer.
		err := raw.Unmarshal(&p.Equal.Value)
		if err != nil {
			return err
		}
		p.Equal.Set = true
		return nil

	case mgobson.ElementDocument:
		// The BSON value is a document. Parse this document as a
		// IntegerConditions.
		err := raw.Unmarshal(&p.IntegerConditions)
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
		return fmt.Errorf("A pattern on an integer should be an integer or an object with the following optional keys: <, <=, >=, >")
	}
}

func (p IntegerPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.Equal.Set {
		return mongobson.MarshalValue(p.Equal.Value)
	}

	resultBson := mongobson.M{}

	if p.Lt.Set {
		bsonFieldName, err := GetFieldBsonName(p, "Lt", "lt")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Lt.Value
	}

	if p.Lte.Set {
		bsonFieldName, err := GetFieldBsonName(p, "Lte", "lte")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Lte.Value
	}

	if p.Gt.Set {
		bsonFieldName, err := GetFieldBsonName(p, "Gt", "gt")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Gt.Value
	}

	if p.Gte.Set {
		bsonFieldName, err := GetFieldBsonName(p, "Gte", "gte")
		if err != nil {
			return bsontype.Undefined, nil, err
		}

		resultBson[bsonFieldName] = p.Gte.Value
	}

	if len(resultBson) > 0 {
		return mongobson.MarshalValue(resultBson)
	}

	return bsontype.Undefined, nil, nil
}

func (p IntegerPattern) IsSet() bool {
	return p.Equal.Set || p.Lt.Set || p.Gt.Set || p.Lte.Set || p.Gte.Set
}

func (p *IntegerPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Int32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		p.Equal.Value = int64(value)
		p.Equal.Set = true
	case bsontype.Int64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		p.Equal.Value = value
		p.Equal.Set = true
	case bsontype.Double:
		value, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		p.Equal.Value = int64(value)
		p.Equal.Set = true
	case bsontype.EmbeddedDocument:
		err := mongobson.Unmarshal(b, &p.IntegerConditions)
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
		return fmt.Errorf("unable to parse int32, int64 or double")
	}

	return nil
}

// IntegerRefPattern is a type representing a pattern that can be applied to the
// value of a field of an event that contains a reference to an integer.
type IntegerRefPattern struct {
	// If EqualNil is true, the field should be nil or not be set to be matched
	// by the pattern.
	EqualNil bool

	IntegerPattern
}

// Empty returns true if the pattern has not been set
func (p IntegerRefPattern) Empty() bool {
	return !p.EqualNil && p.IntegerPattern.Empty()
}

// AsMongoQuery returns a mongodb filter from the IntegerRefPattern for mgo-driver
func (p IntegerRefPattern) AsMongoQuery() mgobson.M {
	if p.EqualNil {
		return nil
	} else {
		return p.IntegerPattern.AsMongoQuery()
	}
}

// AsMongoDriverQuery returns a mongodb filter from the IntegerRefPattern for mongo-driver
func (p IntegerRefPattern) AsMongoDriverQuery() mongobson.M {
	if p.EqualNil {
		return nil
	} else {
		return p.IntegerPattern.AsMongoDriverQuery()
	}
}

// Matches returns true if the value is matched by the pattern.
func (p IntegerRefPattern) Matches(value *types.CpsNumber) bool {
	if value == nil {
		return p.IntegerPattern.Empty()
	} else if p.EqualNil {
		return false
	} else {
		return p.IntegerPattern.Matches(*value)
	}
}

// SetBSON unmarshals a BSON value into an IntegerRefPattern.
func (p *IntegerRefPattern) SetBSON(raw mgobson.Raw) error {
	switch raw.Kind {
	case mgobson.ElementNil:
		// The BSON value is null. The field should be nil or not be set.
		p.EqualNil = true
		return nil

	default:
		return raw.Unmarshal(&p.IntegerPattern)
	}
}

func (p IntegerRefPattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.EqualNil {
		return bsontype.Null, []byte{}, nil
	}

	return p.IntegerPattern.MarshalBSONValue()
}

func (p IntegerRefPattern) IsSet() bool {
	return p.EqualNil || p.IntegerPattern.IsSet()
}

func (p *IntegerRefPattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		p.EqualNil = true
	default:
		return p.IntegerPattern.UnmarshalBSONValue(valueType, b)
	}

	return nil
}
