package flappingrule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	Name           string                    `json:"name" binding:"required,max=255"`
	Description    string                    `json:"description" binding:"max=255"`
	FreqLimit      int                       `json:"freq_limit" binding:"required,gt=0"`
	Duration       types.DurationWithUnit    `json:"duration" binding:"required"`
	AlarmPatterns  pattern.AlarmPatternList  `json:"alarm_patterns"`
	EntityPatterns pattern.EntityPatternList `json:"entity_patterns"`
	Priority       int                       `json:"priority" binding:"required,gt=0"`
	Author         string                    `json:"author" swaggerignore:"true"`
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
	ID             string                    `bson:"_id" json:"_id"`
	Name           string                    `bson:"name" json:"name"`
	Description    string                    `bson:"description" json:"description"`
	FreqLimit      int                       `bson:"freq_limit" json:"freq_limit"`
	Duration       types.DurationWithUnit    `bson:"duration" json:"duration"`
	AlarmPatterns  pattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EntityPatterns pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	Priority       int                       `bson:"priority" json:"priority"`
	Author         common.User               `bson:"author" json:"author"`
	Created        types.CpsTime             `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated        types.CpsTime             `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
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
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id description freq_limit duration author created updated priority"`
}
