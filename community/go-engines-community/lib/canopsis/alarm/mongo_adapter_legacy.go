package alarm

import (
	"fmt"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// Adapter allows you to manipulate alarms in database
type mongoAdapterLegacy struct {
	collection           mongo.Collection
	entityCollectionName string
	mbulk                bulk.Bulk
}

// NewAdapterLegacy gives the correct mongo alarm adapter.
func NewAdapterLegacy(collection mongo.Collection, entityCollectionName string, safeBulk bulk.Bulk) Adapter {
	return mongoAdapterLegacy{
		collection:           collection,
		entityCollectionName: entityCollectionName,
		mbulk:                safeBulk,
	}
}

func DefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C(AlarmCollectionName)
	return mongo.FromMgo(collection)
}

func (a mongoAdapterLegacy) Insert(alarm types.Alarm) error {
	return a.collection.Insert(alarm)
}

func (a mongoAdapterLegacy) Update(alarm types.Alarm) error {
	return a.collection.Update(alarm.ID, alarm)
}

func (a mongoAdapterLegacy) PartialUpdateOpen(alarm *types.Alarm) error {
	panic("not implemented")
}

func (a mongoAdapterLegacy) RemoveId(id string) error {
	return a.collection.Remove(id)
}

func (a mongoAdapterLegacy) RemoveAll() error {
	_, err := a.collection.RemoveAll()
	return err
}

func (a mongoAdapterLegacy) Get(filter map[string]interface{}, alarms *[]types.Alarm) error {
	err := a.collection.Get(filter, alarms)
	if err != nil {
		return err
	}
	for i, al := range *alarms {
		al.Value.Transform()
		(*alarms)[i] = al
	}
	return err
}

func (a mongoAdapterLegacy) GetAlarmsByID(id string) ([]types.Alarm, error) {
	alarms := []types.Alarm{}
	err := a.Get(bson.M{"d": id}, &alarms)
	return alarms, err
}

