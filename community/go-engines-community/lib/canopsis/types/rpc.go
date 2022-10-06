package types

import (
	"encoding/json"
	"errors"
)

type RPCAxeEvent struct {
	EventType  string           `json:"event_type"`
	Parameters RPCAxeParameters `json:"parameters,omitempty"`
	Alarm      *Alarm           `json:"alarm"`
	Entity     *Entity          `json:"entity"`
}

type RPCAxeParameters struct {
	Output string `json:"output"`
	Author string `json:"author"`
	User   string `json:"user"`
	// ChangeState
	State *CpsNumber `json:"state"`
	// AssocTicket
	Ticket string `json:"ticket"`
	// Snooze and Pbehavior
	Duration *DurationWithUnit `json:"duration"`
	// Pbehavior
	Name           string   `json:"name"`
	Reason         string   `json:"reason"`
	Type           string   `json:"type"`
	RRule          string   `json:"rrule"`
	Tstart         *CpsTime `json:"tstart"`
	Tstop          *CpsTime `json:"tstop"`
	StartOnTrigger *bool    `json:"start_on_trigger"`
	// Instruction
	Execution   string `json:"execution"`
	Instruction string `json:"instruction"`
}

type RPCAxeResultEvent struct {
	Alarm           *Alarm          `json:"alarm"`
	AlarmChangeType AlarmChangeType `json:"alarm_change"`
	Error           *RPCError       `json:"error"`
}

type RPCServiceEvent struct {
	Alarm       *Alarm       `json:"alarm"`
	Entity      *Entity      `json:"entity"`
	AlarmChange *AlarmChange `json:"alarm_change"`
}

type RPCServiceResultEvent struct {
	Error *RPCError `json:"error"`
}

type RPCPBehaviorEvent struct {
	Alarm  *Alarm                 `json:"alarm"`
	Entity *Entity                `json:"entity"`
	Params RPCPBehaviorParameters `json:"params"`
}

type RPCPBehaviorParameters struct {
	Author         string            `json:"author"`
	UserID         string            `json:"user"`
	Name           string            `json:"name"`
	Reason         string            `json:"reason"`
	Type           string            `json:"type"`
	RRule          string            `json:"rrule"`
	Tstart         *CpsTime          `json:"tstart,omitempty"`
	Tstop          *CpsTime          `json:"tstop,omitempty"`
	StartOnTrigger *bool             `json:"start_on_trigger,omitempty"`
	Duration       *DurationWithUnit `json:"duration,omitempty"`
}

type RPCPBehaviorResultEvent struct {
	Alarm    *Alarm    `json:"alarm"`
	Entity   *Entity   `json:"entity"`
	PbhEvent Event     `json:"event"`
	Error    *RPCError `json:"error"`
}

type RPCError struct {
	Error error
}

func (e RPCError) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error.Error())
}

func (e *RPCError) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	e.Error = errors.New(s)

	return nil
}

type RPCWebhookEvent struct {
	Parameters   RPCWebhookParameters   `json:"parameters"`
	Alarm        *Alarm                 `json:"alarm"`
	Entity       *Entity                `json:"entity"`
	AckResources bool                   `json:"ack_resources"`
	Header       map[string]string      `json:"header,omitempty"`
	Response     map[string]interface{} `json:"response,omitempty"`
	Message      string                 `json:"message"`
}

type RPCWebhookParameters struct {
	Request       WebhookRequest        `json:"request"`
	DeclareTicket *WebhookDeclareTicket `json:"declare_ticket,omitempty"`
	RetryCount    int64                 `json:"retry_count,omitempty"`
	RetryDelay    *DurationWithUnit     `json:"retry_delay,omitempty"`
	Author        string                `json:"author"`
	User          string                `json:"user"`
}

type RPCWebhookResultEvent struct {
	Alarm           *Alarm                 `json:"alarm"`
	AlarmChangeType AlarmChangeType        `json:"alarm_change_type"`
	Header          map[string]string      `json:"header,omitempty"`
	Response        map[string]interface{} `json:"response,omitempty"`
	Error           *RPCError              `json:"error"`
}

type RPCRemediationEvent struct {
	Alarm       *Alarm      `json:"alarm"`
	Entity      *Entity     `json:"entity"`
	AlarmChange AlarmChange `json:"alarm_change"`
}

type RPCRemediationJobEvent struct {
	JobExecutionID string `json:"job_execution_id"`
}
