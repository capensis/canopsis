package view

import (
	"encoding/json"
	"reflect"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

const (
	WidgetTypeAlarmsList          = "AlarmsList"
	WidgetTypeContextExplorer     = "Context"
	WidgetTypeServiceWeather      = "ServiceWeather"
	WidgetTypeAlarmsCounter       = "Counter"
	WidgetTypeAlarmsStatsCalendar = "StatsCalendar"
	WidgetTypeJunit               = "Junit"
	WidgetTypeMap                 = "Map"

	WidgetInternalParamJunitTestSuites = "test_suites"

	JunitReportFileRegexpSubexpName = "name"
)

const (
	WidgetTemplateTypeAlarmColumns         = "alarm_columns"
	WidgetTemplateTypeEntityColumns        = "entity_columns"
	WidgetTemplateTypeAlarmMoreInfos       = "alarm_more_infos"
	WidgetTemplateTypeAlarmExportToPDF     = "alarm_export_to_pdf"
	WidgetTemplateTypeServiceWeatherItem   = "weather_item"
	WidgetTemplateTypeServiceWeatherModal  = "weather_modal"
	WidgetTemplateTypeServiceWeatherEntity = "weather_entity"
)

type Group struct {
	ID       string           `bson:"_id"`
	Title    string           `bson:"title"`
	Author   string           `bson:"author"`
	Position int64            `bson:"position,omitempty"`
	Created  datetime.CpsTime `bson:"created"`
	Updated  datetime.CpsTime `bson:"updated"`

	IsPrivate bool `bson:"is_private"`
}

type View struct {
	ID              string                        `bson:"_id"`
	Enabled         bool                          `bson:"enabled"`
	Title           string                        `bson:"title"`
	Description     string                        `bson:"description"`
	Group           string                        `bson:"group_id"`
	Tags            []string                      `bson:"tags"`
	PeriodicRefresh *datetime.DurationWithEnabled `bson:"periodic_refresh"`
	Author          string                        `bson:"author"`
	Position        int64                         `bson:"position,omitempty"`
	IsPrivate       bool                          `bson:"is_private"`
	Created         datetime.CpsTime              `bson:"created"`
	Updated         datetime.CpsTime              `bson:"updated"`
}

type Tab struct {
	ID        string           `bson:"_id" json:"_id"`
	Title     string           `bson:"title" json:"title"`
	View      string           `bson:"view,omitempty" json:"-"`
	Author    string           `bson:"author" json:"author"`
	Position  int64            `bson:"position" json:"-"`
	IsPrivate bool             `bson:"is_private" json:"-"`
	Created   datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated   datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type Widget struct {
	ID                 string             `bson:"_id" json:"_id,omitempty"`
	Tab                string             `bson:"tab,omitempty" json:"-"`
	Title              string             `bson:"title" json:"title"`
	Type               string             `bson:"type" json:"type"`
	GridParameters     map[string]any     `bson:"grid_parameters" json:"grid_parameters"`
	Parameters         Parameters         `bson:"parameters" json:"parameters"`
	InternalParameters InternalParameters `bson:"internal_parameters,omitempty" json:"-"`
	Author             string             `bson:"author" json:"author,omitempty"`
	Created            datetime.CpsTime   `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated            datetime.CpsTime   `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	IsPrivate bool `bson:"is_private" json:"is_private"`
}

type Parameters struct {
	MainFilter string `bson:"mainFilter,omitempty" json:"mainFilter,omitempty"`

	// Junit
	IsAPI                 bool     `bson:"is_api,omitempty" json:"is_api,omitempty"`
	Directory             string   `bson:"directory,omitempty" json:"directory,omitempty"`
	ReportFileRegexp      string   `bson:"report_fileregexp,omitempty" json:"report_fileregexp,omitempty"`
	ScreenshotDirectories []string `bson:"screenshot_directories,omitempty" json:"screenshot_directories,omitempty"`
	VideoDirectories      []string `bson:"video_directories,omitempty" json:"video_directories,omitempty"`
	ScreenshotFilemask    string   `bson:"screenshot_filemask,omitempty" json:"screenshot_filemask,omitempty"`
	VideoFilemask         string   `bson:"video_filemask,omitempty" json:"video_filemask,omitempty"`

	// Map
	Map string `bson:"map,omitempty" json:"map,omitempty"`

	RemainParameters map[string]any `bson:",inline" json:"-"`
}

func (p Parameters) MarshalJSON() ([]byte, error) {
	type Alias Parameters
	b, err := json.Marshal(Alias(p))
	if err != nil {
		return nil, err
	}

	m := make(map[string]any)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	for k, v := range p.RemainParameters {
		m[k] = v
	}

	return json.Marshal(m)
}

func (p *Parameters) UnmarshalJSON(b []byte) error {
	type Alias *Parameters
	err := json.Unmarshal(b, Alias(p))
	if err != nil {
		return err
	}
	m := make(map[string]any)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return err
	}

	val := reflect.TypeOf(*p)
	for i := 0; i < val.NumField(); i++ {
		if len(m) == 0 {
			break
		}
		f := val.Field(i)
		tag := f.Tag.Get("json")
		tag = strings.Split(tag, ",")[0]
		delete(m, tag)
	}

	p.RemainParameters = m
	return nil
}

type InternalParameters struct {
	// Junit
	TestSuites []string `bson:"test_suites,omitempty"`

	RemainParameters map[string]any `bson:",inline"`
}

func (p InternalParameters) IsZero() bool {
	return len(p.TestSuites) == 0 &&
		len(p.RemainParameters) == 0
}

type WidgetFilter struct {
	ID        string           `bson:"_id,omitempty"`
	Title     string           `bson:"title"`
	Widget    string           `bson:"widget,omitempty"`
	IsPrivate bool             `bson:"is_private"`
	Author    string           `bson:"author"`
	Position  int64            `bson:"position"`
	Created   datetime.CpsTime `bson:"created,omitempty"`
	Updated   datetime.CpsTime `bson:"updated,omitempty"`

	savedpattern.AlarmPatternFields          `bson:",inline"`
	savedpattern.EntityPatternFields         `bson:",inline"`
	savedpattern.PbehaviorPatternFields      `bson:",inline"`
	savedpattern.WeatherServicePatternFields `bson:",inline"`

	IsUserPreference bool `bson:"is_user_preference"`
}

type WidgetTemplate struct {
	ID      string           `bson:"_id,omitempty"`
	Title   string           `bson:"title"`
	Type    string           `bson:"type"`
	Columns []WidgetColumn   `bson:"columns,omitempty"`
	Content string           `bson:"content,omitempty"`
	Author  string           `bson:"author"`
	Created datetime.CpsTime `bson:"created,omitempty"`
	Updated datetime.CpsTime `bson:"updated,omitempty"`
}

type WidgetColumn struct {
	Value            string `bson:"value," json:"value" binding:"required"`
	Label            string `bson:"label,omitempty" json:"label,omitempty" binding:"max=255"`
	IsHtml           bool   `bson:"isHtml,omitempty" json:"isHtml,omitempty"`
	OnlyIcon         bool   `bson:"onlyIcon,omitempty" json:"onlyIcon,omitempty"`
	ColorIndicator   string `bson:"colorIndicator,omitempty" json:"colorIndicator,omitempty"`
	Template         string `bson:"template,omitempty" json:"template,omitempty"`
	InlineLinksCount int64  `bson:"inlineLinksCount,omitempty" json:"inlineLinksCount,omitempty"`
	LinksInRowCount  int64  `bson:"linksInRowCount,omitempty" json:"linksInRowCount,omitempty"`
}
