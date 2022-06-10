package resolverule

import (
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}
