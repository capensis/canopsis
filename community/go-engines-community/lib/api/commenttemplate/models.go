package commenttemplate

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type Field struct {
	Name     string `bson:"name" json:"name" binding:"required,max=255"`
	Required bool   `bson:"required" json:"required"`
}

type EditRequest struct {
	ID     string  `json:"-"`
	Name   string  `json:"name" binding:"required,max=255"`
	Fields []Field `json:"fields" binding:"required,dive"`

	Author string `json:"author" swaggerignore:"true"`
}

type Response struct {
	ID      string           `bson:"_id" json:"_id"`
	Name    string           `bson:"name" json:"name"`
	Fields  []Field          `bson:"fields" json:"fields"`
	Author  *author.Author   `bson:"author" json:"author"`
	Created datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Document struct {
	ID      string           `bson:"_id,omitempty"`
	Name    string           `bson:"name"`
	Fields  []Field          `bson:"fields"`
	Author  string           `bson:"author"`
	Created datetime.CpsTime `bson:"created,omitempty"`
	Updated datetime.CpsTime `bson:"updated,omitempty"`
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

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string   `json:"sort_by" form:"sort_by" binding:"oneoforempty=name author.display_name updated"`
	IDs    []string `json:"ids[]" form:"ids[]"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}
