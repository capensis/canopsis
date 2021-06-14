package pbehaviorexception

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name created"`
}

type Request struct {
	Name        string          `json:"name" binding:"required,max=255"`
	Description string          `json:"description" binding:"required,max=255"`
	Exdates     []ExdateRequest `json:"exdates" binding:"required,notblank,dive"`
}

type CreateRequest struct {
	Request
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	Request
	ID string `json:"-"`
}

type ExdateRequest struct {
	Begin types.CpsTime `json:"begin" binding:"required" swaggertype:"integer"`
	End   types.CpsTime `json:"end" binding:"required" swaggertype:"integer"`
	Type  string        `json:"type" binding:"required"`
}

type Exception struct {
	ID          string        `bson:"_id" json:"_id"`
	Name        string        `bson:"name" json:"name"`
	Description string        `bson:"description" json:"description"`
	Exdates     []Exdate      `bson:"exdates" json:"exdates"`
	Created     types.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
	Deletable   *bool         `bson:"deletable,omitempty" json:"deletable,omitempty"`
}

type Exdate struct {
	Begin types.CpsTime  `bson:"begin" json:"begin" swaggertype:"integer"`
	End   types.CpsTime  `bson:"end" json:"end" swaggertype:"integer"`
	Type  pbehavior.Type `bson:"type" json:"type"`
}

type AggregationResult struct {
	Data       []Exception `bson:"data" json:"data"`
	TotalCount int64       `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}
