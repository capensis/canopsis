package alarm

import (
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
)

var alarmStepTypes = types.GetAlarmStepTypes()

func ValidateListRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ListRequest)

	if r.SortBy != "" && len(r.MultiSort) != 0 {
		sl.ReportError(r.SortBy, "SortBy", "SortBy", "required_not_both", "MultiSort")
		return
	}

	if len(r.MultiSort) == 0 {
		return
	}

	for _, multiSortValue := range r.MultiSort {
		multiSortData := strings.Split(multiSortValue, ",")
		if len(multiSortData) != 2 {
			sl.ReportError(r.MultiSort, "MultiSort", "MultiSort", "multi_sort_invalid", "MultiSort")
			return
		}

		if multiSortData[1] != common.SortAsc && multiSortData[1] != common.SortDesc {
			sl.ReportError(r.MultiSort, "MultiSort", "MultiSort", "multi_sort_invalid", "MultiSort")
			return
		}
	}
}

func ValidateDetailsRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(DetailsRequest)
	if r.Steps == nil && r.Children == nil {
		sl.ReportError(r.Steps, "Steps", "Steps", "required", "")
		sl.ReportError(r.Children, "Children", "Children", "required", "")
	}

	if r.Steps != nil && r.Steps.Type != "" {
		for _, stepType := range alarmStepTypes {
			if r.Steps.Type == stepType {
				return
			}
		}

		param := strings.Join(alarmStepTypes, " ")
		sl.ReportError(r.Steps.Type, "Steps.Type", "Type", "oneof", param)
	}
}
