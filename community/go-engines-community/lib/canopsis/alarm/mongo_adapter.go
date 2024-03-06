package alarm

import (
	"context"
	"errors"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoAdapter struct {
	mainDbCollection     libmongo.DbCollection
	resolvedDbCollection libmongo.DbCollection
	archivedDbCollection libmongo.DbCollection
}

func NewAdapter(dbClient libmongo.DbClient) Adapter {
	return &mongoAdapter{
		mainDbCollection:     dbClient.Collection(libmongo.AlarmMongoCollection),
		resolvedDbCollection: dbClient.Collection(libmongo.ResolvedAlarmMongoCollection),
		archivedDbCollection: dbClient.Collection(libmongo.ArchivedAlarmMongoCollection),
	}
}

func (a mongoAdapter) GetAlarmsWithCancelMark(ctx context.Context) ([]types.AlarmWithEntity, error) {
	return a.getAlarmsWithEntity(ctx, bson.M{
		"v.canceled": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithSnoozeMark(ctx context.Context) ([]types.AlarmWithEntity, error) {
	return a.getAlarmsWithEntity(ctx, bson.M{
		"v.snooze": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithFlappingStatus(ctx context.Context) ([]types.AlarmWithEntity, error) {
	return a.getAlarmsWithEntity(ctx, bson.M{
		"v.status.val": types.AlarmStatusFlapping,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithoutTicketByComponent(ctx context.Context, component string) ([]types.AlarmWithEntity, error) {
	return a.getAlarmsWithEntity(ctx, bson.M{
		"v.component": component,
		"v.meta":      bson.M{"$exists": false},
		"v.resolved":  nil,
		"v.ticket":    nil,
	})
}

func (a mongoAdapter) GetOpenedAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error) {
	return a.getAlarmWithErr(ctx, bson.M{
		"_id": id,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmByAlarmId(ctx context.Context, id string) (types.Alarm, error) {
	return a.getAlarmWithErr(ctx, bson.M{
		"_id": id,
	})
}

// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
func (a mongoAdapter) GetOpenedAlarmsByIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error {
	var err error
	*alarms, err = a.getAlarms(ctx, bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	})

	return err
}

func (a mongoAdapter) GetOpenedAlarmsWithEntityByIDs(ctx context.Context, ids []string, alarms *[]types.AlarmWithEntity) error {
	filter := bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": bson.A{"", nil}},
	}

	var err error
	*alarms, err = a.getAlarmsWithEntity(ctx, filter)

	return err
}

func (a mongoAdapter) GetOpenedAlarmsWithEntity(ctx context.Context) (libmongo.Cursor, error) {
	filter := bson.M{
		"v.resolved": nil,
	}

	return a.entityAggregateCursor(ctx, filter)
}

func (a mongoAdapter) GetCountOpenedAlarmsByIDs(ctx context.Context, ids []string) (int64, error) {
	return a.getAlarmsCount(ctx, bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	})
}

// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
func (a mongoAdapter) GetOpenedAlarmsByAlarmIDs(ctx context.Context, ids []string, alarms *[]types.Alarm) error {
	var err error
	*alarms, err = a.getAlarms(ctx, bson.M{
		"_id":        bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": bson.A{"", nil}},
	})

	return err
}

func (a mongoAdapter) GetOpenedAlarmsWithLastDatesBefore(
	ctx context.Context,
	time datetime.CpsTime,
) (libmongo.Cursor, error) {
	return a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.status.val": bson.M{"$ne": types.AlarmStatusNoEvents},
			"v.resolved":   nil,
			"$or": []bson.M{
				{"v.last_update_date": bson.M{"$lte": time}},
				{"v.last_event_date": bson.M{"$lte": time}},
			},
		}},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{
			"entity.enabled":              true,
			"entity.last_idle_rule_apply": nil,
			"entity.type": bson.M{"$in": []string{
				types.EntityTypeConnector,
				types.EntityTypeComponent,
				types.EntityTypeResource,
			}},
		}},
	})
}

func (a mongoAdapter) GetOpenedAlarmsByConnectorIdleRules(ctx context.Context) ([]types.Alarm, error) {
	cursor, err := a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"v.state.val":  bson.M{"$ne": types.AlarmStateOK},
			"v.status.val": types.AlarmStatusNoEvents,
			"v.resolved":   nil,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
			"localField":   "d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{"entity.type": types.EntityTypeConnector}},
		{"$match": bson.M{"$expr": bson.M{
			"$gt": bson.A{"$entity.last_event_date", "$v.status.t"},
		}}},
	})

	if err != nil {
		return nil, err
	}

	var alarms []types.Alarm
	err = cursor.All(ctx, &alarms)
	if err != nil {
		return nil, err
	}

	return alarms, nil
}

func (a mongoAdapter) CountResolvedAlarm(ctx context.Context, entityIDs []string) (int, error) {
	res, err := a.getAlarmsCount(ctx, bson.M{
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": bson.M{"$exists": true},
	})

	return int(res), err
}

func (a mongoAdapter) GetLastAlarmByEntityID(ctx context.Context, entityID string) (*types.Alarm, error) {
	cursor, err := a.mainDbCollection.Find(ctx, bson.M{"d": entityID}, options.Find().SetSort(bson.M{"t": -1}))
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		alarm := &types.Alarm{}
		err := cursor.Decode(alarm)
		if err != nil {
			return nil, err
		}

		return alarm, nil
	}

	resolvedCursor, err := a.resolvedDbCollection.Find(ctx, bson.M{"d": entityID}, options.Find().SetSort(bson.M{"t": -1}))
	if err != nil {
		return nil, err
	}

	defer resolvedCursor.Close(ctx)

	if resolvedCursor.Next(ctx) {
		alarm := &types.Alarm{}
		err := resolvedCursor.Decode(alarm)
		if err != nil {
			return nil, err
		}

		return alarm, nil
	}

	return nil, nil
}

func (a mongoAdapter) entityAggregateCursor(ctx context.Context, filter bson.M) (libmongo.Cursor, error) {
	cursor, err := a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{"entity.enabled": true}},
	})

	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func (a mongoAdapter) getAlarmsWithEntity(ctx context.Context, filter bson.M) ([]types.AlarmWithEntity, error) {
	cursor, err := a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
	})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var alarmsWithEntity []types.AlarmWithEntity
	err = cursor.All(ctx, &alarmsWithEntity)
	if err != nil {
		return nil, err
	}

	return alarmsWithEntity, nil
}

