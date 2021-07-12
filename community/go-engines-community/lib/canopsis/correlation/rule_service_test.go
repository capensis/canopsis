package correlation_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/ruleapplicator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	redisV8 "github.com/go-redis/redis/v8"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func testNewMetaAlarmService() (service.MetaAlarmService, entity.Adapter, alarm.Adapter, correlation.RulesAdapter, *redisV8.Client, mongo.DbClient, error) {
	logger := log.NewLogger(true)
	ctx := context.Background()

	redisClient, err := redis.NewSession(ctx, redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	client, err := mongo.NewClient(ctx, 0, 0)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
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
		return nil, nil, nil, nil, nil, nil, err
	}

	err = testFillRulesCollection(rulesCollection)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, err
	}

	rulesAdapter := correlation.NewRuleAdapter(dbClient)

	s := service.NewMetaAlarmService(alarmAdapter, config.NewAlarmConfigProvider(config.CanopsisConf{}, logger), log.NewTestLogger())

	return s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, dbClient, nil
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

	s, entityAdapter, alarmAdapter, rulesAdapter, redisClient, mongoSession, err := testNewMetaAlarmService()
	if err != nil {
		panic(err)
	}

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

		rule := correlation.Rule{}
		err = json.Unmarshal([]byte(`{
			"_id": "testRule",
			"name": "Test",
			"type": "attribute"
		}`), &rule)
		So(err, ShouldBeNil)

		ev, err := s.CreateMetaAlarm(testEvent, []types.AlarmWithEntity{{
			Alarm:  alarm,
			Entity: entity,
		}}, rule)
		So(err, ShouldBeNil)

		fmt.Printf("%v\n", ev)

		container := correlation.NewRuleApplicatorContainer()
		attributeApplicator := ruleapplicator.NewAttributeApplicator(alarmAdapter, log.NewTestLogger(), s, redisClient)
		container.Set(correlation.RuleTypeAttribute, attributeApplicator)

		logger := log.NewLogger(true)

		redisClient, err := redis.NewSession(ctx, redis.RuleTotalEntitiesStorage, logger, 0, 0)
		So(err, ShouldBeNil)

		ruleEntityCounter := correlation.NewRuleEntityCounter(entityAdapter, redisClient, logger)
		valueGroupEntityCounter := correlation.NewValueGroupEntityCounter(mongoSession, redisClient, logger)

		rs := correlation.NewRulesService(rulesAdapter, ruleEntityCounter, valueGroupEntityCounter, container, log.NewTestLogger())
		err = rs.LoadRules(ctx)
		So(err, ShouldBeNil)

		metaAlarms, err := rs.ProcessEvent(ctx, testEvent)
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

		metaAlarms, err = rs.ProcessEvent(ctx, testEvent)
		So(err, ShouldBeNil)
		So(len(metaAlarms), ShouldEqual, 1)
	})
}
