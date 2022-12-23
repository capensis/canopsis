package types

import (
	"encoding/json"
	"errors"
)

// todo move all to rpc package

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

type RPCRemediationEvent struct {
	Alarm       *Alarm      `json:"alarm"`
	Entity      *Entity     `json:"entity"`
	AlarmChange AlarmChange `json:"alarm_change"`
}

type RPCRemediationJobEvent struct {
	JobExecutionID string `json:"job_execution_id"`
}
