package utils

import (
	"errors"
	"reflect"
	"strings"
)

// GetField returns the value of a field of an object or a map, given the
// field's name as a string.
// Multiple field names may be chained, separated by dots:
//
//	GetField(alarm_value, "State.Value")
//
// An error is returned if the field does not exist.
// If the field is a pointer, it will be dereferenced before being returned.
// Note that GetField cannot return map keys that contain a dot.
func GetField(object interface{}, fieldPath string) (interface{}, error) {
	fieldNames := strings.Split(fieldPath, ".")
	value := reflect.ValueOf(object)

	// Dereference the value if its a non-nil pointer
	if value.Kind() == reflect.Ptr && !value.IsNil() {
		value = value.Elem()
	}

	for _, fieldName := range fieldNames {
		// Try to get the field's value
		switch value.Kind() {
		case reflect.Struct:
			value = value.FieldByName(fieldName)
		case reflect.Map:
			value = value.MapIndex(reflect.ValueOf(fieldName))
		default:
			return nil, errors.New("unable to get field from object that is neither a struct nor a map")
		}

		if !value.IsValid() {
			// There is no such field in the struct or the map
			return nil, errors.New("field does not exist")
		}

		// Dereference the value if it is a non-nil pointer
		if (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface) && !value.IsNil() {
			value = value.Elem()
		} else if value.Kind() == reflect.Ptr && value.IsNil() {
			// The field is a nil pointer, so we cannot dereference it
			return nil, nil

		}
	}

	return value.Interface(), nil
}
