package watcher_test

import (
	"context"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/encoding/json"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/mongo/bulk"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/streadway/amqp"
)

type TestPublisher struct {
	Decoder   encoding.Decoder
	TestEvent types.Event
}

func (t *TestPublisher) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	var tmpEvent types.Event
	err := t.Decoder.Decode(msg.Body, &tmpEvent)
	t.TestEvent = tmpEvent
	return err
}

// getServiceAdaptersPublisher returns the watcher service, the entity adapter, the alarm adapter and the test publisher
func getServiceAdaptersPublisher() (watcher.Service, entity.Adapter, alarm.Adapter, *TestPublisher, error) {
	mongo, err := mongo.NewSession(mongo.Timeout)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	redisClient, err := redis.NewSession(redis.CacheWatcher, log.NewTestLogger(), 0, 0)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Resetting collections
	watcherCollection := watcher.DefaultCollection(mongo)
	_, err = watcherCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	entityCollection := entity.DefaultCollection(mongo)
	_, err = entityCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	alarmCollection := alarm.DefaultCollection(mongo)
	_, err = alarmCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// FIXME: This removes all data from the redis cache, and is necessary
	// to make sure that running the tests twice leads to the same result.
	// This prevents from running multiple tests in parallel. A better
	// option might be to used https://github.com/alicebob/miniredis (but
	// this will only fix this problem for redis)
	err = redisClient.FlushDB().Err()
	So(err, ShouldBeNil)

	// Init watcher collection with two watchers
	validWatcherMap := getValidWatcherMap()
	invalidWatcherMap := getInvalidWatcherMap()
	watcherCollection.Insert(validWatcherMap)
	watcherCollection.Insert(invalidWatcherMap)

	// Creating adapters
	entityAdapter := entity.NewAdapter(entityCollection)
	watcherBulk := watcherCollection.NewBulk(bulk.BulkSizeMax)
	watcherAdapter := watcher.NewAdapter(
		watcherCollection,
		watcherBulk,
		entity.EntityCollectionName,
		alarm.AlarmCollectionName,
		log.NewTestLogger())
	alarmAdapter := alarm.NewAdapterLegacy(alarmCollection, entity.EntityCollectionName, alarmCollection.NewBulk(1000))

	countersCache := watcher.NewCountersCache(redisClient, log.NewTestLogger())

	// Creating watcher service and test publisher
	tPub := TestPublisher{Decoder: json.NewDecoder()}
	ws := watcher.NewService(redisClient, &tPub, "", "", json.NewEncoder(), watcherAdapter, alarmAdapter, countersCache, log.NewTestLogger())

	return ws, entityAdapter, alarmAdapter, &tPub, nil
}

func getValidWatcherMap() bson.M {
	mapWatcher := bson.M{
		"_id":     "watcher1",
		"name":    "watcher1",
		"type":    "watcher",
		"infos":   bson.M{},
		"enabled": true,
		"depends": []string{"mockEntity1/mockComponent"},
		"entities": []bson.M{
			bson.M{
				"name": "mockEntity1",
			},
		},
		"state":           bson.M{"method": "worst"},
		"output_template": "{{.State.Critical}}",
	}
	return mapWatcher
}

func getInvalidWatcherMap() bson.M {
	mapWatcher := bson.M{
		"_id":             "watcher2",
		"name":            "watcher2",
		"type":            "watcher",
		"infos":           bson.M{},
		"enabled":         false,
		"entities":        []bson.M{},
		"state":           bson.M{"method": "worst"},
		"output_template": "{{.State.Critical}}",
	}
	return mapWatcher
}

func getWatchedEntityMap() bson.M {
	mapEntity1 := bson.M{
		"_id":     "mockEntity1/mockComponent",
		"name":    "mockEntity1",
		"infos":   bson.M{},
		"state":   bson.M{},
		"enabled": true,
		"impact":  []string{"watcher1"},
		"depends": []string{},
	}
	return mapEntity1
}

func getWatchedEntity() types.Entity {
	var mockEntity1 types.Entity
	mapEntity1 := getWatchedEntityMap()

	bsonEntity1, err := bson.Marshal(mapEntity1)
	So(err, ShouldBeNil)
	So(bson.Unmarshal(bsonEntity1, &mockEntity1), ShouldBeNil)
	mockEntity1.EnsureInitialized()

	return mockEntity1
}

func getEventCritical() types.Event {
	now := types.CpsTime{Time: time.Now()}
	entity := getWatchedEntity()

	alarmChange := types.NewAlarmChange()
	alarmChange.Type = types.AlarmChangeTypeStateIncrease

	return types.Event{
		Component:     "mockComponent",
		Resource:      "mockEntity1",
		EventType:     types.EventTypeCheck,
		SourceType:    types.SourceTypeResource,
		Connector:     "convey",
		ConnectorName: "convey",
		State:         types.AlarmStateCritical,
		Output:        "pong",
		Timestamp:     now,
		Entity:        &entity,
		AlarmChange:   &alarmChange,
	}
}

// Testing

func TestProcess(t *testing.T) {
	Convey("Setup", t, func() {
		ws, ea, aa, tPub, err := getServiceAdaptersPublisher()
		So(err, ShouldBeNil)

		mockEntity := getWatchedEntity()
		err = ea.Insert(mockEntity)
		So(err, ShouldBeNil)

		ctx := context.Background()

		err = ws.ComputeAllWatchers(ctx)
		So(err, ShouldBeNil)

		Convey("Sending an event to process", func() {
			var testAlarm types.Alarm
			event := getEventCritical()
			// Adding an alarm related to the critical event
			testAlarm, err = types.NewAlarm(event, testutils.GetTestConf())
			So(err, ShouldBeNil)
			aa.Insert(testAlarm)
			event.Alarm = &testAlarm

			err = ws.Process(ctx, event)
			So(err, ShouldBeNil)

			Convey("Checking published event", func() {

				// Process sends an event to our test publisher
				So(tPub.TestEvent.State, ShouldEqual, event.State)
				So(tPub.TestEvent.Output, ShouldEqual, "1") // Because there's only one critical alarm opened

				Convey("Processing resolved alarms", func() {
					now := types.CpsTime{Time: time.Now()}
					testAlarm.Resolve(&now)

					ws.ProcessResolvedAlarm(ctx, testAlarm, mockEntity)

					Convey("Checking published event", func() {
						So(tPub.TestEvent.State, ShouldEqual, types.AlarmStateOK) // Because there's no alarm opened
						So(tPub.TestEvent.Output, ShouldEqual, "0")               // Because there's no alarm opened
					})
				})
			})
		})
	})
}
