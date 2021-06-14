package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
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
