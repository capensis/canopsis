package pbehaviorreason

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.FilteredQuery
	WithHidden bool   `form:"with_hidden"`
	SortBy     string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name created"`
}

type EditRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	Description string `json:"description" binding:"required,max=255"`
	Author      string `json:"author" swaggerignore:"true"`

	// Hidden is used in API to hide documents from the list response
	Hidden bool `json:"hidden"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type Response struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`

	Deletable *bool `bson:"deletable,omitempty" json:"deletable,omitempty"`
	Hidden    bool  `bson:"hidden" json:"hidden"`

	Author  *author.Author    `bson:"author,omitempty" json:"author,omitempty"`
	Created *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
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
