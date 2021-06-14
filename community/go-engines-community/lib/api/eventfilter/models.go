package eventfilter

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type EventFilterPayload struct {
	Author      string `bson:"author" json:"author" swaggerignore:"true"`
	Description string `bson:"description" json:"description" binding:"required,max=255"`
	Type        string `bson:"type" json:"type" binding:"required,oneof=break drop enrichment"`

	Patterns *pattern.EventPatternList `bson:"patterns" json:"patterns"`

	Priority int   `bson:"priority" json:"priority"`
	Enabled  *bool `bson:"enabled" json:"enabled" binding:"required"`

	Actions      []eventfilter.Action   `bson:"actions,omitempty" json:"actions,omitempty" binding:"required_if=Type enrichment"`
	ExternalData map[string]interface{} `bson:"external_data,omitempty" json:"external_data,omitempty"`
	OnSuccess    string                 `bson:"on_success,omitempty" json:"on_success,omitempty" binding:"required_if=Type enrichment"`
	OnFailure    string                 `bson:"on_failure,omitempty" json:"on_failure,omitempty" binding:"required_if=Type enrichment"`

	Created *types.CpsTime `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated *types.CpsTime `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type EventFilter struct {
	ID                 string `bson:"_id" json:"_id" binding:"id"`
	EventFilterPayload `bson:",inline"`
}

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id author priority created updated on_success on_failure"`
}
