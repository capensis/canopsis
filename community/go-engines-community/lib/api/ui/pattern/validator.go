package pattern

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)
	switch r.Type {
	case savedpattern.TypeAlarm:
		if len(r.AlarmPattern) == 0 {
			sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required", "")
		}
	case savedpattern.TypeEntity:
		if len(r.EntityPattern) == 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required", "")
		}
	case savedpattern.TypePbehavior:
		if len(r.PbehaviorPattern) == 0 {
			sl.ReportError(r.PbehaviorPattern, "PbehaviorPattern", "PbehaviorPattern", "required", "")
		}
	case savedpattern.TypeWeatherService:
		if len(r.WeatherServicePattern) == 0 {
			sl.ReportError(r.WeatherServicePattern, "WeatherServicePattern", "WeatherServicePattern", "required", "")
		}
	}
}
