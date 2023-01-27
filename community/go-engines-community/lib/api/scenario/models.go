package scenario

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author.name created updated enabled priority"`
}

type EditRequest struct {
	Name     string `json:"name" binding:"required,max=255"`
	Author   string `json:"author" binding:"required,max=255"`
	Enabled  *bool  `json:"enabled" binding:"required"`
	Priority *int   `json:"priority" binding:"gt=0"`

	// Possible trigger values.
	//   * `create` - Alarm creation
	//   * `statedec` - Alarm state decrease
	//   * `changestate` - Alarm state has been changed by "change state" action
	//   * `stateinc` - Alarm state increase
	//   * `changestatus` - Alarm status changes eg. flapping
	//   * `ack` - Alarm has been acked
	//   * `ackremove` - Alarm has been unacked
	//   * `cancel` - Alarm has been cancelled
	//   * `uncancel` - Alarm has been uncancelled
	//   * `comment` - Alarm has been commented
	//   * `declareticketwebhook` - Ticket has been declared by the webhook
	//   * `assocticket` - Ticket has been associated with an alarm
	//   * `snooze` - Alarm has been snoozed
	//   * `unsnooze` - Alarm has been unsnoozed
	//   * `pbhenter` - Alarm enters a periodic behavior
	//   * `activate` - Alarm has been activated
	//   * `resolve` - Alarm has been resolved
	//   * `pbhleave` - Alarm leaves a periodic behavior
	//   * `instructionfail` - Manual instruction has failed
	//   * `autoinstructionfail` - Auto instruction has failed
	//   * `instructionjobfail` - Manual or auto instruction's job is failed
	//   * `instructionjobcomplete` - Manual or auto instruction's job is completed
	//   * `instructioncomplete` - Manual instruction is completed
	//   * `autoinstructioncomplete` - Auto instruction is completed
	Triggers             []string                `json:"triggers" binding:"required,notblank,dive,oneof=create statedec stateinc changestate changestatus ack ackremove cancel uncancel comment declareticketwebhook assocticket snooze unsnooze resolve activate pbhenter pbhleave instructionfail autoinstructionfail instructionjobfail instructionjobcomplete instructioncomplete autoinstructioncomplete"`
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
	OldAlarmPatterns         oldpattern.AlarmPatternList  `json:"old_alarm_patterns"`
	OldEntityPatterns        oldpattern.EntityPatternList `json:"old_entity_patterns"`
	DropScenarioIfNotMatched *bool                        `json:"drop_scenario_if_not_matched" binding:"required"`
	EmitTrigger              *bool                        `json:"emit_trigger" binding:"required"`

	common.EntityPatternFieldsRequest `bson:",inline"`
	common.AlarmPatternFieldsRequest  `bson:",inline"`
}

type Scenario struct {
	ID                   string                  `bson:"_id" json:"_id"`
	Name                 string                  `bson:"name" json:"name"`
	Author               *author.Author          `bson:"author" json:"author"`
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
	OldAlarmPatterns         oldpattern.AlarmPatternList  `bson:"old_alarm_patterns,omitempty" json:"old_alarm_patterns,omitempty"`
	OldEntityPatterns        oldpattern.EntityPatternList `bson:"old_entity_patterns,omitempty" json:"old_entity_patterns,omitempty"`
	DropScenarioIfNotMatched bool                         `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool                         `bson:"emit_trigger" json:"emit_trigger"`

	savedpattern.EntityPatternFields `bson:",inline"`
	savedpattern.AlarmPatternFields  `bson:",inline"`
}

type Parameters struct {
	Output string `json:"output,omitempty" bson:"output"`

	ForwardAuthor *bool  `json:"forward_author,omitempty" bson:"forward_author"`
	Author        string `json:"author,omitempty" bson:"author"`

	// ChangeState
	State *types.CpsNumber `json:"state,omitempty" bson:"state"`
	// AssocTicket
	Ticket     string            `json:"ticket,omitempty" bson:"ticket"`
	TicketURL  string            `json:"ticket_url,omitempty" bson:"ticket_url"`
	TicketData map[string]string `json:"ticket_data,omitempty" bson:"ticket_data"`
	// AssocTicket and Webhook
	TicketSystemName string `json:"ticket_system_name,omitempty" bson:"ticket_system_name"`
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
	Request       *request.Parameters           `json:"request,omitempty" bson:"request"`
	SkipForChild  *bool                         `json:"skip_for_child,omitempty" bson:"skip_for_child"`
	DeclareTicket *request.WebhookDeclareTicket `json:"declare_ticket,omitempty" bson:"declare_ticket"`
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
