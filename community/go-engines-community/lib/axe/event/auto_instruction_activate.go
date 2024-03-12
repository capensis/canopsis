package event

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewAutoInstructionActivateProcessor(
	client mongo.DbClient,
) Processor {
	return &autoInstructionActivateProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type autoInstructionActivateProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *autoInstructionActivateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.activation_date"] = nil
	match["auto_instruction_in_progress"] = true
	update := []bson.M{
		{"$unset": "auto_instruction_in_progress"},
		{"$set": bson.M{
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$subtract": bson.A{
					event.Parameters.Timestamp,
					"$v.inactive_start",
				}},
			}},
			"v.inactive_start": bson.M{"$cond": bson.M{
				"if": bson.M{"$and": []bson.M{
					{"$eq": bson.A{"$v.snooze", nil}},
					{"$in": bson.A{"$v.pbehavior_info", bson.A{nil, "", pbehavior.TypeActive}}},
					{"$ne": bson.A{"$inactive_delay_meta_alarm_in_progress", true}},
				}},
				"then": nil,
				"else": event.Parameters.Timestamp,
			}},
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
	alarmChange.Type = types.AlarmChangeTypeAutoInstructionActivate
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
