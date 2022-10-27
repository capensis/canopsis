package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
)

type ListRequest struct {
	pagination.FilteredQuery
	OnlyDefault bool     `form:"default"`
	SortBy      string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=name priority"`
	Types       []string `form:"types[]" json:"types"`
}

type EditRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
	Type        string `json:"type" binding:"required,oneof=active inactive maintenance pause"`
	Priority    *int   `json:"priority" binding:"required"`
	IconName    string `json:"icon_name" binding:"required,max=255"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type Type struct {
	ID          string `bson:"_id,omitempty" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
	Priority    int    `bson:"priority" json:"priority"`
	IconName    string `bson:"icon_name" json:"icon_name"`
	Editable    *bool  `bson:"editable,omitempty" json:"editable,omitempty"`
	Deletable   *bool  `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type AggregationResult struct {
	Data       []Type `bson:"data" json:"data"`
	TotalCount int64  `bson:"total_count" json:"total_count"`
}

// GetTotal implementation PaginatedData interface
func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

// GetData implementation PaginatedData interface
func (r AggregationResult) GetData() interface{} {
	return r.Data
}
