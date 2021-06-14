package entity

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type ListRequestWithPagination struct {
	pagination.Query
	ListRequest
}

type ListRequest struct {
	Search   string   `form:"search" json:"search"`
	Filter   string   `form:"filter" json:"filter"`
	Sort     string   `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy   string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=name type"`
	SearchBy []string `form:"active_columns[]" json:"active_columns[]"`
}

type ExportRequest struct {
	ListRequest
	Separator string `form:"separator" json:"separator" binding:"oneoforempty=comma semicolon tab space"`
}

type ExportResponse struct {
	ID     string `json:"_id"`
	Status int    `json:"status"`
}

type AggregationResult struct {
	Data       []types.Entity `bson:"data" json:"data"`
	TotalCount int64          `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}

type Entity struct {
	ID             string          `bson:"_id" json:"_id"`
	Name           string          `bson:"name" json:"name"`
	Impacts        []string        `bson:"impact" json:"impact"`
	Depends        []string        `bson:"depends" json:"depends"`
	EnableHistory  []types.CpsTime `bson:"enable_history" json:"enable_history" swaggertype:"array,integer"`
	Measurements   interface{}     `bson:"measurements" json:"measurements"`
	Enabled        bool            `bson:"enabled" json:"enabled"`
	Infos          map[string]Info `bson:"infos" json:"infos"`
	ComponentInfos map[string]Info `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string          `bson:"type" json:"type"`
	Component      string          `bson:"component,omitempty" json:"component,omitempty"`
}

type Info struct {
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}
