package view

import (
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"reflect"
	"strings"
)

const (
	WidgetTypeJunit = "Junit"

	WidgetInternalParamJunitTestSuites = "test_suites"

	JunitReportFileRegexpSubexpName = "name"
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
	Tags            []string                   `bson:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `bson:"periodic_refresh"`
	Author          string                     `bson:"author"`
	Position        int64                      `bson:"position"`
	Created         types.CpsTime              `bson:"created"`
	Updated         types.CpsTime              `bson:"updated"`
}

type Tab struct {
	ID       string        `bson:"_id" json:"_id"`
	Title    string        `bson:"title" json:"title"`
	View     string        `bson:"view" json:"-"`
	Author   string        `bson:"author" json:"author"`
	Position int64         `bson:"position" json:"-"`
	Created  types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated  types.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type Widget struct {
	ID                 string                 `bson:"_id" json:"_id,omitempty"`
	Tab                string                 `bson:"tab" json:"-"`
	Title              string                 `bson:"title" json:"title"`
	Type               string                 `bson:"type" json:"type"`
	GridParameters     map[string]interface{} `bson:"grid_parameters" json:"grid_parameters"`
	Parameters         Parameters             `bson:"parameters" json:"parameters"`
	InternalParameters InternalParameters     `bson:"internal_parameters,omitempty" json:"-"`
	Author             string                 `bson:"author" json:"author,omitempty"`
	Created            types.CpsTime          `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated            types.CpsTime          `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Parameters struct {
	MainFilter string `bson:"main_filter,omitempty" json:"main_filter,omitempty"`

	// Junit
	IsAPI                 bool     `bson:"is_api,omitempty" json:"is_api,omitempty"`
	Directory             string   `bson:"directory,omitempty" json:"directory,omitempty"`
	ReportFileRegexp      string   `bson:"report_fileregexp,omitempty" json:"report_fileregexp,omitempty"`
	ScreenshotDirectories []string `bson:"screenshot_directories,omitempty" json:"screenshot_directories,omitempty"`
	VideoDirectories      []string `bson:"video_directories,omitempty" json:"video_directories,omitempty"`
	ScreenshotFilemask    string   `bson:"screenshot_filemask,omitempty" json:"screenshot_filemask,omitempty"`
	VideoFilemask         string   `bson:"video_filemask,omitempty" json:"video_filemask,omitempty"`

	RemainParameters map[string]interface{} `bson:",inline" json:"-"`
}

func (p Parameters) MarshalJSON() ([]byte, error) {
	type Alias Parameters
	b, err := json.Marshal(Alias(p))
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
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
	m := make(map[string]interface{})
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

	RemainParameters map[string]interface{} `bson:",inline"`
}

func (p InternalParameters) IsZero() bool {
	return len(p.TestSuites) == 0 &&
		len(p.RemainParameters) == 0
}

type WidgetFilter struct {
	ID        string        `bson:"_id"`
	Title     string        `bson:"title"`
	Widget    string        `bson:"widget"`
	IsPrivate bool          `bson:"is_private"`
	Author    string        `bson:"author"`
	Created   types.CpsTime `bson:"created,omitempty"`
	Updated   types.CpsTime `bson:"updated,omitempty"`

	savedpattern.AlarmPatternFields     `bson:",inline"`
	savedpattern.EntityPatternFields    `bson:",inline"`
	savedpattern.PbehaviorPatternFields `bson:",inline"`

	// Deprecated : contains old mongo query which cannot be migrated to pattern.
	OldMongoQuery map[string]interface{} `bson:"old_mongo_query,omitempty"`
}
