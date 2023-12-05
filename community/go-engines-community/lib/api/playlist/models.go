package playlist

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type ListRequest struct {
	pagination.FilteredQuery
	SortBy string   `json:"sort_by" form:"sort_by" binding:"oneoforempty=name enabled fullscreen"`
	Ids    []string `form:"-"`
}

type EditRequest struct {
	ID         string                    `json:"-"`
	Author     string                    `json:"author" swaggerignore:"true"`
	Name       string                    `json:"name" binding:"required,max=255"`
	Enabled    *bool                     `json:"enabled" binding:"required"`
	Fullscreen *bool                     `json:"fullscreen" binding:"required"`
	TabsList   []string                  `json:"tabs_list" binding:"required,notblank"`
	Interval   datetime.DurationWithUnit `json:"interval" binding:"required"`
}

type Response struct {
	ID         string                    `bson:"_id,omitempty" json:"_id"`
	Author     *author.Author            `bson:"author" json:"author,omitempty"`
	Name       string                    `bson:"name" json:"name"`
	Enabled    bool                      `bson:"enabled" json:"enabled"`
	Fullscreen bool                      `bson:"fullscreen" json:"fullscreen"`
	TabsList   []string                  `bson:"tabs_list" json:"tabs_list"`
	Interval   datetime.DurationWithUnit `bson:"interval" json:"interval"`
	Created    datetime.CpsTime          `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated    datetime.CpsTime          `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Playlist struct {
	ID         string                    `bson:"_id,omitempty" json:"_id"`
	Author     string                    `bson:"author" json:"author,omitempty"`
	Name       string                    `bson:"name" json:"name"`
	Enabled    bool                      `bson:"enabled" json:"enabled"`
	Fullscreen bool                      `bson:"fullscreen" json:"fullscreen"`
	TabsList   []string                  `bson:"tabs_list" json:"tabs_list"`
	Interval   datetime.DurationWithUnit `bson:"interval" json:"interval"`
	Created    datetime.CpsTime          `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated    datetime.CpsTime          `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
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
