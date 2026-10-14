package commenttemplate

import (
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)

	if len(r.Fields) == 0 {
		sl.ReportError(r.Fields, "Fields", "Fields", "required", "")
	}
}
