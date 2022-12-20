package executor

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewSnoozeExecutor creates new executor.
func NewSnoozeExecutor(configProvider config.AlarmConfigProvider) operation.Executor {
	return &snoozeExecutor{configProvider: configProvider}
}

type snoozeExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new snooze step for alarm.
func (e *snoozeExecutor) Exec(
	_ context.Context,
	op types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := op.Parameters
	if params.Duration == nil {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	if alarm.Value.Snooze != nil {
		return "", nil
	}

	err := alarm.PartialUpdateSnooze(
		time,
		*params.Duration,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeSnooze, nil
}
