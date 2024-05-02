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

func NewJunitProcessor(
	client mongo.DbClient,
) Processor {
	return &junitProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
		alarmStepTypeMap: map[string]string{
			types.EventTypeJunitTestSuiteUpdated: types.AlarmStepJunitTestSuiteUpdate,
			types.EventTypeJunitTestCaseUpdated:  types.AlarmStepJunitTestCaseUpdate,
		},
		alarmChangeTypeMap: map[string]types.AlarmChangeType{
			types.EventTypeJunitTestSuiteUpdated: types.AlarmChangeTypeJunitTestSuiteUpdate,
			types.EventTypeJunitTestCaseUpdated:  types.AlarmChangeTypeJunitTestCaseUpdate,
		},
	}
}

type junitProcessor struct {
	alarmCollection    mongo.DbCollection
	alarmStepTypeMap   map[string]string
	alarmChangeTypeMap map[string]types.AlarmChangeType
}

func (p *junitProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	alarmStepType, ok := p.alarmStepTypeMap[event.EventType]
	if !ok {
		return result, nil
	}

	alarmChangeType, ok := p.alarmChangeTypeMap[event.EventType]
	if !ok {
		return result, nil
	}

	match := getOpenAlarmMatchWithStepsLimit(event)
	newStepQuery := stepUpdateQueryWithInPbhInterval(alarmStepType, event.Parameters.Output, event.Parameters)
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
	alarmChange.Type = alarmChangeType
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
