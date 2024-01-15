package rpc

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type AxeEvent struct {
	EventType   string        `json:"event_type"`
	Parameters  AxeParameters `json:"parameters,omitempty"`
	Alarm       *types.Alarm  `json:"alarm,omitempty"`
	AlarmID     string        `json:"alarm_id,omitempty"`
	Entity      *types.Entity `json:"entity,omitempty"`
	Healthcheck bool          `json:"healthcheck,omitempty"`
}

type AxeParameters struct {
	Output    string           `json:"output,omitempty"`
	Author    string           `json:"author,omitempty"`
	User      string           `json:"user,omitempty"`
	Role      string           `json:"role,omitempty"`
	Initiator string           `json:"initiator,omitempty"`
	Timestamp datetime.CpsTime `json:"timestamp,omitempty"`
	// ChangeState
	State *types.CpsNumber `json:"state,omitempty"`
	// AssocTicket and Webhook
	types.TicketInfo
	// Webhook
	TicketResources   bool              `json:"ticket_resources,omitempty"`
	WebhookRequest    bool              `json:"webhook_request,omitempty"`
	WebhookHeader     map[string]string `json:"webhook_header,omitempty"`
	WebhookResponse   map[string]any    `json:"webhook_response,omitempty"`
	WebhookFailReason string            `json:"webhook_fail_reason,omitempty"`
	WebhookError      *Error            `json:"webhook_error,omitempty"`
	EmitTrigger       bool              `json:"emit_trigger,omitempty"`
	// Snooze and Pbehavior
	Duration *datetime.DurationWithUnit `json:"duration,omitempty"`
	// Pbehavior enter
	PbehaviorInfo types.PbehaviorInfo `json:"pbehavior_info,omitempty"`
	// Pbehavior create
	Name           string            `json:"name,omitempty"`
	Reason         string            `json:"reason,omitempty"`
	Type           string            `json:"type,omitempty"`
	RRule          string            `json:"rrule,omitempty"`
	Tstart         *datetime.CpsTime `json:"tstart,omitempty"`
	Tstop          *datetime.CpsTime `json:"tstop,omitempty"`
	StartOnTrigger *bool             `json:"start_on_trigger,omitempty"`
	// Instruction
	Execution   string `json:"execution,omitempty"`
	Instruction string `json:"instruction,omitempty"`
	// Trigger
	Trigger string `json:"trigger,omitempty"`
	// Check
	LongOutput    string            `json:"long_output,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	Connector     string            `json:"connector,omitempty"`
	ConnectorName string            `json:"connector_name,omitempty"`
	// Idle events
	IdleRuleApply string `json:"idle_rule_apply,omitempty"`
	// Meta alarm create
	MetaAlarmRuleID     string   `json:"meta_alarm_rule_id,omitempty"`
	MetaAlarmValuePath  string   `json:"meta_alarm_value_path,omitempty"`
	DisplayName         string   `json:"display_name,omitempty"`
	MetaAlarmChildren   []string `json:"meta_alarm_children,omitempty"`
	StateSettingUpdated bool     `json:"state_setting_updated,omitempty"`
}

type AxeResultEvent struct {
	Alarm           *types.Alarm          `json:"alarm"`
	AlarmChangeType types.AlarmChangeType `json:"alarm_change"`
	WebhookHeader   map[string]string     `json:"webhook_header,omitempty"`
	WebhookResponse map[string]any        `json:"webhook_response,omitempty"`
	Error           *Error                `json:"error"`
}

type WebhookEvent struct {
	Execution   string `json:"execution"`
	Healthcheck bool   `json:"healthcheck"`
}

type PbehaviorRecomputeEvent struct {
	Ids []string `json:"ids"`
}

type PbehaviorEvent struct {
	Alarm       *types.Alarm        `json:"alarm"`
	Entity      *types.Entity       `json:"entity"`
	Params      PbehaviorParameters `json:"params"`
	Healthcheck bool                `json:"healthcheck"`
}

type PbehaviorParameters struct {
	Author         string                     `json:"author"`
	UserID         string                     `json:"user"`
	Name           string                     `json:"name"`
	Reason         string                     `json:"reason"`
	Type           string                     `json:"type"`
	RRule          string                     `json:"rrule"`
	Tstart         *datetime.CpsTime          `json:"tstart,omitempty"`
	Tstop          *datetime.CpsTime          `json:"tstop,omitempty"`
	StartOnTrigger *bool                      `json:"start_on_trigger,omitempty"`
	Duration       *datetime.DurationWithUnit `json:"duration,omitempty"`
}

type PbehaviorResultEvent struct {
	Alarm    *types.Alarm  `json:"alarm"`
	Entity   *types.Entity `json:"entity"`
	PbhEvent types.Event   `json:"event"`
	Error    *Error        `json:"error"`
}

type Error struct {
	Error error
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error.Error())
}

func (e *Error) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	e.Error = errors.New(s)

	return nil
}
