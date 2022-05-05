package scenario

import (
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/mitchellh/mapstructure"
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
	Type                     string                    `json:"type" binding:"required,oneof=ack ackremove assocticket cancel changestate pbehavior snooze webhook"`
	Parameters               interface{}               `json:"parameters,omitempty"`
	Comment                  string                    `json:"comment"`
	AlarmPatterns            pattern.AlarmPatternList  `json:"alarm_patterns"`
	EntityPatterns           pattern.EntityPatternList `json:"entity_patterns"`
	DropScenarioIfNotMatched *bool                     `json:"drop_scenario_if_not_matched" binding:"required"`
	EmitTrigger              *bool                     `json:"emit_trigger" binding:"required"`
}

func (r *ActionRequest) UnmarshalJSON(b []byte) error {
	type Alias ActionRequest
	tmp := struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch r.Type {
	case types.ActionTypeSnooze:
		var params SnoozeParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypeChangeState:
		var params ChangeStateParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypeAssocTicket:
		var params AssocTicketParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypePbehavior:
		var params PbehaviorParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	case types.ActionTypeWebhook:
		var params WebhookParameterRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	default:
		var params ParametersRequest
		err := mapstructure.Decode(r.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		r.Parameters = params
	}

	return nil
}

type SnoozeParametersRequest struct {
	Duration *types.DurationWithUnit `json:"duration" binding:"required"`
	Output   string                  `json:"output" binding:"max=255"`
	Author   string                  `json:"author"`
	User     string                  `json:"user" swaggerignore:"true"`
}

type PbehaviorParametersRequest struct {
	Name           string                  `json:"name" binding:"required,max=255"`
	Author         string                  `json:"author" swaggerignore:"true"`
	User           string                  `json:"user" swaggerignore:"true"`
	Reason         string                  `json:"reason" binding:"required"`
	Type           string                  `json:"type" binding:"required"`
	RRule          string                  `json:"rrule"`
	Tstart         *int                    `json:"tstart" binding:"required_with=Tstop"`
	Tstop          *int                    `json:"tstop" binding:"required_with=Tstart"`
	StartOnTrigger *bool                   `json:"start_on_trigger" binding:"required_with=Duration" mapstructure:"start_on_trigger"`
	Duration       *types.DurationWithUnit `json:"duration" binding:"required_with=StartOnTrigger"`
}

type ChangeStateParametersRequest struct {
	State  *types.CpsNumber `json:"state" binding:"required"`
	Output string           `json:"output" binding:"required,max=255"`
	Author string           `json:"author"`
	User   string           `json:"user" swaggerignore:"true"`
}

type AssocTicketParametersRequest struct {
	Ticket string `json:"ticket" binding:"required,max=255"`
	Output string `json:"output" binding:"max=255"`
	Author string `json:"author"`
	User   string `json:"user" swaggerignore:"true"`
}

type WebhookParameterRequest struct {
	Request       WebhookRequest          `json:"request" binding:"required"`
	DeclareTicket *WebhookDeclareTicket   `json:"declare_ticket" mapstructure:"declare_ticket"`
	RetryCount    int64                   `json:"retry_count" mapstructure:"retry_count" binding:"min=0"`
	RetryDelay    *types.DurationWithUnit `json:"retry_delay" mapstructure:"retry_delay"`
}

type WebhookRequest struct {
	URL    string `json:"url" binding:"required,url"`
	Method string `json:"method" binding:"required"`
	Auth   *struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	} `json:"auth"`
	Headers    map[string]string `json:"headers"`
	Payload    string            `json:"payload"`
	SkipVerify *bool             `json:"skip_verify" mapstructure:"skip_verify"`
}

type WebhookDeclareTicket struct {
	EmptyResponse bool              `mapstructure:"empty_response"`
	IsRegexp      bool              `mapstructure:"is_regexp"`
	Fields        map[string]string `mapstructure:",remain"`
}

func (t WebhookDeclareTicket) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"empty_response": t.EmptyResponse,
		"is_regexp":      t.IsRegexp,
	}

	for k, v := range t.Fields {
		m[k] = v
	}

	return json.Marshal(m)
}

type ParametersRequest struct {
	Output string `json:"output" binding:"max=255"`
	Author string `json:"author"`
	User   string `json:"user" swaggerignore:"true"`
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
	Type                     string                    `bson:"type" json:"type"`
	Comment                  string                    `bson:"comment" json:"comment"`
	Parameters               map[string]interface{}    `bson:"parameters,omitempty" json:"parameters,omitempty"`
	AlarmPatterns            pattern.AlarmPatternList  `bson:"alarm_patterns" json:"alarm_patterns"`
	EntityPatterns           pattern.EntityPatternList `bson:"entity_patterns" json:"entity_patterns"`
	DropScenarioIfNotMatched bool                      `bson:"drop_scenario_if_not_matched" json:"drop_scenario_if_not_matched"`
	EmitTrigger              bool                      `bson:"emit_trigger" json:"emit_trigger"`
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
