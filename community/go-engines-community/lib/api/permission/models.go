package permission

import "git.canopsis.net/canopsis/go-engines/lib/api/pagination"

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=name description"`
}

type Permission struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Type        string `bson:"type" json:"type"`
}

type AggregationResult struct {
	Data       []Permission `bson:"data" json:"data"`
	TotalCount int64        `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
