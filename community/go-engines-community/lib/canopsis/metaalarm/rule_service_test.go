package metaalarm_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/ruleapplicator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redisV8 "github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func testNewMetaAlarmService() (service.MetaAlarmService, entity.Adapter, alarm.Adapter, metaalarm.RulesAdapter, *redisV8.Client, redis.LockClient, mongo.DbClient, error) {
	logger := log.NewLogger(true)
	ctx := context.Background()

	redisClient, err := redis.NewSession(ctx, redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	redisClient2, err := redis.NewSession(ctx, redis.CorrelationLockStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	redisLockClient := redis.NewLockClient(redisClient2)

	client, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	alarmAdapter := alarm.NewAdapter(client)
	entityAdapter := entity.NewAdapter(client)

	dbClient, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rulesCollection := dbClient.Collection(mongo.MetaAlarmRulesMongoCollection)
	_, err = rulesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	err = testFillRulesCollection(rulesCollection)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	rulesAdapter := metaalarm.NewRuleAdapter(dbClient)

	s := service.NewMetaAlarmService(alarmAdapter, rulesAdapter,
		config.NewAlarmConfigProvider(config.CanopsisConf{}, logger), log.NewTestLogger())

	return s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, redisLockClient, dbClient, nil
}

func testFillRulesCollection(rulesCollection mongo.DbCollection) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	_, err := rulesCollection.InsertOne(ctx, rule)

	return err
}

func TestProcessAttributes(t *testing.T) {
	ctx := context.Background()

	s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, redlockClient, mongoSession, err := testNewMetaAlarmService()
	displayNameScheme, err := config.CreateDisplayNameTpl(config.AlarmDefaultNameScheme)
	if err != nil {
		panic(err)
	}
	c := config.AlarmConfig{
		FlappingFreqLimit:    0,
		FlappingInterval:     0,
		StealthyInterval:     0,
		BaggotTime:           time.Second,
		EnableLastEventDate:  true,
		CancelAutosolveDelay: time.Hour,
		DisplayNameScheme:    displayNameScheme,
		OutputLength:         10,
	}

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

		err = alarmAdapter.Insert(ctx, alarm)
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

		redisClient, err := redis.NewSession(ctx, redis.RuleTotalEntitiesStorage, logger, 0, 0)
		So(err, ShouldBeNil)

		ruleEntityCounter := metaalarm.NewRuleEntityCounter(entityAdapter, redisClient, logger)
		valueGroupEntityCounter := metaalarm.NewValueGroupEntityCounter(mongoSession, redisClient, logger)

		rs := metaalarm.NewRulesService(rulesAdapter, ruleEntityCounter, valueGroupEntityCounter, container, log.NewTestLogger())
		err = rs.LoadRules(ctx)
		So(err, ShouldBeNil)

		metaAlarms, err := rs.ProcessEvent(ctx, &testEvent)
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

		err = alarmAdapter.Insert(ctx, alarm)
		So(err, ShouldBeNil)

		testEvent.Alarm = &alarm
		testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

		metaAlarms, err = rs.ProcessEvent(ctx, &testEvent)
		So(err, ShouldBeNil)
		So(len(metaAlarms), ShouldEqual, 1)
	})
}
