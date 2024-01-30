package alarm

//go:generate easyjson -no_std_marshalers

import (
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/export"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pbehaviorcomment"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/statesettings"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/link"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/db"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
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
	WithInstructions   bool `form:"with_instructions" json:"with_instructions"`
	WithDeclareTickets bool `form:"with_declare_tickets" json:"with_declare_tickets"`
	WithLinks          bool `form:"with_links" json:"with_links"`
	WithDependencies   bool `form:"with_dependencies" json:"with_dependencies"`
}

type FilterRequest struct {
	BaseFilterRequest
	SearchBy []string `form:"active_columns[]" json:"active_columns[]"`
}

func (r BaseFilterRequest) GetOpenedFilter() int {
	if r.Opened == nil {
		return OpenedAndRecentResolved
	}

	if *r.Opened {
		return OnlyOpened
	}

	return OnlyResolved
}

type BaseFilterRequest struct {
	Filters     []string          `form:"filters[]" json:"filters"`
	Search      string            `form:"search" json:"search"`
	TimeField   string            `form:"time_field" json:"time_field" binding:"oneoforempty=t v.creation_date v.resolved v.last_update_date v.last_event_date"`
	StartFrom   *datetime.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo     *datetime.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
	Opened      *bool             `form:"opened" json:"opened"`
	OnlyParents bool              `form:"correlation" json:"correlation"`
	Category    string            `form:"category" json:"category"`
	Tag         string            `form:"tag" json:"tag"`

	AlarmPattern     string `form:"alarm_pattern" json:"alarm_pattern"`
	EntityPattern    string `form:"entity_pattern" json:"entity_pattern"`
	PbehaviorPattern string `form:"pbehavior_pattern" json:"pbehavior_pattern"`

	Instructions  []InstructionFilterRequest `form:"instructions[]" json:"instructions"`
	OnlyBookmarks bool                       `form:"only_bookmarks" json:"only_bookmarks"`
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
	ID        string            `form:"_id" json:"_id" binding:"required"`
	StartFrom *datetime.CpsTime `form:"tstart" json:"tstart" swaggertype:"integer"`
	StartTo   *datetime.CpsTime `form:"tstop" json:"tstop" swaggertype:"integer"`
}

type SortRequest struct {
	MultiSort []string `form:"multi_sort[]" json:"multi_sort"`
	Sort      string   `form:"sort" json:"sort" binding:"oneoforempty=asc desc"`
	SortBy    string   `form:"sort_by" json:"sort_by"`
}

type DetailsRequest struct {
	ID                 string               `json:"_id" binding:"required"`
	Search             string               `json:"search"`
	SearchBy           []string             `json:"search_by"`
	Opened             *bool                `json:"opened"`
	WithInstructions   bool                 `json:"with_instructions"`
	WithDeclareTickets bool                 `json:"with_declare_tickets"`
	WithDependencies   bool                 `json:"with_dependencies"`
	Steps              *StepsRequest        `json:"steps"`
	Children           *ChildDetailsRequest `json:"children"`
	PerfData           []string             `json:"perf_data"`
}

func (r *DetailsRequest) Format() {
	defaultQuery := pagination.GetDefaultQuery()

	if r.Steps != nil {
		r.Steps.Paginate = true
		if r.Steps.Page == 0 {
			r.Steps.Page = defaultQuery.Page
		}
		if r.Steps.Limit == 0 {
			r.Steps.Limit = defaultQuery.Limit
		}
	}

	if r.Children != nil {
		r.Children.Paginate = true
		if r.Children.Page == 0 {
			r.Children.Page = defaultQuery.Page
		}
		if r.Children.Limit == 0 {
			r.Children.Limit = defaultQuery.Limit
		}
	}
}

func (r *DetailsRequest) GetOpenedFilter() int {
	if r.Opened == nil {
		return OpenedAndRecentResolved
	}

	if *r.Opened {
		return OnlyOpened
	}

	return OnlyResolved
}

type StepsRequest struct {
	pagination.Query
	Reversed bool   `json:"reversed"`
	Type     string `json:"type"`
}

type ChildDetailsRequest struct {
	pagination.Query
	SortRequest
}

