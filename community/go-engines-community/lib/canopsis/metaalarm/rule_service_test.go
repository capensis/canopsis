package metaalarm_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm/ruleapplicator"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"

	"github.com/bsm/redislock"
	"github.com/globalsign/mgo"
	redisV7 "github.com/go-redis/redis/v7"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func testNewMetaAlarmService() (service.MetaAlarmService, entity.Adapter, alarm.Adapter, metaalarm.RulesAdapter, *redisV7.Client, *redislock.Client, *mgo.Session, error) {
	logger := log.NewLogger(true)

	redisClient, err := redis.NewSession(redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	redisClient2, err := redis.NewSession(redis.CorrelationLockStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	redisLockClient := redislock.New(redisClient2)

	session, err := mongo.NewSession(mongo.Timeout)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	alarmCollection := testAlarmDefaultCollection(session)
	_, err = alarmCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	alarmAdapter := alarm.NewAdapterLegacy(alarmCollection, entity.EntityCollectionName, alarmCollection.NewBulk(1000))
	entityAdapter := entity.NewAdapter(alarmCollection)

	rulesCollection := metaalarm.DefaultRulesCollection(session)
	_, err = rulesCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	err = testFillRulesCollection(rulesCollection)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	rulesAdapter := metaalarm.NewRuleAdapter(rulesCollection)

	s := service.NewMetaAlarmService(alarmAdapter, entityAdapter, rulesAdapter, testutils.GetTestConf(), log.NewTestLogger())

	return s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, redisLockClient, session, nil
}

func testAlarmDefaultCollection(session *mgo.Session) mongo.Collection {
	collection := session.DB(canopsis.DbName).C("test_periodical_alarm")
	return mongo.FromMgo(collection)
}

func testFillRulesCollection(rulesCollection mongo.Collection) error {
	rule := bson.M{
		"_id":  "testRule",
		"name": "Test",
		"type": "attribute",
		"config": bson.M{
			"alarm_patterns": []bson.M{
				{"v": bson.M{
					"connector_name": bson.M{
						"regex_match": "test31$*",
					},
				}},
			},
		},
	}

	return rulesCollection.Insert(rule)
}

func TestProcessAttributes(t *testing.T) {
	s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, redlockClient, mongoSession, err := testNewMetaAlarmService()
	c := testutils.GetTestConf()

	Convey("Create Process Attributes test", t, func() {
		So(err, ShouldBeNil)

		now := types.CpsTime{Time: time.Now()}
		testEvent := types.Event{
			Component:     "testComponent",
			Resource:      "testResource",
			EventType:     types.EventTypeCheck,
			SourceType:    types.SourceTypeResource,
			Connector:     "test",
			ConnectorName: "test",
			State:         types.AlarmStateCritical,
			Output:        "",
			Timestamp:     now,
		}

		alarm, err := types.NewAlarm(testEvent, c)
		entity := types.Entity{ID: testEvent.GetEID()}
		So(err, ShouldBeNil)

		err = alarmAdapter.Insert(alarm)
		So(err, ShouldBeNil)

		testEvent.Alarm = &alarm

		rule := metaalarm.Rule{}
		err = json.Unmarshal([]byte(`{
			"_id": "testRule",
			"name": "Test",
			"type": "attribute"
		}`), &rule)
		So(err, ShouldBeNil)

		_, err = s.CreateMetaAlarm(&testEvent, []types.AlarmWithEntity{{
			Alarm:  alarm,
			Entity: entity,
		}}, rule)
		So(err, ShouldBeNil)

		container := metaalarm.NewRuleApplicatorContainer()
		attributeApplicator := ruleapplicator.NewAttributeApplicator(alarmAdapter, log.NewTestLogger(), s, redisClient, redlockClient)
		container.Set(metaalarm.RuleTypeAttribute, attributeApplicator)

		logger := log.NewLogger(true)

		redisClient, err := redis.NewSession(redis.RuleTotalEntitiesStorage, logger, 0, 0)
		So(err, ShouldBeNil)

		ruleEntityCounter := metaalarm.NewRuleEntityCounter(entityAdapter, redisClient, logger)
		valueGroupEntityCounter := metaalarm.NewValueGroupEntityCounter(mongoSession, redisClient, logger)

		rs := metaalarm.NewRulesService(rulesAdapter, ruleEntityCounter, valueGroupEntityCounter, container, log.NewTestLogger())
		err = rs.LoadRules()
		So(err, ShouldBeNil)

		metaAlarms, err := rs.ProcessEvent(&testEvent)
		So(err, ShouldBeNil)
		So(len(metaAlarms), ShouldEqual, 0)

		testEvent = types.Event{
			Component:     "testComponent",
			Resource:      "testResource",
			EventType:     types.EventTypeCheck,
			SourceType:    types.SourceTypeResource,
			Connector:     "test31",
			ConnectorName: "test31",
			State:         types.AlarmStateMajor,
			Output:        "",
			Timestamp:     now,
		}

		alarm, err = types.NewAlarm(testEvent, c)
		So(err, ShouldBeNil)

		err = alarmAdapter.Insert(alarm)
		So(err, ShouldBeNil)

		testEvent.Alarm = &alarm
		testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

		metaAlarms, err = rs.ProcessEvent(&testEvent)
		So(err, ShouldBeNil)
		So(len(metaAlarms), ShouldEqual, 1)
	})
}
