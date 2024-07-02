package goja

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/dop251/goja"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mappingFunc func(v any) error

func transformOptions(vm *goja.Runtime, opts goja.Value, dbOpts any, mappingArg ...map[string]mappingFunc) error {
	var mapping map[string]mappingFunc
	if len(mappingArg) > 0 {
		mapping = mappingArg[0]
	}

	if opts == nil {
		return nil
	}

	vt := opts.ExportType()
	if vt == nil {
		return nil
	}

	if vt.Kind() != reflect.Map {
		return errors.New("invalid type for options: " + vt.String())
	}

	obj := opts.ToObject(vm)
	rv := reflect.ValueOf(dbOpts)
	rt := rv.Type()
	for _, k := range obj.Keys() {
		v, err := transformValue(vm, obj.Get(k))
		if err != nil {
			return fmt.Errorf("invalid option %q: %w", k, err)
		}

		if f, ok := mapping[k]; ok {
			err := f(v)
			if err != nil {
				return err
			}

			continue
		}

		methodName := "Set" + strings.ToUpper(k[:1]) + k[1:]
		if _, ok := rt.MethodByName(methodName); !ok {
			return fmt.Errorf("unknown option %q for %T", k, dbOpts)
		}

		m := rv.MethodByName(methodName)
		if m.Type().NumIn() != 1 {
			return fmt.Errorf("set method for option %q has %d arguments instead of one", k, m.Type().NumIn())
		}

		arg := reflect.ValueOf(v)
		if !arg.Type().AssignableTo(m.Type().In(0)) {
			return fmt.Errorf("set method for option %q accepts %s instead of %s", k, m.Type().In(0).String(), arg.Type().String())
		}

		m.Call([]reflect.Value{arg})
	}

	return nil
}

func transformValue(vm *goja.Runtime, v goja.Value) (any, error) {
	if v == nil {
		return nil, nil
	}

	vt := v.ExportType()
	if vt == nil {
		return nil, nil
	}

	switch vt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String,
		reflect.Bool:

		return v.Export(), nil
	case reflect.Array, reflect.Slice:
		res := bson.A{}
		var err error
		vm.ForOf(v, func(curValue goja.Value) (continueIteration bool) {
			var tv any
			tv, err = transformValue(vm, curValue)
			if err != nil {
				return false
			}

			res = append(res, tv)

			return true
		})
		if err != nil {
			return nil, err
		}

		return res, nil
		//obj := v.ToObject(vm)
		//keys := obj.Keys()
		//res := make(bson.A, len(keys))
		//var err error
		//for i := range keys {
		//	res[i], err = transformValue(vm, obj.Get(strconv.Itoa(i)))
		//	if err != nil {
		//		return nil, fmt.Errorf("invalid item %d: %w", i, err)
		//	}
		//}
		//
		//return res, nil
	case reflect.Map:
		obj := v.ToObject(vm)
		switch obj.ClassName() {
		case "Object":
			keys := obj.Keys()
			res := make(bson.D, len(keys))
			for i, k := range keys {
				tv, err := transformValue(vm, obj.Get(k))
				if err != nil {
					return nil, fmt.Errorf("invalid field %q: %w", k, err)
				}

				res[i] = bson.E{
					Key:   k,
					Value: tv,
				}
			}

			return res, nil
		case "RegExp":
			str := obj.String()
			if len(str) == 0 || str[0] != '/' {
				return nil, errors.New("invalid regular expression: " + str)
			}

			p, o, ok := strings.Cut(str[1:], "/")
			if !ok {
				return nil, errors.New("invalid regular expression: " + str)
			}

			return primitive.Regex{
				Pattern: p,
				Options: o,
			}, nil
		default:
			return nil, errors.New("unsupported object class: " + obj.ClassName())
		}
	case reflect.Struct:
		if t, ok := v.Export().(time.Time); ok {
			return t, nil
		}
	}

	return nil, errors.New("unsupported value type: " + vt.String())
}