func (a mongoAdapter) getAlarms(ctx context.Context, filter bson.M) ([]types.Alarm, error) {
	cursor, err := a.mainDbCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	alarms := make([]types.Alarm, 0)
	err = cursor.All(ctx, &alarms)
	if err != nil {
		return nil, err
	}

	return alarms, nil
}

func (a mongoAdapter) getAlarmsCount(ctx context.Context, filter bson.M) (int64, error) {
	return a.mainDbCollection.CountDocuments(ctx, filter)
}

func (a mongoAdapter) getAlarmWithErr(ctx context.Context, filter bson.M) (types.Alarm, error) {
	alarm := types.Alarm{}
	err := a.mainDbCollection.FindOne(ctx, filter).Decode(&alarm)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return alarm, err
		}

		return alarm, err
	}

	alarm.Value.Transform()

	return alarm, nil
}

func (a mongoAdapter) DeleteResolvedAlarms(ctx context.Context, duration time.Duration) error {
	_, err := a.mainDbCollection.DeleteMany(ctx, bson.M{
		"v.resolved": bson.M{"$lte": time.Now().Unix() - int64(duration.Seconds())},
	})

	return err
}

func (a *mongoAdapter) CopyAlarmToResolvedCollection(ctx context.Context, alarm types.Alarm) error {
	_, err := a.resolvedDbCollection.UpdateOne(
		ctx,
		bson.M{"_id": alarm.ID},
		bson.M{"$set": alarm},
		options.Update().SetUpsert(true),
	)

	return err
}

func (a *mongoAdapter) FindToCheckPbehaviorInfo(ctx context.Context, createdBefore datetime.CpsTime, idsWithPbehaviors []string) (libmongo.Cursor, error) {
	filter := bson.M{
		"v.resolved": nil,
		"t":          bson.M{"$lt": createdBefore},
	}

	if len(idsWithPbehaviors) > 0 {
		filter["$or"] = []bson.M{
			{"d": bson.M{"$in": idsWithPbehaviors}},
			{"v.pbehavior_info": bson.M{"$ne": nil}},
		}
	} else {
		filter["v.pbehavior_info"] = bson.M{"$ne": nil}
	}

	return a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": filter},
		{"$project": bson.M{
			"alarm": "$$ROOT",
			"_id":   0,
		}},
		{"$lookup": bson.M{
			"from":         libmongo.EntityMongoCollection,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		{"$unwind": "$entity"},
		{"$match": bson.M{"entity.enabled": true}},
	})
}

func (a *mongoAdapter) GetWorstAlarmStateAndMaxLastEventDate(ctx context.Context, entityIds []string) (int64, int64, error) {
	cursor, err := a.mainDbCollection.Aggregate(ctx, []bson.M{
		{"$match": bson.M{
			"d":          bson.M{"$in": entityIds},
			"v.resolved": nil,
		}},
		{"$group": bson.M{
			"_id":             nil,
			"state":           bson.M{"$max": "$v.state.val"},
			"last_event_date": bson.M{"$max": "$v.last_event_date"},
		}},
	})
	if err != nil {
		return 0, 0, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		res := struct {
			State         int64 `bson:"state"`
			LastEventDate int64 `bson:"last_event_date"`
		}{}

		err := cursor.Decode(&res)

		return res.State, res.LastEventDate, err
	}

	return 0, 0, nil
}

func (a *mongoAdapter) UpdateLastEventDate(ctx context.Context, entityIds []string, t datetime.CpsTime) error {
	_, err := a.mainDbCollection.UpdateMany(ctx, bson.M{
		"d":          bson.M{"$in": entityIds},
		"v.resolved": nil,
	}, bson.M{
		"$set": bson.M{"v.last_event_date": t},
	})

	return err
}
