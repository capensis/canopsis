package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
)

type instructionExecutor struct {
	alarmStepTypeMap   map[string]string
	alarmChangeTypeMap map[string]types.AlarmChangeType
}

func NewInstructionExecutor() operation.Executor {
	return &instructionExecutor{
		alarmStepTypeMap: map[string]string{
			types.EventTypeInstructionStarted:      types.AlarmStepInstructionStart,
			types.EventTypeInstructionPaused:       types.AlarmStepInstructionPause,
			types.EventTypeInstructionResumed:      types.AlarmStepInstructionResume,
			types.EventTypeInstructionCompleted:    types.AlarmStepInstructionComplete,
			types.EventTypeInstructionAborted:      types.AlarmStepInstructionAbort,
			types.EventTypeInstructionFailed:       types.AlarmStepInstructionFail,
			types.EventTypeInstructionJobStarted:   types.AlarmStepInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmStepInstructionJobComplete,
			types.EventTypeInstructionJobAborted:   types.AlarmStepInstructionJobAbort,
			types.EventTypeInstructionJobFailed:    types.AlarmStepInstructionJobFail,
		},
		alarmChangeTypeMap: map[string]types.AlarmChangeType{
			types.EventTypeInstructionStarted:      types.AlarmChangeTypeInstructionStart,
			types.EventTypeInstructionPaused:       types.AlarmChangeTypeInstructionPause,
			types.EventTypeInstructionResumed:      types.AlarmChangeTypeInstructionResume,
			types.EventTypeInstructionCompleted:    types.AlarmChangeTypeInstructionComplete,
			types.EventTypeInstructionAborted:      types.AlarmChangeTypeInstructionAbort,
			types.EventTypeInstructionFailed:       types.AlarmChangeTypeInstructionFail,
			types.EventTypeInstructionJobStarted:   types.AlarmChangeTypeInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmChangeTypeInstructionJobComplete,
			types.EventTypeInstructionJobAborted:   types.AlarmChangeTypeInstructionJobAbort,
			types.EventTypeInstructionJobFailed:    types.AlarmChangeTypeInstructionJobFail,
		},
	}
}

func (e *instructionExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	alarmStepType, ok := e.alarmStepTypeMap[operation.Type]
	if !ok {
		return "", nil
	}

	alarmChangeType, ok := e.alarmChangeTypeMap[operation.Type]
	if !ok {
		return "", nil
	}

	err := alarm.PartialUpdateAddStep(
		alarmStepType,
		time,
		params.Author,
		params.Output,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return alarmChangeType, nil
}
