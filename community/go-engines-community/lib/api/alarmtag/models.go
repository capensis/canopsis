package alarmtag

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.Query
	Search string `form:"search" json:"search"`
	Sort   string `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=value created"`
}

type Response struct {
	ID      string        `bson:"_id" json:"_id"`
	Value   string        `bson:"value" json:"value"`
	Color   string        `bson:"color" json:"color"`
	Created types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
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
