package broadcastmessage

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type Payload struct {
	Color   string        `bson:"color" json:"color" binding:"required,iscolor"`
	Message string        `bson:"message" json:"message" binding:"required"`
	Start   types.CpsTime `bson:"start" json:"start" binding:"required" swaggertype:"integer"`
	End     types.CpsTime `bson:"end" json:"end" binding:"required" swaggertype:"integer"`

	Created *types.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *types.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type BroadcastMessage struct {
	ID      string `bson:"_id" json:"_id" binding:"id"`
	Payload `bson:",inline"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id message"`
}

type AggregationResult struct {
	Data       []BroadcastMessage `bson:"data" json:"data"`
	TotalCount int64              `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
