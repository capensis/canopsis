package utils

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/globalsign/mgo/bson"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type OptionalStringArray struct {
	Set   bool
	Value []string
}

func (a *OptionalStringArray) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementArray {
		return fmt.Errorf("unable to parse array")
	}

	err := raw.Unmarshal(&a.Value)
	if err != nil {
		return fmt.Errorf("unable to parse array: %v", err)
	}

	if len(a.Value) == 0 {
		a.Set = false
	} else {
		a.Set = true
	}

	return nil
}

func (a *OptionalStringArray) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Array:
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
			stringValue, ok := v.StringValueOK()
			if !ok {
				return fmt.Errorf("unable to parse string element")
			}

			a.Value = append(a.Value, stringValue)
		}
	default:
		return fmt.Errorf("unable to parse array")
	}

	if len(a.Value) == 0 {
		a.Set = false
	} else {
		a.Set = true
	}

	return nil
}

// OptionalInt64 is a wrapper around int64 that implements the bson.Setter
// interface.
//
// Using this type instead of int64 in a struct allows to :
//  - check whether the value was set or not in the bson document.
//  - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not an integer, SetBSON
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the SetBSON method of an ancestor, calls to MongoDB
// queries (e.g. mgo.Collection.Find) may fail.
type OptionalInt64 struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the int64. It should only be taken into
	// account if Set is true.
	Value int64
}

func (i *OptionalInt64) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementInt32 && raw.Kind != bson.ElementInt64 {
		return fmt.Errorf("unable to parse integer")
	}

	err := raw.Unmarshal(&i.Value)
	if err != nil {
		return fmt.Errorf("unable to parse integer: %v", err)
	}

	i.Set = true
	return nil
}

func (i *OptionalInt64) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Int32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		i.Value = int64(value)
	case bsontype.Int64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		i.Value = value
	case bsontype.Double:
		value, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		i.Value = int64(value)
	default:
		return fmt.Errorf("unable to parse integer")
	}

	i.Set = true
	return nil
}

// OptionalBool is a wrapper around bool that implements the bson.Setter
// interface.
//
// Using this type instead of bool in a struct allows to :
//  - check whether the value was set or not in the bson document.
//  - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not a bool, SetBSON
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the SetBSON method of an ancestor, calls to MongoDB
// queries (e.g. mgo.Collection.Find) may fail.
type OptionalBool struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the bool. It should only be taken into
	// account if Set is true.
	Value bool
}

func (s *OptionalBool) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementBool {
		return fmt.Errorf("unable to parse bool")
	}

	err := raw.Unmarshal(&s.Value)
	if err != nil {
		return fmt.Errorf("unable to parse bool: %v", err)
	}

	s.Set = true
	return nil
}

func (i *OptionalBool) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Boolean:
		value, _, ok := bsoncore.ReadBoolean(b)
		if !ok {
			return errors.New("invalid value, expected bool")
		}

		i.Value = value
	default:
		return fmt.Errorf("unable to parse bool")
	}

	i.Set = true
	return nil
}

// OptionalString is a wrapper around string that implements the bson.Setter
// interface.
//
// Using this type instead of string in a struct allows to :
//  - check whether the value was set or not in the bson document.
//  - raise an error when trying to unmarshal a value that is not an integer.
//
// Note that when trying to unmarshal a value that is not a string, SetBSON
// will raise an error that will not be handled by bson.Unmarshal. If this
// error is not handled in the SetBSON method of an ancestor, calls to MongoDB
// queries (e.g. mgo.Collection.Find) may fail.
type OptionalString struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the string. It should only be taken into
	// account if Set is true.
	Value string
}

func (s *OptionalString) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementString {
		return fmt.Errorf("unable to parse string")
	}

	err := raw.Unmarshal(&s.Value)
	if err != nil {
		return fmt.Errorf("unable to parse string: %v", err)
	}

	s.Set = true
	return nil
}

func (i *OptionalString) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.String:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("invalid value, expected string")
		}

		i.Value = value
	default:
		return fmt.Errorf("unable to parse string")
	}

	i.Set = true
	return nil
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

func (i *OptionalInterface) SetBSON(raw bson.Raw) error {
	err := raw.Unmarshal(&i.Value)
	if err != nil {
		return fmt.Errorf("unable to parse interface{}: %v", err)
	}

	i.Set = true
	return nil
}

