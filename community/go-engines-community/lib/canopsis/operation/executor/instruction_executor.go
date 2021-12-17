package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type instructionExecutor struct {
	alarmStepTypeMap   map[string]string
	alarmChangeTypeMap map[string]types.AlarmChangeType
	metricsSender      metrics.Sender
}

func NewInstructionExecutor(metricsSender metrics.Sender) operation.Executor {
	return &instructionExecutor{
		alarmStepTypeMap: map[string]string{
			// Manual instruction
			types.EventTypeInstructionStarted:   types.AlarmStepInstructionStart,
			types.EventTypeInstructionPaused:    types.AlarmStepInstructionPause,
			types.EventTypeInstructionResumed:   types.AlarmStepInstructionResume,
			types.EventTypeInstructionCompleted: types.AlarmStepInstructionComplete,
			types.EventTypeInstructionFailed:    types.AlarmStepInstructionFail,
			// Auto instruction
			types.EventTypeAutoInstructionStarted:        types.AlarmStepAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted:      types.AlarmStepAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:         types.AlarmStepAutoInstructionFail,
			types.EventTypeAutoInstructionAlreadyRunning: types.AlarmStepAutoInstructionAlreadyRunning,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmStepInstructionAbort,
			// Job
			types.EventTypeInstructionJobStarted:   types.AlarmStepInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmStepInstructionJobComplete,
			types.EventTypeInstructionJobAborted:   types.AlarmStepInstructionJobAbort,
			types.EventTypeInstructionJobFailed:    types.AlarmStepInstructionJobFail,
		},
		alarmChangeTypeMap: map[string]types.AlarmChangeType{
			// Manual instruction
			types.EventTypeInstructionStarted:   types.AlarmChangeTypeInstructionStart,
			types.EventTypeInstructionPaused:    types.AlarmChangeTypeInstructionPause,
			types.EventTypeInstructionResumed:   types.AlarmChangeTypeInstructionResume,
			types.EventTypeInstructionCompleted: types.AlarmChangeTypeInstructionComplete,
			types.EventTypeInstructionFailed:    types.AlarmChangeTypeInstructionFail,
			// Auto instruction
			types.EventTypeAutoInstructionStarted:        types.AlarmChangeTypeAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted:      types.AlarmChangeTypeAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:         types.AlarmChangeTypeAutoInstructionFail,
			types.EventTypeAutoInstructionAlreadyRunning: types.AlarmChangeTypeAutoInstructionAlreadyRunning,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmChangeTypeInstructionAbort,
			// Job
			types.EventTypeInstructionJobStarted:   types.AlarmChangeTypeInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmChangeTypeInstructionJobComplete,
			types.EventTypeInstructionJobAborted:   types.AlarmChangeTypeInstructionJobAbort,
			types.EventTypeInstructionJobFailed:    types.AlarmChangeTypeInstructionJobFail,
		},
		metricsSender: metricsSender,
	}
}

func (e *instructionExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationInstructionParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationInstructionParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	alarmStepType, ok := e.alarmStepTypeMap[operation.Type]
	if !ok {
		return "", nil
	}

	alarmChangeType, ok := e.alarmChangeTypeMap[operation.Type]
	if !ok {
		return "", nil
	}

	err := alarm.PartialUpdateAddInstructionStep(
		alarmStepType,
		time,
		params.Execution,
		params.Author,
		params.Output,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	switch alarmChangeType {
	case types.AlarmStepAutoInstructionStart, types.AlarmStepAutoInstructionAlreadyRunning:
		go e.metricsSender.SendAutoInstructionStart(context.Background(), *alarm, time.Time)
	}

	return alarmChangeType, nil
}
