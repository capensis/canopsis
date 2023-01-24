package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

// NewDeclareTicketExecutor creates new executor.
func NewDeclareTicketExecutor() operation.Executor {
	return &declareTicketExecutor{}
}

type declareTicketExecutor struct{}

// Exec triggers manual declareticket trigger.
func (e *declareTicketExecutor) Exec(
	_ context.Context,
	_ types.Operation,
	_ *types.Alarm,
	_ *types.Entity,
	_ types.CpsTime,
	_, _, _ string,
) (types.AlarmChangeType, error) {
	return types.AlarmChangeTypeDeclareTicket, nil
}
