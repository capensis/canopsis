package eventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	Author       string                                        `bson:"author" json:"author" swaggerignore:"true"`
	Description  string                                        `bson:"description" json:"description" binding:"required,max=255"`
	Type         string                                        `bson:"type" json:"type" binding:"required,oneof=break drop enrichment change_entity"`
	Priority     int                                           `bson:"priority" json:"priority"`
	Enabled      bool                                          `bson:"enabled" json:"enabled"`
	Config       eventfilter.RuleConfig                        `bson:"config" json:"config" binding:"dive"`
	ExternalData map[string]eventfilter.ExternalDataParameters `bson:"external_data" json:"external_data,omitempty" binding:"dive"`
	Created      *types.CpsTime                                `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated      *types.CpsTime                                `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	common.EntityPatternFieldsRequest
	EventPattern pattern.Event `json:"event_pattern" bson:"event_pattern" binding:"event_pattern"`

	RRule      string           `json:"rrule,omitempty"`
	Start      *types.CpsTime   `json:"start,omitempty" swaggertype:"integer"`
	Stop       *types.CpsTime   `json:"stop,omitempty" swaggertype:"integer"`
	Exdates    []exdate.Request `json:"exdates" binding:"dive"`
	Exceptions []string         `json:"exceptions"`
}

type Response struct {
	ID                               string                                        `bson:"_id" json:"_id"`
	Author                           *author.Author                                `bson:"author" json:"author" swaggerignore:"true"`
	Description                      string                                        `bson:"description" json:"description"`
	Type                             string                                        `bson:"type" json:"type"`
	Priority                         int                                           `bson:"priority" json:"priority"`
	Enabled                          bool                                          `bson:"enabled" json:"enabled"`
	Config                           eventfilter.RuleConfig                        `bson:"config" json:"config"`
	ExternalData                     map[string]eventfilter.ExternalDataParameters `bson:"external_data" json:"external_data,omitempty"`
	Created                          *types.CpsTime                                `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated                          *types.CpsTime                                `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	RRule                            string                                        `bson:"rrule" json:"rrule"`
	Start                            *types.CpsTime                                `bson:"start,omitempty" json:"start,omitempty" swaggertype:"integer"`
	Stop                             *types.CpsTime                                `bson:"stop,omitempty" json:"stop,omitempty" swaggertype:"integer"`
	Exdates                          []types.Exdate                                `bson:"exdates" json:"exdates"`
	Exceptions                       []Exception                                   `bson:"exceptions" json:"exceptions"`
	OldPatterns                      oldpattern.EventPatternList                   `bson:"old_patterns,omitempty" json:"old_patterns,omitempty"`
	EventPattern                     pattern.Event                                 `bson:"event_pattern" json:"event_pattern"`
	savedpattern.EntityPatternFields `bson:",inline"`
}

type Exception struct {
	ID          string         `bson:"_id" json:"_id"`
	Name        string         `bson:"name" json:"name"`
	Description string         `bson:"description" json:"description"`
	Exdates     []types.Exdate `bson:"exdates" json:"exdates"`
	Created     types.CpsTime  `bson:"created" json:"created" swaggertype:"integer"`
}

type CreateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"-" json:"-"`
}

type BulkUpdateRequestItem struct {
	EditRequest `bson:",inline"`
	ID          string `bson:"_id" json:"_id" binding:"required"`
}

type BulkDeleteRequestItem struct {
	ID string `bson:"_id" json:"_id" binding:"required"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id author.name priority created updated on_success on_failure"`
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
