package colortheme

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type EditRequest struct {
	ID       string `json:"-"`
	Name     string `json:"name" binding:"required"`
	Colors   Colors `json:"colors"`
	FontSize int    `json:"font_size" binding:"required,oneof=1 2 3"`
	Author   string `json:"author" swaggerignore:"true"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type Response struct {
	ID        string            `bson:"_id" json:"_id"`
	Name      string            `bson:"name" json:"name"`
	Colors    Colors            `bson:"colors" json:"colors"`
	FontSize  int               `bson:"font_size" json:"font_size"`
	Author    *author.Author    `bson:"author" json:"author"`
	Created   *datetime.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated   *datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	Deletable bool              `bson:"deletable" json:"deletable"`
}

type Colors struct {
	Main struct {
		Primary           string `bson:"primary" json:"primary" binding:"required,iscolor"`
		Secondary         string `bson:"secondary" json:"secondary" binding:"required,iscolor"`
		Accent            string `bson:"accent" json:"accent" binding:"required,iscolor"`
		Error             string `bson:"error" json:"error" binding:"required,iscolor"`
		Info              string `bson:"info" json:"info" binding:"required,iscolor"`
		Success           string `bson:"success" json:"success" binding:"required,iscolor"`
		Warning           string `bson:"warning" json:"warning" binding:"required,iscolor"`
		Background        string `bson:"background" json:"background" binding:"required,iscolor"`
		ActiveColor       string `bson:"active_color" json:"active_color" binding:"required,iscolor"`
		ErrorBackground   string `bson:"error_background" json:"error_background" binding:"required,iscolor"`
		WarningBackground string `bson:"warning_background" json:"warning_background" binding:"required,iscolor"`
		InfoBackground    string `bson:"info_background" json:"info_background" binding:"required,iscolor"`
		SuccessBackground string `bson:"success_background" json:"success_background" binding:"required,iscolor"`
	} `bson:"main" json:"main"`
	Table struct {
		Background    string `bson:"background" json:"background" binding:"required,iscolor"`
		RowColor      string `bson:"row_color" json:"row_color" binding:"required,iscolor"`
		ShiftRow      *bool  `bson:"shift_row" json:"shift_row" binding:"required"`
		ShiftRowColor string `bson:"shift_row_color,omitempty" json:"shift_row_color,omitempty" binding:"required_if=ShiftRow true,iscolororempty"`
		HoverRow      *bool  `bson:"hover_row" json:"hover_row" binding:"required"`
		HoverRowColor string `bson:"hover_row_color,omitempty" json:"hover_row_color,omitempty" binding:"required_if=HoverRow true,iscolororempty"`
	} `bson:"table" json:"table"`
	State struct {
		OK       string `bson:"ok" json:"ok" binding:"required,iscolor"`
		Minor    string `bson:"minor" json:"minor" binding:"required,iscolor"`
		Major    string `bson:"major" json:"major" binding:"required,iscolor"`
		Critical string `bson:"critical" json:"critical" binding:"required,iscolor"`
	} `bson:"state" json:"state"`
}

type Document struct {
	ID        string           `bson:"_id,omitempty"`
	Name      string           `bson:"name"`
	Colors    Colors           `bson:"colors"`
	FontSize  int              `bson:"font_size"`
	Author    string           `bson:"author"`
	Created   datetime.CpsTime `bson:"created,omitempty"`
	Updated   datetime.CpsTime `bson:"updated,omitempty"`
	Deletable bool             `bson:"deletable"`
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

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=name updated"`
}
