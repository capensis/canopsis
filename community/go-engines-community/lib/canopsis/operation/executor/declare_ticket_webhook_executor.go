package executor

import (
	"context"
	"fmt"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
)

// NewDeclareTicketWebhookExecutor creates new executor.
func NewDeclareTicketWebhookExecutor(configProvider config.AlarmConfigProvider) operationlib.Executor {
	return &declareTicketWebhookExecutor{configProvider: configProvider}
}

type declareTicketWebhookExecutor struct {
	configProvider config.AlarmConfigProvider
}

// Exec creates new declare ticket step for alarm.
func (e *declareTicketWebhookExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationDeclareTicketParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationDeclareTicketParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	err := alarm.PartialUpdateDeclareTicket(
		time,
		params.Author,
		utils.TruncateString(params.Output, e.configProvider.Get().OutputLength),
		params.Ticket,
		params.Data,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeDeclareTicketWebhook, nil
}
