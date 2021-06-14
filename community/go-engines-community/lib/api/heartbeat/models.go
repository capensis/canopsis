package heartbeat

import (
	"encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author expected_interval created updated"`
}

type BaseEditRequest struct {
	Name             string                 `json:"name" binding:"required,max=255"`
	Description      string                 `json:"description" binding:"required,max=255"`
	Author           string                 `json:"author" binding:"required,max=255"`
	Pattern          map[string]string      `json:"pattern" binding:"required"`
	ExpectedInterval types.CpsShortDuration `json:"expected_interval" binding:"required"`
	Output           string                 `json:"output" binding:"max=255"`
}

type CreateRequest struct {
	BaseEditRequest
	ID string `json:"_id" binding:"max=255"`
}

type UpdateRequest struct {
	BaseEditRequest
	ID string `json:"-"`
}

type BulkCreateRequest struct {
	Items []CreateRequest `binding:"required,notblank,dive"`
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

type Heartbeat struct {
	ID               string                 `bson:"_id,omitempty" json:"_id"`
	Name             string                 `bson:"name" json:"name"`
	Description      string                 `bson:"description" json:"description"`
	Author           string                 `bson:"author" json:"author"`
	Pattern          map[string]string      `bson:"pattern" json:"pattern"`
	ExpectedInterval types.CpsShortDuration `bson:"expected_interval" json:"expected_interval"`
	Output           string                 `bson:"output" json:"output"`
	Created          *types.CpsTime         `bson:"created,omitempty" json:"created" swaggertype:"integer"`
	Updated          *types.CpsTime         `bson:"updated,omitempty" json:"updated" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Heartbeat `bson:"data" json:"data"`
	TotalCount int64       `bson:"total_count" json:"total_count"`
}

// GetTotal implementation PaginatedData interface
func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

// GetData implementation PaginatedData interface
func (r AggregationResult) GetData() interface{} {
	return r.Data
}
