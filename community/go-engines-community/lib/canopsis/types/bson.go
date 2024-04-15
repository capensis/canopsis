package types

import (
	"errors"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type OptionalStringArray struct {
	Set   bool
	Value []string
}

func (a *OptionalStringArray) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	if valueType != bson.TypeArray {
		return errors.New("unable to parse array")
	}
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
		stringValue, ok := v.StringValueOK()
		if !ok {
			return errors.New("unable to parse string element")
		}

		a.Value = append(a.Value, stringValue)
	}

	a.Set = len(a.Value) > 0

	return nil
}

func (a OptionalStringArray) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if a.Set {
		return bson.MarshalValue(a.Value)
	}

	return bson.TypeUndefined, nil, nil
}

// OptionalInt64 is a wrapper around int64 that implements the bson.Setter
// interface.
//
// Using this type instead of int64 in a struct allows to :
//   - check whether the value was set or not in the bson document.
//   - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not an integer, UnmarshalBSONValue
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the UnmarshalBSONValue method of an ancestor, calls to MongoDB
// queries may fail.
type OptionalInt64 struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the int64. It should only be taken into
	// account if Set is true.
	Value int64
}

func (i *OptionalInt64) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeInt32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		i.Value = int64(value)
	case bson.TypeInt64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		i.Value = value
	case bson.TypeDouble:
		value, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}
		if intVal := int64(value); value == float64(intVal) {
			i.Value = intVal
		} else {
			return errors.New("invalid value, double is not whole number")
		}

	default:
		return errors.New("unable to parse integer")
	}

	i.Set = true
	return nil
}

func (i OptionalInt64) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if i.Set {
		return bson.MarshalValue(i.Value)
	}

	return bson.TypeUndefined, nil, nil
}

// OptionalBool is a wrapper around bool that implements the bson.Setter
// interface.
//
// Using this type instead of bool in a struct allows to :
//   - check whether the value was set or not in the bson document.
//   - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not a bool, UnmarshalBSONValue
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the UnmarshalBSONValue method of an ancestor, calls to MongoDB
// queries may fail.
type OptionalBool struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the bool. It should only be taken into
	// account if Set is true.
	Value bool
}

func (s *OptionalBool) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeBoolean:
		value, _, ok := bsoncore.ReadBoolean(b)
		if !ok {
			return errors.New("invalid value, expected bool")
		}

		s.Value = value
	default:
		return errors.New("unable to parse bool")
	}

	s.Set = true
	return nil
}

func (s OptionalBool) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if s.Set {
		return bson.MarshalValue(s.Value)
	}

	return bson.TypeUndefined, nil, nil
}

// OptionalString is a wrapper around string that implements the bson.Setter
// interface.
//
// Using this type instead of string in a struct allows to :
//   - check whether the value was set or not in the bson document.
//   - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not a string, UnmarshalBSONValue
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the UnmarshalBSONValue method of an ancestor, calls to MongoDB
// queries may fail.
type OptionalString struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the string. It should only be taken into
	// account if Set is true.
	Value string
}

func (s *OptionalString) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeString:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		s.Value = value
	default:
		return errors.New("unable to parse string")
	}

	s.Set = true
	return nil
}

func (s OptionalString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if s.Set {
		return bson.MarshalValue(s.Value)
	}

	return bson.TypeUndefined, nil, nil
}

// OptionalInterface is a wrapper around interface{} that implements the
// bson.Setter interface.
//
// Using this type instead of interface{} in a struct allows to check whether
// the value was set or not in the bson document.
type OptionalInterface struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the interface{}. It should only be taken
	// into account if Set is true.
	Value interface{}
}

func (i *OptionalInterface) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeInt32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		i.Value = int64(value)
	case bson.TypeInt64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		i.Value = value
	case bson.TypeDouble:
		value, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		i.Value = int64(value)
	case bson.TypeBoolean:
		value, _, ok := bsoncore.ReadBoolean(b)
		if !ok {
			return errors.New("invalid value, expected bool")
		}

		i.Value = value
	case bson.TypeString:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		i.Value = value
	default:
		return fmt.Errorf("unsupported type = %s", valueType.String())
	}

	i.Set = true
	return nil
}

func (i OptionalInterface) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if i.Set {
		return bson.MarshalValue(i.Value)
	}

	return bson.TypeUndefined, nil, nil
}

// OptionalRegexp is a wrapper around regexp.Regexp that implements the
// bson.Setter interface.
//
// Using this type in a struct allows to :
//   - check whether the value was set or not in the bson document.
//   - automatically compile a regular expression.
//   - raise an error when trying to unmarshal a value that is not a valid
//     regular expression.
//
// Note that when trying to unmarshal a value that is not a valid regular
// expression, UnmarshalBSONValue will raise an error that will not be handled by
// bson.Unmarshal. If this error is not handled in the UnmarshalBSONValue method of an
// ancestor, calls to MongoDB queries may fail.
type OptionalRegexp struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the regular expression. It should only be
	// taken into account if Set is true.
	Value utils.RegexExpression
}

func (r *OptionalRegexp) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bson.TypeString:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("unable to parse regular expression")
		}

		if re, err := utils.NewRegexExpression(value); err != nil {
			return err
		} else {
			r.Value = re
		}
	default:
		return errors.New("unable to parse regular expression")
	}

	r.Set = true
	return nil
}

func (r OptionalRegexp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if r.Set {
		return bson.MarshalValue(r.Value)
	}

	return bson.TypeUndefined, nil, nil
}
