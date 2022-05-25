package scenario

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author expected_interval created updated enabled priority delay"`
}

type EditRequest struct {
	Name                 string                  `json:"name" binding:"required,max=255"`
	Author               string                  `json:"author" binding:"required,max=255"`
	Enabled              *bool                   `json:"enabled" binding:"required"`
	Priority             *int                    `json:"priority" binding:"gt=0"`
	Triggers             []string                `json:"triggers" binding:"required,notblank,dive,oneof=create statedec stateinc changestate changestatus ack ackremove cancel uncancel comment done declareticket declareticketwebhook assocticket snooze unsnooze resolve activate pbhenter pbhleave"`
	DisableDuringPeriods []string                `json:"disable_during_periods" binding:"dive,oneof=maintenance pause inactive"`
	Delay                *types.DurationWithUnit `json:"delay"`
	Actions              []ActionRequest         `json:"actions" binding:"required,notblank,dive"`
}

type CreateRequest struct {
	EditRequest
	ID string `json:"_id" binding:"id"`
}

type UpdateRequest struct {
	EditRequest
	ID string `json:"-"`
}

type BulkUpdateRequestItem struct {
	EditRequest
	ID string `json:"_id" binding:"required"`
}

type BulkDeleteRequestItem struct {
	ID string `json:"_id" binding:"required"`
}

// for swagger
type BulkCreateResponseItem struct {
	ID     string            `json:"id,omitempty"`
	Item   CreateRequest     `json:"item"`
	Status int               `json:"status"`
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// for swagger
type BulkUpdateResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkUpdateRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
}

// for swagger
type BulkDeleteResponseItem struct {
	ID     string                `json:"id,omitempty"`
	Item   BulkDeleteRequestItem `json:"item"`
	Status int                   `json:"status"`
	Error  string                `json:"error,omitempty"`
	Errors map[string]string     `json:"errors,omitempty"`
}

type GetMinimalPriorityResponse struct {
	Priority int `json:"priority"`
}

type CheckPriorityRequest struct {
	Priority int `json:"priority" binding:"required,gt=0"`
}

type CheckPriorityResponse struct {
	Valid               bool `json:"valid"`
	RecommendedPriority int  `json:"recommended_priority,omitempty"`
}

type ActionRequest struct {
	Type                     string                       `json:"type" binding:"required,oneof=ack ackremove assocticket cancel changestate pbehavior snooze webhook"`
	Parameters               action.Parameters            `json:"parameters,omitempty"`
	Comment                  string                       `json:"comment"`
	AlarmPatterns            oldpattern.AlarmPatternList  `json:"alarm_patterns"`
	EntityPatterns           oldpattern.EntityPatternList `json:"entity_patterns"`
	DropScenarioIfNotMatched *bool                        `json:"drop_scenario_if_not_matched" binding:"required"`
	EmitTrigger              *bool                        `json:"emit_trigger" binding:"required"`
}

type Scenario struct {
	ID                   string                  `bson:"_id" json:"_id"`
	Name                 string                  `bson:"name" json:"name"`
	Author               string                  `bson:"author" json:"author"`
	Enabled              bool                    `bson:"enabled" json:"enabled"`
	DisableDuringPeriods []string                `bson:"disable_during_periods" json:"disable_during_periods"`
	Triggers             []string                `bson:"triggers" json:"triggers"`
	Actions              []Action                `bson:"actions" json:"actions"`
	Priority             int                     `bson:"priority" json:"priority"`
	Delay                *types.DurationWithUnit `bson:"delay" json:"delay"`
	Created              types.CpsTime           `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated              types.CpsTime           `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Action struct {
	Type                     string                       `bson:"type" json:"type"`
	Comment                  string                       `bson:"comment" json:"comment"`
	Parameters               Parameters                   `bson:"parameters,omitempty" json:"parameters,omitempty"`
	AlarmPatterns            oldpattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EntityPatterns           oldpattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	DropScenarioIfNotMatched bool                         `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool                         `bson:"emit_trigger" json:"emit_trigger"`
}

type Parameters struct {
	Output string `json:"output,omitempty" bson:"output"`

	ForwardAuthor *bool  `json:"forward_author,omitempty" bson:"forward_author"`
	Author        string `json:"author,omitempty" bson:"author"`

	// ChangeState
	State *types.CpsNumber `json:"state,omitempty" bson:"state"`
	// AssocTicket
	Ticket string `json:"ticket,omitempty" bson:"ticket"`
	// Snooze and Pbehavior
	Duration *types.DurationWithUnit `json:"duration,omitempty" bson:"duration"`
	// Pbehavior
	Name           string            `json:"name,omitempty" bson:"name"`
	Reason         *pbehavior.Reason `json:"reason,omitempty" bson:"reason"`
	Type           *pbehavior.Type   `json:"type,omitempty" bson:"type"`
	RRule          string            `json:"rrule,omitempty" bson:"rrule"`
	Tstart         *int64            `json:"tstart,omitempty" bson:"tstart"`
	Tstop          *int64            `json:"tstop,omitempty" bson:"tstop"`
	StartOnTrigger *bool             `json:"start_on_trigger,omitempty" bson:"start_on_trigger"`
	// Webhook
	Request       *types.WebhookRequest       `json:"request,omitempty" bson:"request"`
	DeclareTicket *types.WebhookDeclareTicket `json:"declare_ticket,omitempty" bson:"declare_ticket"`
	RetryCount    int64                       `json:"retry_count,omitempty" bson:"retry_count"`
	RetryDelay    *types.DurationWithUnit     `json:"retry_delay,omitempty" bson:"retry_delay"`
}

type AggregationResult struct {
	Data       []Scenario `bson:"data" json:"data"`
	TotalCount int64      `bson:"total_count" json:"total_count"`
}

// GetTotal implementation PaginatedData interface
func (r AggregationResult) GetTotal() int64 {
	return r.TotalCount
}

// GetData implementation PaginatedData interface
func (r AggregationResult) GetData() interface{} {
	return r.Data
}
