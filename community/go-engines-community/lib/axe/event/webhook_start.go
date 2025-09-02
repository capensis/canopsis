package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewWebhookStartProcessor(
	client mongo.DbClient,
) Processor {
	return &webhookStartProcessor{
		alarmCollection:         client.Collection(mongo.AlarmMongoCollection),
		resolvedAlarmCollection: client.Collection(mongo.ResolvedAlarmMongoCollection),
	}
}

type webhookStartProcessor struct {
	alarmCollection         mongo.DbCollection
	resolvedAlarmCollection mongo.DbCollection
}

func (p *webhookStartProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil || !event.Entity.Enabled {
		return result, nil
	}

	match := getExactAlarmMatchWithStepsLimit(event)
	if match == nil {
		return result, nil
	}

	match["v.steps"] = bson.M{"$not": bson.M{"$elemMatch": bson.M{
		"exec": event.Parameters.Execution,
		"_t":   types.AlarmStepWebhookStart,
	}}}
	newStepQuery := execStepUpdateQueryWithInPbhInterval(types.AlarmStepWebhookStart, event.Parameters.RuleExecution,
		event.Parameters.Output, event.Parameters)
	update := []bson.M{
		{"$set": bson.M{
			"v.steps":            addStepUpdateQuery(newStepQuery),
			"v.last_update_date": event.Parameters.Timestamp,
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

	if alarm.IsResolved() {
		_, err = p.resolvedAlarmCollection.UpdateOne(ctx, match, update)
		if err != nil {
			return result, err
		}
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeWebhookStart
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
