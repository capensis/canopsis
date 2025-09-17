package webhook

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
)

const (
	StatusCreated = iota
	StatusRunning
	StatusSucceeded
	StatusFailed
	StatusAborted
)

const MultipleURLsDelimiter = ","

type History struct {
	BaseHistory `bson:",inline"`

	Execution         string   `bson:"execution" json:"execution"`
	Alarms            []string `bson:"alarms,omitempty" json:"alarms,omitempty"`
	Scenario          string   `bson:"scenario,omitempty" json:"scenario,omitempty"`
	DeclareTicketRule string   `bson:"declare_ticket_rule,omitempty" json:"declare_ticket_rule,omitempty"`
	Name              string   `bson:"name" json:"name"`

	Index              int64  `bson:"index" json:"index"`
	NextExec           string `bson:"next_exec,omitempty" json:"next_exec,omitempty"`
	StopOnFail         bool   `bson:"stop_on_fail,omitempty" json:"stop_on_fail,omitempty"`
	StopOnSuccess      bool   `bson:"stop_on_success,omitempty" json:"stop_on_success,omitempty"`
	MultipleURLs       bool   `bson:"multiple_urls,omitempty" json:"multiple_urls,omitempty"`
	ResolvedRequestURL string `bson:"resolved_request_url,omitempty" json:"resolved_request_url,omitempty"`

	SystemName      string                        `bson:"system_name,omitempty" json:"system_name,omitempty"`
	EmitTrigger     bool                          `bson:"emit_trigger,omitempty" json:"emit_trigger,omitempty"`
	Comment         string                        `bson:"comment,omitempty" json:"comment,omitempty"`
	AuthToken       *request.WebhookAuthToken     `bson:"auth_token,omitempty" json:"auth_token,omitempty"`
	DeclareTicket   *request.WebhookDeclareTicket `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	TicketResources bool                          `bson:"ticket_resources,omitempty" json:"ticket_resources,omitempty"`
	UserID          string                        `bson:"user,omitempty" json:"user,omitempty"`
	Username        string                        `bson:"username,omitempty" json:"username,omitempty"`
	Initiator       string                        `bson:"initiator,omitempty" json:"initiator,omitempty"`
	EventInitiator  string                        `bson:"event_initiator,omitempty" json:"event_initiator,omitempty"`
	EventOutput     string                        `bson:"event_output,omitempty" json:"event_output,omitempty"`
	Trigger         string                        `bson:"trigger,omitempty" json:"trigger,omitempty"`

	ResponseCode   int64             `bson:"response_code,omitempty" json:"response_code,omitempty"`
	ResponseHeader map[string]string `bson:"response_header,omitempty" json:"response_header,omitempty"`
	ResponseBody   map[string]any    `bson:"response_body,omitempty" json:"response_body,omitempty"`

	TicketID   string            `bson:"ticket_id,omitempty" json:"ticket_id,omitempty"`
	TicketURL  string            `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketData map[string]string `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`

	IsTest bool `bson:"is_test,omitempty" json:"is_test,omitempty"`
}

type TokenHistory struct {
	BaseHistory        `bson:",inline"`
	Rule               string                    `bson:"rule" json:"rule"`
	ResponseField      string                    `bson:"response_field,omitempty" json:"response_field,omitempty"`
	Template           string                    `bson:"template,omitempty" json:"template,omitempty"`
	ExpirationDuration datetime.DurationWithUnit `bson:"expiration_duration" json:"expiration_duration"`
	Token              string                    `bson:"token,omitempty" json:"token,omitempty"`
	ExpiredAt          datetime.MicroTime        `bson:"expired_at,omitempty" json:"expired_at,omitempty"`
}

type BaseHistory struct {
	ID string `bson:"_id" json:"_id"`

	Status     int64  `bson:"status" json:"status"`
	FailReason string `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`

	Request     request.Parameters `bson:"request" json:"request"`
	RawRequest  string             `bson:"raw_request,omitempty" json:"raw_request,omitempty"`
	RawResponse string             `bson:"raw_response,omitempty" json:"raw_response,omitempty"`

	CreatedAt   datetime.MicroTime `bson:"created_at" json:"created_at"`
	LaunchedAt  datetime.MicroTime `bson:"launched_at,omitempty" json:"launched_at,omitempty"`
	CompletedAt datetime.MicroTime `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	LastPing    datetime.MicroTime `bson:"last_ping,omitempty" json:"last_ping,omitempty"`
}
