package widgettemplate

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	for i, column := range r.Columns {
		if column.Value != "" && !view.IsValidWidgetColumn(r.Type, column.Value) {
			sl.ReportError(column, fmt.Sprintf("Columns.%d.Value", i), "Value", "invalid", "")
		}
	}
}
