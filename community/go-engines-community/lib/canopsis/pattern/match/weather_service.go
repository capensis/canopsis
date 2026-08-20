package match

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

func ValidateWeatherServicePattern(p pattern.WeatherServicePattern) bool {
	return len(WeatherServicePatternErrors(p)) == 0
}

func WeatherServicePatternErrors(p pattern.WeatherServicePattern) []ConditionError {
	var errs []ConditionError
	for gidx, group := range p {
		if len(group) == 0 {
			errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: -1, Err: pattern.ErrEmptyGroup})
			continue
		}

		for cidx, v := range group {
			f := v.Field
			cond := v.Condition
			var err error

			switch f {
			case "is_grey":
				_, err = cond.MatchBool(true)
			case "icon", "secondary_icon":
				_, err = cond.MatchString("")
			case "state.val":
				_, err = cond.MatchInt(0)
			default:
				err = pattern.ErrUnsupportedField
			}

			if err != nil {
				errs = append(errs, ConditionError{GroupIdx: gidx, CondIdx: cidx, Err: err})
			}
		}
	}

	return errs
}
