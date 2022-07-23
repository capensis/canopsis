package eventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type EditRequest struct {
	Author       string                                        `bson:"author" json:"author" swaggerignore:"true"`
	Description  string                                        `bson:"description" json:"description" binding:"required,max=255"`
	Type         string                                        `bson:"type" json:"type" binding:"required,oneof=break drop enrichment change_entity"`
	Priority     int                                           `bson:"priority" json:"priority"`
	Enabled      bool                                          `bson:"enabled" json:"enabled"`
	Config       eventfilter.RuleConfig                        `bson:"config" json:"config" binding:"dive"`
	ExternalData map[string]eventfilter.ExternalDataParameters `bson:"external_data" json:"external_data,omitempty"`
	Created      *types.CpsTime                                `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated      *types.CpsTime                                `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`

	common.EntityPatternFieldsRequest
	EventPattern pattern.Event `json:"event_pattern" bson:"event_pattern" binding:"event_pattern"`
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
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id author priority created updated on_success on_failure"`
}
