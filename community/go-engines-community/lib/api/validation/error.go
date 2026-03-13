package validation

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var reBrackets = regexp.MustCompile(`\[([^\]]+)\]`)

func NewError(errors validator.ValidationErrors, validatedStruct any) *Error {
	return &Error{
		errors:            errors,
		rvValidatedStruct: reflect.ValueOf(validatedStruct),
	}
}

func NewSingleError(tag, field, namespace string, validatedStruct any) *Error {
	return NewError(
		validator.ValidationErrors{
			NewFieldError(tag, field, namespace),
		},
		validatedStruct,
	)
}

func NewSingleErrorWithParam(tag, field, namespace, param string, validatedStruct any) *Error {
	return NewError(
		validator.ValidationErrors{
			NewFieldErrorWithParam(tag, field, namespace, param),
		},
		validatedStruct,
	)
}

type Error struct {
	errors validator.ValidationErrors
	// rvValidatedStruct is used to transform struct namespace to json|form tag namespace.
	// For example, "EditRequest.Duration.Value" to "duration.value"
	// There is no such possibility in github.com/go-playground/validator library.
	rvValidatedStruct reflect.Value
}

func (e *Error) Error() string {
	return e.errors.Error()
}

// TransformNamespace transforms struct namespace to json|form tag namespace.
// for example:
// - Username -> username
// - Items[0] -> items.0
// - Items[0].Name -> items.0.name
// - EditRequest.Duration.Value -> duration.value
func (e *Error) TransformNamespace(ns string) string {
	if ns == "" || !e.rvValidatedStruct.IsValid() {
		return ns
	}

	// remove brackets
	ns = reBrackets.ReplaceAllString(ns, ".$1")
	// replace name to json tag name
	path := strings.Split(ns, ".")
	rv := e.rvValidatedStruct
loop:
	for i := range path {
		rk := rv.Kind()

		switch rk {
		case reflect.Interface, reflect.Pointer:
			rv = rv.Elem()
			rk = rv.Kind()
		}

		switch rk {
		case reflect.Struct:
			rt := rv.Type()
			if path[i] == rt.Name() {
				path[i] = ""
				continue
			}

			if f, ok := rt.FieldByName(path[i]); ok {
				tag := f.Tag.Get("json")
				if tag == "" {
					tag = f.Tag.Get("form")
				}

				tags := strings.Split(tag, ",")
				if len(tags) > 1 && tags[len(tags)-1] == "omitempty" {
					tag = strings.Join(tags[:len(tags)-1], ",")
				}

				if tag == "-" {
					tag = ""
				}

				rv = rv.FieldByName(path[i])
				path[i] = tag
			} else {
				path[i] = ""
			}
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(path[i])
			if err != nil {
				break loop
			}

			rv = rv.Index(index)
		case reflect.Map:
			kv := reflect.ValueOf(path[i])
			if !kv.Type().AssignableTo(rv.Type().Key()) {
				break loop
			}

			rv = rv.MapIndex(kv)
		default:
			path[i] = ""
		}
	}

	k := 0
	for _, p := range path {
		if p != "" {
			path[k] = p
			k++
		}
	}

	path = path[:k]

	return strings.Join(path, ".")
}
