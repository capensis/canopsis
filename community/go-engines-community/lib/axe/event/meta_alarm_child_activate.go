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

func NewMetaAlarmChildActivateProcessor(
	client mongo.DbClient,
) Processor {
	return &metaAlarmChildActivateProcessor{
		alarmCollection: client.Collection(mongo.AlarmMongoCollection),
	}
}

type metaAlarmChildActivateProcessor struct {
	alarmCollection mongo.DbCollection
}

func (p *metaAlarmChildActivateProcessor) Process(ctx context.Context, event rpc.AxeEvent) (Result, error) {
	result := Result{}
	if event.Entity == nil {
		return result, nil
	}

	match := getOpenAlarmMatch(event)
	match["inactive_delay_meta_alarm_in_progress"] = true
	update := []bson.M{
		{"$unset": bson.A{
			"inactive_delay_meta_alarm_in_progress",
			"meta_alarm_inactive_delay",
		}},
		{"$set": bson.M{
			"v.inactive_duration": bson.M{"$sum": bson.A{
				"$v.inactive_duration",
				bson.M{"$subtract": bson.A{
					event.Parameters.Timestamp,
					"$v.inactive_start",
				}},
			}},
			"v.inactive_start": updateInactiveStart(event.Parameters.Timestamp, true, true, true, false),
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
	alarmChange.Type = types.AlarmChangeTypeMetaAlarmChildActivate
	result.Forward = true
	result.Alarm = alarm
	result.AlarmChange = alarmChange

	return result, nil
}
