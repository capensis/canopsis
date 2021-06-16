package alarm

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type ListRequestWithPagination struct {
	pagination.Query
	ListRequest
}

type ListRequest struct {
	FilterRequest
	WithSteps    bool   `form:"with_steps" json:"with_steps"`
	WithChildren bool   `form:"with_consequences" json:"with_consequences"`
	MultiSort    string `form:"multi_sort" json:"multi_sort"`
	Sort         string `form:"sort_dir" json:"sort_dir" binding:"oneoforempty=asc desc"`
	SortBy       string `form:"sort_key" json:"sort_key"`
}

type FilterRequest struct {
	Filter              string         `form:"filter" json:"filter"`
	Search              string         `form:"search" json:"search"`
	SearchBy            []string       `form:"active_columns[]" json:"active_columns[]"`
	StartFrom           *types.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo             *types.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
	OnlyOpened          bool           `form:"opened" json:"opened"`
	OnlyResolved        bool           `form:"resolved" json:"resolved"`
	OnlyParents         bool           `form:"correlation" json:"correlation"`
	OnlyManual          bool           `form:"manual" json:"manual"`
	WithInstructions    string         `form:"with_instructions" json:"with_instructions"`
	WithoutInstructions string         `form:"without_instructions" json:"without_instructions"`
	Category            string         `form:"category" json:"category"`
}

type ExportRequest struct {
	ListRequest
	Separator  string `form:"separator" json:"separator" binding:"oneoforempty=comma semicolon tab space"`
	TimeFormat string `form:"time_format" json:"time_format" binding:"time_format"`
}

type ExportResponse struct {
	ID     string `json:"_id"`
	Status int    `json:"status"`
}

type Alarm struct {
	ID            string                            `bson:"_id" json:"_id"`
	Time          types.CpsTime                     `bson:"t" json:"t" swaggertype:"integer"`
	Entity        entity.Entity                     `bson:"entity" json:"entity"`
	Value         AlarmValue                        `bson:"v" json:"v"`
	Infos         map[string]map[string]interface{} `bson:"infos" json:"infos"`
	Pbehavior     *Pbehavior                        `bson:"pbehavior,omitempty" json:"pbehavior,omitempty"`
	MetaAlarmRule *MetaAlarmRule                    `bson:"meta_alarm_rule,omitempty" json:"rule,omitempty"`
	IsMetaAlarm   *bool                             `bson:"is_meta_alarm,omitempty" json:"metaalarm,omitempty"`
	ChildrenIDs   *struct {
		Data  []string `bson:"data"`
		Total int      `bson:"total"`
	} `bson:"children_ids,omitempty" json:"-"`
	Children             *Children               `bson:"children,omitempty" json:"consequences,omitempty"`
	Causes               *Causes                 `bson:"causes,omitempty" json:"causes,omitempty"`
	FilteredChildrenIDs  []string                `bson:"filtered_children_ids,omitempty" json:"filtered_children,omitempty"`
	AssignedInstructions []InstructionWithAlarms `bson:"assigned_instructions,omitempty" json:"assigned_instructions,omitempty"`
	Links                map[string]interface{}  `json:"links"`
	ImpactState          int64                   `bson:"impact_state" json:"impact_state"`
}

type MetaAlarmRule struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type Causes struct {
	Rules []MetaAlarmRule `bson:"rules" json:"rules"`
	Total int             `bson:"total" json:"total"`
}

type Children struct {
	Data  []Alarm `bson:"data,omitempty" json:"data,omitempty"`
	Total int     `bson:"total" json:"total"`
}

