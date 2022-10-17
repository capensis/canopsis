package serviceweather

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statistics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	pagination.Query
	Filters  []string `form:"filters[]" json:"filters"`
	Category string   `form:"category" json:"category"`
	Sort     string   `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy   string   `form:"sort_by" json:"sort_by" binding:"oneoforempty=name state infos.* impact_state"`
}

type EntitiesListRequest struct {
	pagination.Query
	WithInstructions bool   `form:"with_instructions"`
	Sort             string `form:"sort"`
	SortBy           string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name state infos.* impact_state"`
}

type Service struct {
	ID             string                   `json:"_id" bson:"_id"`
	Name           string                   `json:"name" bson:"name"`
	Infos          map[string]Info          `json:"infos" bson:"infos"`
	Connector      string                   `json:"connector" bson:"connector"`
	ConnectorName  string                   `json:"connector_name" bson:"connector_name"`
	Component      string                   `json:"component" bson:"component"`
	Resource       string                   `json:"resource" bson:"resource"`
	HasOpenAlarm   bool                     `json:"is_action_required" bson:"has_open_alarm"`
	State          common.AlarmStep         `json:"state" bson:"state"`
	Status         common.AlarmStep         `json:"status" bson:"status"`
	Snooze         *common.AlarmStep        `json:"snooze" bson:"snooze"`
	Ack            *common.AlarmStep        `json:"ack" bson:"ack"`
	Icon           string                   `json:"icon" bson:"icon"`
	SecondaryIcon  string                   `json:"secondary_icon" bson:"secondary_icon"`
	Output         string                   `json:"output" bson:"output"`
	LastUpdateDate *types.CpsTime           `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	Counters       Counters                 `json:"counters" bson:"counters"`
	PbehaviorInfo  *entity.PbehaviorInfo    `json:"pbehavior_info" bson:"pbehavior_info"`
	Pbehaviors     []alarm.Pbehavior        `json:"pbehaviors" bson:"pbehaviors"`
	ImpactLevel    int                      `json:"impact_level" bson:"impact_level"`
	ImpactState    int                      `json:"impact_state" bson:"impact_state"`
	Category       *entitycategory.Category `json:"category" bson:"category"`
	IsGrey         bool                     `json:"is_grey" bson:"is_grey"`
	IdleSince      *types.CpsTime           `json:"idle_since,omitempty" bson:"idle_since,omitempty" swaggertype:"integer"`

	Links []WeatherLink `json:"linklist" bson:"-"`
}

type Info struct {
	Description string      `bson:"description,omitempty" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

type Counters struct {
	All    int64 `bson:"all" json:"all"`
	Alarms int64 `bson:"active" json:"active"`
	State  struct {
		Critical int64 `bson:"critical" json:"critical"`
		Major    int64 `bson:"major" json:"major"`
		Minor    int64 `bson:"minor" json:"minor"`
		Ok       int64 `bson:"ok" json:"ok"`
	} `bson:"state" json:"state"`
	Acknowledged            int64            `bson:"acked" json:"acked"`
	NotAcknowledged         int64            `bson:"unacked" json:"unacked"`
	NotAcknowledgedUnderPbh int64            `bson:"unacked_under_pbh" json:"unacked_under_pbh"`
	UnderPbehavior          int64            `bson:"under_pbh" json:"under_pbh"`
	Depends                 int64            `bson:"depends" json:"depends"`
	PbhTypeCounters         []PbhTypeCounter `bson:"pbh_types" json:"pbh_types"`
}

type PbhTypeCounter struct {
	Count int64              `json:"count" bson:"count"`
	Type  pbehaviortype.Type `json:"type" bson:"type"`
}

type AggregationResult struct {
	Data       []Service `bson:"data" json:"data"`
	TotalCount int64     `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type Entity struct {
	ID      string `json:"_id" bson:"_id"`
	AlarmID string `json:"alarm_id" bson:"alarm_id"`

	AssignedInstructions             *[]alarm.AssignedInstruction `bson:"-" json:"assigned_instructions,omitempty"`
	IsAutoInstructionRunning         *bool                        `bson:"-" json:"is_auto_instruction_running,omitempty"`
	IsAllAutoInstructionsCompleted   *bool                        `bson:"-" json:"is_all_auto_instructions_completed,omitempty"`
	IsAutoInstructionFailed          *bool                        `bson:"-" json:"is_auto_instruction_failed,omitempty"`
	IsManualInstructionRunning       *bool                        `bson:"-" json:"is_manual_instruction_running,omitempty"`
	IsManualInstructionWaitingResult *bool                        `bson:"-" json:"is_manual_instruction_waiting_result,omitempty"`

	Name           string                     `json:"name" bson:"name"`
	Infos          map[string]Info            `json:"infos" bson:"infos"`
	Type           string                     `json:"source_type" bson:"type"`
	Category       *entitycategory.Category   `json:"category" bson:"category"`
	Connector      string                     `json:"connector" bson:"connector"`
	ConnectorName  string                     `json:"connector_name" bson:"connector_name"`
	Component      string                     `json:"component" bson:"component"`
	Resource       string                     `json:"resource" bson:"resource"`
	State          common.AlarmStep           `json:"state" bson:"state"`
	Status         common.AlarmStep           `json:"status" bson:"status"`
	Snooze         *common.AlarmStep          `json:"snooze" bson:"snooze"`
	Ack            *common.AlarmStep          `json:"ack" bson:"ack"`
	Ticket         *alarm.AlarmTicket         `json:"ticket" bson:"ticket"`
	LastUpdateDate *types.CpsTime             `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	CreationDate   *types.CpsTime             `json:"alarm_creation_date" bson:"creation_date" swaggertype:"integer"`
	DisplayName    string                     `json:"alarm_display_name" bson:"display_name"`
	Icon           string                     `json:"icon" bson:"icon"`
	Pbehaviors     []alarm.Pbehavior          `json:"pbehaviors" bson:"pbehaviors"`
	PbehaviorInfo  *entity.PbehaviorInfo      `json:"pbehavior_info" bson:"pbehavior_info"`
	IsGrey         bool                       `json:"is_grey" bson:"is_grey"`
	ImpactLevel    int                        `json:"impact_level" bson:"impact_level"`
	ImpactState    int                        `json:"impact_state" bson:"impact_state"`
	IdleSince      *types.CpsTime             `json:"idle_since,omitempty" bson:"idle_since,omitempty" swaggertype:"integer"`
	Stats          statistics.EventStatistics `json:"stats" bson:"stats"`

	Links []WeatherLink `json:"linklist" bson:"-"`

	DependsCount int `bson:"depends_count" json:"depends_count"`
}

type WeatherLink struct {
	Category string      `json:"cat_name"`
	Links    interface{} `json:"links"`
}

type EntityAggregationResult struct {
	Data       []Entity `bson:"data" json:"data"`
	TotalCount int64    `bson:"total_count" json:"total_count"`
}

func (r *EntityAggregationResult) GetData() interface{} {
	return r.Data
}

func (r *EntityAggregationResult) GetTotal() int64 {
	return r.TotalCount
}
