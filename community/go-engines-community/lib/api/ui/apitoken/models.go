package apitoken

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=name created expired accessed"`
	Active *bool  `form:"active"`
}

type EditRequest struct {
	Name     string                     `json:"name" binding:"required,max=255"`
	Duration *datetime.DurationWithUnit `json:"duration"`
	Author   string                     `json:"author" swaggerignore:"true"`
}

type Token struct {
	ID       string            `bson:"_id"`
	Name     string            `bson:"name"`
	Value    string            `bson:"value"`
	Hash     string            `bson:"hash"`
	User     string            `bson:"user"`
	Expired  *datetime.CpsTime `bson:"expired"`
	Created  *datetime.CpsTime `bson:"created"`
	Accessed *datetime.CpsTime `bson:"accessed"`
}

type Response struct {
	ID       string            `bson:"_id" json:"_id"`
	Name     string            `bson:"name" json:"name"`
	Value    string            `bson:"value" json:"value"`
	User     *author.Author    `bson:"user" json:"user"`
	Roles    []author.Role     `bson:"roles" json:"roles"`
	Expired  *datetime.CpsTime `bson:"expired" json:"expired" swaggertype:"integer"`
	Created  *datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Accessed *datetime.CpsTime `bson:"accessed" json:"accessed" swaggertype:"integer"`
}

type AggregationResult struct {
	Data       []Response `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() any {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
