package messageratestats

import (
	"github.com/go-playground/validator/v10"
)

func ValidateListRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ListRequest)
	if !r.To.IsZero() && !r.From.IsZero() && r.To.Before(r.From) {
		sl.ReportError(r.To, "To", "To", "gtfield", "From")
	}
}
