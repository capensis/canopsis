package externaldatatable

import (
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"github.com/go-playground/validator/v10"
)

func ValidateRefParameters(sl validator.StructLevel, templateExecutor template.Executor, r []externaldata.RefParameters, availableTypes []string) {
	for i, params := range r {
		structNs := "ExternalData." + strconv.Itoa(i)
		if !slices.Contains(availableTypes, params.Type) {
			sl.ReportError(params.Type, "Type", structNs+".Type", "oneof", strings.Join(availableTypes, " "))
			continue
		}

		switch params.Type {
		case externaldata.RefTypeTable:
			if params.Table == "" {
				sl.ReportError(params.Table, "Table", structNs+".Table", "required", "")
			}

			for k, v := range params.Regexp {
				if v != "" {
					parsedValue := templateExecutor.Parse(v)
					if parsedValue.Err != nil {
						sl.ReportError(v, k, structNs+".Regexp."+k, "template", "")
					}
				}
			}

			for k, v := range params.Select {
				if v != "" {
					parsedValue := templateExecutor.Parse(v)
					if parsedValue.Err != nil {
						sl.ReportError(v, k, structNs+".Select."+k, "template", "")
					}
				}
			}
		case externaldata.RefTypeAPI:
			if params.Request == nil {
				sl.ReportError(params.Request, "Request", structNs+".Request", "required", "")
				return
			}
		}
	}
}
