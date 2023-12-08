package rpc

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type AxeEvent struct {
	EventType  string        `json:"event_type"`
	Parameters AxeParameters `json:"parameters,omitempty"`
	Alarm      *types.Alarm  `json:"alarm"`
	Entity     *types.Entity `json:"entity"`
}

type AxeParameters struct {
	Output string `json:"output,omitempty"`
	Author string `json:"author,omitempty"`
	User   string `json:"user,omitempty"`
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
	Duration *types.DurationWithUnit `json:"duration,omitempty"`
	// Pbehavior
	Name           string         `json:"name,omitempty"`
	Reason         string         `json:"reason,omitempty"`
	Type           string         `json:"type,omitempty"`
	RRule          string         `json:"rrule,omitempty"`
	Tstart         *types.CpsTime `json:"tstart,omitempty"`
	Tstop          *types.CpsTime `json:"tstop,omitempty"`
	StartOnTrigger *bool          `json:"start_on_trigger,omitempty"`
	// Instruction
	Execution   string `json:"execution,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}

type AxeResultEvent struct {
	Alarm           *types.Alarm          `json:"alarm"`
	AlarmChangeType types.AlarmChangeType `json:"alarm_change"`
	WebhookHeader   map[string]string     `json:"webhook_header,omitempty"`
	WebhookResponse map[string]any        `json:"webhook_response,omitempty"`
	Error           *Error                `json:"error"`
}

type WebhookEvent struct {
	Execution string `json:"execution"`
}

type PbehaviorRecomputeEvent struct {
	Ids []string `json:"ids"`
}

type PbehaviorEvent struct {
	Alarm  *types.Alarm        `json:"alarm"`
	Entity *types.Entity       `json:"entity"`
	Params PbehaviorParameters `json:"params"`
}

type PbehaviorParameters struct {
	Author         string                  `json:"author"`
	UserID         string                  `json:"user"`
	Name           string                  `json:"name"`
	Reason         string                  `json:"reason"`
	Type           string                  `json:"type"`
	RRule          string                  `json:"rrule"`
	Tstart         *types.CpsTime          `json:"tstart,omitempty"`
	Tstop          *types.CpsTime          `json:"tstop,omitempty"`
	StartOnTrigger *bool                   `json:"start_on_trigger,omitempty"`
	Duration       *types.DurationWithUnit `json:"duration,omitempty"`
}

type PbehaviorResultEvent struct {
	Alarm    *types.Alarm  `json:"alarm"`
	Entity   *types.Entity `json:"entity"`
	PbhEvent types.Event   `json:"event"`
	Error    *Error        `json:"error"`
}

type DynamicInfosEvent struct {
	// On update rule
	ID     string `json:"_id"`
	Action int    `json:"action"`
	// On update alarm
	Alarm           *types.Alarm          `json:"alarm"`
	Entity          *types.Entity         `json:"entity"`
	AlarmChangeType types.AlarmChangeType `json:"alarm_change"`
}

type DynamicInfosResultEvent struct {
	Alarm           *types.Alarm          `json:"alarm"`
	AlarmChangeType types.AlarmChangeType `json:"alarm_change"`
	Error           *Error                `json:"error"`
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
