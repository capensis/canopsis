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

func NewWebhookFailProcessor(
	client mongo.DbClient,
) Processor {
	return &webhookFailProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type webhookFailProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *webhookFailProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
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
	var update bson.M
	outputBuilder := strings.Builder{}
	outputBuilder.WriteString(event.Parameters.Output)
	if event.Parameters.WebhookFailReason != "" {
		outputBuilder.WriteString(". Fail reason: ")
		outputBuilder.WriteString(event.Parameters.WebhookFailReason)
		outputBuilder.WriteRune('.')
	}

	if event.Parameters.TicketInfo.TicketRuleID == "" {
		newStep := types.NewAlarmStep(types.AlarmStepWebhookFail, event.Parameters.Timestamp, event.Parameters.Author,
			outputBuilder.String(), event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
		newStep.Execution = event.Parameters.Execution
		update = bson.M{
			"$push": bson.M{
				"v.steps": newStep,
			},
		}
		alarmChange.Type = types.AlarmChangeTypeWebhookFail
	} else {
		ticketOutput := outputBuilder.String()
		requestOutput := ticketOutput
		stepType := types.AlarmStepWebhookFail
		if event.Parameters.WebhookRequest {
			requestOutput = event.Parameters.Output
			stepType = types.AlarmStepWebhookComplete
		}

		newStep := types.NewAlarmStep(stepType, event.Parameters.Timestamp, event.Parameters.Author, requestOutput,
			event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator)
		newStep.Execution = event.Parameters.Execution
		newTicketStep := types.NewTicketStep(types.AlarmStepDeclareTicketFail, event.Parameters.Timestamp, event.Parameters.Author,
			ticketOutput, event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, event.Parameters.TicketInfo)
		newTicketStep.Execution = event.Parameters.Execution
		update = bson.M{
			"$push": bson.M{
				"v.tickets": newTicketStep,
				"v.steps":   bson.M{"$each": bson.A{newStep, newTicketStep}},
			},
		}
		alarmChange.Type = types.AlarmChangeTypeDeclareTicketWebhookFail
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
