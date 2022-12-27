package executor

import (
	"context"

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
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := operation.Parameters

	if userID == "" {
		userID = params.User
	}

	err := alarm.PartialUpdateAssocTicket(
		time,
		nil,
		params.Author,
		params.Ticket,
		params.TicketMetaAlarmID,
		params.TicketRuleName,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	metricsUserID := ""
	if initiator == types.InitiatorUser {
		metricsUserID = userID
	}
	e.metricsSender.SendTicket(*alarm, metricsUserID, time.Time)

	return types.AlarmChangeTypeAssocTicket, nil
}
