package view

import (
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"github.com/mitchellh/mapstructure"
)

type ListRequest struct {
	pagination.Query
	Search string   `form:"search"`
	Ids    []string `form:"-"`
}

type EditRequest struct {
	BaseEditRequest
	ID string `json:"-"`
}

type BaseEditRequest struct {
	Enabled         *bool                      `json:"enabled" binding:"required"`
	Title           string                     `json:"title" binding:"required,max=255"`
	Description     string                     `json:"description" binding:"max=255"`
	Group           string                     `json:"group" binding:"required"`
	Tabs            []TabRequest               `json:"tabs" binding:"dive"`
	Tags            []string                   `json:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `json:"periodic_refresh"`
	Author          string                     `json:"author" swaggerignore:"true"`
}

type TabRequest struct {
	ID      string          `json:"_id" binding:"required,max=255"`
	Title   string          `json:"title" binding:"required,max=255"`
	Widgets []WidgetRequest `json:"widgets" binding:"dive"`
}

type WidgetRequest struct {
	ID             string                 `json:"_id" binding:"required,max=255"`
	Title          string                 `json:"title" binding:"max=255"`
	Type           string                 `json:"type" binding:"required,max=255"`
	GridParameters map[string]interface{} `json:"grid_parameters"`
	Parameters     interface{}            `json:"parameters"`
}

type WidgetParametersJunitRequest struct {
	IsAPI                 bool                   `json:"is_api" mapstructure:"is_api"`
	Directory             string                 `json:"directory" mapstructure:"directory"`
	ReportFileRegexp      string                 `json:"report_fileregexp" mapstructure:"report_fileregexp"`
	ScreenshotDirectories []string               `json:"screenshot_directories" mapstructure:"screenshot_directories"`
	VideoDirectories      []string               `json:"video_directories" mapstructure:"video_directories"`
	ScreenshotFilemask    string                 `json:"screenshot_filemask" mapstructure:"screenshot_filemask"`
	VideoFilemask         string                 `json:"video_filemask" mapstructure:"video_filemask"`
	RemainParameters      map[string]interface{} `json:"-" mapstructure:",remain"`
}

func (w *WidgetRequest) UnmarshalJSON(b []byte) error {
	type Alias WidgetRequest
	var tmp Alias

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	*w = WidgetRequest(tmp)

	switch w.Type {
	case view.WidgetTypeJunit:
		var params WidgetParametersJunitRequest
		err := mapstructure.Decode(w.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		w.Parameters = params
	}

	return nil
}

func (w WidgetParametersJunitRequest) MarshalJSON() ([]byte, error) {
	type Alias WidgetParametersJunitRequest
	b, err := json.Marshal(Alias(w))
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	for k, v := range w.RemainParameters {
		m[k] = v
	}

	return json.Marshal(m)
}

type EditPositionRequest struct {
	Items []EditPositionItemRequest `json:"items" binding:"required,notblank,dive"`
}

func (r EditPositionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *EditPositionRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}

type EditPositionItemRequest struct {
	ID    string   `json:"_id" binding:"required"`
	Views []string `json:"views" binding:"required"`
}

type AggregationResult struct {
	Data       []viewgroup.View `bson:"data" json:"data"`
	TotalCount int64            `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type BulkCreateRequest struct {
	Items []EditRequest `binding:"required,notblank,dive"`
}

func (r BulkCreateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *BulkCreateRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}

type BulkUpdateRequest struct {
	Items []BulkUpdateRequestItem `binding:"required,notblank,dive"`
}

type BulkUpdateRequestItem struct {
	BaseEditRequest
	ID string `json:"_id" binding:"required"`
}

func (r BulkUpdateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *BulkUpdateRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}

type BulkDeleteRequest struct {
	IDs []string `form:"ids[]" binding:"required,notblank"`
}
