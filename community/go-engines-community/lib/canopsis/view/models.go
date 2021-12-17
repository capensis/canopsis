package view

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"reflect"
)

const (
	WidgetTypeJunit = "Junit"
)

const (
	WidgetParamJunitIsApi              = "is_api"
	WidgetParamJunitDir                = "directory"
	WidgetParamJunitReportFileRegexp   = "report_fileregexp"
	WidgetParamJunitScreenshotDirs     = "screenshot_directories"
	WidgetParamJunitScreenshotFilemask = "screenshot_filemask"
	WidgetParamJunitVideoDirs          = "video_directories"
	WidgetParamJunitVideoFilemask      = "video_filemask"
	WidgetInternalParamJunitTestSuites = "test_suites"
)

const JunitReportFileRegexpSubexpName = "name"

type Group struct {
	ID       string        `bson:"_id"`
	Title    string        `bson:"title"`
	Author   string        `bson:"author"`
	Position int64         `bson:"position"`
	Created  types.CpsTime `bson:"created"`
	Updated  types.CpsTime `bson:"updated"`
}

type View struct {
	ID              string                     `bson:"_id"`
	Enabled         bool                       `bson:"enabled"`
	Title           string                     `bson:"title"`
	Description     string                     `bson:"description"`
	Group           string                     `bson:"group_id"`
	Tabs            []Tab                      `bson:"tabs"`
	Tags            []string                   `bson:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `bson:"periodic_refresh"`
	Author          string                     `bson:"author"`
	Position        int64                      `bson:"position"`
	Created         types.CpsTime              `bson:"created"`
	Updated         types.CpsTime              `bson:"updated"`
}

type Tab struct {
	ID      string   `bson:"_id" json:"_id"`
	Title   string   `bson:"title" json:"title"`
	Widgets []Widget `bson:"widgets" json:"widgets"`
}

type Widget struct {
	ID                 string                 `bson:"_id" json:"_id"`
	Title              string                 `bson:"title" json:"title"`
	Type               string                 `bson:"type" json:"type"`
	GridParameters     map[string]interface{} `bson:"grid_parameters" json:"grid_parameters"`
	Parameters         map[string]interface{} `bson:"parameters" json:"parameters"`
	InternalParameters map[string]interface{} `bson:"internal_parameters,omitempty" json:"-"`
}

func (w Widget) GetStringParameter(k, defaultVal string) string {
	if v, ok := w.Parameters[k]; ok {
		if str, ok := v.(string); ok {
			return str
		}
	}

	return defaultVal
}

func (w Widget) GetBoolParameter(k string, defaultVal bool) bool {
	if v, ok := w.Parameters[k]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}

	return defaultVal
}

func (w Widget) GetStringsParameter(k string, defaultVal []string) []string {
	if v, ok := w.Parameters[k]; ok {
		typeVal := reflect.ValueOf(v)

		switch typeVal.Kind() {
		case reflect.Array, reflect.Slice:
			val := make([]string, 0)
			for i := 0; i < typeVal.Len(); i++ {
				if str, ok := typeVal.Index(i).Interface().(string); ok {
					val = append(val, str)
				}
			}
			return val
		}
	}

	return defaultVal
}

func (w Widget) GetStringsInternalParameter(k string, defaultVal []string) []string {
	if v, ok := w.InternalParameters[k]; ok {
		typeVal := reflect.ValueOf(v)

		switch typeVal.Kind() {
		case reflect.Array, reflect.Slice:
			val := make([]string, 0)
			for i := 0; i < typeVal.Len(); i++ {
				if str, ok := typeVal.Index(i).Interface().(string); ok {
					val = append(val, str)
				}
			}
			return val
		}
	}

	return defaultVal
}
