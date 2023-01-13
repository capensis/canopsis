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
	Color       string `json:"color" binding:"required,iscolor"`
}

type CreateRequest struct {
	EditRequest
	ID       string `json:"_id" binding:"id"`
	IconName string `json:"icon_name" binding:"required,max=255"`
}

type UpdateRequest struct {
	EditRequest
	ID       string `json:"-"`
	IconName string `json:"icon_name" binding:"max=255"`
}

type Type struct {
	ID          string `bson:"_id,omitempty" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
	Priority    int    `bson:"priority" json:"priority"`
	IconName    string `bson:"icon_name" json:"icon_name"`
	Color       string `bson:"color" json:"color"`
	Default     *bool  `bson:"default,omitempty" json:"default,omitempty"`
	Deletable   *bool  `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type AggregationResult struct {
	Data       []Type `bson:"data" json:"data"`
	TotalCount int64  `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}
