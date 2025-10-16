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
		if !slices.Contains(availableTypes, params.Type) {
			sl.ReportError(params.Type, "ExternalData."+strconv.Itoa(i)+".Type", "Type", "oneof", strings.Join(availableTypes, " "))
			continue
		}

		switch params.Type {
		case externaldata.RefTypeTable:
			if params.Table == "" {
				sl.ReportError(params.Table, "ExternalData."+strconv.Itoa(i)+".Table", "Table", "required", "")
			}

			for k, v := range params.Regexp {
				if v != "" {
					parsedValue := templateExecutor.Parse(v)
					if parsedValue.Err != nil {
						sl.ReportError(v, "ExternalData."+strconv.Itoa(i)+".Regexp."+k, k, "template", "")
					}
				}
			}

			for k, v := range params.Select {
				if v != "" {
					parsedValue := templateExecutor.Parse(v)
					if parsedValue.Err != nil {
						sl.ReportError(v, "ExternalData."+strconv.Itoa(i)+".Select."+k, k, "template", "")
					}
				}
			}
		case externaldata.RefTypeAPI:
			if params.Request == nil {
				sl.ReportError(params.Request, "ExternalData."+strconv.Itoa(i)+".Request", "Request", "required", "")
				return
			}
		}
	}
}
