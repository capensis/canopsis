package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// NewDeclareTicketWebhookExecutor creates new executor.
func NewDeclareTicketWebhookExecutor(configProvider config.AlarmConfigProvider) operation.Executor {
	return &declareTicketWebhookExecutor{configProvider: configProvider}
}

type declareTicketWebhookExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new declare ticket step for alarm.
func (e *declareTicketWebhookExecutor) Exec(
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

	err := alarm.PartialUpdateDeclareTicket(
		time,
		params.Author,
		userID,
		role,
		initiator,
		params.TicketInfo,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeDeclareTicketWebhook, nil
}
