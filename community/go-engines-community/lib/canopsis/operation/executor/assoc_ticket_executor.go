package executor

import (
	"context"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// NewAssocTicketExecutor creates new executor.
func NewAssocTicketExecutor(metricsSender metrics.Sender) operationlib.Executor {
	return &assocTicketExecutor{metricsSender: metricsSender}
}

type assocTicketExecutor struct {
	metricsSender metrics.Sender
}

// Exec creates new assoc ticket step for alarm.
func (e *assocTicketExecutor) Exec(
	ctx context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ types.Entity,
	time types.CpsTime,
	role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationAssocTicketParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationAssocTicketParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	err := alarm.PartialUpdateAssocTicket(
		time,
		nil,
		params.Author,
		params.Ticket,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	go e.metricsSender.SendAssocTicket(ctx, *alarm, params.Author, time.Time)

	return types.AlarmChangeTypeAssocTicket, nil
}
