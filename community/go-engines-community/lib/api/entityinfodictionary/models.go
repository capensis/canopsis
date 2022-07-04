package entityinfodictionary

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"

type ListValuesRequest struct {
	pagination.Query
	Key    string `form:"key" json:"key" binding:"required"`
	Search string `form:"search" json:"search"`
}

type ListKeysRequest struct {
	pagination.Query
	Search string `form:"search" json:"search"`
}

type Result struct {
	Value string `bson:"value" json:"value"`
}

type AggregationResult struct {
	Data       []Result `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}