func (a mongoAdapterLegacy) GetAlarmsWithCancelMark() ([]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	filter := bson.M{
		"v.canceled": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.Get(filter, &alarms)
	return alarms, err
}

func (a mongoAdapterLegacy) GetAlarmsWithDoneMark() ([]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	filter := bson.M{
		"v.done": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.Get(filter, &alarms)
	return alarms, err
}

func (a mongoAdapterLegacy) GetAlarmsWithSnoozeMark() ([]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	filter := bson.M{
		"v.snooze": bson.M{"$ne": nil},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.Get(filter, &alarms)
	return alarms, err
}

func (a mongoAdapterLegacy) GetAlarmsWithFlappingStatus() ([]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	filter := bson.M{
		"v.status.val": types.AlarmStatusFlapping,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.Get(filter, &alarms)
	return alarms, err
}

func (a mongoAdapterLegacy) GetAllOpenedResourceAlarmsByComponent(component string) ([]types.AlarmWithEntity, error) {
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

	alarms := make([]types.AlarmWithEntity, 0)
	err := a.getAlarmsWithEntity(req, &alarms)

	return alarms, err
}

func (a mongoAdapterLegacy) GetUnacknowledgedAlarmsByComponent(component string) ([]types.Alarm, error) {
	req := bson.M{
		"v.component": component,
		"v.meta":      bson.M{"$exists": false},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
		"v.ack": nil,
	}

	alarms := make([]types.Alarm, 0)
	err := a.Get(req, &alarms)

	return alarms, err
}

func (a mongoAdapterLegacy) GetAlarmsWithoutTicketByComponent(component string) ([]types.Alarm, error) {
	req := bson.M{
		"v.component": component,
		"v.meta":      bson.M{"$exists": false},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
		"v.ticket": nil,
	}

	alarms := make([]types.Alarm, 0)
	err := a.Get(req, &alarms)

	return alarms, err
}

func (a mongoAdapterLegacy) GetOpenedAlarmByAlarmId(id string) (types.Alarm, error) {
	al := types.Alarm{}
	query := bson.M{
		"_id": id,
		"$or": []bson.M{
			bson.M{"v.resolved": nil},
			bson.M{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.collection.GetOne(query, nil, &al)

	return al, err
}

func (a mongoAdapterLegacy) GetOpenedAlarm(connector, connectorName, id string) (types.Alarm, error) {
	al := types.Alarm{}
	query := bson.M{
		"d":                id,
		"v.connector":      connector,
		"v.connector_name": connectorName,
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	err := a.collection.GetOne(query, nil, &al)
	if err == nil {
		al.Value.Transform()
	}

	return al, err
}

func (a mongoAdapterLegacy) GetOpenedMetaAlarm(ruleId string, valuePath string) (types.Alarm, error) {
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

	err := a.collection.GetFirst(query, "-v.creation_date", &al)
	if err == nil {
		al.Value.Transform()
	}

	return al, err
}

func (a mongoAdapterLegacy) GetLastAlarm(connector, connectorName, id string) (types.Alarm, error) {
	al := types.Alarm{}
	query := bson.M{
		"d":                id,
		"v.connector":      connector,
		"v.connector_name": connectorName,
	}
	err := a.collection.GetFirst(query, "-v.creation_date", &al)
	if err == nil {
		al.Value.Transform()
	}

	return al, err
}

func (a mongoAdapterLegacy) GetUnresolved() ([]types.Alarm, error) {
	alarms := make([]types.Alarm, 0)
	filter := bson.M{"v.resolved": nil}
	err := a.Get(filter, &alarms)
	return alarms, err
}

// getAlarmsWithEntity gets the alarms matching the filter given in parameter and adds the entity associated to each alarm
// It writes in the AlarmWithEntity array the result of the mongo aggregation
func (a mongoAdapterLegacy) getAlarmsWithEntity(filter bson.M, alarmsWithEntity *[]types.AlarmWithEntity) error {
	pipeline := []bson.M{
		// Filters the alarms according the given parameter
		bson.M{"$match": filter},
		// Defines alarm the root document, i.e. the top-level document, in the aggregation pipeline
		bson.M{"$project": bson.M{
			"alarm": "$$ROOT",
		}},
		// JOIN between periodical_alarm and default_entities based on entities id
		bson.M{"$lookup": bson.M{
			"from":         a.entityCollectionName,
			"localField":   "alarm.d",
			"foreignField": "_id",
			"as":           "entity",
		}},
		// Transforms the array of entity to a document
		// Also removes the alarms without the entities
		bson.M{"$unwind": "$entity"},
	}

	return a.collection.Aggregate(pipeline, alarmsWithEntity)
}

// GetOpenedAlarmsByIDs gets ongoing alarms related the provided entity ids
func (a mongoAdapterLegacy) GetOpenedAlarmsByIDs(ids []string, alarms *[]types.Alarm) error {
	filter := bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	}
	return a.Get(filter, alarms)
}

func (a mongoAdapterLegacy) GetOpenedAlarmsWithEntityByIDs(ids []string, alarms *[]types.AlarmWithEntity) error {
	filter := bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	}
	return a.getAlarmsWithEntity(filter, alarms)
}

func (a mongoAdapterLegacy) GetCountOpenedAlarmsByIDs(ids []string) (int, error) {
	filter := bson.M{
		"d":          bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	}

	return a.collection.Count(filter)
}

// GetOpenedAlarmsByIDs gets ongoing alarms related the provided alarm ids
func (a mongoAdapterLegacy) GetOpenedAlarmsByAlarmIDs(ids []string, alarms *[]types.Alarm) error {
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"v.resolved": bson.M{"$in": []interface{}{"", nil}},
	}
	return a.Get(filter, alarms)
}

func (a mongoAdapterLegacy) GetOpenedAlarmsWithEntityByAlarmIDs(ids []string, alarms *[]types.AlarmWithEntity) error {
	filter := bson.M{
		"_id":        bson.M{"$in": ids},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}
	return a.getAlarmsWithEntity(filter, alarms)
}

func (a mongoAdapterLegacy) MassUpdate(alarms []types.Alarm, notUpdateResolved bool) error {
	for _, alarm := range alarms {
		var op bulk.OpUpdate
		if notUpdateResolved {
			op = bulk.NewUpdate(bson.M{
				"_id":        alarm.ID,
				"v.resolved": bson.M{"$in": []interface{}{"", nil}},
			}, alarm)
		} else {
			op = bulk.NewUpdate(bson.M{"_id": alarm.ID}, alarm)
		}
		_, err := a.mbulk.AddUpdate(op)
		if err != nil {
			return fmt.Errorf("alarm.Adapter.MassUpdate: %v", err)
		}
	}

	_, err := a.mbulk.PerformUpdates()
	if err != nil {
		return fmt.Errorf("alarm.Adapter.MussUpdate - flush: %v", err)
	}

	return nil
}

func (a mongoAdapterLegacy) MassUpdateWithEntity(alarmsWithEntity []types.AlarmWithEntity) error {
	alarms := make([]types.Alarm, 0)
	for _, alarmWithEntity := range alarmsWithEntity {
		alarms = append(alarms, alarmWithEntity.Alarm)
	}
	return a.MassUpdate(alarms, false)
}

func (a mongoAdapterLegacy) GetOpenedAlarmsWithLastEventDateBefore(
	date time.Time,
) ([]types.AlarmWithEntity, error) {
	filter := bson.M{
		"v.last_event_date": bson.M{"$lte": date.Unix()},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}

	alarms := make([]types.AlarmWithEntity, 0)
	err := a.getAlarmsWithEntity(filter, &alarms)

	return alarms, err
}

func (a mongoAdapterLegacy) GetOpenedAlarmsWithLastUpdateDateBefore(
	date time.Time,
) ([]types.AlarmWithEntity, error) {
	filter := bson.M{
		"v.last_update_date": bson.M{"$lte": date.Unix()},
		"$or": []bson.M{
			{"v.resolved": nil},
			{"v.resolved": bson.M{"$exists": false}},
		},
	}

	alarms := make([]types.AlarmWithEntity, 0)
	err := a.getAlarmsWithEntity(filter, &alarms)

	return alarms, err
}

func (a mongoAdapterLegacy) CountResolvedAlarm(alarmList []string) (int, error) {
	filter := bson.M{
		"d":          bson.M{"$in": alarmList},
		"v.resolved": bson.M{"$exists": true},
	}
	alarms := make([]types.Alarm, 0)
	err := a.collection.Get(filter, &alarms)
	if err != nil {
		return 0, err
	}
	return len(alarms), nil
}
