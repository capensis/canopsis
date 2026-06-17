package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"github.com/go-playground/validator/v10"
)

func ValidateAlarmPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.Alarm)

	return ok && match.ValidateAlarmPattern(p, nil, nil)
}

func ValidateEntityPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.Entity)

	return ok && match.ValidateEntityPattern(p, nil)
}

func ValidateEventPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.Event)

	return ok && match.ValidateEventPattern(p)
}

func ValidatePbehaviorPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.PbehaviorInfo)

	return ok && match.ValidatePbehaviorInfoPattern(p)
}

func ValidateWeatherServicePattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}

	p, ok := i.(pattern.WeatherServicePattern)

	return ok && match.ValidateWeatherServicePattern(p)
}
