package view

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"reflect"
	"time"
)

const (
	WidgetTypeJunit = "Junit"
)

const (
	WidgetParamJunitIsApi              = "is_api"
	WidgetParamJunitDir                = "directory"
	WidgetParamJunitScreenshotDirs     = "screenshot_directories"
	WidgetParamJunitScreenshotFilemask = "screenshot_filemask"
	WidgetParamJunitVideoDirs          = "video_directories"
	WidgetParamJunitVideoFilemask      = "video_filemask"
	WidgetInternalParamJunitTestSuites = "test_suites"
)

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

func (w Widget) GetDurationParameter(k string, defaultVal time.Duration) time.Duration {
	if v, ok := w.Parameters[k]; ok {
		if d, ok := v.(types.DurationWithUnit); ok {
			if d.Seconds > 0 {
				return time.Duration(d.Seconds) * time.Second
			}

			return defaultVal
		}

		typeVal := reflect.ValueOf(v)
		var seconds int64

		switch typeVal.Kind() {
		case reflect.Map:
			for _, mk := range typeVal.MapKeys() {
				if mk.String() == "seconds" {
					sv := typeVal.MapIndex(mk).Interface()
					switch value := sv.(type) {
					case float64:
						seconds = int64(value)
					case float32:
						seconds = int64(value)
					case int64:
						seconds = value
					case int32:
						seconds = int64(value)
					case int:
						seconds = int64(value)
					}

					break
				}
			}

			if seconds > 0 {
				return time.Duration(seconds) * time.Second
			}
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
