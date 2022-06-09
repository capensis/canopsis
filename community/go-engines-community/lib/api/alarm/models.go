package alarm

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	OnlyOpened = iota
	OpenedAndRecentResolved
	OnlyResolved
)

type ListRequestWithPagination struct {
	pagination.Query
	ListRequest
}

type ListRequest struct {
	FilterRequest
	SortRequest
	WithInstructions bool `form:"with_instructions" json:"with_instructions"`
	WithLinks        bool `form:"with_links" json:"with_links"`
}

type FilterRequest struct {
	BaseFilterRequest
	SearchBy []string `form:"active_columns[]" json:"active_columns[]"`
}

type BaseFilterRequest struct {
	Filter      string         `form:"filter" json:"filter"`
	Search      string         `form:"search" json:"search"`
	TimeField   string         `form:"time_field" json:"time_field" binding:"oneoforempty=t v.creation_date v.resolved v.last_update_date v.last_event_date"`
	StartFrom   *types.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo     *types.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
	Opened      *bool          `form:"opened" json:"opened"`
	OnlyParents bool           `form:"correlation" json:"correlation"`
	Category    string         `form:"category" json:"category"`

	IncludeInstructionTypes []int    `form:"include_instruction_types[]" json:"include_instruction_types"`
	ExcludeInstructionTypes []int    `form:"exclude_instruction_types[]" json:"exclude_instruction_types"`
	IncludeInstructions     []string `form:"include_instructions[]" json:"include_instructions"`
	ExcludeInstructions     []string `form:"exclude_instructions[]" json:"exclude_instructions"`
}

func (r FilterRequest) GetOpenedFilter() int {
	if r.Opened == nil {
		return OpenedAndRecentResolved
	}

	if *r.Opened {
		return OnlyOpened
	}

	return OnlyResolved
}

type ListByServiceRequest struct {
	pagination.Query
	SortRequest
}

type ListByComponentRequest struct {
	pagination.Query
	SortRequest
	ID string `form:"_id" json:"_id" binding:"required"`
}

type ResolvedListRequest struct {
	pagination.Query
	SortRequest
	ID        string         `form:"_id" json:"_id" binding:"required"`
	StartFrom *types.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo   *types.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
}

type SortRequest struct {
	MultiSort []string `form:"multi_sort[]" json:"multi_sort"`
	Sort      string   `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy    string   `form:"sort_by" json:"sort_by"`
}

type ManualRequest struct {
	Search string `form:"search" json:"search"`
}

type ManualResponse struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type DetailsRequest struct {
	ID               string               `json:"_id" binding:"required"`
	Opened           *bool                `json:"opened"`
	WithInstructions bool                 `json:"with_instructions"`
	Steps            *pagination.Query    `json:"steps"`
	Children         *ChildDetailsRequest `json:"children"`
}

type ChildDetailsRequest struct {
	pagination.Query
	SortRequest
}

func (r DetailsRequest) GetOpenedFilter() int {
	if r.Opened == nil {
		return OpenedAndRecentResolved
	}

	if *r.Opened {
		return OnlyOpened
	}

	return OnlyResolved
}

type DetailsResponse struct {
	ID     string            `json:"_id"`
	Status int               `json:"status"`
	Data   Details           `json:"data,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
	Error  string            `json:"error,omitempty"`
}

type Details struct {
	Steps    *StepDetails     `bson:"steps" json:"steps,omitempty"`
	Children *ChildrenDetails `bson:"children" json:"children,omitempty"`

	IsMetaAlarm bool   `json:"-" bson:"is_meta_alarm"`
	EntityID    string `json:"-" bson:"d"`
	StepsCount  int64  `json:"-" bson:"steps_count"`
}

type StepDetails struct {
	Data []AlarmStep          `json:"data"`
	Meta common.PaginatedMeta `json:"meta"`
}

type ChildrenDetails struct {
	Data []Alarm              `json:"data"`
	Meta common.PaginatedMeta `json:"meta"`
}

type ExportRequest struct {
	BaseFilterRequest
	Fields     export.Fields `json:"fields"`
	Separator  string        `json:"separator" binding:"oneoforempty=comma semicolon tab space"`
	TimeFormat string        `json:"time_format" binding:"time_format"`
}

type ExportResponse struct {
	ID     string `json:"_id"`
	Status int    `json:"status"`
}

