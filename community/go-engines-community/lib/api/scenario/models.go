package scenario

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/author"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/savedpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/x/bsonx/bsoncore"
)

type FilteredQuery struct {
	pagination.FilteredQuery
	SortBy string `json:"sort_by" form:"sort_by" binding:"oneoforempty=_id name author.name author.display_name created updated enabled priority"`
}

type EditRequest struct {
	Name                 string                     `json:"name" binding:"required,max=255"`
	Author               string                     `json:"author" swaggerignore:"true"`
	Enabled              *bool                      `json:"enabled" binding:"required"`
	Priority             int64                      `json:"priority" binding:"min=0"`
	Triggers             []Trigger                  `json:"triggers" binding:"required,notblank,dive"`
	DisableDuringPeriods []string                   `json:"disable_during_periods" binding:"dive,oneof=maintenance pause inactive"`
	Delay                *datetime.DurationWithUnit `json:"delay"`
	Actions              []ActionRequest            `json:"actions" binding:"required,notblank,dive"`
}

type Trigger struct {
	// Possible trigger type values.
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
	//   * `pbhleave` - Alarm leaves a periodic behavior
	//   * `activate` - Alarm has been activated
	//   * `resolve` - Alarm has been resolved
	//   * `instructionfail` - Manual instruction has failed
	//   * `autoinstructionfail` - Auto instruction has failed
	//   * `instructionjobfail` - Manual or auto instruction's job is failed
	//   * `instructionjobcomplete` - Manual or auto instruction's job is completed
	//   * `instructioncomplete` - Manual instruction is completed
	//   * `autoinstructioncomplete` - Auto instruction is completed
	//   * `autoinstructionresultok` - Alarm is in OK state after all auto instructions
	//   * `autoinstructionresultfail` - Alarm is in not in OK state after all auto instructions
	//   * `eventscount` - Alarm check events count
	Type      string `json:"type" binding:"required,oneof=create statedec stateinc changestate changestatus ack ackremove cancel uncancel comment declareticketwebhook assocticket snooze unsnooze resolve activate pbhenter pbhleave instructionfail autoinstructionfail instructionjobfail instructionjobcomplete instructioncomplete autoinstructioncomplete autoinstructionresultok autoinstructionresultfail eventscount"`
	Threshold int    `json:"threshold,omitempty" binding:"required_if=Type eventscount,excluded_unless=Type eventscount,omitempty,gt=1"`
}

func (t *Trigger) UnmarshalBSONValue(valueType byte, b []byte) error {
	if bson.Type(valueType) != bson.TypeString {
		return errors.New("trigger should be a string")
	}
	value, _, ok := bsoncore.ReadString(b)
	if !ok {
		return errors.New("invalid trigger value, expected string")
	}

	thresholdStr, ok := strings.CutPrefix(value, string(types.AlarmChangeEventsCount))
	if !ok {
		t.Type = value
		return nil
	}

	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		return fmt.Errorf("cannot decode an %s trigger threshold value: %w", types.AlarmChangeEventsCount, err)
	}

	t.Type = string(types.AlarmChangeEventsCount)
	t.Threshold = threshold

	return nil
}

