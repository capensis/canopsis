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
