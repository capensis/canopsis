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

	if userID == "" {
		userID = params.User
	}

	if params.TicketInfo.TicketRuleID != "" {
		err := alarm.PartialUpdateWebhookDeclareTicketFail(
			params.DeclareTicketRequest,
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

		return types.AlarmChangeTypeDeclareTicketWebhookFail, nil
	}

	err := alarm.PartialUpdateAddStep(
		types.AlarmStepWebhookFail,
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

	return types.AlarmChangeTypeWebhookFail, nil
}
