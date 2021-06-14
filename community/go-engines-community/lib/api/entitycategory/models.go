package entitycategory

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
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
	ID      string         `bson:"_id" json:"_id"`
	Name    string         `bson:"name" json:"name"`
	Author  string         `bson:"author" json:"author"`
	Created *types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Updated *types.CpsTime `bson:"updated" json:"updated" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Category `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