type Alarm struct {
	ID     string                            `bson:"_id" json:"_id"`
	Time   types.CpsTime                     `bson:"t" json:"t" swaggertype:"integer"`
	Entity entity.Entity                     `bson:"entity" json:"entity"`
	Value  AlarmValue                        `bson:"v" json:"v"`
	Infos  map[string]map[string]interface{} `bson:"infos" json:"infos"`

	Pbehavior *Pbehavior `bson:"pbehavior,omitempty" json:"pbehavior,omitempty"`

	// Meta alarm fields
	MetaAlarmRule        *MetaAlarmRule `bson:"meta_alarm_rule,omitempty" json:"meta_alarm_rule,omitempty"`
	IsMetaAlarm          *bool          `bson:"is_meta_alarm,omitempty" json:"is_meta_alarm,omitempty"`
	Children             *int64         `bson:"children,omitempty" json:"children,omitempty"`
	ChildrenInstructions *bool          `bson:"children_instructions" json:"children_instructions,omitempty"`
	FilteredChildrenIDs  []string       `bson:"filtered_children,omitempty" json:"filtered_children,omitempty"`
	// Meta alarm child fields
	Parents        *int64          `bson:"parents" json:"parents,omitempty"`
	MetaAlarmRules []MetaAlarmRule `bson:"meta_alarm_rules" json:"meta_alarm_rules,omitempty"`

	AssignedInstructions             *[]AssignedInstruction `bson:"assigned_instructions,omitempty" json:"assigned_instructions,omitempty"`
	IsAutoInstructionRunning         *bool                  `bson:"-" json:"is_auto_instruction_running,omitempty"`
	IsAllAutoInstructionsCompleted   *bool                  `bson:"-" json:"is_all_auto_instructions_completed,omitempty"`
	IsAutoInstructionFailed          *bool                  `bson:"-" json:"is_auto_instruction_failed,omitempty"`
	IsManualInstructionRunning       *bool                  `bson:"-" json:"is_manual_instruction_running,omitempty"`
	IsManualInstructionWaitingResult *bool                  `bson:"-" json:"is_manual_instruction_waiting_result,omitempty"`

	Links       map[string]interface{} `bson:"-" json:"links,omitempty"`
	ImpactState int64                  `bson:"impact_state" json:"impact_state"`
}

type MetaAlarmRule struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type AlarmValue struct {
	ACK               *AlarmStep            `bson:"ack,omitempty" json:"ack,omitempty"`
	Canceled          *AlarmStep            `bson:"canceled,omitempty" json:"canceled,omitempty"`
	Done              *AlarmStep            `bson:"done,omitempty" json:"done,omitempty"`
	Snooze            *AlarmStep            `bson:"snooze,omitempty" json:"snooze,omitempty"`
	State             *AlarmStep            `bson:"state,omitempty" json:"state,omitempty"`
	Status            *AlarmStep            `bson:"status,omitempty" json:"status,omitempty"`
	Ticket            *AlarmTicket          `bson:"ticket,omitempty" json:"ticket,omitempty"`
	LastComment       *AlarmStep            `bson:"last_comment,omitempty" json:"last_comment,omitempty"`
	Steps             []AlarmStep           `bson:"steps,omitempty" json:"steps,omitempty"`
	Component         string                `bson:"component" json:"component"`
	Connector         string                `bson:"connector" json:"connector"`
	ConnectorName     string                `bson:"connector_name" json:"connector_name"`
	CreationDate      types.CpsTime         `bson:"creation_date" json:"creation_date" swaggertype:"integer"`
	ActivationDate    *types.CpsTime        `bson:"activation_date,omitempty" json:"activation_date,omitempty" swaggertype:"integer"`
	DisplayName       string                `bson:"display_name" json:"display_name"`
	InitialOutput     string                `bson:"initial_output" json:"initial_output"`
	Output            string                `bson:"output" json:"output"`
	InitialLongOutput string                `bson:"initial_long_output" json:"initial_long_output"`
	LongOutput        string                `bson:"long_output" json:"long_output"`
	LongOutputHistory []string              `bson:"long_output_history" json:"long_output_history"`
	LastUpdateDate    types.CpsTime         `bson:"last_update_date" json:"last_update_date" swaggertype:"integer"`
	LastEventDate     types.CpsTime         `bson:"last_event_date" json:"last_event_date" swaggertype:"integer"`
	Resource          string                `bson:"resource,omitempty" json:"resource,omitempty"`
	Resolved          *types.CpsTime        `bson:"resolved,omitempty" json:"resolved,omitempty" swaggertype:"integer"`
	PbehaviorInfo     *entity.PbehaviorInfo `bson:"pbehavior_info,omitempty" json:"pbehavior_info,omitempty"`
	Tags              []string              `bson:"tags" json:"tags"`
	Meta              string                `bson:"meta,omitempty" json:"meta,omitempty"`
	Parents           []string              `bson:"parents" json:"parents"`
	Children          []string              `bson:"children" json:"children"`

	StateChangesSinceStatusUpdate types.CpsNumber `bson:"state_changes_since_status_update,omitempty" json:"state_changes_since_status_update,omitempty"`
	TotalStateChanges             types.CpsNumber `bson:"total_state_changes,omitempty" json:"total_state_changes,omitempty"`

	Duration                  int   `bson:"duration" json:"duration"`
	CurrentStateDuration      int   `bson:"current_state_duration" json:"current_state_duration"`
	SnoozeDuration            int64 `bson:"snooze_duration" json:"snooze_duration"`
	PbehaviorInactiveDuration int64 `bson:"pbh_inactive_duration" json:"pbh_inactive_duration"`
	ActiveDuration            int64 `bson:"active_duration" json:"active_duration"`

	EventsCount types.CpsNumber `bson:"events_count,omitempty" json:"events_count,omitempty"`

	RuleVersion map[string]string                 `bson:"infos_rule_version" json:"infos_rule_version"`
	Infos       map[string]map[string]interface{} `bson:"infos" json:"infos"`
}

