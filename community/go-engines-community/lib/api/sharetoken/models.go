package sharetoken

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.Query
	Search string `form:"search"`
	Sort   string `form:"sort" binding:"oneoforempty=asc desc"`
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=user.name description created accessed expired"`
}

type EditRequest struct {
	Duration    *types.DurationWithUnit `json:"duration"`
	Description string                  `json:"description" binding:"max=500"`
}

type Response struct {
	ID          string         `bson:"_id" json:"_id"`
	Value       string         `bson:"value" json:"value"`
	User        *author.Author `bson:"user" json:"user"`
	Role        *author.Role   `bson:"role" json:"role"`
	Description string         `bson:"description" json:"description"`
	Created     types.CpsTime  `bson:"created" json:"created" swaggertype:"integer"`
	Accessed    types.CpsTime  `bson:"accessed" json:"accessed" swaggertype:"integer"`
	Expired     *types.CpsTime `bson:"expired" json:"expired" swaggertype:"integer"`
}

type Model struct {
	ID          string         `bson:"_id"`
	Value       string         `bson:"value"`
	User        string         `bson:"user"`
	Description string         `bson:"description"`
	Created     types.CpsTime  `bson:"created"`
	Accessed    types.CpsTime  `bson:"accessed"`
	Expired     *types.CpsTime `bson:"expired,omitempty"`
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
