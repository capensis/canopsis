package datastorage

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"github.com/go-playground/validator/v10"
)

func ValidateConfig(sl validator.StructLevel) {
	r := sl.Current().Interface().(datastorage.Config)

	if r.Remediation.AccumulateAfter != nil && r.Remediation.DeleteAfter != nil &&
		r.Remediation.AccumulateAfter.Seconds >= r.Remediation.DeleteAfter.Seconds {
		sl.ReportError(r.Remediation.DeleteAfter, "Remediation.DeleteAfter", "DeleteAfter", "gtfield", "AccumulateAfter")
	}
}
