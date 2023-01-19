package executor

import (
	"context"

	operationlib "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/operation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
)

func NewWebhookFailExecutor() operationlib.Executor {
	return &webhookFailExecutor{}
}

type webhookFailExecutor struct {
}

func (e *webhookFailExecutor) Exec(
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

	if params.TicketInfo.TicketRuleID != "" {
		err := alarm.PartialUpdateWebhookDeclareTicketFail(
			params.WebhookRequest,
			time,
			params.Execution,
			params.Author,
			params.Output,
			params.WebhookFailReason,
			userID,
			role,
			initiator,
			params.TicketInfo,
		)
		if err != nil {
			return "", err
		}

		return types.AlarmChangeTypeDeclareTicketWebhookFail, nil
	}

	err := alarm.PartialUpdateWebhookFail(
		time,
		params.Execution,
		params.Author,
		params.Output,
		params.WebhookFailReason,
		userID,
		role,
		initiator,
	)
	if err != nil {
		return "", err
	}

	return types.AlarmChangeTypeWebhookFail, nil
}
