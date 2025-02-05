package externaldatatable

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"github.com/go-playground/validator/v10"
)

func ValidateRefParameters(sl validator.StructLevel, r []externaldata.RefParameters, availableTypes []string) {
	for i, params := range r {
		found := false
		for _, t := range availableTypes {
			if params.Type == t {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(params.Type, "ExternalData."+strconv.Itoa(i)+".Type", "Type", "oneof", strings.Join(availableTypes, " "))
			continue
		}

		switch params.Type {
		case externaldata.RefTypeTable:
			if params.Table == "" {
				sl.ReportError(params.Table, "ExternalData."+strconv.Itoa(i)+".Table", "Table", "required", "")
			}
		case externaldata.RefTypeAPI:
			if params.Request == nil {
				sl.ReportError(params.Request, "ExternalData."+strconv.Itoa(i)+".Request", "Request", "required", "")
			}
		}
	}
}
