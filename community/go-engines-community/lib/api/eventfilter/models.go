package eventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/exdate"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	Author       string                                        `json:"author" swaggerignore:"true"`
	Description  string                                        `json:"description" binding:"required,max=255"`
	Type         string                                        `json:"type" binding:"required,oneof=break drop enrichment change_entity"`
	Priority     int64                                         `json:"priority" binding:"min=0"`
	Enabled      bool                                          `json:"enabled"`
	Config       eventfilter.RuleConfig                        `json:"config"`
	ExternalData map[string]eventfilter.ExternalDataParameters `json:"external_data,omitempty" binding:"dive"`

	common.EntityPatternFieldsRequest
	EventPattern pattern.Event `json:"event_pattern" binding:"event_pattern"`

	RRule      string            `json:"rrule,omitempty"`
	Start      *datetime.CpsTime `json:"start,omitempty" swaggertype:"integer"`
	Stop       *datetime.CpsTime `json:"stop,omitempty" swaggertype:"integer"`
	Exdates    []exdate.Request  `json:"exdates" binding:"dive"`
	Exceptions []string          `json:"exceptions"`
}

type Response struct {
	ID           string                                        `bson:"_id" json:"_id"`
	Author       *author.Author                                `bson:"author" json:"author" swaggerignore:"true"`
	Description  string                                        `bson:"description" json:"description"`
	Type         string                                        `bson:"type" json:"type"`
	Priority     int64                                         `bson:"priority" json:"priority"`
	Enabled      bool                                          `bson:"enabled" json:"enabled"`
	Config       eventfilter.RuleConfig                        `bson:"config" json:"config"`
	ExternalData map[string]eventfilter.ExternalDataParameters `bson:"external_data" json:"external_data,omitempty"`
	Created      *datetime.CpsTime                             `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated      *datetime.CpsTime                             `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
	RRule        string                                        `bson:"rrule" json:"rrule"`
	Start        *datetime.CpsTime                             `bson:"start,omitempty" json:"start,omitempty" swaggertype:"integer"`
	Stop         *datetime.CpsTime                             `bson:"stop,omitempty" json:"stop,omitempty" swaggertype:"integer"`
	Exdates      []types.Exdate                                `bson:"exdates" json:"exdates"`
	Exceptions   []Exception                                   `bson:"exceptions" json:"exceptions"`

	EventsCount         int64 `bson:"events_count" json:"events_count"`
	FailuresCount       int64 `bson:"failures_count" json:"failures_count"`
	UnreadFailuresCount int64 `bson:"unread_failures_count" json:"unread_failures_count"`

	EventPattern                     pattern.Event `bson:"event_pattern" json:"event_pattern"`
	savedpattern.EntityPatternFields `bson:",inline"`
}

type Exception struct {
	ID          string           `bson:"_id" json:"_id"`
	Name        string           `bson:"name" json:"name"`
	Description string           `bson:"description" json:"description"`
	Exdates     []types.Exdate   `bson:"exdates" json:"exdates"`
	Created     datetime.CpsTime `bson:"created" json:"created" swaggertype:"integer"`
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
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id author.name author.display_name priority created updated on_success on_failure"`
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

type FailureRequest struct {
	pagination.Query
	Type *int `json:"type" form:"type"`
}

type FailureResponse struct {
	ID        string           `bson:"_id" json:"_id"`
	Type      int64            `bson:"type" json:"type"`
	Timestamp datetime.CpsTime `bson:"t" json:"t" swaggertype:"integer"`
	Message   string           `bson:"message" json:"message"`
	Event     map[string]any   `bson:"event" json:"event" swaggertype:"object"`
	Unread    bool             `bson:"unread" json:"unread"`
}

type AggregationFailureResult struct {
	Data       []FailureResponse `bson:"data" json:"data"`
	TotalCount int64             `bson:"total_count" json:"total_count"`
}

func (r *AggregationFailureResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationFailureResult) GetTotal() int64 {
	return r.TotalCount
}
