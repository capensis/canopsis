package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"github.com/go-playground/validator/v10"
)

func ValidateEventPattern(fl validator.FieldLevel) bool {
	i := fl.Field().Interface()
	if i == nil {
		return true
	}
	p, ok := i.(pattern.Event)
	if !ok {
		return false
	}

	return match.ValidateEventPattern(p)
}
