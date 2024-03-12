package alarmtag

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=value created"`
	Values []string `form:"values[]" json:"values"`
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
	ID        string           `bson:"_id" json:"_id"`
	Type      int64            `bson:"type" json:"type"`
	Value     string           `bson:"value" json:"value"`
	Color     string           `bson:"color" json:"color"`
	Author    *author.Author   `bson:"author" json:"author"`
	Updated   datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
	Created   datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Deletable *bool            `bson:"deletable,omitempty" json:"deletable,omitempty"`

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

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}
