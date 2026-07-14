package event

import (
	"context"
	"errors"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
	var result Result
	alarmChange := types.NewAlarmChange()

	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getExactAlarmMatchWithStepsLimit(event)
	if match == nil {
		return result, nil
	}

	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   bson.M{"$in": bson.A{types.AlarmStepWebhookComplete, types.AlarmStepWebhookFail}},
	}}}

	var update []bson.M

	stopOutputBuilder := strings.Builder{}
	stopOutputBuilder.WriteString(event.Parameters.Output)
	stopOutputBuilder.WriteString(". ")
	stopOutputBuilder.WriteString("Stopped by: ")
	stopOutputBuilder.WriteString(event.Parameters.Author)
	stopOutputBuilder.WriteString(".")

	if event.Parameters.TicketInfo.TicketRuleID == "" {
		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookStop, event.Parameters.Execution,
			stopOutputBuilder.String(), event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps":            addStepUpdateQuery(newStepQuery),
				"v.last_update_date": event.Parameters.Timestamp,
			}},
		}

		alarmChange.Type = types.AlarmChangeTypeAutoWebhookFail
	} else {
		outputBuilder := strings.Builder{}
		outputBuilder.WriteString(event.Parameters.Output)
		outputBuilder.WriteString(". ")
		outputBuilder.WriteString(types.OutputFailReasonPrefix)
		outputBuilder.WriteString("stopped by ")
		outputBuilder.WriteString(event.Parameters.Author)
		outputBuilder.WriteString(".")

		newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookStop, event.Parameters.Execution,
			stopOutputBuilder.String(), event.Parameters)
		newTicketStepQuery := ticketStepUpdateQueryWithInPbhInterval(types.AlarmStepDeclareTicketFail,
			event.Parameters.Execution, outputBuilder.String(), event.Parameters)
		update = []bson.M{
			{"$set": bson.M{
				"v.steps":            addStepUpdateQuery(newStepQuery, newTicketStepQuery),
				"v.failed_ticket":    newTicketStepQuery,
				"v.last_update_date": event.Parameters.Timestamp,
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

	if alarm.IsResolved() {
		_, err = p.resolvedAlarmCollection.UpdateOne(ctx, match, update)
		if err != nil {
			return result, err
		}
	}

	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
