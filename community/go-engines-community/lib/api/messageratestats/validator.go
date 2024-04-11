package messageratestats

import (
	"github.com/go-playground/validator/v10"
)

func ValidateListRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ListRequest)
	if !r.To.IsZero() && !r.From.IsZero() && r.To.Before(r.From) {
		sl.ReportError(r.To, "To", "To", "gtfield", "From")
	}

	if r.From.IsZero() && r.Interval == IntervalHour {
		sl.ReportError(r.From, "From", "From", "required_if", "Interval hour")
	}

	if r.To.IsZero() && r.Interval == IntervalHour {
		sl.ReportError(r.To, "To", "To", "required_if", "Interval hour")
	}
}