func (t *Trigger) String() string {
	if t == nil {
		return ""
	}

	if t.Type == string(types.AlarmChangeEventsCount) {
		return t.Type + strconv.Itoa(t.Threshold)
	}

	return t.Type
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

type ActionRequest struct {
	Type                     string            `json:"type" binding:"required,oneof=ack ackremove assocticket cancel changestate pbehavior pbehaviorremove snooze unsnooze webhook"`
	Parameters               action.Parameters `json:"parameters,omitempty"`
	Comment                  string            `json:"comment"`
	DropScenarioIfNotMatched *bool             `json:"drop_scenario_if_not_matched" binding:"required"`
	EmitTrigger              *bool             `json:"emit_trigger" binding:"required"`

	common.EntityPatternFieldsRequest `bson:",inline"`
	common.AlarmPatternFieldsRequest  `bson:",inline"`
}

type Scenario struct {
	ID                   string                     `bson:"_id" json:"_id"`
	Name                 string                     `bson:"name" json:"name"`
	Author               *author.Author             `bson:"author" json:"author"`
	Enabled              bool                       `bson:"enabled" json:"enabled"`
	DisableDuringPeriods []string                   `bson:"disable_during_periods" json:"disable_during_periods"`
	Triggers             []Trigger                  `bson:"triggers" json:"triggers"`
	Actions              []Action                   `bson:"actions" json:"actions"`
	Priority             int64                      `bson:"priority" json:"priority"`
	Delay                *datetime.DurationWithUnit `bson:"delay" json:"delay"`
	Created              datetime.CpsTime           `bson:"created,omitempty" json:"created,omitempty" swaggertype:"integer"`
	Updated              datetime.CpsTime           `bson:"updated,omitempty" json:"updated,omitempty" swaggertype:"integer"`
}

type Action struct {
	Type                     string     `bson:"type" json:"type"`
	Comment                  string     `bson:"comment" json:"comment"`
	Parameters               Parameters `bson:"parameters,omitempty" json:"parameters,omitempty"`
	DropScenarioIfNotMatched bool       `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool       `bson:"emit_trigger" json:"emit_trigger"`

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
	Ticket         string            `json:"ticket,omitempty" bson:"ticket"`
	TicketURL      string            `json:"ticket_url,omitempty" bson:"ticket_url"`
	TicketURLTitle string            `json:"ticket_url_title,omitempty" bson:"ticket_url_title"`
	TicketData     map[string]string `json:"ticket_data,omitempty" bson:"ticket_data"`
	// AssocTicket and Webhook
	TicketSystemName string `json:"ticket_system_name,omitempty" bson:"ticket_system_name"`
	// Snooze and Pbehavior
	Duration *datetime.DurationWithUnit `json:"duration,omitempty" bson:"duration"`
	// Pbehavior
	Name           string            `json:"name,omitempty" bson:"name"`
	Reason         *pbehavior.Reason `json:"reason,omitempty" bson:"reason"`
	Type           *pbehavior.Type   `json:"type,omitempty" bson:"type"`
	RRule          string            `json:"rrule,omitempty" bson:"rrule"`
	Tstart         *int64            `json:"tstart,omitempty" bson:"tstart"`
	Tstop          *int64            `json:"tstop,omitempty" bson:"tstop"`
	StartOnTrigger *bool             `json:"start_on_trigger,omitempty" bson:"start_on_trigger"`
	// Webhook
	Request            *request.Parameters           `json:"request,omitempty" bson:"request"`
	SkipForChild       *bool                         `json:"skip_for_child,omitempty" bson:"skip_for_child"`
	SkipForInstruction *bool                         `json:"skip_for_instruction,omitempty" bson:"skip_for_instruction,omitempty"`
	DeclareTicket      *request.WebhookDeclareTicket `json:"declare_ticket,omitempty" bson:"declare_ticket"`
	StopOnFail         *bool                         `json:"stop_on_fail,omitempty" bson:"stop_on_fail,omitempty"`
	StopOnSuccess      *bool                         `json:"stop_on_success,omitempty" bson:"stop_on_success,omitempty"`
	MultipleURLs       *bool                         `json:"multiple_urls,omitempty" bson:"multiple_urls,omitempty"`
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

type TemplateRequest struct {
	Rule struct {
		TemplateRuleRequest
		ID string `json:"_id" binding:"id"`
	} `json:"rule"`
	TestData struct {
		Test  string `json:"test"`
		Event string `json:"event"`
		// TestData.Responses keys correspond with Rule.Actions keys
		Responses map[int]string `json:"responses"`
	} `json:"testdata"`
}

type TemplateRuleRequest struct {
	Name     string                  `json:"name" binding:"required"`
	Triggers []Trigger               `json:"triggers" binding:"required,notblank,dive"`
	Actions  []TemplateActionRequest `json:"actions" binding:"required,notblank,dive"`
}

type TemplateActionRequest struct {
	Type       string                   `json:"type" binding:"required"`
	Parameters TemplateActionParameters `json:"parameters"`
}

type TemplateActionParameters struct {
	Output        string                                 `json:"output"`
	Author        string                                 `json:"author"`
	ForwardAuthor *bool                                  `json:"forward_author"`
	Request       *template.TemplateParameters           `json:"request"`
	DeclareTicket *template.TemplateWebhookDeclareTicket `json:"declare_ticket"`
}

type TemplateVarsResponse struct {
	Output       []template.VarResponse `json:"output"`
	Author       []template.VarResponse `json:"author"`
	FirstWebhook []template.VarResponse `json:"first_webhook"`
	Webhook      []template.VarResponse `json:"webhook"`
	Ticket       []template.VarResponse `json:"ticket"`
}
