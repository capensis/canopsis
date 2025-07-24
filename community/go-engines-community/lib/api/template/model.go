package template

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
)

const (
	TypeEvent = iota
	TypeResponse
)

type ValidateResponse struct {
	IsValid bool                 `json:"is_valid"`
	Err     *validator.ErrReport `json:"err"`
	Result  string               `json:"result"`
}

type EditDataRequest struct {
	ID          string            `json:"-"`
	Type        *int              `json:"type" binding:"required,oneof=0 1"`
	Name        string            `json:"name" binding:"required,max=255"`
	Description string            `json:"description" binding:"max=500"`
	Body        map[string]any    `json:"body" binding:"required"`
	Headers     map[string]string `json:"headers" binding:"dive,max=500"`
	Author      string            `json:"author" swaggerignore:"true"`
}

type ListDataRequest struct {
	pagination.FilteredQuery
	SortBy string `form:"sort_by" binding:"oneoforempty=_id name description type"`
	Type   *int   `form:"type"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type VarResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DataResponse struct {
	ID          string            `json:"_id" bson:"_id"`
	Type        int               `json:"type" bson:"type"`
	Name        string            `json:"name" bson:"name"`
	Description string            `json:"description" bson:"description"`
	Body        map[string]any    `json:"body" bson:"body"`
	Headers     map[string]string `json:"headers,omitempty" bson:"headers"`
	Created     datetime.CpsTime  `json:"created" bson:"created" swaggertype:"integer"`
	Updated     datetime.CpsTime  `json:"updated" bson:"updated" swaggertype:"integer"`
}

type AggregationDataResult struct {
	Data       []DataResponse `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

func (r *AggregationDataResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationDataResult) GetTotal() int64 {
	return r.TotalCount
}

type DataModel struct {
	ID          string            `bson:"_id,omitempty"`
	Type        int               `json:"type"`
	Name        string            `bson:"name"`
	Description string            `bson:"description"`
	Body        map[string]any    `bson:"body"`
	Headers     map[string]string `bson:"headers"`
	Author      string            `bson:"author"`
	Created     *datetime.CpsTime `bson:"created,omitempty"`
	Updated     *datetime.CpsTime `bson:"updated"`
}
