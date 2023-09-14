package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func NewForwardWithAlarmProcessor(
	client mongo.DbClient,
) Processor {
	return &forwardWithAlarmProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type forwardWithAlarmProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *forwardWithAlarmProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	alarm := types.Alarm{}
	err := p.alarmCollection.FindOne(ctx, match).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			return result, nil
		}

		return result, err
	}

	alarmChange := types.NewAlarmChange()
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
