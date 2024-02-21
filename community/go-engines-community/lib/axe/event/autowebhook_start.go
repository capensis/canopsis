package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAutoWebhookStartProcessor(
	client mongo.DbClient,
) Processor {
	return &autoWebhookStartProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type autoWebhookStartProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *autoWebhookStartProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   types.AlarmStepWebhookStart,
	}}}
	newStep := types.NewAlarmStep(types.AlarmStepWebhookStart, event.Parameters.Timestamp, event.Parameters.Author,
		event.Parameters.Output, event.Parameters.User, event.Parameters.Role, event.Parameters.Initiator, false)
	newStep.Execution = event.Parameters.Execution
	newStepQuery := stepUpdateQuery(newStep)
	update := []bson.M{
		{"$set": bson.M{
			"v.steps": addStepUpdateQuery(newStepQuery),
		}},
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

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeAutoWebhookStart
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
