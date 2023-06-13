package alarmtag

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.Query
	Search string `form:"search" json:"search"`
	Sort   string `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=value created"`
}

type CreateRequest struct {
	Value  string `json:"value" binding:"required,max=255"`
	Color  string `json:"color" binding:"required,iscolor"`
	Author string `json:"author" swaggerignore:"true"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
}

type UpdateRequest struct {
	ID     string `json:"-"`
	Color  string `json:"color" binding:"required,iscolor"`
	Author string `json:"author" swaggerignore:"true"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
}

type Response struct {
	ID      string         `bson:"_id" json:"_id"`
	Type    int64          `bson:"type" json:"type"`
	Value   string         `bson:"value" json:"value"`
	Color   string         `bson:"color" json:"color"`
	Author  *author.Author `bson:"author" json:"author"`
	Updated types.CpsTime  `bson:"updated" json:"updated" swaggertype:"integer"`
	Created types.CpsTime  `bson:"created" json:"created" swaggertype:"integer"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
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
