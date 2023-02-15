package alarm

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorcomment"
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

const (
	NoIcon = iota
	IconManualInProgress
	IconAutoInProgress
	IconAutoFailed
	IconManualFailed
	IconManualFailedOtherInProgress
	IconAutoFailedOtherInProgress
	IconManualFailedManualAvailable
	IconAutoFailedManualAvailable
	IconManualAvailable
	IconAutoSuccessful
	IconManualSuccessful
	IconAutoSuccessfulOtherInProgress
	IconManualSuccessfulOtherInProgress
	IconAutoSuccessfulManualAvailable
	IconManualSuccessfulManualAvailable
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

func (r FilterRequest) GetOpenedFilter() int {
	if r.Opened == nil {
		return OpenedAndRecentResolved
	}

	if *r.Opened {
		return OnlyOpened
	}

	return OnlyResolved
}

type BaseFilterRequest struct {
	Filters     []string       `form:"filters[]" json:"filters"`
	Search      string         `form:"search" json:"search"`
	TimeField   string         `form:"time_field" json:"time_field" binding:"oneoforempty=t v.creation_date v.resolved v.last_update_date v.last_event_date"`
	StartFrom   *types.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo     *types.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
	Opened      *bool          `form:"opened" json:"opened"`
	OnlyParents bool           `form:"correlation" json:"correlation"`
	Category    string         `form:"category" json:"category"`
	Tag         string         `form:"tag" json:"tag"`

	AlarmPattern     string `form:"alarm_pattern" json:"alarm_pattern"`
	EntityPattern    string `form:"entity_pattern" json:"entity_pattern"`
	PbehaviorPattern string `form:"pbehavior_pattern" json:"pbehavior_pattern"`

	Instructions []InstructionFilterRequest `form:"instructions[]" json:"instructions"`
}

type InstructionFilterRequest struct {
	Running      *bool    `form:"running" json:"running"`
	IncludeTypes []int    `form:"include_types[]" json:"include_types"`
	ExcludeTypes []int    `form:"exclude_types[]" json:"exclude_types"`
	Include      []string `form:"include[]" json:"include"`
	Exclude      []string `form:"exclude[]" json:"exclude"`
}

type ListByServiceRequest struct {
	pagination.Query
	SortRequest
	Search      string `form:"search" json:"search"`
	Category    string `form:"category" json:"category"`
	WithService bool   `form:"with_service" json:"with_service"`
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
	Steps            *StepsRequest        `json:"steps"`
	Children         *ChildDetailsRequest `json:"children"`
}

type StepsRequest struct {
	pagination.Query
	Reversed bool `json:"reversed"`
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
	Data []common.AlarmStep   `json:"data"`
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
	Tags   []string                          `bson:"tags" json:"tags"`
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

	AssignedInstructions         *[]AssignedInstruction `bson:"assigned_instructions,omitempty" json:"assigned_instructions,omitempty"`
	InstructionExecutionIcon     int                    `bson:"-" json:"instruction_execution_icon,omitempty"`
	RunningManualInstructions    []string               `bson:"-" json:"running_manual_instructions,omitempty"`
	RunningAutoInstructions      []string               `bson:"-" json:"running_auto_instructions,omitempty"`
	FailedManualInstructions     []string               `bson:"-" json:"failed_manual_instructions,omitempty"`
	FailedAutoInstructions       []string               `bson:"-" json:"failed_auto_instructions,omitempty"`
	SuccessfulManualInstructions []string               `bson:"-" json:"successful_manual_instructions,omitempty"`
	SuccessfulAutoInstructions   []string               `bson:"-" json:"successful_auto_instructions,omitempty"`

	Links       map[string]interface{} `bson:"-" json:"links,omitempty"`
	ImpactState int64                  `bson:"impact_state" json:"impact_state"`
}

type MetaAlarmRule struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type AlarmValue struct {
	ACK         *common.AlarmStep  `bson:"ack,omitempty" json:"ack,omitempty"`
	Canceled    *common.AlarmStep  `bson:"canceled,omitempty" json:"canceled,omitempty"`
	Done        *common.AlarmStep  `bson:"done,omitempty" json:"done,omitempty"`
	Snooze      *common.AlarmStep  `bson:"snooze,omitempty" json:"snooze,omitempty"`
	State       *common.AlarmStep  `bson:"state,omitempty" json:"state,omitempty"`
	Status      *common.AlarmStep  `bson:"status,omitempty" json:"status,omitempty"`
	Ticket      *AlarmTicket       `bson:"ticket,omitempty" json:"ticket,omitempty"`
	LastComment *common.AlarmStep  `bson:"last_comment,omitempty" json:"last_comment,omitempty"`
	Steps       []common.AlarmStep `bson:"steps,omitempty" json:"steps,omitempty"`

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
	ID     string            `bson:"_id" json:"_id"`
	Author *author.Author    `bson:"author" json:"author"`
	Name   string            `bson:"name" json:"name"`
	RRule  string            `bson:"rrule" json:"rrule"`
	Start  *types.CpsTime    `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop   *types.CpsTime    `bson:"tstop" json:"tstop" swaggertype:"integer"`
	Type   *pbehavior.Type   `bson:"type" json:"type"`
	Reason *pbehavior.Reason `bson:"reason" json:"reason"`

	LastComment *pbehaviorcomment.Response `bson:"last_comment" json:"last_comment"`
}

type Instruction struct {
	ID            string   `bson:"_id"`
	Name          string   `bson:"name"`
	Type          int      `bson:"type"`
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
	Type      int        `json:"type"`
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
	ID                           string   `bson:"_id"`
	Icon                         int      `bson:"-"`
	RunningManualInstructions    []string `bson:"running_manual_instructions"`
	RunningAutoInstructions      []string `bson:"running_auto_instructions"`
	FailedManualInstructions     []string `bson:"failed_manual_instructions"`
	FailedAutoInstructions       []string `bson:"failed_auto_instructions"`
	SuccessfulManualInstructions []string `bson:"successful_manual_instructions"`
	SuccessfulAutoInstructions   []string `bson:"successful_auto_instructions"`
	LastFailed                   *int     `bson:"last_failed"`
	LastSuccessful               *int     `bson:"last_successful"`
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

type GetOpenRequest struct {
	ID string `form:"_id" json:"_id" binding:"required"`
}
