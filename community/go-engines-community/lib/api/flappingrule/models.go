package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
)

type EditRequest struct {
	Name        string                    `json:"name" binding:"required,max=255"`
	Description string                    `json:"description" binding:"max=255"`
	FreqLimit   int                       `json:"freq_limit" binding:"required,gt=0"`
	Duration    datetime.DurationWithUnit `json:"duration" binding:"required"`
	Priority    int64                     `json:"priority" binding:"min=0"`
	Author      string                    `json:"author" swaggerignore:"true"`

	common.AlarmPatternFieldsRequest
	common.EntityPatternFieldsRequest
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
	ID          string                    `bson:"_id" json:"_id"`
	Name        string                    `bson:"name" json:"name"`
	Description string                    `bson:"description" json:"description"`
	FreqLimit   int                       `bson:"freq_limit" json:"freq_limit"`
	Duration    datetime.DurationWithUnit `bson:"duration" json:"duration"`
	Priority    int64                     `bson:"priority" json:"priority"`
	Author      *author.Author            `bson:"author" json:"author"`
	Created     datetime.CpsTime          `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated     datetime.CpsTime          `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
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
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name description freq_limit author.name author.display_name created updated priority"`
}
