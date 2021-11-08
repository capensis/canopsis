package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type Overridable interface {
	Clone() interface{}
}

// Override calls Clone method to get new deep-copy of dest and merge values from src map.
// In case of error, pointer of original dest will be return with error.
func Override(dest Overridable, src map[string]interface{}) (interface{}, error) {
	confVal := reflect.ValueOf(dest)
	if confVal.Type().Kind() != reflect.Ptr {
		return dest, errors.New("dest value must be pointer")
	}

	temp := dest.Clone()
	if reflect.ValueOf(temp).Type().Kind() != reflect.Ptr {
		return dest, errors.New("clone method of dest doesn't return pointer")
	}

	err := OverrideInPlace(temp, src)
	if err != nil {
		return dest, err
	}

	return temp, nil
}

// OverrideInPlace merges values from src to dest. Value updates are made directly on dest.
// When error occurs, values could be not fully merged and might cause inconsistence.
func OverrideInPlace(dest interface{}, src map[string]interface{}) error {
	var dstVal reflect.Value
	var ok bool

	if dstVal, ok = dest.(reflect.Value); !ok {
		dstVal = reflect.ValueOf(dest)
		if dstVal.Type().Kind() == reflect.Ptr {
			dstVal = reflect.Indirect(dstVal)
		}
	}

	for key, val := range src {
		fieldName, fieldPtr := getFieldByTag(dstVal, key)
		if fieldPtr == nil {
			return fmt.Errorf("type %s doesn't have field having toml tag '%s'", dstVal.Type().Name(), key)
		}

		fieldVal := *fieldPtr
		if !fieldVal.CanSet() {
			return fmt.Errorf("field %s of %s cannot be set", fieldVal.Type().Name(), dstVal.Type().Name())
		}

		switch value := val.(type) {
		case map[string]interface{}:
			if err := handleObject(fieldName, fieldVal, value); err != nil {
				return err
			}

		case []map[string]interface{}:
			if err := handleObjectArray(fieldName, fieldVal, value); err != nil {
				return err
			}

		case []interface{}:
			if err := handleArray(fieldName, fieldVal, value); err != nil {
				return err
			}

		case string:
			if fieldVal.Type().Kind() == reflect.String {
				fieldVal.Set(reflect.ValueOf(value))
			} else {
				return fmt.Errorf("value of string is not assignable to field %s of type %s", fieldName, fieldVal.Type().String())
			}
		case int64:
			if fieldVal.Type().Kind() == reflect.Int {
				fieldVal.Set(reflect.ValueOf(int(value)))

			} else {
				return fmt.Errorf("value of int64 is not assignable to field %s of type %s", fieldName, fieldVal.Type().String())
			}
		case bool:
			if fieldVal.Type().Kind() == reflect.Bool {
				fieldVal.Set(reflect.ValueOf(value))
			} else {
				return fmt.Errorf("value of bool is not assignable to field %s of type %s", fieldName, fieldVal.Type().String())
			}
		default:
			return fmt.Errorf("field %s has unhandled type: %s", fieldName, fieldVal.Type().String())
		}
	}

	return nil
}

func getFieldByTag(obj reflect.Value, tag string) (string, *reflect.Value) {
	typ := obj.Type()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Tag.Get("toml") == tag {
			fieldVal := obj.Field(i)
			return typ.Field(i).Name, &fieldVal
		}
	}
	return "", nil
}

// handleObject
func handleObject(fieldName string, fieldVal reflect.Value, value map[string]interface{}) error {
	if fieldVal.Type().Kind() != reflect.Struct {
		return fmt.Errorf("value of map[string]interface{} is not assignable to field %s of type %s", fieldName, fieldVal.Type().String())
	}

	if err := OverrideInPlace(fieldVal, value); err != nil {
		return err
	}

	return nil
}

// handleObjectArray
func handleObjectArray(fieldName string, fieldVal reflect.Value, value []map[string]interface{}) error {
	if fieldVal.Type().Kind() != reflect.Slice {
		return fmt.Errorf("value of []map[string]interface{} is not assignable to field %s of type %s", fieldName, fieldVal.Type().String())
	}

	elemTyp := fieldVal.Type().Elem()
	if elemTyp.Kind() == reflect.Ptr {
		elemTyp = elemTyp.Elem()
	}

	for _, item := range value {
		// initialize new native golang struct
		nativeItem := reflect.New(elemTyp)
		err := mapstructure.Decode(item, nativeItem.Interface())
		if err != nil {
			return fmt.Errorf("failed calling FromMap")
		}

		// nativeItem is a pointer
		nativeItem = reflect.Indirect(nativeItem)

		// append to array
		newArr := reflect.Append(fieldVal, nativeItem)
		fieldVal.Set(newArr)
	}

	return nil
}

// handleArray
func handleArray(fieldName string, fieldVal reflect.Value, value []interface{}) error {
	if len(value) > 0 {
		itemInf := value[0]
		switch itemInf.(type) {
		case string:
			if fieldVal.Type() == reflect.TypeOf([]string{}) {
				stringArr := make([]string, len(value))
				for i := 0; i < len(value); i++ {
					stringArr[i] = value[i].(string)
				}
				fieldVal.Set(reflect.ValueOf(stringArr))
				return nil
			}
		case int64:
			if fieldVal.Type() == reflect.TypeOf([]int{}) {
				intArr := make([]int, len(value))
				for i := 0; i < len(value); i++ {
					intArr[i] = int(value[i].(int64))
				}
				fieldVal.Set(reflect.ValueOf(intArr))
				return nil
			}
		case bool:
			if fieldVal.Type() == reflect.TypeOf([]bool{}) {
				boolArr := make([]bool, len(value))
				for i := 0; i < len(value); i++ {
					boolArr[i] = value[i].(bool)
				}
				fieldVal.Set(reflect.ValueOf(boolArr))
				return nil
			}
		default:
			return fmt.Errorf("field %s has unhandled type: %s", fieldName, fieldVal.Type().String())
		}
	}

	// Set to an empty array
	fieldVal.Set(reflect.Indirect(reflect.New(fieldVal.Type())))
	return nil
}
