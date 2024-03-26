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

func NewMetaAlarmChildDeactivateProcessor(
	client mongo.DbClient,
) Processor {
	return &metaAlarmChildDeactivateProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type metaAlarmChildDeactivateProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *metaAlarmChildDeactivateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["v.activation_date"] = nil
	match["inactive_delay_meta_alarm_in_progress"] = nil
	update := []bson.M{
		{"$set": bson.M{
			"inactive_delay_meta_alarm_in_progress": true,
			"v.inactive_start":                      event.Parameters.Timestamp,
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$cond": bson.M{
					"if": bson.M{"$gt": bson.A{"$v.inactive_start", 0}},
					"then": bson.M{"$subtract": bson.A{
						event.Parameters.Timestamp,
						"$v.inactive_start",
					}},
					"else": 0,
				}},
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
	alarmChange.Type = types.AlarmChangeTypeMetaAlarmChildDeactivate
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
