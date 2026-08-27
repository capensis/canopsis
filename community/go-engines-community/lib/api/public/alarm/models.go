package alarm

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	ResolutionOpen     = "open"
	ResolutionResolved = "resolved"
)

const (
	TypeStandalone = "standalone"
	TypeMeta       = "meta"
	TypeGrouped    = "grouped"
)

type PaginatedListRequest struct {
	pagination.Query
	ListRequest
}

type ListRequest struct {
	Search   string   `form:"search" binding:"max=10000"`
	SearchBy []string `form:"search_by"`
	Sort     string   `form:"sort"`

	Resolution string            `form:"resolution" binding:"required,oneof=open resolved"`
	Types      []string          `form:"type" binding:"unique,dive,oneof=meta grouped standalone"`
	MetaAlarms []string          `form:"meta_alarm" binding:"dive,max=10000"`
	TimeField  string            `form:"time_field" binding:"required_with=From To,oneoforempty=creation_date resolved activation_date last_event_date last_update_date last_st_upd_dt"`
	From       *datetime.CpsTime `form:"from" swaggertype:"integer"`
	To         *datetime.CpsTime `form:"to"  swaggertype:"integer"`
	Tags       []string          `form:"tag"`

	EntityTypes      []string `form:"entity_type" binding:"dive,oneof=connector component resource service"`
	EntityCategories []string `form:"entity_category" binding:"dive,max=50"`
	EntityConnectors []string `form:"entity_connector" binding:"dive,max=10000"`
	EntityComponents []string `form:"entity_component" binding:"dive,max=10000"`
	EntityUpstreams  []string `form:"entity_upstream" binding:"dive,max=10000"`
}

type PaginatedSearchRequest struct {
	pagination.Query
	SearchRequest
}

type SearchRequest struct {
	AlarmPattern     pattern.Alarm         `json:"alarm_pattern" bson:"alarm_pattern" binding:"alarm_pattern"`
	EntityPattern    pattern.Entity        `json:"entity_pattern" bson:"entity_pattern" binding:"entity_pattern"`
	PbehaviorPattern pattern.PbehaviorInfo `json:"pbehavior_pattern" bson:"pbehavior_pattern" binding:"pbehavior_pattern"`
	Sort             string                `json:"sort" bson:"sort"`
	Resolution       string                `json:"resolution" binding:"required,oneof=open resolved"`
	Types            []string              `json:"types" binding:"unique,dive,oneof=meta grouped standalone"`
	MetaAlarms       []string              `form:"meta_alarm" binding:"dive,max=10000"`
}

type StepRequest struct {
	pagination.Query
	Sort string `form:"sort"`

	Types      []string          `form:"type"`
	Initiators []string          `form:"initiator" binding:"unique,dive,oneof=user external system"`
	From       *datetime.CpsTime `form:"from" swaggertype:"integer"`
	To         *datetime.CpsTime `form:"to"  swaggertype:"integer"`
}

type Response struct {
	ID            string          `bson:"_id" json:"_id"`
	DisplayName   string          `bson:"display_name" json:"display_name"`
	Connector     string          `bson:"connector" json:"connector"`
	ConnectorName string          `bson:"connector_name" json:"connector_name"`
	Component     string          `bson:"component" json:"component"`
	Resource      string          `bson:"resource" json:"resource,omitempty"`
	Entity        entity.Response `bson:"entity" json:"entity"`

	ACK         *Step `bson:"ack" json:"ack"`
	Canceled    *Step `bson:"canceled" json:"canceled"`
	Snooze      *Step `bson:"snooze" json:"snooze"`
	State       *Step `bson:"state" json:"state"`
	Status      *Step `bson:"status" json:"status"`
	ChangeState *Step `bson:"change_state" json:"change_state"`
	Ticket      *Step `bson:"ticket" json:"ticket"`
	LastComment *Step `bson:"last_comment" json:"last_comment"`

	CreationDate                datetime.CpsTime  `bson:"creation_date" json:"creation_date" swaggertype:"integer"`
	Resolved                    *datetime.CpsTime `bson:"resolved" json:"resolved" swaggertype:"integer"`
	ActivationDate              *datetime.CpsTime `bson:"activation_date" json:"activation_date" swaggertype:"integer"`
	LastUpdateDate              datetime.CpsTime  `bson:"last_update_date" json:"last_update_date" swaggertype:"integer"`
	LastEventDate               *datetime.CpsTime `bson:"last_event_date" json:"last_event_date" swaggertype:"integer"`
	LastStateOrStatusUpdateDate datetime.CpsTime  `bson:"last_st_upd_dt" json:"last_st_upd_dt" swaggertype:"integer"`

	InitialOutput     string `bson:"initial_output" json:"initial_output"`
	Output            string `bson:"output" json:"output"`
	InitialLongOutput string `bson:"initial_long_output" json:"initial_long_output"`
	LongOutput        string `bson:"long_output" json:"long_output"`

	StateChangesSinceStatusUpdate types.CpsNumber `bson:"state_changes_since_status_update" json:"state_changes_since_status_update"`
	TotalStateChanges             types.CpsNumber `bson:"total_state_changes" json:"total_state_changes"`

	Duration             int   `bson:"duration" json:"duration"`
	CurrentStateDuration int   `bson:"current_state_duration" json:"current_state_duration"`
	ActiveDuration       int64 `bson:"active_duration" json:"active_duration"`

	MaxState     *types.CpsNumber `bson:"max_state" json:"max_state"`
	InitialState *types.CpsNumber `bson:"initial_state" json:"initial_state"`

	Tags  []string                  `bson:"tags" json:"tags"`
	Infos map[string]map[string]any `bson:"infos" json:"infos" swaggertype:"object"`

	PbehaviorInfo *entity.PbehaviorInfo `bson:"pbehavior_info" json:"pbehavior_info"`

	MetaAlarmRule *response.Reference `bson:"meta_alarm_rule" json:"meta_alarm_rule"`
	Type          string              `bson:"type" json:"type"`
}

func (r *Response) setDefaults() {
	r.Entity.SetDefaults()
}

type Step struct {
	Type      string            `bson:"_t" json:"_t"`
	Timestamp *datetime.CpsTime `bson:"t" json:"t" swaggertype:"integer"`
	Author    string            `bson:"a" json:"a"`
	UserID    string            `bson:"user_id,omitempty" json:"user_id"`
	Message   string            `bson:"m" json:"m"`
	Value     types.CpsNumber   `bson:"val" json:"val"`
	Initiator string            `bson:"initiator" json:"initiator"`

	Ticket             string            `bson:"ticket" json:"ticket,omitempty"`
	TicketURL          string            `bson:"ticket_url" json:"ticket_url,omitempty"`
	TicketComment      string            `bson:"ticket_comment" json:"ticket_comment,omitempty"`
	TicketSystemName   string            `bson:"ticket_system_name" json:"ticket_system_name,omitempty"`
	TicketData         map[string]string `bson:"ticket_data" json:"ticket_data,omitempty"`
	TicketStatus       int               `bson:"ticket_status" json:"ticket_status,omitempty"`
	TicketSourceStatus string            `bson:"ticket_source_status" json:"ticket_source_status,omitempty"`

	StructuredMessage []types.StructuredMessage `bson:"struct_m" json:"struct_m,omitempty"`
}