func (i *OptionalInterface) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.Int32:
		value, _, ok := bsoncore.ReadInt32(b)
		if !ok {
			return errors.New("invalid value, expected int32")
		}

		i.Value = int64(value)
	case bsontype.Int64:
		value, _, ok := bsoncore.ReadInt64(b)
		if !ok {
			return errors.New("invalid value, expected int64")
		}

		i.Value = value
	case bsontype.Double:
		value, _, ok := bsoncore.ReadDouble(b)
		if !ok {
			return errors.New("invalid value, expected double")
		}

		i.Value = int64(value)
	case bsontype.Boolean:
		value, _, ok := bsoncore.ReadBoolean(b)
		if !ok {
			return errors.New("invalid value, expected bool")
		}

		i.Value = value
	case bsontype.String:
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

// OptionalRegexp is a wrapper around regexp.Regexp that implements the
// bson.Setter interface.
//
// Using this type in a struct allows to :
//  - check whether the value was set or not in the bson document.
//  - automatically compile a regular expression.
//  - raise an error when trying to unmarshal a value that is not a valid
//    regular expression.
//
// Note that when trying to unmarshal a value that is not a valid regular
// expression, SetBSON will raise an error that will not be handled by
// bson.Unmarshal. If this error is not handled in the SetBSON method of an
// ancestor, calls to MongoDB queries (e.g. mgo.Collection.Find) may fail.
type OptionalRegexp struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the regular expression. It should only be
	// taken into account if Set is true.
	Value RegexExpression
}

func (r *OptionalRegexp) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementString {
		return fmt.Errorf("unable to parse regular expression")
	}

	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse regex: %v", err)
	}

	if re, err := NewRegexExpression(s); err != nil {
		return err
	} else {
		r.Value = re
	}

	r.Set = true
	return nil
}

func (r *OptionalRegexp) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	switch valueType {
	case bsontype.String:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("unable to parse regular expression")
		}

		if re, err := NewRegexExpression(value); err != nil {
			return err
		} else {
			r.Value = re
		}
	default:
		return fmt.Errorf("unable to parse regular expression")
	}

	r.Set = true
	return nil
}

// OptionalTemplate is a wrapper around template.Template that implements the
// bson.Setter interface.
//
// Using this type in a struct allows to :
//  - check whether the value was set or not in the bson document.
//  - automatically compile a template.
//  - raise an error when trying to unmarshal a value that is not a valid
//    template.
//
// Note that when trying to unmarshal a value that is not a valid template,
// SetBSON will raise an error that will not be handled by bson.Unmarshal. If
// this error is not handled in the SetBSON method of an ancestor, calls to
// MongoDB queries (e.g. mgo.Collection.Find) may fail.
type OptionalTemplate struct {
	// Set is a boolean indicating whether the value was set or not.
	Set bool

	// Value contains the value of the regular expression. It should only be
	// taken into account if Set is true.
	Value *template.Template
}

func (t *OptionalTemplate) SetBSON(raw bson.Raw) error {
	if raw.Kind != bson.ElementString {
		return fmt.Errorf("unable to parse template")
	}

	var s string
	err := raw.Unmarshal(&s)
	if err != nil {
		return fmt.Errorf("unable to parse template: %v", err)
	}

	t.Value, err = template.New("value").Parse(s)
	if err != nil {
		return fmt.Errorf("unable to parse template: %v", err)
	}

	// This makes the template return an error when a missing key is used,
	// instead of replacing its value with "<no value>". This is necessary for
	// the event-filter's actions, which should fail in this case.
	t.Value.Option("missingkey=error")

	t.Set = true
	return nil
}

func (t *OptionalTemplate) UnmarshalBSONValue(valueType bsontype.Type, b []byte) error {
	var err error
	switch valueType {
	case bsontype.String:
		value, _, ok := bsoncore.ReadString(b)
		if !ok {
			return errors.New("unable to parse template")
		}

		t.Value, err = template.New("value").Parse(value)
		if err != nil {
			return fmt.Errorf("unable to parse template: %v", err)
		}

		// This makes the template return an error when a missing key is used,
		// instead of replacing its value with "<no value>". This is necessary for
		// the event-filter's actions, which should fail in this case.
		t.Value.Option("missingkey=error")
	default:
		return fmt.Errorf("unable to parse template")
	}

	t.Set = true
	return nil
}
