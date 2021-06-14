package pbehaviorreason

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name created"`
}

type Request struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
}

type CreateRequest struct {
	Request
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	Request
	ID string `json:"-"`
}

type Reason struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Deletable   *bool  `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type AggregationResult struct {
	Data       []Reason `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
