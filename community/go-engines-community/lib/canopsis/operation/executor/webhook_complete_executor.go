package executor

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewWebhookCompleteExecutor(metricsSender metrics.Sender) operationlib.Executor {
	return &webhookCompleteExecutor{metricsSender: metricsSender}
}

type webhookCompleteExecutor struct {
	metricsSender metrics.Sender
}

func (e *webhookCompleteExecutor) Exec(
	_ context.Context,
	operation types.Operation,
	alarm *types.Alarm,
	_ *types.Entity,
	time types.CpsTime,
	userID, role, initiator string,
) (types.AlarmChangeType, error) {
	params := operation.Parameters

	for i := len(alarm.Value.Steps) - 1; i >= 0; i-- {
		step := alarm.Value.Steps[i]
		if step.Execution == params.Execution && (step.Type == types.AlarmStepWebhookComplete || step.Type == types.AlarmStepWebhookFail) {
			return types.AlarmChangeTypeNone, nil
		}
	}

	if userID == "" {
		userID = params.User
	}

	if params.TicketInfo.Ticket == "" {
		err := alarm.PartialUpdateWebhookComplete(
			time,
			params.Execution,
			params.Author,
			params.Output,
			userID,
			role,
			initiator,
		)
		if err != nil {
			return "", err
		}

		return types.AlarmChangeTypeWebhookComplete, nil
	}

	err := alarm.PartialUpdateWebhookDeclareTicket(
		time,
		params.Execution,
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

	return types.AlarmChangeTypeDeclareTicketWebhook, nil
}
