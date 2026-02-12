package widget

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/filemask"
	"github.com/go-playground/validator/v10"
)

const WidgetExportColumnsParam = "widgetExportColumns"

type Validator interface {
	ValidateEditRequest(sl validator.StructLevel)
	ValidateFilterRequest(sl validator.StructLevel)
}

func NewValidator(templateExecutor template.Executor) Validator {
	return &baseValidator{
		filterValidator:  widgetfilter.NewValidator(),
		templateExecutor: templateExecutor,
	}
}

type baseValidator struct {
	filterValidator  *widgetfilter.Validator
	templateExecutor template.Executor
}

func (v *baseValidator) ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	switch r.Type {
	case view.WidgetTypeJunit:
		validateJunitParametersRequest(sl, r.Parameters)
	case view.WidgetTypeMap:
		validateMapParametersRequest(sl, r.Parameters)
	case view.WidgetTypeExternalData:
		validateExternalDataParametersRequest(sl, r.Parameters)
	}
	validateTemplateParametersRequest(sl, r)

	if exportColumnsRaw, ok := r.Parameters.RemainParameters[WidgetExportColumnsParam]; ok {
		if exportColumns, ok := exportColumnsRaw.([]any); ok {
			for idx, columnRaw := range exportColumns {
				column, ok := columnRaw.(map[string]any)
				if !ok {
					continue
				}

				tplRaw, ok := column["template"]
				if !ok {
					continue
				}

				tpl, ok := tplRaw.(string)
				if !ok {
					continue
				}

				if tpl != "" {
					parsedValue := v.templateExecutor.Parse(tpl)
					if parsedValue.Err != nil {
						sl.ReportError(tpl, "template", "Parameters.RemainParameters."+WidgetExportColumnsParam+"."+strconv.Itoa(idx)+".template", "template", "")
					}
				}
			}
		}
	}
}

func (v *baseValidator) ValidateFilterRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(FilterRequest)
	v.filterValidator.ValidatePatterns(sl, r.BaseEditRequest)
}

func validateJunitParametersRequest(sl validator.StructLevel, r view.Parameters) {
	isAPI := r.IsAPI

	if r.Directory == "" {
		if !isAPI {
			sl.ReportError(r.Directory, "Directory", "Parameters.Directory", "required", "")
		}
	} else if isAPI {
		sl.ReportError(r.Directory, "Directory", "Parameters.Directory", "must_be_empty", "")
	}

	if len(r.ScreenshotDirectories) > 0 && isAPI {
		sl.ReportError(r.ScreenshotDirectories, "ScreenshotDirectories", "Parameters.ScreenshotDirectories", "must_be_empty", "")
	}

	if len(r.VideoDirectories) > 0 && isAPI {
		sl.ReportError(r.VideoDirectories, "VideoDirectories", "Parameters.VideoDirectories", "must_be_empty", "")
	}

	if r.ScreenshotFilemask != "" {
		_, err := filemask.NewFileMask(r.ScreenshotFilemask)
		if err != nil {
			sl.ReportError(r.ScreenshotFilemask, "ScreenshotFilemask", "Parameters.ScreenshotFilemask", "filemask", "")
		}
	}

	if r.VideoFilemask != "" {
		_, err := filemask.NewFileMask(r.VideoFilemask)
		if err != nil {
			sl.ReportError(r.VideoFilemask, "VideoFilemask", "Parameters.VideoFilemask", "filemask", "")
		}
	}

	if r.ReportFileRegexp != "" {
		re, err := regexp.Compile(r.ReportFileRegexp)
		if err != nil || re.SubexpIndex(view.JunitReportFileRegexpSubexpName) < 0 {
			sl.ReportError(r.ReportFileRegexp, "ReportFileRegexp", "Parameters.ReportFileRegexp", "regexp", "")
		}
	}
}

func validateMapParametersRequest(sl validator.StructLevel, r view.Parameters) {
	if r.Map == "" {
		sl.ReportError(r.Map, "Map", "Parameters.Map", "required", "")
	}
}

func validateExternalDataParametersRequest(sl validator.StructLevel, r view.Parameters) {
	if r.Table == "" {
		sl.ReportError(r.Table, "Table", "Parameters.Table", "required", "")
	}
}

func validateTemplateParametersRequest(sl validator.StructLevel, r EditRequest) {
	widgetParametersByType := view.GetWidgetTemplateParameters()[r.Type]
	for tplType, widgetParameters := range widgetParametersByType {
		for _, parameter := range widgetParameters {
			parameters := r.Parameters.RemainParameters
			key := parameter
			parts := strings.Split(parameter, ".")
			if len(parts) > 1 {
				key = parts[len(parts)-1]
				var ok bool
				for i := 0; i < len(parts)-1; i++ {
					parameters, ok = parameters[parts[i]].(map[string]any)
					if !ok {
						break
					}
				}
				if !ok {
					continue
				}
			}

			tplId, ok := parameters[key+"Template"].(string)
			if ok && tplId != "" {
				continue
			}

			switch tplType {
			case view.WidgetTemplateTypeAlarmColumns,
				view.WidgetTemplateTypeEntityColumns:
			default:
				continue
			}

			columns, ok := parameters[key].([]any)
			if !ok {
				continue
			}

			for i, column := range columns {
				if m, ok := column.(map[string]any); ok {
					val, _ := m["value"].(string)
					structNs := fmt.Sprintf("Parameters.RemainParameters.%s.%d.value", parameter, i)
					if val == "" {
						sl.ReportError(val, "value", structNs, "required", "")
					}
				}
			}
		}
	}
}
