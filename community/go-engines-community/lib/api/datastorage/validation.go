package datastorage

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datastorage"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

func ValidateConfig(sl validator.StructLevel) {
	r := sl.Current().Interface().(datastorage.Config)

	if !durationGt(r.Remediation.DeleteStatsAfter, r.Remediation.DeleteAfter) {
		sl.ReportError(r.Remediation.DeleteStatsAfter, "Remediation.DeleteStatsAfter", "DeleteStatsAfter", "gtfield", "DeleteAfter")
	}

	if !durationGt(r.Remediation.DeleteModStatsAfter, r.Remediation.DeleteAfter) {
		sl.ReportError(r.Remediation.DeleteModStatsAfter, "Remediation.DeleteModStatsAfter", "DeleteModStatsAfter", "gtfield", "DeleteAfter")
	}

	if !durationGt(r.Remediation.DeleteModStatsAfter, r.Remediation.DeleteStatsAfter) {
		sl.ReportError(r.Remediation.DeleteModStatsAfter, "Remediation.DeleteModStatsAfter", "DeleteModStatsAfter", "gtfield", "DeleteStatsAfter")
	}

	if !durationGt(r.Alarm.DeleteAfter, r.Alarm.ArchiveAfter) {
		sl.ReportError(r.Alarm.DeleteAfter, "Alarm.DeleteAfter", "DeleteAfter", "gtfield", "ArchiveAfter")
	}

	if types.IsDurationEnabledAndValid(r.Alarm.DeleteAfter) && !types.IsDurationEnabledAndValid(r.Alarm.ArchiveAfter) {
		sl.ReportError(r.Alarm.ArchiveAfter, "Alarm.ArchiveAfter", "ArchiveAfter", "required_if", "DeleteAfter")
	}
}

func durationGt(left, right *types.DurationWithEnabled) bool {
	if left != nil && right == nil {
		return false
	}

	if types.IsDurationEnabledAndValid(left) && types.IsDurationEnabledAndValid(right) {
		now := types.NewCpsTime()
		leftAt := left.AddTo(now)
		rightAt := right.AddTo(now)

		return leftAt.After(rightAt)
	}

	return true
}
