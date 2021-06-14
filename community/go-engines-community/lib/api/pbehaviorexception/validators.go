package pbehaviorexception

import (
	"github.com/go-playground/validator/v10"
)

func ValidateExdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ExdateRequest)

	if r.End.Before(r.Begin.Time) {
		sl.ReportError(r.End, "End", "End", "gtfield", "Begin")
	}
}
