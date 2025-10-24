package entityinfosproperty

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
)

type EditRequest struct {
	Description string `bson:"description" json:"description" binding:"max=255"`
	Alias       string `bson:"alias,omitempty" json:"alias,omitempty" binding:"max=255"`

	// Possible type values.
	//   * `0` - type boolean
	//   * `1` - type number
	//   * `2` - type timestamp
	//   * `3` - type string
	//   * `4` - type string_array
	Type *int `bson:"type" json:"type" binding:"required,oneof=0 1 2 3 4"`

	Author  string           `bson:"author,omitempty" json:"author,omitempty" swaggerignore:"true"`
	Created datetime.CpsTime `bson:"created,omitempty" json:"-" swaggerignore:"true"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"-" swaggerignore:"true"`
}

type CreateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"_id" binding:"id"`
	Name        string `bson:"name" json:"name" binding:"required"`
}

type UpdateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"-"`
}

type InfoProperty struct {
	ID          string `bson:"_id" json:"_id"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Alias       string `bson:"alias" json:"alias"`

	// Possible type values.
	//   * `0` - type boolean
	//   * `1` - type number
	//   * `2` - type timestamp
	//   * `3` - type string
	//   * `4` - type string_array
	Type int `bson:"type" json:"type"`

	Created datetime.CpsTime `bson:"created,omitempty" json:"created,omitzero" swaggertype:"integer"`
	Updated datetime.CpsTime `bson:"updated,omitempty" json:"updated,omitzero" swaggertype:"integer"`
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
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=name alias type"`

	// Possible type values.
	//   * `0` - type boolean
	//   * `1` - type number
	//   * `2` - type timestamp
	//   * `3` - type string
	//   * `4` - type string_array
	Type *int `json:"type" form:"type" binding:"omitempty,oneof=0 1 2 3 4"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}
