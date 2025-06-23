package entityinfosproperty

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

const EntityInfosTypeTimestamp = "timestamp"

type EditRequest struct {
	Description string `bson:"description,omitempty" json:"description,omitempty" binding:"max=255"`
	Alias       string `bson:"alias" json:"alias" binding:"max=255"`
	Type        string `bson:"type" json:"type" binding:"required,oneof=string number boolean string_array timestamp"`

	Author  string           `bson:"author,omitempty" json:"author,omitempty" swaggerignore:"true"`
	Created datetime.CpsTime `bson:"created,omitempty" json:"-" swaggerignore:"true"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"-" swaggerignore:"true"`
}

type CreateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"_id" binding:"id"`
	Key         string `bson:"key" json:"key" binding:"required"`
}

type UpdateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"-"`
}

type InfoProperty struct {
	ID          string `bson:"_id" json:"_id"`
	Key         string `bson:"key" json:"key"`
	Description string `bson:"description,omitempty" json:"description,omitempty"`
	Alias       string `bson:"alias" json:"alias"`
	Type        string `bson:"type" json:"type"`

	Created datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Response struct {
	InfoProperty `bson:",inline"`
	Author       *author.Author `bson:"author" json:"author"`
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

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=key alias type"`
	Type   string `json:"type" form:"type" binding:"oneoforempty=string number boolean string_array timestamp"`
}
