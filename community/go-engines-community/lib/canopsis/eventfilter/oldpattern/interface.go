package oldpattern

import (
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// InterfacePattern is a type representing a pattern that can be applied to the
// value of a field of an event that contains an interface{}. It only allows to
// define patterns on integers and strings.
type InterfacePattern struct {
	// If EqualNil is true, the field should be nil or not be set to be matched
	// by the pattern.
	EqualNil bool

	IntegerConditions
	StringConditions
	StringArrayConditions
}

func (p *InterfacePattern) AsMongoDriverQuery() bson.M {
	var query bson.M
	query = make(bson.M)
	if p.EqualNil {
		return nil
	} else {
		if !p.IntegerConditions.Empty() && !p.StringConditions.Empty() ||
			(!p.IntegerConditions.Empty() && !p.StringArrayConditions.Empty()) ||
			(!p.StringArrayConditions.Empty() && !p.StringConditions.Empty()) {
			// Returns a mongo query that always fails, as the conditions are conflicting
			query["$in"] = nil
		} else if !p.IntegerConditions.Empty() {
			query = p.IntegerConditions.AsMongoDriverQuery()
		} else if !p.StringConditions.Empty() {
			query = p.StringConditions.AsMongoDriverQuery()
		} else if !p.StringArrayConditions.Empty() {
			query = p.StringArrayConditions.AsMongoDriverQuery()
		}
	}
	return query
}

func (p InterfacePattern) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if p.EqualNil {
		return bsontype.Null, []byte{}, nil
	}

	iPattern := IntegerPattern{IntegerConditions: p.IntegerConditions}
	marshalledType, marshalledData, err := iPattern.MarshalBSONValue()
	if err != nil {
		return bsontype.Undefined, nil, err
	}

	if marshalledType != bsontype.Undefined {
		return marshalledType, marshalledData, err
	}

	sPattern := StringPattern{StringConditions: p.StringConditions}
	marshalledType, marshalledData, err = sPattern.MarshalBSONValue()
	if err != nil {
		return bsontype.Undefined, nil, err
	}

	if marshalledType != bsontype.Undefined {
		return marshalledType, marshalledData, err
	}

	saPattern := StringArrayPattern{StringArrayConditions: p.StringArrayConditions}
	marshalledType, marshalledData, err = saPattern.MarshalBSONValue()
	if err != nil {
		return bsontype.Undefined, nil, err
	}

	if marshalledType != bsontype.Undefined {
		return marshalledType, marshalledData, err
	}

	return bsontype.Undefined, nil, nil
}

func (p *InterfacePattern) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Null:
		p.EqualNil = true
	case bsontype.Int32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		p.IntegerConditions.Equal.Value = int64(value)
		p.IntegerConditions.Equal.Set = true
	case bsontype.Int64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		p.IntegerConditions.Equal.Value = value
		p.IntegerConditions.Equal.Set = true
	case bsontype.String:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		p.StringConditions.Equal.Value = value
		p.StringConditions.Equal.Set = true
	case bsontype.EmbeddedDocument:
		err := bson.Unmarshal(b, &p.IntegerConditions)
		if err != nil {
			return err
		}
		err = bson.Unmarshal(b, &p.StringConditions)
		if err != nil {
			return err
		}
		err = bson.Unmarshal(b, &p.StringArrayConditions)
		if err != nil {
			return err
		}

		if len(p.IntegerConditions.UnexpectedFields) != 0 &&
			len(p.StringConditions.UnexpectedFields) != 0 &&
			len(p.StringArrayConditions.UnexpectedFields) != 0 {
			return UnexpectedFieldsError{
				Err: fmt.Errorf("interface patterns should only contain conditions on a string (regex_match) or conditions on an integer (>, <, >=, <=) or conditions on a string array (has_every, has_one_of, has_not, is_empty); those 4 types of conditions cannot be mixed"),
			}
		}
		return nil
	default:
		return fmt.Errorf("a pattern on a interface should be an integer, a string or an object")
	}

	return nil
}

// Matches returns true if the value is matched by the pattern. If the pattern
// contains a regular expression with sub-expressions, the values of the
// sub-expressions are written in the matches argument.
func (p InterfacePattern) Matches(value interface{}, matches *RegexMatches) bool {
	if value == nil {
		return p.IntegerConditions.Empty() && p.StringConditions.Empty() &&
			(p.StringArrayConditions.Empty() || p.StringArrayConditions.OnlyHasNotCondition() || p.StringArrayConditions.OnlyIsEmpty())
	}
	if p.EqualNil {
		return false
	}

	intValue, success := types.AsInteger(value)
	if success {
		return p.StringConditions.Empty() && p.StringArrayConditions.Empty() && p.IntegerConditions.Matches(types.CpsNumber(intValue))
	}

	stringValue, success := value.(string)
	if success {
		return p.IntegerConditions.Empty() && p.StringArrayConditions.Empty() && p.StringConditions.Matches(stringValue, matches)
	}

	arrayValue, err := types.InterfaceToStringSlice(value)
	if err != nil {
		log.Error().Err(err)
		return false
	}
	return p.StringConditions.Empty() && p.IntegerConditions.Empty() && p.StringArrayConditions.Matches(arrayValue)
}
