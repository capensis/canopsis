package widgettemplate

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
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
	Type    string              `json:"type" binding:"required,oneof=alarm_columns entity_columns alarm_more_infos alarm_export_to_pdf weather_item weather_modal weather_entity"`
	Columns []view.WidgetColumn `json:"columns" binding:"dive"`
	Content string              `json:"content"`
	Author  string              `json:"author" swaggerignore:"true"`
}

type Response struct {
	ID      string              `bson:"_id" json:"_id"`
	Title   string              `bson:"title" json:"title"`
	Type    string              `bson:"type" json:"type"`
	Columns []view.WidgetColumn `bson:"columns" json:"columns,omitempty"`
	Content string              `bson:"content" json:"content,omitempty"`
	Author  *author.Author      `bson:"author" json:"author"`
	Created *datetime.CpsTime   `bson:"created" json:"created" swaggertype:"integer"`
	Updated *datetime.CpsTime   `bson:"updated" json:"updated" swaggertype:"integer"`
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
