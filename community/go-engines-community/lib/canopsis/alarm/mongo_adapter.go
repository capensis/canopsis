package alarm

import (
	"context"
	"errors"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/errt"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AlarmCollectionName = libmongo.AlarmMongoCollection
)

type mongoAdapter struct {
	dbClient     libmongo.DbClient
	dbCollection libmongo.DbCollection
}

func NewAdapter(dbClient libmongo.DbClient) Adapter {
	return &mongoAdapter{
		dbClient:     dbClient,
		dbCollection: dbClient.Collection(AlarmCollectionName),
	}
}

func (a mongoAdapter) Insert(alarm types.Alarm) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := a.dbCollection.InsertOne(ctx, alarm)
	if err != nil {
		return err
	}

	return nil
}

func (a mongoAdapter) Update(alarm types.Alarm) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.dbCollection.UpdateOne(ctx, bson.M{"_id": alarm.ID}, bson.M{"$set": alarm})
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return errors.New("alarm not modified")
	}

	return nil
}

func (a mongoAdapter) PartialUpdateOpen(ctx context.Context, alarm *types.Alarm) error {
	update := alarm.GetUpdate()
	if len(update) == 0 {
		return nil
	}

	res, err := a.dbCollection.UpdateOne(ctx, bson.M{
		"_id": alarm.ID,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}, update)
	if err != nil {
		return err
	}

	if res.MatchedCount > 0 && res.ModifiedCount > 0 {
		alarm.CleanUpdate()
	}

	return nil
}

func (a mongoAdapter) RemoveId(id string) error {
	panic("not implemented")
}

func (a mongoAdapter) RemoveAll() error {
	panic("not implemented")
}

func (a mongoAdapter) Get(filter map[string]interface{}, alarms *[]types.Alarm) error {
	panic("not implemented")
}

func (a mongoAdapter) GetAlarmsByID(id string) ([]types.Alarm, error) {
	return a.getAlarms(bson.M{"d": id})
}

func (a mongoAdapter) GetAlarmsWithCancelMark() ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.canceled": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithDoneMark() ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.done": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithSnoozeMark() ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.snooze": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmsWithFlappingStatus() ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.status.val": types.AlarmStatusFlapping,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAllOpenedResourceAlarmsByComponent(component string) ([]types.AlarmWithEntity, error) {
	req := bson.M{
		"v.component":  component,
		"v.resource":   bson.M{"$exists": true},
		"v.status.val": bson.M{"$ne": 0},
		"v.meta":       bson.M{"$exists": false},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}

	return a.getAlarmsWithEntity(req)
}

func (a mongoAdapter) GetUnacknowledgedAlarmsByComponent(component string) ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.component": component,
		"v.meta":      bson.M{"$exists": false},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
		"v.ack": nil,
	})
}

func (a mongoAdapter) GetAlarmsWithoutTicketByComponent(component string) ([]types.Alarm, error) {
	return a.getAlarms(bson.M{
		"v.component": component,
		"v.meta":      bson.M{"$exists": false},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
		"v.ticket": nil,
	})
}

func (a mongoAdapter) GetOpenedAlarmByAlarmId(id string) (types.Alarm, error) {
	return a.getAlarmWithErr(bson.M{
		"_id": id,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetAlarmByAlarmId(id string) (types.Alarm, error) {
	return a.getAlarmWithErr(bson.M{"_id": id})
}

func (a mongoAdapter) GetOpenedAlarm(connector, connectorName, id string) (types.Alarm, error) {
	return a.getAlarmWithErr(bson.M{
		"d":                id,
		"v.connector":      connector,
		"v.connector_name": connectorName,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	})
}

func (a mongoAdapter) GetOpenedMetaAlarm(ruleId string, valuePath string) (types.Alarm, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	al := types.Alarm{}
	query := bson.M{
		"v.meta": ruleId,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}

	if valuePath != "" {
		query["v.meta_value_path"] = valuePath
	}

	err := a.dbCollection.FindOne(ctx, query, options.FindOne().SetSort(bson.M{"v.creation_date": -1})).Decode(&al)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return al, errt.NewNotFound(err)
		}

		return al, err
	}

	al.Value.Transform()
	return al, nil
}

func (a mongoAdapter) GetLastAlarm(connector, connectorName, id string) (types.Alarm, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alarm := types.Alarm{}
	query := bson.M{
		"d":                id,
		"v.connector":      connector,
		"v.connector_name": connectorName,
	}
	err := a.dbCollection.FindOne(ctx, query, options.FindOne().SetSort(bson.M{"v.creation_date": -1})).Decode(&alarm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return alarm, errt.NewNotFound(err)
		}

		return alarm, err
	}

	alarm.Value.Transform()
	return alarm, nil
}

