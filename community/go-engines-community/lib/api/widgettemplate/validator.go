package widgettemplate

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"github.com/go-playground/validator/v10"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	switch r.Type {
	case view.WidgetTemplateTypeAlarmColumns,
		view.WidgetTemplateTypeEntityColumns:
		if len(r.Columns) == 0 {
			sl.ReportError(r.Columns, "Columns", "Columns", "required", "")
		}

		for i, column := range r.Columns {
			if column.Value != "" && !view.IsValidWidgetColumn(r.Type, column.Value) {
				sl.ReportError(column, fmt.Sprintf("Columns.%d.Value", i), "Value", "invalid", "")
			}
		}

		if r.Content != "" {
			sl.ReportError(r.Content, "Content", "Content", "must_be_empty", "")
		}
	case view.WidgetTemplateTypeAlarmMoreInfos,
		view.WidgetTemplateTypeServiceWeatherItem,
		view.WidgetTemplateTypeServiceWeatherModal,
		view.WidgetTemplateTypeServiceWeatherEntity:
		if r.Content == "" {
			sl.ReportError(r.Content, "Content", "Content", "required", "")
		}
		if len(r.Columns) > 0 {
			sl.ReportError(r.Columns, "Columns", "Columns", "must_be_empty", "")
		}
	}
}
