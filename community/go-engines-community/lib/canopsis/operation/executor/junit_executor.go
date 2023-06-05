package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

type junitExecutor struct {
	alarmStepTypeMap   map[string]string
	alarmChangeTypeMap map[string]types.AlarmChangeType
}

func NewJunitExecutor() operation.Executor {
	return &junitExecutor{
		alarmStepTypeMap: map[string]string{
			types.EventTypeJunitTestSuiteUpdated: types.AlarmStepJunitTestSuiteUpdate,
			types.EventTypeJunitTestCaseUpdated:  types.AlarmStepJunitTestCaseUpdate,
		},
		alarmChangeTypeMap: map[string]types.AlarmChangeType{
			types.EventTypeJunitTestSuiteUpdated: types.AlarmChangeTypeJunitTestSuiteUpdate,
			types.EventTypeJunitTestCaseUpdated:  types.AlarmChangeTypeJunitTestCaseUpdate,
		},
	}
}

func (e *junitExecutor) Exec(
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

	err := alarm.PartialUpdateAddStep(
		alarmStepType,
		time,
		params.Author,
		params.Output,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return alarmChangeType, nil
}
