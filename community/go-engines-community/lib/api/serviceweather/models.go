package serviceweather

import (
	alarmapi "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entitycategory"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviortype"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequest struct {
	EntitiesListRequest
	Category string `form:"category" json:"category"`
}

type EntitiesListRequest struct {
	pagination.FilteredQuery
	WithInstructions bool   `form:"with_instructions"`
	SortBy           string `form:"sort_by" json:"sort_by" binding:"oneoforempty=name state infos.* impact_state"`
}

type Service struct {
	ID             string              `json:"_id" bson:"_id"`
	Name           string              `json:"name" bson:"name"`
	Infos          map[string]Info     `json:"infos" bson:"infos"`
	Connector      string              `json:"connector" bson:"connector"`
	ConnectorName  string              `json:"connector_name" bson:"connector_name"`
	Component      string              `json:"component" bson:"component"`
	Resource       string              `json:"resource" bson:"resource"`
	HasOpenAlarm   bool                `json:"is_action_required" bson:"has_open_alarm"`
	State          alarmapi.AlarmStep  `json:"state" bson:"state"`
	Status         alarmapi.AlarmStep  `json:"status" bson:"status"`
	Snooze         *alarmapi.AlarmStep `json:"snooze" bson:"snooze"`
	Ack            *alarmapi.AlarmStep `json:"ack" bson:"ack"`
	Icon           string              `json:"icon" bson:"icon"`
	SecondaryIcon  string              `json:"secondary_icon" bson:"secondary_icon"`
	Output         string              `json:"output" bson:"output"`
	LastUpdateDate *types.CpsTime      `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	AlarmCounters  []AlarmCounter      `json:"alarm_counters" bson:"alarm_counters"`
	Links          []struct {
		Name  string      `json:"cat_name" bson:"cat_name"`
		Links interface{} `json:"links" bson:"links"`
	} `json:"linklist" bson:"links"`
	Pbehaviors  []pbehavior.Response     `json:"pbehaviors" bson:"-"`
	PbehaviorID string                   `json:"-" bson:"pbehavior_id"`
	ImpactLevel int                      `json:"impact_level" bson:"impact_level"`
	ImpactState int                      `json:"impact_state" bson:"impact_state"`
	Category    *entitycategory.Category `json:"category" bson:"category"`
	IsGrey      bool                     `json:"is_grey" bson:"is_grey"`
	IdleSince   *types.CpsTime           `json:"idle_since,omitempty" bson:"idle_since,omitempty" swaggertype:"integer"`
}

type Info struct {
	Description string      `bson:"description,omitempty" json:"description"`
	Value       interface{} `bson:"value" json:"value"`
}

type AlarmCounter struct {
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
	ID                               string                            `json:"_id" bson:"_id"`
	AlarmID                          string                            `json:"alarm_id" bson:"alarm_id"`
	AssignedInstructions             *[]alarmapi.InstructionWithAlarms `bson:"-" json:"assigned_instructions,omitempty"`
	IsAutoInstructionRunning         *bool                             `bson:"-" json:"is_auto_instruction_running,omitempty"`
	IsAllAutoInstructionsCompleted   *bool                             `bson:"-" json:"is_all_auto_instructions_completed,omitempty"`
	IsAutoInstructionFailed          *bool                             `bson:"-" json:"is_auto_instruction_failed,omitempty"`
	IsManualInstructionWaitingResult *bool                             `bson:"-" json:"is_manual_instruction_waiting_result,omitempty"`
	Name                             string                            `json:"name" bson:"name"`
	Infos                            map[string]Info                   `json:"infos" bson:"infos"`
	Type                             string                            `json:"source_type" bson:"type"`
	Category                         *entitycategory.Category          `json:"category" bson:"category"`
	Connector                        string                            `json:"connector" bson:"connector"`
	ConnectorName                    string                            `json:"connector_name" bson:"connector_name"`
	Component                        string                            `json:"component" bson:"component"`
	Resource                         string                            `json:"resource" bson:"resource"`
	State                            alarmapi.AlarmStep                `json:"state" bson:"state"`
	Status                           alarmapi.AlarmStep                `json:"status" bson:"status"`
	Snooze                           *alarmapi.AlarmStep               `json:"snooze" bson:"snooze"`
	Ack                              *alarmapi.AlarmStep               `json:"ack" bson:"ack"`
	Ticket                           *alarmapi.AlarmTicket             `json:"ticket" bson:"ticket"`
	LastUpdateDate                   *types.CpsTime                    `json:"last_update_date" bson:"last_update_date" swaggertype:"integer"`
	CreationDate                     *types.CpsTime                    `json:"alarm_creation_date" bson:"creation_date" swaggertype:"integer"`
	DisplayName                      string                            `json:"alarm_display_name" bson:"display_name"`
	Icon                             string                            `json:"icon" bson:"icon"`
	Pbehaviors                       []pbehavior.Response              `json:"pbehaviors" bson:"-"`
	PbehaviorInfo                    types.PbehaviorInfo               `json:"-" bson:"pbehavior_info"`
	Links                            []WeatherLink                     `json:"linklist" bson:"-"`
	Stats                            Stats                             `json:"stats" bson:"-"`
	IsGrey                           bool                              `json:"is_grey"`
	ImpactLevel                      int                               `json:"impact_level" bson:"impact_level"`
	ImpactState                      int                               `json:"impact_state" bson:"impact_state"`
	IdleSince                        *types.CpsTime                    `json:"idle_since,omitempty" bson:"idle_since,omitempty" swaggertype:"integer"`
}

type WeatherLink struct {
	Category string      `json:"cat_name"`
	Links    interface{} `json:"links"`
}

type Stats struct {
	FailEventsCount int            `json:"ko" bson:"ko"`
	OKEventsCount   int            `json:"ok" bson:"ok"`
	LastFailEvent   *types.CpsTime `json:"last_ko" bson:"last_ko" swaggertype:"integer"`
	LastEvent       *types.CpsTime `json:"last_event" bson:"last_event" swaggertype:"integer"`
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
