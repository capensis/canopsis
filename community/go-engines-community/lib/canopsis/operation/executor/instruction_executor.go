package executor

import (
	"context"

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
			types.EventTypeAutoInstructionStarted:   types.AlarmStepAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted: types.AlarmStepAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:    types.AlarmStepAutoInstructionFail,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmStepInstructionAbort,
			// Job
			types.EventTypeInstructionJobStarted:   types.AlarmStepInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmStepInstructionJobComplete,
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
			types.EventTypeAutoInstructionStarted:   types.AlarmChangeTypeAutoInstructionStart,
			types.EventTypeAutoInstructionCompleted: types.AlarmChangeTypeAutoInstructionComplete,
			types.EventTypeAutoInstructionFailed:    types.AlarmChangeTypeAutoInstructionFail,
			// Manual and auto instruction
			types.EventTypeInstructionAborted: types.AlarmChangeTypeInstructionAbort,
			// Job
			types.EventTypeInstructionJobStarted:   types.AlarmChangeTypeInstructionJobStart,
			types.EventTypeInstructionJobCompleted: types.AlarmChangeTypeInstructionJobComplete,
			types.EventTypeInstructionJobFailed:    types.AlarmChangeTypeInstructionJobFail,
		},
		metricsSender: metricsSender,
	}
}

func (e *instructionExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters

	if userID == "" {
		userID = params.User
	}

	alarmStepType, ok := e.alarmStepTypeMap[op.Type]
	if !ok {
		return "", nil
	}

	alarmChangeType, ok := e.alarmChangeTypeMap[op.Type]
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
	case types.AlarmStepInstructionComplete, types.AlarmStepInstructionFail:
		instrID := params.Instruction

		assigned := false
		for _, assignedInstr := range alarm.KPIAssignedInstructions {
			if assignedInstr == instrID {
				assigned = true
				break
			}
		}

		if !assigned {
			break
		}

		executed := false
		for _, executedInstr := range alarm.KPIExecutedInstructions {
			if executedInstr == instrID {
				executed = true
				break
			}
		}

		if executed {
			break
		}

		alarm.PartialUpdateAddExecutedInstruction(instrID)

		if len(alarm.KPIExecutedInstructions) == 0 {
			e.metricsSender.SendInstructionExecutionForAlarm(alarm.EntityID, time.Time)
		}

		e.metricsSender.SendInstructionExecutionForInstruction(instrID, time.Time)
	case types.AlarmStepAutoInstructionStart:
		e.metricsSender.SendAutoInstructionStart(*alarm, time.Time)
	}

	return alarmChangeType, nil
}
