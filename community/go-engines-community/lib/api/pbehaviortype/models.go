package pbehaviortype

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
)

type ListRequest struct {
	pagination.FilteredQuery
	OnlyDefault bool     `form:"default"`
	WithHidden  bool     `form:"with_hidden"`
	SortBy      string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=name priority"`
	Types       []string `form:"types[]" json:"types"`
}

type EditRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
	Type        string `json:"type" binding:"required,oneof=active inactive maintenance pause"`
	Priority    int64  `json:"priority" binding:"required,min=1"`
	Color       string `json:"color" binding:"required,iscolor"`
	Author      string `json:"author" swaggerignore:"true"`
	Hidden      bool   `json:"hidden"`
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

type Response struct {
	pbehavior.Type `bson:",inline"`
	Author         *author.Author `bson:"author" json:"author"`
	Default        *bool          `bson:"default,omitempty" json:"default,omitempty"`
	Deletable      *bool          `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}

type PriorityResponse struct {
	Priority int64 `json:"priority"`
}
