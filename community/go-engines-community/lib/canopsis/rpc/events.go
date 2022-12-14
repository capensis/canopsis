package rpc

import (
	"encoding/json"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/request"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

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
	Author        string                        `json:"author"`
	User          string                        `json:"user"`
}

type WebhookResultEvent struct {
	Alarm           *types.Alarm           `json:"alarm"`
	AlarmChangeType types.AlarmChangeType  `json:"alarm_change_type"`
	Header          map[string]string      `json:"header,omitempty"`
	Response        map[string]interface{} `json:"response,omitempty"`
	Error           *Error                 `json:"error"`
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
