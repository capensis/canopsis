package alarm

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidateListRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ListRequest)

	if r.SortBy != "" && r.MultiSort != "" {
		sl.ReportError(r.SortBy, "SortBy", "SortBy", "required_not_both", "MultiSort")
		return
	}

	if r.MultiSort == "" {
		return
	}

	multiSortData := strings.Split(r.MultiSort, ",")
	if len(multiSortData)%2 != 0 {
		sl.ReportError(r.MultiSort, "MultiSort", "MultiSort", "multi_sort_invalid", "MultiSort")
		return
	}

	for i := 0; i < len(multiSortData); i += 2 {
		if multiSortData[i+1] != common.SortAsc && multiSortData[i+1] != common.SortDesc {
			sl.ReportError(r.MultiSort, "MultiSort", "MultiSort", "multi_sort_invalid", "MultiSort")
		}
	}
}
