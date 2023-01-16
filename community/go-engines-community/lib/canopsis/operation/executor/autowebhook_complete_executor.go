package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewAutoWebhookCompleteExecutor(metricsSender metrics.Sender) operationlib.Executor {
	return &autoWebhookCompleteExecutor{metricsSender: metricsSender}
}

type autoWebhookCompleteExecutor struct {
	metricsSender metrics.Sender
}

func (e *autoWebhookCompleteExecutor) Exec(
	_ context.Context,
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

	if params.Ticket == "" {
		err := alarm.PartialUpdateAddStep(
			types.AlarmStepWebhookComplete,
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

		return types.AlarmChangeTypeAutoWebhookComplete, nil
	}

	err := alarm.PartialUpdateWebhookDeclareTicket(
		time,
		params.Author,
		params.Output,
		userID,
		role,
		initiator,
		params.TicketInfo,
	)
	if err != nil {
		return "", err
	}

	metricsUserID := ""
	if initiator == types.InitiatorUser {
		metricsUserID = userID
	}
	e.metricsSender.SendTicket(*alarm, metricsUserID, time.Time)

	return types.AlarmChangeTypeAutoDeclareTicketWebhook, nil
}
