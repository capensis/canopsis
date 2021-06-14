package serviceweather

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ListRequest)
	// Validate sort
	if r.Sort != "" {
		sorts := []string{
			common.SortAsc,
			common.SortDesc,
		}

		found := false
		for _, sort := range sorts {
			if sort == r.Sort {
				found = true
			}
		}

		if !found {
			param := strings.Join(sorts, " ")
			sl.ReportError(r.Sort, "Sort", "sort", "oneof", param)
		}
	}
}
