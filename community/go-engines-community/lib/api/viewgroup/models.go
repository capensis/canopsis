package viewgroup

import (
	"encoding/json"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type ListRequest struct {
	pagination.Query
	Search    string `form:"search"`
	WithViews bool   `form:"with_views"`
	WithFlags bool   `form:"with_flags"`
}

type EditRequest struct {
	BaseEditRequest
	ID string `json:"-"`
}

type BaseEditRequest struct {
	Title  string `json:"title" binding:"required,max=255"`
	Author string `json:"author" swaggerignore:"true"`
}

type ViewGroup struct {
	ID        string        `bson:"_id" json:"_id"`
	Title     string        `bson:"title" json:"title"`
	Author    string        `bson:"author" json:"author"`
	Views     *[]View       `bson:"views,omitempty" json:"views,omitempty"`
	Created   types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated   types.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
	Deletable *bool         `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type View struct {
	ID              string                     `bson:"_id" json:"_id"`
	Enabled         bool                       `bson:"enabled" json:"enabled"`
	Title           string                     `bson:"title" json:"title"`
	Description     string                     `bson:"description" json:"description"`
	Group           ViewGroup                  `bson:"group" json:"group"`
	Tabs            []view.Tab                 `bson:"tabs" json:"tabs"`
	Tags            []string                   `bson:"tags" json:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `bson:"periodic_refresh" json:"periodic_refresh"`
	Author          string                     `bson:"author" json:"author"`
	Created         *types.CpsTime             `bson:"created" json:"created" swaggertype:"integer"`
	Updated         *types.CpsTime             `bson:"updated" json:"updated" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []ViewGroup `bson:"data" json:"data"`
	TotalCount int64       `bson:"total_count" json:"total_count"`
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
	IDs []string `form:"ids[]" json:"ids" binding:"required,unique,notblank"`
}
