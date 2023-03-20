package entity

import (
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type ListRequestWithPagination struct {
	pagination.Query
	ListRequest
	WithFlags bool `form:"with_flags" json:"with_flags"`
}

type ListRequest struct {
	BaseFilterRequest
	SortRequest
	SearchBy []string `form:"active_columns[]" json:"active_columns[]"`
}

type SortRequest struct {
	Sort   string `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy string `form:"sort_by" json:"sort_by"`
}

type BaseFilterRequest struct {
	Search   string   `form:"search" json:"search"`
	Filters  []string `form:"filters[]" json:"filters"`
	Category string   `form:"category" json:"category"`
	Type     []string `form:"type[]" json:"type"`
	NoEvents bool     `form:"no_events" json:"no_events"`
}

type ExportRequest struct {
	BaseFilterRequest
	Fields    export.Fields `json:"fields"`
	Separator string        `json:"separator" binding:"oneoforempty=comma semicolon tab space"`
}

type CleanRequest struct {
	Archive             *bool `json:"archive" bson:"archive" binding:"required"`
	ArchiveDependencies bool  `json:"archive_dependencies" bson:"archive_dependencies"`
}

type CleanTask struct {
	Archive             *bool
	ArchiveDependencies bool
	UserID              string
}

type ExportResponse struct {
	ID string `json:"_id"`
	// Possible values.
	//   * `0` - Running
	//   * `1` - Succeeded
	//   * `2` - Failed
	Status int `json:"status"`
}

type Entity struct {
	ID             string         `bson:"_id" json:"_id"`
	Name           string         `bson:"name" json:"name"`
	Enabled        bool           `bson:"enabled" json:"enabled"`
	Infos          Infos          `bson:"infos" json:"infos"`
	ComponentInfos Infos          `bson:"component_infos,omitempty" json:"component_infos,omitempty"`
	Type           string         `bson:"type" json:"type"`
	ImpactLevel    int64          `bson:"impact_level" json:"impact_level"`
	Category       *Category      `bson:"category" json:"category"`
	IdleSince      *types.CpsTime `bson:"idle_since,omitempty" json:"idle_since,omitempty" swaggertype:"integer"`
	LastEventDate  *types.CpsTime `bson:"last_event_date,omitempty" json:"last_event_date,omitempty" swaggertype:"integer"`

	PbehaviorInfo     *PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	LastPbehaviorDate *types.CpsTime `bson:"last_pbehavior_date,omitempty" json:"last_pbehavior_date,omitempty" swaggertype:"integer"`

	Connector string `bson:"connector,omitempty" json:"connector,omitempty"`
	Component string `bson:"component,omitempty" json:"component,omitempty"`

	// ConnectorType contains a part before "/" of connector id.
	ConnectorType string `bson:"-" json:"connector_type,omitempty"`

	// Flags
	Deletable *bool `bson:"deletable,omitempty" json:"deletable,omitempty"`

	// Stats
	OKEvents *int `bson:"ok_events" json:"ok_events,omitempty"`
	KOEvents *int `bson:"ko_events" json:"ko_events,omitempty"`

	// Alarm fields
	State               *int              `bson:"state" json:"state,omitempty"`
	ImpactState         *int              `bson:"impact_state" json:"impact_state,omitempty"`
	Status              *int              `bson:"status" json:"status,omitempty"`
	Ack                 *common.AlarmStep `bson:"ack" json:"ack,omitempty"`
	Snooze              *common.AlarmStep `bson:"snooze" json:"snooze,omitempty"`
	AlarmLastUpdateDate *types.CpsTime    `bson:"alarm_last_update_date" json:"alarm_last_update_date,omitempty" swaggertype:"integer"`

	// Links
	Links link.LinksByCategory `bson:"-" json:"links,omitempty"`

	// DependsCount contains only service's dependencies
	DependsCount *int `bson:"depends_count" json:"depends_count,omitempty"`
	// ImpactsCount contains only services
	ImpactsCount *int `bson:"impacts_count" json:"impacts_count,omitempty"`

	Coordinates *types.Coordinates `bson:"coordinates,omitempty" json:"coordinates,omitempty"`

	OldEntityPatterns                oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty" json:"old_entity_patterns,omitempty"`
	savedpattern.EntityPatternFields `bson:",inline"`

	Resources []string `bson:"resources,omitempty" json:"-"`

	ImportSource string         `bson:"import_source,omitempty" json:"import_source,omitempty"`
	Imported     *types.CpsTime `bson:"imported,omitempty" json:"imported,omitempty" swaggertype:"integer"`
}

func (e *Entity) fillConnectorType() {
	if e.Type == types.EntityTypeConnector {
		e.ConnectorType = strings.TrimSuffix(e.ID, "/"+e.Name)
	}
}

type Category struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type Infos map[string]Info

type Info struct {
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

func (i *Infos) UnmarshalBSONValue(_ bsontype.Type, b []byte) error {
	var tmp map[string]Info
	err := bson.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	for k, info := range tmp {
		if info.Name == "" {
			info.Name = k
			tmp[k] = info
		}
	}

	*i = tmp

	return nil
}

type PbehaviorInfo struct {
	types.PbehaviorInfo `bson:",inline"`

	IconName string `bson:"icon_name" json:"icon_name"`
}

type AggregationResult struct {
	Data       []Entity `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

func (r AggregationResult) GetData() interface{} {
	return r.Data
}

type BulkToggleRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

type SimplifiedEntity struct {
	ID        string   `bson:"_id"`
	Type      string   `bson:"type"`
	Enabled   bool     `bson:"enabled"`
	Resources []string `bson:"resources"`
}

type ContextGraphRequest struct {
	ID string `form:"_id" binding:"required"`
}

type ContextGraphResponse struct {
	Impacts []string `bson:"impact" json:"impact"`
	Depends []string `bson:"depends" json:"depends"`
}
