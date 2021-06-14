package pbehaviortype

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	validateType(sl, r.Type)
}

func validateType(sl validator.StructLevel, rType string) {
	if rType == "" {
		return
	}

	switch rType {
	case pbehavior.TypeActive,
		pbehavior.TypeInactive,
		pbehavior.TypeMaintenance,
		pbehavior.TypePause:
	/*do nothing*/
	default:
		types := []string{
			pbehavior.TypeActive,
			pbehavior.TypeInactive,
			pbehavior.TypeMaintenance,
			pbehavior.TypePause,
		}
		param := strings.Join(types, " ")
		sl.ReportError(rType, "Type", "type", "oneof", param)
	}
}
