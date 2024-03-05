package widgetfilter

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateCreateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CreateRequest)
	v.ValidatePatterns(sl, r.BaseEditRequest)
}

func (v *Validator) ValidateUpdateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(UpdateRequest)
	v.ValidatePatterns(sl, r.BaseEditRequest)
}

func (v *Validator) ValidatePatterns(sl validator.StructLevel, r BaseEditRequest) {
	if len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" &&
		len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.PbehaviorPattern) == 0 && r.CorporatePbehaviorPattern == "" &&
		len(r.WeatherServicePattern) == 0 {

		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required", "")
		sl.ReportError(r.CorporateAlarmPattern, "CorporateAlarmPattern", "CorporateAlarmPattern", "required", "")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		sl.ReportError(r.CorporateEntityPattern, "CorporateEntityPattern", "CorporateEntityPattern", "required", "")
		sl.ReportError(r.PbehaviorPattern, "PbehaviorPattern", "PbehaviorPattern", "required", "")
		sl.ReportError(r.CorporatePbehaviorPattern, "CorporatePbehaviorPattern", "CorporatePbehaviorPattern", "required", "")

		sl.ReportError(r.WeatherServicePattern, "WeatherServicePattern", "WeatherServicePattern", "required", "")
	}

	if !r.WeatherServicePattern.Validate() {
		sl.ReportError(r.WeatherServicePattern, "WeatherServicePattern", "WeatherServicePattern", "weather_service_pattern", "")
	}
}