type DetailsResponse struct {
	ID     string            `json:"_id"`
	Status int               `json:"status"`
	Data   Details           `json:"data,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
	Error  string            `json:"error,omitempty"`
}

type EntityDetails struct {
	types.Entity `bson:",inline" json:",inline"`
	// DependsCount contains only service's dependencies
	DependsCount int `bson:"depends_count" json:"depends_count"`
	// ImpactsCount contains only services
	ImpactsCount int                        `bson:"impacts_count" json:"impacts_count"`
	StateSetting *CheckStateSettingResponse `bson:"state_setting,omitempty" json:"state_setting,omitempty"`
}

type Details struct {
	// Only for websocket
	ID string `bson:"-" json:"_id,omitempty"`

	Steps    *StepDetails     `bson:"steps" json:"steps,omitempty"`
	Children *ChildrenDetails `bson:"children" json:"children,omitempty"`

	FilteredPerfData []string `bson:"filtered_perf_data" json:"filtered_perf_data,omitempty"`

	IsMetaAlarm bool  `json:"-" bson:"is_meta_alarm"`
	StepsCount  int64 `json:"-" bson:"steps_count"`
	// Entity isn't the same as Entity of Alarm, but have counts in response as well
	Entity EntityDetails `json:"entity" bson:"entity"`
}

type StepDetails struct {
	Data []common.AlarmStep   `json:"data"`
	Meta common.PaginatedMeta `json:"meta"`
}

type ChildrenDetails struct {
	Data []Alarm              `json:"data"`
	Meta common.PaginatedMeta `json:"meta"`
}

type CheckStateSettingResponse struct {
	ID                     string                         `bson:"_id" json:"_id"`
	Title                  string                         `json:"title"`
	Method                 string                         `json:"method"`
	InheritedEntityPattern *pattern.Entity                `bson:"inherited_entity_pattern,omitempty" json:"inherited_entity_pattern,omitempty"`
	StateThresholds        *statesettings.StateThresholds `bson:"state_thresholds,omitempty" json:"state_thresholds,omitempty"`
}

type ExportRequest struct {
	ExportFetchParameters
	Fields    export.Fields `json:"fields"`
	Separator string        `json:"separator" binding:"oneoforempty=comma semicolon tab space"`
}

// ExportFetchParameters
// easyjson:json
type ExportFetchParameters struct {
	BaseFilterRequest
	TimeFormat string `json:"time_format" binding:"time_format"`
}

type ExportResponse struct {
	ID string `json:"_id"`
	// Possible values.
	//   * `0` - Running
	//   * `1` - Succeeded
	//   * `2` - Failed
	Status int64 `json:"status"`
}

type Alarm struct {
	ID     string                            `bson:"_id" json:"_id"`
	Time   datetime.CpsTime                  `bson:"t" json:"t" swaggertype:"integer"`
	Entity entity.Entity                     `bson:"entity" json:"entity"`
	Value  AlarmValue                        `bson:"v" json:"v"`
	Tags   []string                          `bson:"tags" json:"tags"`
	Infos  map[string]map[string]interface{} `bson:"infos" json:"infos"`

	Pbehavior *Pbehavior `bson:"pbehavior,omitempty" json:"pbehavior,omitempty"`

	// Meta alarm fields
	MetaAlarmRule        *MetaAlarmRule `bson:"meta_alarm_rule,omitempty" json:"meta_alarm_rule,omitempty"`
	IsMetaAlarm          *bool          `bson:"is_meta_alarm,omitempty" json:"is_meta_alarm,omitempty"`
	Children             *int64         `bson:"children,omitempty" json:"children,omitempty"`
	OpenedChildren       *int64         `bson:"opened_children,omitempty" json:"opened_children,omitempty"`
	ClosedChildren       *int64         `bson:"closed_children,omitempty" json:"closed_children,omitempty"`
	ChildrenInstructions *bool          `bson:"children_instructions" json:"children_instructions,omitempty"`
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

	Links       link.LinksByCategory `bson:"-" json:"links,omitempty"`
	ImpactState int64                `bson:"impact_state" json:"impact_state"`

	AssignedDeclareTicketRules []AssignedDeclareTicketRule `bson:"-" json:"assigned_declare_ticket_rules,omitempty"`

	// Only for details request
	Filtered *bool `bson:"filtered" json:"filtered,omitempty"`

	Bookmark bool `bson:"bookmark" json:"bookmark"`
}

type MetaAlarmRule struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
	Type string `bson:"type" json:"type"`
}

type AlarmValue struct {
	ACK         *common.AlarmStep  `bson:"ack,omitempty" json:"ack,omitempty"`
	Canceled    *common.AlarmStep  `bson:"canceled,omitempty" json:"canceled,omitempty"`
	Snooze      *common.AlarmStep  `bson:"snooze,omitempty" json:"snooze,omitempty"`
	State       *common.AlarmStep  `bson:"state,omitempty" json:"state,omitempty"`
	Status      *common.AlarmStep  `bson:"status,omitempty" json:"status,omitempty"`
	Tickets     []common.AlarmStep `bson:"tickets,omitempty" json:"tickets,omitempty"`
	Ticket      *common.AlarmStep  `bson:"ticket,omitempty" json:"ticket,omitempty"`
	LastComment *common.AlarmStep  `bson:"last_comment,omitempty" json:"last_comment,omitempty"`
	Steps       []common.AlarmStep `bson:"steps,omitempty" json:"steps,omitempty"`

	Component         string                `bson:"component" json:"component"`
	Connector         string                `bson:"connector" json:"connector"`
	ConnectorName     string                `bson:"connector_name" json:"connector_name"`
	CreationDate      datetime.CpsTime      `bson:"creation_date" json:"creation_date" swaggertype:"integer"`
	ActivationDate    *datetime.CpsTime     `bson:"activation_date,omitempty" json:"activation_date,omitempty" swaggertype:"integer"`
	DisplayName       string                `bson:"display_name" json:"display_name"`
	InitialOutput     string                `bson:"initial_output" json:"initial_output"`
	Output            string                `bson:"output" json:"output"`
	InitialLongOutput string                `bson:"initial_long_output" json:"initial_long_output"`
	LongOutput        string                `bson:"long_output" json:"long_output"`
	LongOutputHistory []string              `bson:"long_output_history" json:"long_output_history"`
	LastUpdateDate    datetime.CpsTime      `bson:"last_update_date" json:"last_update_date" swaggertype:"integer"`
	LastEventDate     datetime.CpsTime      `bson:"last_event_date" json:"last_event_date" swaggertype:"integer"`
	Resource          string                `bson:"resource,omitempty" json:"resource,omitempty"`
	Resolved          *datetime.CpsTime     `bson:"resolved,omitempty" json:"resolved,omitempty" swaggertype:"integer"`
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

type Pbehavior struct {
	ID     string            `bson:"_id" json:"_id"`
	Author *author.Author    `bson:"author" json:"author"`
	Name   string            `bson:"name" json:"name"`
	RRule  string            `bson:"rrule" json:"rrule"`
	Start  *datetime.CpsTime `bson:"tstart" json:"tstart" swaggertype:"integer"`
	Stop   *datetime.CpsTime `bson:"tstop" json:"tstop" swaggertype:"integer"`
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

type AssignedDeclareTicketRule struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`
}

