package view

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.Query
	Search string   `form:"search"`
	Ids    []string `form:"-"`
}

type EditRequest struct {
	ID          string   `json:"-"`
	Enabled     *bool    `json:"enabled" binding:"required"`
	Title       string   `json:"title" binding:"required,max=255"`
	Description string   `json:"description" binding:"max=255"`
	Group       string   `json:"group" binding:"required"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author" swaggerignore:"true"`

	PeriodicRefresh *types.DurationWithEnabled `json:"periodic_refresh"`
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

type Response struct {
	ID              string                     `bson:"_id" json:"_id,omitempty"`
	Enabled         bool                       `bson:"enabled" json:"enabled"`
	Title           string                     `bson:"title" json:"title"`
	Description     string                     `bson:"description" json:"description"`
	Tabs            *[]viewtab.Response        `bson:"tabs" json:"tabs,omitempty"`
	Tags            []string                   `bson:"tags" json:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `bson:"periodic_refresh" json:"periodic_refresh"`
	Group           *ViewGroup                 `bson:"group" json:"group,omitempty"`
	Author          *author.Author             `bson:"author" json:"author,omitempty"`
	Created         *types.CpsTime             `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated         *types.CpsTime             `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`
}

type ViewGroup struct {
	ID      string         `bson:"_id" json:"_id,omitempty"`
	Title   string         `bson:"title" json:"title"`
	Author  *author.Author `bson:"author" json:"author,omitempty"`
	Created *types.CpsTime `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated *types.CpsTime `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type ImportItemRequest struct {
	ViewGroup
	Views []Response `json:"views"`
}

type ImportRequest struct {
	Items []ImportItemRequest `json:"items" binding:"required,notblank,dive"`
}

func (r ImportRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Items)
}

func (r *ImportRequest) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &r.Items)
}

type ExportRequest struct {
	Groups []struct {
		ID    string   `json:"_id" binding:"required"`
		Views []string `json:"views"`
	} `json:"groups"`
	Views []string `json:"views"`
}

type ExportViewGroupResponse struct {
	ViewGroup `bson:",inline"`
	Views     []Response `json:"views"`
}

type ExportResponse struct {
	Groups []ExportViewGroupResponse `json:"groups"`
	Views  []Response                `json:"views"`
}
