package entitycategory

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=name created"`
}

type EditRequest struct {
	ID     string `json:"-"`
	Name   string `json:"name" binding:"required,max=255"`
	Author string `json:"author" swaggerignore:"true"`
}

type Category struct {
	ID      string            `bson:"_id" json:"_id"`
	Name    string            `bson:"name" json:"name"`
	Author  string            `bson:"author" json:"author"`
	Created *datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated *datetime.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type Response struct {
	ID      string            `bson:"_id" json:"_id"`
	Name    string            `bson:"name" json:"name"`
	Author  *author.Author    `bson:"author" json:"author"`
	Created *datetime.CpsTime `bson:"created" json:"created"`
	Updated *datetime.CpsTime `bson:"updated" json:"updated"`
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
