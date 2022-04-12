package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type RPCAxeEvent struct {
	EventType  string      `json:"event_type"`
	Parameters interface{} `json:"parameters,omitempty"`
	Alarm      *Alarm      `json:"alarm"`
	Entity     *Entity     `json:"entity"`
}

func (e *RPCAxeEvent) UnmarshalJSON(b []byte) error {
	type Alias RPCAxeEvent
	tmp := struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	switch e.EventType {
	case ActionTypeSnooze:
		var params OperationSnoozeParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	case ActionTypeChangeState:
		var params OperationChangeStateParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	case ActionTypeAssocTicket:
		var params OperationAssocTicketParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	case ActionTypePbehavior:
		var params ActionPBehaviorParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	case ActionTypeWebhook:
		var params WebhookParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	case EventTypeInstructionStarted, EventTypeInstructionPaused,
		EventTypeInstructionResumed, EventTypeInstructionCompleted,
		EventTypeInstructionFailed, EventTypeInstructionAborted,
		EventTypeAutoInstructionStarted, EventTypeAutoInstructionCompleted,
		EventTypeAutoInstructionFailed,
		EventTypeInstructionJobStarted, EventTypeInstructionJobCompleted,
		EventTypeInstructionJobAborted, EventTypeInstructionJobFailed:
		var params OperationInstructionParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	default:
		var params OperationParameters
		err := mapstructure.Decode(e.Parameters, &params)
		if err != nil {
			return fmt.Errorf("cannot decode map struct : %v", err)
		}
		e.Parameters = params
	}

	return nil
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
	Alarm  *Alarm                    `json:"alarm"`
	Entity *Entity                   `json:"entity"`
	Params ActionPBehaviorParameters `json:"params"`
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
	Parameters   WebhookParameters      `json:"parameters"`
	Alarm        *Alarm                 `json:"alarm"`
	Entity       *Entity                `json:"entity"`
	AckResources bool                   `json:"ack_resources"`
	Header       map[string]string      `json:"header,omitempty"`
	Response     map[string]interface{} `json:"response,omitempty"`
	Message      string                 `json:"message"`
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
