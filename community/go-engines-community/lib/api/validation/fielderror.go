package validation

import (
	"reflect"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewFieldError(tag, field, namespace string) validator.FieldError {
	return &fieldError{
		tag:   tag,
		field: field,
		ns:    namespace,
		param: "",
		value: nil,
	}
}

func NewFieldErrorWithParam(tag, field, namespace, param string) validator.FieldError {
	return &fieldError{
		tag:   tag,
		field: field,
		ns:    namespace,
		param: param,
		value: nil,
	}
}

type fieldError struct {
	tag   string
	field string
	ns    string
	param string
	value any
}

func (e *fieldError) Tag() string {
	return e.tag
}

func (e *fieldError) ActualTag() string {
	return e.tag
}

func (e *fieldError) Namespace() string {
	return e.ns
}

func (e *fieldError) StructNamespace() string {
	return e.ns
}

func (e *fieldError) Field() string {
	return e.field
}

func (e *fieldError) StructField() string {
	return e.field
}

func (e *fieldError) Value() any {
	return e.value
}

func (e *fieldError) Param() string {
	return e.param
}

func (e *fieldError) Kind() reflect.Kind {
	return reflect.TypeOf(e.value).Kind()
}

func (e *fieldError) Type() reflect.Type {
	return reflect.TypeOf(e.value)
}

func (e *fieldError) Translate(ut ut.Translator) string {
	if ut != nil {
		if t, err := ut.T(e.Tag(), e.Field(), e.Param()); err == nil {
			return t
		}
	}

	return e.Error()
}

func (e *fieldError) Error() string {
	return "validation for " + e.field + " failed on the " + e.tag + " tag"
}
