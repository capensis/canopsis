package alarm

import (
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/go-playground/validator/v10"
)

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
	if r.Steps == nil && len(r.GroupedSteps.IDs) == 0 && r.Children == nil {
		sl.ReportError(r.Steps, "Steps", "Steps", "required", "")
		sl.ReportError(r.GroupedSteps, "GroupedSteps", "GroupedSteps", "required", "")
		sl.ReportError(r.Children, "Children", "Children", "required", "")
	}
}
