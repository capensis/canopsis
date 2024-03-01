package event

import (
	"context"
	"errors"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAutoWebhookFailProcessor(
	client mongo.DbClient,
) Processor {
	return &autoWebhookFailProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type autoWebhookFailProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *autoWebhookFailProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   bson.M{"$in": bson.A{types.AlarmStepWebhookComplete, types.AlarmStepWebhookFail}},
	}}}
	alarmChange := types.NewAlarmChange()
	var update []bson.M
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(event.Parameters.Output)
	if event.Parameters.WebhookFailReason != "" {
		outputBuilder.WriteString(". Fail reason: ")
		outputBuilder.WriteString(event.Parameters.WebhookFailReason)
		outputBuilder.WriteRune('.')
	}

	if event.Parameters.TicketInfo.TicketRuleID == "" {
		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookFail, event.Parameters.Execution,
			outputBuilder.String(), event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps": addStepUpdateQuery(newStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeAutoWebhookFail
	} else {
		ticketOutput := outputBuilder.String()
		requestOutput := ticketOutput
		stepType := types.AlarmStepWebhookFail
		if event.Parameters.WebhookRequest {
			requestOutput = event.Parameters.Output
			stepType = types.AlarmStepWebhookComplete
		}

		newStepQuery := execStepUpdateQueryWithInPbhInterval(stepType, event.Parameters.Execution, requestOutput, event.Parameters)
		newTicketStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepDeclareTicketFail,
			event.Parameters.Execution, ticketOutput, event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.tickets": addTicketUpdateQuery(newTicketStepQuery),
				"v.steps":   addStepUpdateQuery(newStepQuery, newTicketStepQuery),
			}},
		}
		alarmChange.Type = types.AlarmChangeTypeAutoDeclareTicketWebhookFail
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	alarm := types.Alarm{}
	err := p.alarmCollection.FindOneAndUpdate(ctx, match, update, opts).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