type AlarmStep struct {
	Type         string             `bson:"_t" json:"_t"`
	Timestamp    *types.CpsTime     `bson:"t" json:"t" swaggertype:"integer"`
	Author       string             `bson:"a" json:"a"`
	UserID       string             `bson:"user_id,omitempty" json:"user_id"`
	Message      string             `bson:"m" json:"m"`
	Value        types.CpsNumber    `bson:"val" json:"val"`
	Initiator    string             `bson:"initiator" json:"initiator"`
	Execution    string             `bson:"exec,omitempty" json:"-"`
	StateCounter *types.CropCounter `bson:"statecounter,omitempty" json:"statecounter,omitempty"`
}

type AlarmTicket struct {
	Type      string            `bson:"_t" json:"_t"`
	Timestamp types.CpsTime     `bson:"t" json:"t" swaggertype:"integer"`
	Author    string            `bson:"a" json:"a"`
	UserID    string            `bson:"user_id,omitempty" json:"user_id"`
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

	LastComment *pbehavior.Comment `bson:"last_comment" json:"last_comment"`
}

type Instruction struct {
	ID            string   `bson:"_id"`
	Name          string   `bson:"name"`
	ActiveOnPbh   []string `bson:"active_on_pbh"`
	DisabledOnPbh []string `bson:"disabled_on_pbh"`

	savedpattern.AlarmPatternFields  `bson:",inline"`
	savedpattern.EntityPatternFields `bson:",inline"`

	OldAlarmPatterns  oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty"`
	OldEntityPatterns oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty"`
}

type AssignedInstruction struct {
	ID        string     `json:"_id"`
	Name      string     `json:"name"`
	Execution *Execution `json:"execution"`
}

type InstructionWithExecutions struct {
	Instruction `bson:",inline"`

	Executions []struct {
		Execution `bson:",inline"`
		Alarm     string `bson:"alarm"`
	} `bson:"executions"`
}

func (i InstructionWithExecutions) GetExecution(alarmId string) *Execution {
	for _, v := range i.Executions {
		if v.Alarm == alarmId {
			return &v.Execution
		}
	}

	return nil
}

type ExecutionStatus struct {
	ID                  string `bson:"_id"`
	AutoRunning         *bool  `bson:"auto_running"`
	ManualRunning       *bool  `bson:"manual_running"`
	ManualWaitingResult *bool  `bson:"manual_waiting_result"`
	AutoFailed          *bool  `bson:"auto_failed"`
	AutoAllCompleted    *bool  `bson:"auto_all_completed"`
}

type Execution struct {
	ID     string `bson:"_id" json:"_id"`
	Status int    `bson:"status" json:"status"`
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
