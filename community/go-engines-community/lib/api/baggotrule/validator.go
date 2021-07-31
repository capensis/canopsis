package baggotrule

import (
	"github.com/go-playground/validator/v10"
)

type baggotRuleValidator struct {
}

func (v *baggotRuleValidator) Validate(sl validator.StructLevel) {
	var payload = sl.Current().Interface().(Payload)

	if payload.PRaw.AlarmPatterns != nil && sl.Validator().Var(payload.PRaw.AlarmPatterns, "alarmpatterns") != nil {
		sl.ReportError(payload.AlarmPatterns, "alarm_patterns", "AlarmPattern",
			"alarmpatterns", "")
		return
	}

	if (payload.AlarmPatterns != nil && !payload.AlarmPatterns.IsValid()) || (payload.PRaw.AlarmPatterns != nil && payload.AlarmPatterns == nil) {
		sl.ReportError(payload.AlarmPatterns, "alarm_patterns", "AlarmPatterns",
			"alarmpatterns", "")
	}

	if (payload.EntityPatterns != nil && !payload.EntityPatterns.IsValid()) || (payload.PRaw.EntityPatterns != nil && payload.EntityPatterns == nil) {
		sl.ReportError(payload.EntityPatterns, "entity_patterns", "EntityPatterns",
			"entitypatterns", "")
	}
}

func NewBaggotRuleValidator() *baggotRuleValidator {
	return &baggotRuleValidator{}
}
