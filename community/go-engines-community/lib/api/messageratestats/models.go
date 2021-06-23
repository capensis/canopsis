package messageratestats

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	IntervalMinute = "minute"
	IntervalHour   = "hour"
)

type ListRequest struct {
	pagination.Query
	Interval string        `form:"interval" json:"interval" binding:"required,oneof=minute hour"`
	From     types.CpsTime `form:"from" json:"from" binding:"required" swaggertype:"integer"`
	To       types.CpsTime `form:"to" json:"to" binding:"required" swaggertype:"integer"`
	Sort     string        `form:"sort" json:"sort"`
	SortBy   string        `form:"sort_by" json:"sort_by" binding:"oneoforempty=_id received"`
}

type StatsResponse struct {
	ID       int64 `bson:"_id" json:"_id"`
	Received int64 `bson:"received" json:"received"`
}

type AggregationResult struct {
	Data       []StatsResponse `bson:"data" json:"data"`
	TotalCount int64           `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}
