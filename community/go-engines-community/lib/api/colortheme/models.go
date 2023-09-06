package colortheme

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	Name         string        `bson:"name" json:"name" binding:"required"`
	Colors       Colors        `bson:"colors" json:"colors"`
	LastModified types.CpsTime `bson:"last_modified" json:"-"`
}

type CreateRequest struct {
	EditRequest
	ID string `bson:"_id" json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type Theme struct {
	ID           string        `bson:"_id" json:"_id"`
	Name         string        `bson:"name" json:"name"`
	Colors       Colors        `bson:"colors" json:"colors"`
	LastModified types.CpsTime `bson:"last_modified" json:"last_modified" swaggertype:"integer"`
}

type Colors struct {
	Main struct {
		Primary    string `bson:"primary" json:"primary" binding:"required,iscolor"`
		Secondary  string `bson:"secondary" json:"secondary" binding:"required,iscolor"`
		Accent     string `bson:"accent" json:"accent" binding:"required,iscolor"`
		Error      string `bson:"error" json:"error" binding:"required,iscolor"`
		Info       string `bson:"info" json:"info" binding:"required,iscolor"`
		Success    string `bson:"success" json:"success" binding:"required,iscolor"`
		Warning    string `bson:"warning" json:"warning" binding:"required,iscolor"`
		Background string `bson:"background" json:"background" binding:"required,iscolor"`
	} `bson:"main" json:"main"`
	Table struct {
		Background    string `bson:"background" json:"background" binding:"required,iscolor"`
		ActiveColor   string `bson:"active_color" json:"active_color" binding:"required,iscolor"`
		RowColor      string `bson:"row_color" json:"row_color" binding:"required,iscolor"`
		ShiftRowColor string `bson:"shift_row_color,omitempty" json:"shift_row_color,omitempty" binding:"iscolororempty"`
		HoverRowColor string `bson:"hover_row_color,omitempty" json:"hover_row_color,omitempty" binding:"iscolororempty"`
	} `bson:"table" json:"table"`
	State struct {
		OK       string `bson:"ok" json:"ok" binding:"required,iscolor"`
		Minor    string `bson:"minor" json:"minor" binding:"required,iscolor"`
		Major    string `bson:"major" json:"major" binding:"required,iscolor"`
		Critical string `bson:"critical" json:"critical" binding:"required,iscolor"`
	} `bson:"state" json:"state"`
}

type AggregationResult struct {
	Data       []Theme `bson:"data" json:"data"`
	TotalCount int64   `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=name last_modified"`
}
