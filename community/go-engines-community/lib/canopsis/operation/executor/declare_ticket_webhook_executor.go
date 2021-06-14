package executor

import (
	"fmt"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/config"
	operationlib "git.canopsis.net/canopsis/go-engines/lib/canopsis/operation"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
)

// NewDeclareTicketWebhookExecutor creates new executor.
func NewDeclareTicketWebhookExecutor(cfg config.CanopsisConf) operationlib.Executor {
	return &declareTicketWebhookExecutor{cfg: cfg}
}

type declareTicketWebhookExecutor struct {
	cfg config.CanopsisConf
}

// Exec creates new declare ticket step for alarm.
func (e *declareTicketWebhookExecutor) Exec(
	operation types.Operation,
	alarm *types.Alarm,
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
		utils.TruncateString(params.Output, e.cfg.Alarm.OutputLength),
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
