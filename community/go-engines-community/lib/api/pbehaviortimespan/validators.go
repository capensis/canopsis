package pbehaviortimespan

import (
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
)

func ValidateTimespansRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(TimespansRequest)

	if r.EndAt != nil && r.EndAt.Before(r.StartAt) {
		sl.ReportError(r.EndAt, "EndAt", "EndAt", "gtfield", "StartAt")
	}

	if r.ViewTo.Before(r.ViewFrom) {
		sl.ReportError(r.ViewTo, "ViewTo", "ViewTo", "gtfield", "ViewFrom")
	}

	if r.RRule != "" {
		_, err := rrule.StrToROption(r.RRule)
		if err != nil {
			sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
		}
	}
}

func ValidateExdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ExdateRequest)

	if r.End.Before(r.Begin) {
		sl.ReportError(r.End, "End", "End", "gtfield", "Begin")
	}
}
