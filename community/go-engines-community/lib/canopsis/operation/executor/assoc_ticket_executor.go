package executor

import (
	"context"
	"fmt"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// NewAssocTicketExecutor creates new executor.
func NewAssocTicketExecutor() operationlib.Executor {
	return &assocTicketExecutor{}
}

type assocTicketExecutor struct {
}

// Exec creates new assoc ticket step for alarm.
func (e *assocTicketExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	var params types.OperationAssocTicketParameters
	var ok bool
	if params, ok = operation.Parameters.(types.OperationAssocTicketParameters); !ok {
		return "", fmt.Errorf("invalid parameters")
	}

	if userID == "" {
		userID = params.User
	}

	err := alarm.PartialUpdateAssocTicket(
		time,
		nil,
		params.Author,
		params.Ticket,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeAssocTicket, nil
}
