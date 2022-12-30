package webhook

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

const (
	StatusRunning = iota
	StatusStarted
	StatusSucceeded
	StatusFailed
)

type History struct {
	ID            string                        `bson:"_id" json:"_id"`
	Alarms        []string                      `bson:"alarms" json:"alarms"`
	Scenario      string                        `bson:"scenario,omitempty" json:"scenario,omitempty"`
	Action        int64                         `bson:"action,omitempty" json:"action,omitempty"`
	Execution     string                        `bson:"execution" json:"execution"`
	Name          string                        `bson:"name" json:"name"`
	SystemName    string                        `bson:"system_name" json:"system_name"`
	Status        int64                         `bson:"status" json:"status"`
	Comment       string                        `bson:"comment,omitempty" json:"comment,omitempty"`
	Request       request.Parameters            `bson:"request" json:"request"`
	DeclareTicket *request.WebhookDeclareTicket `bson:"declare_ticket,omitempty" json:"declare_ticket,omitempty"`
	NextExec      string                        `bson:"next_exec,omitempty" json:"next_exec,omitempty"`
	Index         *int64                        `bson:"index,omitempty" json:"index,omitempty"`
	StopOnFail    bool                          `bson:"stop_on_fail,omitempty" json:"stop_on_fail,omitempty"`
	FailReason    string                        `bson:"fail_reason,omitempty" json:"fail_reason,omitempty"`
	RawRequest    string                        `bson:"raw_request,omitempty" json:"raw_request,omitempty"`
	RawResponse   string                        `bson:"raw_response,omitempty" json:"raw_response,omitempty"`
	TicketID      string                        `bson:"ticket_id,omitempty" json:"ticket_id,omitempty"`
	TicketURL     string                        `bson:"ticket_url,omitempty" json:"ticket_url,omitempty"`
	TicketData    map[string]string             `bson:"ticket_data,omitempty" json:"ticket_data,omitempty"`
	UserID        string                        `bson:"user,omitempty" json:"user,omitempty"`
	Username      string                        `bson:"username,omitempty" json:"username,omitempty"`
	CreatedAt     types.CpsTime                 `bson:"created_at,omitempty" json:"created_at,omitempty"`
	LaunchedAt    types.CpsTime                 `bson:"launched_at,omitempty" json:"launched_at,omitempty"`
	CompletedAt   types.CpsTime                 `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
}
