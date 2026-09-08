package entity

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type PaginatedListRequest struct {
	pagination.Query
	ListRequest
}

type ListRequest struct {
	Search       string   `form:"search" binding:"max=10000"`
	SearchBy     []string `form:"search_by"`
	Sort         string   `form:"sort"`
	Enabled      *bool    `form:"enabled"`
	Types        []string `form:"type" binding:"dive,oneof=connector component resource service"`
	Categories   []string `form:"category" binding:"dive,max=50"`
	Connectors   []string `form:"connector" binding:"dive,max=10000"`
	Components   []string `form:"component" binding:"dive,max=10000"`
	Upstreams    []string `form:"upstream" binding:"dive,max=10000"`
	ImportSource string   `form:"import_source" binding:"max=10000"`
}

type PaginatedSearchRequest struct {
	pagination.Query
	SearchRequest
}

type SearchRequest struct {
	EntityPattern    pattern.Entity        `json:"entity_pattern" bson:"entity_pattern" binding:"entity_pattern"`
	PbehaviorPattern pattern.PbehaviorInfo `json:"pbehavior_pattern" bson:"pbehavior_pattern" binding:"pbehavior_pattern"`
	Sort             string                `json:"sort" bson:"sort"`
	Enabled          *bool                 `json:"enabled" bson:"enabled"`
	ImportSource     string                `json:"import_source" bson:"import_source" binding:"max=10000"`
}

type Response struct {
	ID            string              `bson:"_id" json:"_id"`
	Name          string              `bson:"name" json:"name"`
	Description   string              `bson:"description" json:"description"`
	Enabled       bool                `bson:"enabled" json:"enabled"`
	Infos         map[string]Info     `bson:"infos" json:"infos"`
	Type          string              `bson:"type" json:"type"`
	Category      *response.Reference `bson:"category" json:"category"`
	Created       datetime.CpsTime    `bson:"created" json:"created" swaggertype:"integer"`
	Updated       *datetime.CpsTime   `bson:"updated,omitempty" json:"updated" swaggertype:"integer"`
	LastEventDate *datetime.CpsTime   `bson:"last_event_date" json:"last_event_date" swaggertype:"integer"`
	ImpactLevel   int                 `json:"impact_level" bson:"impact_level"`

	PbehaviorInfo *PbehaviorInfo `bson:"pbehavior_info" json:"pbehavior_info"`

	Connector string `bson:"connector,omitempty" json:"connector,omitempty"`

	Component      string           `bson:"component,omitempty" json:"component,omitempty"`
	ComponentInfos *map[string]Info `bson:"component_infos" json:"component_infos,omitempty"`

	Upstream *string `bson:"upstream,omitempty" json:"upstream,omitempty"`

	OutputTemplate *string `bson:"output_template,omitempty" json:"output_template,omitempty"`

	Coordinates *types.Coordinates `bson:"coordinates" json:"coordinates"`

	EntityPattern pattern.Entity `bson:"entity_pattern" json:"entity_pattern,omitempty"`

	ImportSource string            `bson:"import_source,omitempty" json:"import_source"`
	Imported     *datetime.CpsTime `bson:"imported,omitempty" json:"imported" swaggertype:"integer"`
}

func (r *Response) SetDefaults() {
	if r.Type == types.EntityTypeResource && r.ComponentInfos == nil {
		r.ComponentInfos = new(map[string]Info{})
	}

	if r.Type == types.EntityTypeService && r.OutputTemplate == nil {
		r.OutputTemplate = new(string)
	}
}

type Info struct {
	Description string `bson:"description" json:"description"`
	Value       any    `bson:"value" json:"value"`
}

type PbehaviorInfo struct {
	Timestamp *datetime.CpsTime  `bson:"timestamp" json:"timestamp" swaggertype:"integer"`
	Rule      response.Reference `bson:"rule" json:"rule"`
	Reason    response.Reference `bson:"reason" json:"reason"`
	Type      struct {
		response.Reference `bson:",inline"`
		CanonicalType      string `bson:"canonical_type" json:"canonical_type"`
	} `bson:"type" json:"type"`
}
