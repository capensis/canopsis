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

	if r.Alarm.ArchiveAfter != nil && r.Alarm.DeleteAfter != nil &&
		r.Alarm.ArchiveAfter.Seconds >= r.Alarm.DeleteAfter.Seconds {
		sl.ReportError(r.Remediation.DeleteAfter, "Alarm.DeleteAfter", "DeleteAfter", "gtfield", "ArchiveAfter")
	}

	if r.Alarm.ArchiveAfter == nil && r.Alarm.DeleteAfter != nil {
		sl.ReportError(r.Remediation.DeleteAfter, "Alarm.ArchiveAfter", "ArchiveAfter", "required_if", "DeleteAfter")
	}
}
