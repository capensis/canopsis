package view

import (
	"encoding/json"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewgroup"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
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
	Tags            []string                   `json:"tags"`
	PeriodicRefresh *types.DurationWithEnabled `json:"periodic_refresh"`
	Author          string                     `json:"author" swaggerignore:"true"`
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
