package healthcheck

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/go-playground/validator/v10"
)

func ValidateLimitParameters(sl validator.StructLevel) {
	p := sl.Current().Interface().(config.LimitParameters)

	if p.Enabled && p.Limit < 1 {
		sl.ReportError(p.Limit, "Limit", "Limit", "gt", "0")
	}
}

func ValidateEngineParameters(sl validator.StructLevel) {
	p := sl.Current().Interface().(config.EngineParameters)

	if p.Enabled && p.Minimal < 1 {
		sl.ReportError(p.Minimal, "Minimal", "Minimal", "gt", "0")
	}
}
