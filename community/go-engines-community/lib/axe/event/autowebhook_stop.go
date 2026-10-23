package event

import (
	"context"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

func NewAutoWebhookStopProcessor(
	client mongo.DbClient,
) Processor {
	return &autoWebhookStopProcessor{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
	}
}

type autoWebhookStopProcessor struct {
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
}

func (p *autoWebhookStopProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	return processWebhookStop(ctx, event, p.alarmCollection, p.resolvedAlarmCollection,
		event.Parameters.Execution, types.AlarmChangeTypeAutoWebhookFail, types.AlarmChangeTypeAutoDeclareTicketWebhookFail)
}
