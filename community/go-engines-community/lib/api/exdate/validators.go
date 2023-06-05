package exdate

import (
	"github.com/go-playground/validator/v10"
)

func ValidateExdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(Request)

	if r.End.Before(r.Begin) {
		sl.ReportError(r.End, "End", "End", "gtfield", "Begin")
	}
}
