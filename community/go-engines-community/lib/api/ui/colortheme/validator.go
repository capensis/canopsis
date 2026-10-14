package colortheme

import "github.com/go-playground/validator/v10"

func ValidateColors(sl validator.StructLevel) {
	r := sl.Current().Interface().(Colors)
	if r.Widget.BorderRadius != nil && *r.Widget.BorderRadius < 0 {
		sl.ReportError(r.Widget.BorderRadius, "Widget.BorderRadius", "Widget.BorderRadius", "gte", "0")
	}
}