type AlarmValue struct {
	ACK                           *AlarmStep             `bson:"ack,omitempty" json:"ack,omitempty"`
	Canceled                      *AlarmStep             `bson:"canceled,omitempty" json:"canceled,omitempty"`
	Done                          *AlarmStep             `bson:"done,omitempty" json:"done,omitempty"`
	Snooze                        *AlarmStep             `bson:"snooze,omitempty" json:"snooze,omitempty"`
	State                         *AlarmStep             `bson:"state,omitempty" json:"state,omitempty"`
	Status                        *AlarmStep             `bson:"status,omitempty" json:"status,omitempty"`
	Ticket                        *AlarmTicket           `bson:"ticket,omitempty" json:"ticket,omitempty"`
	LastComment                   *AlarmStep             `bson:"last_comment,omitempty" json:"lastComment,omitempty"`
	Steps                         []AlarmStep            `bson:"steps,omitempty" json:"steps,omitempty"`
	Component                     string                 `bson:"component" json:"component"`
	Connector                     string                 `bson:"connector" json:"connector"`
	ConnectorName                 string                 `bson:"connector_name" json:"connector_name"`
	CreationDate                  types.CpsTime          `bson:"creation_date" json:"creation_date" swaggertype:"integer"`
	ActivationDate                *types.CpsTime         `bson:"activation_date,omitempty" json:"activation_date,omitempty" swaggertype:"integer"`
	DisplayName                   string                 `bson:"display_name" json:"display_name"`
	InitialOutput                 string                 `bson:"initial_output" json:"initial_output"`
	Output                        string                 `bson:"output" json:"output"`
	InitialLongOutput             string                 `bson:"initial_long_output" json:"initial_long_output"`
	LongOutput                    string                 `bson:"long_output" json:"long_output"`
	LongOutputHistory             []string               `bson:"long_output_history" json:"long_output_history"`
	LastUpdateDate                types.CpsTime          `bson:"last_update_date" json:"last_update_date" swaggertype:"integer"`
	LastEventDate                 types.CpsTime          `bson:"last_event_date" json:"last_event_date" swaggertype:"integer"`
	Resource                      string                 `bson:"resource,omitempty" json:"resource,omitempty"`
	Resolved                      *types.CpsTime         `bson:"resolved,omitempty" json:"resolved,omitempty" swaggertype:"integer"`
	PbehaviorInfo                 *types.PbehaviorInfo   `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	Tags                          []string               `bson:"tags" json:"tags"`
	Meta                          string                 `bson:"meta,omitempty" json:"meta,omitempty"`
	Parents                       []string               `bson:"parents" json:"parents"`
	Children                      []string               `bson:"children" json:"children"`
	StateChangesSinceStatusUpdate types.CpsNumber        `bson:"state_changes_since_status_update,omitempty" json:"state_changes_since_status_update,omitempty"`
	TotalStateChanges             types.CpsNumber        `bson:"total_state_changes,omitempty" json:"total_state_changes,omitempty"`
	Extra                         map[string]interface{} `bson:"extra" json:"extra"`
	RuleVersion                   map[string]string      `bson:"infos_rule_version" json:"infos_rule_version"`
	Duration                      int                    `bson:"duration" json:"duration"`
	CurrentStateDuration          int                    `bson:"current_state_duration" json:"current_state_duration"`
	EventsCount                   types.CpsNumber        `bson:"events_count,omitempty" json:"events_count,omitempty"`

	Infos map[string]map[string]interface{} `bson:"infos" json:"infos"`
}

type AlarmStep struct {
	Type      string          `bson:"_t" json:"_t"`
	Timestamp types.CpsTime   `bson:"t" json:"t" swaggertype:"integer"`
	Author    string          `bson:"a" json:"a"`
	Message   string          `bson:"m" json:"m"`
	Value     types.CpsNumber `bson:"val" json:"val"`
	Initiator string          `bson:"initiator" json:"initiator"`
}

type AlarmTicket struct {
	Type      string            `bson:"_t" json:"_t"`
	Timestamp types.CpsTime     `bson:"t" json:"t" swaggertype:"integer"`
	Author    string            `bson:"a" json:"a"`
	Message   string            `bson:"m" json:"m"`
	Value     string            `bson:"val" json:"val"`
	Data      map[string]string `bson:"data" json:"data"`
}

type Pbehavior struct {
	ID     string          `bson:"_id" json:"_id"`
	Author string          `bson:"author" json:"author"`
	Name   string          `bson:"name" json:"name"`
	RRule  string          `bson:"rrule" json:"rrule"`
	Start  *types.CpsTime  `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop   *types.CpsTime  `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Type   *pbehavior.Type `bson:"type" json:"type"`
}

type InstructionWithAlarms struct {
	ID                   string                    `bson:"_id" json:"_id"`
	AlarmPatterns        pattern.AlarmPatternList  `bson:"alarm_patterns" json:"-"`
	EntityPatterns       pattern.EntityPatternList `bson:"entity_patterns" json:"-"`
	Name                 string                    `bson:"name" json:"name"`
	ActiveOnPbh          []string                  `bson:"active_on_pbh,omitempty" json:"active_on_pbh,omitempty"`
	DisabledOnPbh        []string                  `bson:"disabled_on_pbh,omitempty" json:"disabled_on_pbh,omitempty"`
	Execution            *Execution                `bson:"-" json:"execution"`
	AlarmsWithExecutions []Execution               `bson:"alarms_with_executions" json:"-"`
	Created              types.CpsTime             `bson:"created,omitempty" json:"-"`
}

type Execution struct {
	ID     string `bson:"_id" json:"_id"`
	Status int    `bson:"status" json:"status"`
	Alarm  string `bson:"alarm" json:"-"`
}

func (i InstructionWithAlarms) GetExecution(alarmId string) *Execution {
	for _, v := range i.AlarmsWithExecutions {
		if v.Alarm == alarmId {
			return &v
		}
	}

	return nil
}

type AggregationResult struct {
	Data       []Alarm `bson:"data"`
	TotalCount int64   `bson:"total_count" json:"total_count"`
}

func (r *AggregationResult) GetData() interface{} {
	return r.Data
}

func (r *AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

type Count struct {
	Total          int `bson:"total" json:"total"`
	TotalActive    int `bson:"total_active" json:"total_active"`
	TotalSnooze    int `bson:"total_snooze" json:"snooze"`
	TotalAck       int `bson:"total_ack" json:"ack"`
	TotalTicket    int `bson:"total_ticket" json:"ticket"`
	TotalPbehavior int `bson:"total_pbehavior" json:"pbehavior_active"`
}
