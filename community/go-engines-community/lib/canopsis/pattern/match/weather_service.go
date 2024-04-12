package match

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"

func ValidateWeatherServicePattern(p pattern.WeatherServicePattern) bool {
	for _, group := range p {
		if len(group) == 0 {
			return false
		}

		for _, v := range group {
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
				return false
			}
		}
	}

	return true
}
