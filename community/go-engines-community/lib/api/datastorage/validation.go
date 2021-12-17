package datastorage

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

func ValidateConfig(sl validator.StructLevel) {
	r := sl.Current().Interface().(datastorage.Config)
	now := types.NewCpsTime()

	if r.Remediation.AccumulateAfter != nil && r.Remediation.DeleteAfter != nil &&
		r.Remediation.AccumulateAfter.Value > 0 && r.Remediation.DeleteAfter.Value > 0 {
		accumulateAt := r.Remediation.AccumulateAfter.AddTo(now)
		deleteAt := r.Remediation.DeleteAfter.AddTo(now)

		if !accumulateAt.Before(deleteAt) {
			sl.ReportError(r.Remediation.DeleteAfter, "Remediation.DeleteAfter", "DeleteAfter", "gtfield", "AccumulateAfter")
		}
	}

	if r.Alarm.ArchiveAfter != nil && r.Alarm.DeleteAfter != nil &&
		r.Alarm.ArchiveAfter.Value > 0 && r.Alarm.DeleteAfter.Value > 0 {
		archiveAt := r.Alarm.ArchiveAfter.AddTo(now)
		deleteAt := r.Alarm.DeleteAfter.AddTo(now)

		if !archiveAt.Before(deleteAt) {
			sl.ReportError(r.Remediation.DeleteAfter, "Alarm.DeleteAfter", "DeleteAfter", "gtfield", "ArchiveAfter")
		}
	}

	if r.Alarm.ArchiveAfter == nil && r.Alarm.DeleteAfter != nil {
		sl.ReportError(r.Remediation.DeleteAfter, "Alarm.ArchiveAfter", "ArchiveAfter", "required_if", "DeleteAfter")
	}
}
