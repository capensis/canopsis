package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewTriggerProcessor(
	client mongo.DbClient,
) Processor {
	return &triggerProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type triggerProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *triggerProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := bson.M{
		"d":          event.Entity.ID,
		"v.resolved": nil,
	}

	if event.Alarm != nil {
		match = bson.M{"_id": event.Alarm.ID}
	} else if event.AlarmID != "" {
		match = bson.M{"_id": event.AlarmID}
	}

	alarm := types.Alarm{}
	err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeType(event.Parameters.Trigger)
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
