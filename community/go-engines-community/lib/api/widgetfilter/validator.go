package widgetfilter

import (
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" &&
		len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.PbehaviorPattern) == 0 && r.CorporatePbehaviorPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required", "")
		sl.ReportError(r.CorporateAlarmPattern, "CorporateAlarmPattern", "CorporateAlarmPattern", "required", "")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		sl.ReportError(r.CorporateEntityPattern, "CorporateEntityPattern", "CorporateEntityPattern", "required", "")
		sl.ReportError(r.PbehaviorPattern, "PbehaviorPattern", "PbehaviorPattern", "required", "")
		sl.ReportError(r.CorporatePbehaviorPattern, "CorporatePbehaviorPattern", "CorporatePbehaviorPattern", "required", "")
	}
}