func (a mongoAdapter) GetUnresolved() ([]types.Alarm, error) {
	return a.getAlarms(bson.M{"v.resolved": nil})
}

// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
func (a mongoAdapter) GetOpenedAlarmsByIDs(ids []string, alarms *[]types.Alarm) error {
	var err error
	*alarms, err = a.getAlarms(bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	})

	return err
}

func (a mongoAdapter) GetOpenedAlarmsWithEntityByIDs(ids []string, alarms *[]types.AlarmWithEntity) error {
	filter := bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": bson.A{"", nil}},
	}

	var err error
	*alarms, err = a.getAlarmsWithEntity(filter)

	return err
}

func (a mongoAdapter) GetCountOpenedAlarmsByIDs(ids []string) (int64, error) {
	return a.getAlarmsCount(bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	})
}

// GetOpenedAlarmsByAlarmIDs gets ongoing alarms related the provided alarm ids
func (a mongoAdapter) GetOpenedAlarmsByAlarmIDs(ids []string, alarms *[]types.Alarm) error {
	var err error
	*alarms, err = a.getAlarms(bson.M{
		"_id":        bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": bson.A{"", nil}},
	})

	return err
}

func (a mongoAdapter) GetOpenedAlarmsWithEntityByAlarmIDs(ids []string, alarms *[]types.AlarmWithEntity) error {
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": bson.A{"", nil}},
	}

	var err error
	*alarms, err = a.getAlarmsWithEntity(filter)

	return err

}

func (a mongoAdapter) MassUpdate(alarms []types.Alarm, notUpdateResolved bool) error {
	if len(alarms) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	models := make([]mongo.WriteModel, len(alarms))

	for i, alarm := range alarms {
		filter := bson.M{"_id": alarm.ID}
		if notUpdateResolved {
			filter["v.resolved"] = bson.M{"$in": []interface{}{"", nil}}
		}

		models[i] = mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(bson.M{"$set": alarm})
	}

	_, err := a.dbCollection.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (a mongoAdapter) MassUpdateWithEntity(alarmsWithEntity []types.AlarmWithEntity) error {
	panic("not implemented")
}

func (a mongoAdapter) MassPartialUpdateOpen(ctx context.Context, updatedAlarm *types.Alarm, alarmID []string) error {
	update := updatedAlarm.GetUpdate()
	if len(update) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := a.dbCollection.UpdateMany(ctx, bson.M{
		"_id": bson.M{"$in": alarmID},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}, update)
	if err != nil {
		return err
	}

	if res.MatchedCount > 0 && res.ModifiedCount > 0 {
		updatedAlarm.CleanUpdate()
	}

	return nil
}

func (a mongoAdapter) GetOpenedAlarmsWithLastDatesBefore(
	ctx context.Context,
	time types.CpsTime,
) (libmongo.Cursor, error) {
	return a.dbCollection.Aggregate(ctx, []bson.M{
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
	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
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

func (a mongoAdapter) CountResolvedAlarm(entityIDs []string) (int, error) {
	res, err := a.getAlarmsCount(bson.M{
		"d":          bson.M{"$in": entityIDs},
		"v.resolved": bson.M{"$exists": true},
	})

	return int(res), err
}

func (a mongoAdapter) GetLastAlarmByEntityID(ctx context.Context, entityID string) (*types.Alarm, error) {
	cursor, err := a.dbCollection.Find(ctx, bson.M{"d": entityID}, options.Find().SetSort(bson.M{"t": -1}))
	if err != nil {
		return nil, err
	}
	if cursor.Next(ctx) {
		alarm := &types.Alarm{}
		err := cursor.Decode(alarm)
		if err != nil {
			return nil, err
		}

		return alarm, nil
	}

	return nil, cursor.Close(ctx)
}

func (a mongoAdapter) getAlarmsWithEntity(filter bson.M) ([]types.AlarmWithEntity, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Aggregate(ctx, []bson.M{
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

func (a mongoAdapter) getAlarms(filter bson.M) ([]types.Alarm, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := a.dbCollection.Find(ctx, filter)
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

func (a mongoAdapter) getAlarmsCount(filter bson.M) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return a.dbCollection.CountDocuments(ctx, filter)
}

func (a mongoAdapter) getAlarmWithErr(filter bson.M) (types.Alarm, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	alarm := types.Alarm{}
	err := a.dbCollection.FindOne(ctx, filter).Decode(&alarm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return alarm, errt.NewNotFound(err)
		}

		return alarm, err
	}

	alarm.Value.Transform()

	return alarm, nil
}
