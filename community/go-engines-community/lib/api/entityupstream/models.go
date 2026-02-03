package entityupstream

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
)

type DownstreamsRequest struct {
	pagination.Query
	entity.SortRequest
	ID string `form:"_id" json:"_id" binding:"required"`
}

type UpstreamRequest struct {
	ID string `form:"_id" json:"_id" binding:"required"`
}

type Response struct {
	entity.Entity `bson:",inline"`
}

type AggregationResult struct {
	Data       []Response `bson:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
