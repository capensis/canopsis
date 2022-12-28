package rpc

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
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
	DeclareTicket        bool              `json:"declare_ticket,omitempty"`
	DeclareTicketRequest bool              `json:"declare_ticket_request,omitempty"`
	TicketResources      bool              `json:"ticket_resources,omitempty"`
	WebhookHeader        map[string]string `json:"webhook_header,omitempty"`
	WebhookResponse      map[string]any    `json:"webhook_response,omitempty"`
	WebhookError         *Error            `json:"webhook_error,omitempty"`
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
	Error           *types.RPCError       `json:"error"`
}

type WebhookEvent struct {
	Parameters   WebhookParameters      `json:"parameters"`
	Alarm        *types.Alarm           `json:"alarm"`
	Entity       *types.Entity          `json:"entity"`
	AckResources bool                   `json:"ack_resources"`
	Header       map[string]string      `json:"header,omitempty"`
	Response     map[string]interface{} `json:"response,omitempty"`
	Message      string                 `json:"message"`
}

type WebhookParameters struct {
	Request       request.Parameters            `json:"request"`
	DeclareTicket *request.WebhookDeclareTicket `json:"declare_ticket,omitempty"`
	Scenario      string                        `json:"scenario"`
	Author        string                        `json:"author"`
	User          string                        `json:"user"`
}

type WebhookResultEvent struct {
	Alarm           *types.Alarm          `json:"alarm"`
	AlarmChangeType types.AlarmChangeType `json:"alarm_change_type"`
	Header          map[string]string     `json:"header,omitempty"`
	Response        map[string]any        `json:"response,omitempty"`
	Error           *Error                `json:"error"`
}

type Error struct {
	Error error
}

func (e Error) MarshalJSON() ([]byte, error) {
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
