package entitycomment

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Request struct {
	Message string `json:"message" binding:"required,max=200"`
	Entity  string `json:"entity" binding:"required"`
}

type UpdateRequest struct {
	Request
	ID string `json:"-"`
}

type Comment struct {
	Timestamp datetime.CpsTime `bson:"t" json:"t" swaggertype:"integer"`
	Author    *types.Author    `bson:"a" json:"author"`
	Message   string           `bson:"m" json:"message"`
}

type Response struct {
	ID      string `bson:"_id" json:"_id"`
	Entity  string `json:"entity,omitempty"`
	Comment `bson:",inline" json:",inline"`
}

type ListRequest struct {
	pagination.Query
	Entity string `form:"entity" json:"entity" binding:"required"`
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
