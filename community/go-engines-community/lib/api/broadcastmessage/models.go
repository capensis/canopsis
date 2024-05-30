package broadcastmessage

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type EditRequest struct {
	Color   string           `bson:"color" json:"color" binding:"required,iscolor"`
	Message string           `bson:"message" json:"message" binding:"required"`
	Start   datetime.CpsTime `bson:"start" json:"start" binding:"required" swaggertype:"integer"`
	End     datetime.CpsTime `bson:"end" json:"end" binding:"required" swaggertype:"integer"`

	Author  string            `bson:"author,omitempty" json:"author,omitempty" swaggerignore:"true"`
	Created *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggerignore:"true"`
	Updated *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggerignore:"true"`
}

type CreateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"-"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id message"`
}

type Response struct {
	ID      string           `bson:"_id" json:"_id"`
	Color   string           `bson:"color" json:"color"`
	Message string           `bson:"message" json:"message"`
	Start   datetime.CpsTime `bson:"start" json:"start" swaggertype:"integer"`
	End     datetime.CpsTime `bson:"end" json:"end" swaggertype:"integer"`

	Author  *author.Author    `bson:"author,omitempty" json:"author,omitempty"`
	Created *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	Maintenance bool `bson:"-" json:"maintenance,omitempty"`
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
