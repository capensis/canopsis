package pbehaviorexception

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name created"`
}

type EditRequest struct {
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description" binding:"required,max=255"`
	Author      string          `json:"author" swaggerignore:"true"`
	Exdates     []ExdateRequest `json:"exdates" binding:"required,notblank,dive"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type ExdateRequest struct {
	exdate.Request
	Type string `json:"type" binding:"required"`
}

type Response struct {
	ID          string            `bson:"_id" json:"_id"`
	Name        string            `bson:"name" json:"name"`
	Description string            `bson:"description" json:"description"`
	Exdates     []Exdate          `bson:"exdates" json:"exdates"`
	Author      *author.Author    `bson:"author,omitempty" json:"author,omitempty"`
	Created     *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated     *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	Deletable   *bool             `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type Exdate struct {
	Begin datetime.CpsTime `bson:"begin" json:"begin" swaggertype:"integer"`
	End   datetime.CpsTime `bson:"end" json:"end" swaggertype:"integer"`
	Type  pbehavior.Type   `bson:"type" json:"type"`
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