type DeclareTicketRule struct {
	ID   string `bson:"_id" json:"_id"`
	Name string `bson:"name" json:"name"`

	savedpattern.AlarmPatternFields     `bson:",inline"`
	savedpattern.EntityPatternFields    `bson:",inline"`
	savedpattern.PbehaviorPatternFields `bson:",inline"`
}

func (r DeclareTicketRule) getDeclareTicketQuery() (bson.M, error) {
	alarmPatternQuery, err := db.AlarmPatternToMongoQuery(r.AlarmPattern, "")
	if err != nil {
		return nil, fmt.Errorf("invalid alarm pattern in declare ticket rule id=%q: %w", r.ID, err)
	}

	entityPatternQuery, err := db.EntityPatternToMongoQuery(r.EntityPattern, "entity")
	if err != nil {
		return nil, fmt.Errorf("invalid entity pattern in declare ticket rule id=%q: %w", r.ID, err)
	}

	pbhPatternQuery, err := db.PbehaviorInfoPatternToMongoQuery(r.PbehaviorPattern, "v")
	if err != nil {
		return nil, fmt.Errorf("invalid pbehavior pattern in declare ticket rule id=%q: %w", r.ID, err)
	}

	if len(alarmPatternQuery) == 0 && len(entityPatternQuery) == 0 && len(pbhPatternQuery) == 0 {
		return nil, nil
	}

	var and []bson.M
	if len(alarmPatternQuery) > 0 {
		and = append(and, alarmPatternQuery)
	}

	if len(entityPatternQuery) > 0 {
		and = append(and, entityPatternQuery)
	}

	if len(pbhPatternQuery) > 0 {
		and = append(and, pbhPatternQuery)
	}

	return bson.M{"$and": and}, nil
}

type LinksRequest struct {
	Ids []string `form:"ids[]" json:"ids" binding:"required,notblank"`
}
