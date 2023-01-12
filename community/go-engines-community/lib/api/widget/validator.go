package widget

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widgetfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/filemask"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	ValidateEditRequest(sl validator.StructLevel)
	ValidateFilterRequest(ctx context.Context, sl validator.StructLevel)
}

func NewValidator(dbClient mongo.DbClient) Validator {
	return &baseValidator{
		filterValidator: widgetfilter.NewValidator(dbClient),
	}
}

type baseValidator struct {
	filterValidator *widgetfilter.Validator
}

func (v *baseValidator) ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	switch r.Type {
	case view.WidgetTypeJunit:
		validateJunitParametersRequest(sl, r.Parameters)
	case view.WidgetTypeMap:
		validateMapParametersRequest(sl, r.Parameters)
	}
	validateColumnParametersRequest(sl, r)
}

func (v *baseValidator) ValidateFilterRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(FilterRequest)
	v.filterValidator.ValidatePatterns(ctx, sl, r.BaseEditRequest, r.ID)
}

func validateJunitParametersRequest(sl validator.StructLevel, r view.Parameters) {
	isAPI := r.IsAPI

	if r.Directory == "" {
		if !isAPI {
			sl.ReportError(r.Directory, "parameters.directory", "Directory", "required", "")
		}
	} else if isAPI {
		sl.ReportError(r.Directory, "parameters.directory", "Directory", "must_be_empty", "")
	}

	if len(r.ScreenshotDirectories) > 0 && isAPI {
		sl.ReportError(r.ScreenshotDirectories, "parameters.screenshot_directories", "ScreenshotDirectories", "must_be_empty", "")
	}

	if len(r.VideoDirectories) > 0 && isAPI {
		sl.ReportError(r.VideoDirectories, "parameters.video_directories", "VideoDirectories", "must_be_empty", "")
	}

	if r.ScreenshotFilemask != "" {
		_, err := filemask.NewFileMask(r.ScreenshotFilemask)
		if err != nil {
			sl.ReportError(r.ScreenshotFilemask, "parameters.screenshot_filemask", "ScreenshotFilemask", "filemask", "")
		}
	}

	if r.VideoFilemask != "" {
		_, err := filemask.NewFileMask(r.VideoFilemask)
		if err != nil {
			sl.ReportError(r.VideoFilemask, "parameters.video_filemask", "VideoFilemask", "filemask", "")
		}
	}

	if r.ReportFileRegexp != "" {
		re, err := regexp.Compile(r.ReportFileRegexp)
		if err != nil || re.SubexpIndex(view.JunitReportFileRegexpSubexpName) < 0 {
			sl.ReportError(r.ReportFileRegexp, "parameters.report_fileregexp", "ReportFileRegexp", "regexp", "")
		}
	}
}

func validateMapParametersRequest(sl validator.StructLevel, r view.Parameters) {
	if r.Map == "" {
		sl.ReportError(r.Map, "parameters.map", "Map", "required", "")
	}
}

func validateColumnParametersRequest(sl validator.StructLevel, r EditRequest) {
	columnParametersByType := view.GetWidgetColumnParameters()[r.Type]
	for tplType, columnParameters := range columnParametersByType {
		for _, parameter := range columnParameters {
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

			columns, ok := parameters[key].([]any)
			if !ok {
				continue
			}

			for i, column := range columns {
				if m, ok := column.(map[string]any); ok {
					val, _ := m["value"].(string)
					fieldName := fmt.Sprintf("Parameters.%s.%d.value", parameter, i)
					if val == "" {
						sl.ReportError(column, fieldName, "Value", "required", "")
					} else if !view.IsValidWidgetColumn(tplType, val) {
						sl.ReportError(column, fieldName, "Value", "invalid", "")
					}
				}
			}
		}
	}
}
