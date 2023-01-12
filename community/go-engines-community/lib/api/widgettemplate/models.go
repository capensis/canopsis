package widgettemplate

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id title type author.name created updated"`
	Type   string `json:"type" form:"type"`
}

type EditRequest struct {
	ID      string              `json:"-"`
	Title   string              `json:"title" binding:"required,max=255"`
	Type    string              `json:"type" binding:"required,oneof=alarm entity"`
	Columns []view.WidgetColumn `json:"columns" binding:"required,notblank,dive"`
	Author  string              `json:"author" swaggerignore:"true"`
}

type Response struct {
	ID      string              `bson:"_id" json:"_id"`
	Title   string              `bson:"title" json:"title"`
	Type    string              `bson:"type" json:"type"`
	Columns []view.WidgetColumn `bson:"columns" json:"columns"`
	Author  *author.Author      `bson:"author" json:"author,omitempty"`
	Created *types.CpsTime      `bson:"created" json:"created,omitempty" swaggertype:"integer"`
	Updated *types.CpsTime      `bson:"updated" json:"updated,omitempty" swaggertype:"integer"`
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
