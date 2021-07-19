package ruleapplicator_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/testutils"
	"go.mongodb.org/mongo-driver/bson"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/ruleapplicator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewComplexApplicator() (*ruleapplicator.ComplexApplicator, alarm.Adapter, entity.Adapter, error) {
	logger := log.NewLogger(true)
	ctx := context.Background()

	redisClient, err := redis.NewSession(ctx, redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}
	redisClient.FlushDB(ctx)

	redisClient2, err := redis.NewSession(ctx, redis.CorrelationLockStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}
	redisClient2.FlushDB(ctx)

	redisLockClient := redis.NewLockClient(redisClient2)

	dbClient, err := libmongo.NewClient(ctx, 0, 0)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	entityCollection := dbClient.Collection(libmongo.EntityMongoCollection)
	_, err = entityCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, err
	}

	alarmAdapter := alarm.NewAdapter(dbClient)
	entityAdapter := entity.NewAdapter(dbClient)

	rulesCollection := dbClient.Collection(libmongo.MetaAlarmRulesMongoCollection)
	_, err = rulesCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return nil, nil, nil, err
	}

	rulesAdapter := metaalarm.NewRuleAdapter(dbClient)

	metaAlarmService := service.NewMetaAlarmService(alarmAdapter, rulesAdapter,
		config.NewAlarmConfigProvider(testutils.GetTestConf(), logger), log.NewLogger(true))

	redisClient3, err := redis.NewSession(ctx, redis.RuleTotalEntitiesStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	ruleEntityCounter := metaalarm.NewRuleEntityCounter(entityAdapter, redisClient3, logger)

	applicator := ruleapplicator.NewComplexApplicator(alarmAdapter, metaAlarmService, redisClient, redisLockClient, ruleEntityCounter, logger)
	return &applicator, alarmAdapter, entityAdapter, nil
}

type alarmPatternListWrapper struct {
	PatternList pattern.AlarmPatternList `bson:"list"`
}

func TestRuleMatch(t *testing.T) {
	ctx := context.Background()
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
	applicator, alarmAdapter, entityAdapter, err := testNewComplexApplicator()

	Convey("Create RuleMatch test", t, func() {
		So(err, ShouldBeNil)

		testEvent := types.Event{}
		rule := metaalarm.Rule{}

		Convey("empty event and rule", func() {
			metaAlarmEvents, _ := applicator.Apply(ctx, &testEvent, rule)
			So(len(metaAlarmEvents), ShouldEqual, 0)
		})

		Convey("empty event complex rule", func() {
			Convey("Init rule", func() {
				err = json.Unmarshal([]byte(`{
					"_id": "testRule",
					"name": "Test",
					"type": "complex",
					"config": {
						"time_interval": 3,
						"threshold_count": 1
					}
				}`), &rule)
				So(err, ShouldBeNil)
				So(rule.Config.TimeInterval, ShouldEqual, 3)

				// So(err, ShouldBeNil)
				// rule.Config.AttributePatterns = w.PatternList

				mapPattern := bson.M{
					"list": []bson.M{
						{
							"v": bson.M{
								"connector_name": "test899",
							},
						},
					},
				}
				bsonPattern, err := bson.Marshal(mapPattern)
				So(err, ShouldBeNil)

				var w alarmPatternListWrapper
				err = bson.Unmarshal(bsonPattern, &w)
				So(err, ShouldBeNil)

				rule.Config.AlarmPatterns = w.PatternList

				So(fmt.Sprintf("%+v", rule.Config.AlarmPatterns.AsMongoDriverQuery()["$or"]), ShouldEqual, "[map[v.connector_name:map[$eq:test899]]]")

				Convey("match event complex rule", func() {

					now := types.CpsTime{Time: time.Now()}
					testEvent = types.Event{
						Component:     "testComponent",
						Resource:      "testResource1",
						EventType:     types.EventTypeCheck,
						SourceType:    types.SourceTypeResource,
						Connector:     "test",
						ConnectorName: "test899",
						State:         types.AlarmStateCritical,
						Output:        "",
						Timestamp:     now,
					}

					entity1 := types.NewEntity("testResource1/testComponent", "testResource", "resource", nil, nil, nil)
					err = entityAdapter.Insert(ctx, entity1)
					So(err, ShouldBeNil)

					alarm1, err := types.NewAlarm(testEvent, c)
					So(err, ShouldBeNil)

					err = alarmAdapter.Insert(ctx, alarm1)
					So(err, ShouldBeNil)

					testEvent.Alarm = &alarm1
					testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

					Convey("Test RuleMatched", func() {
						metaAlarmEvents, _ := applicator.Apply(ctx, &testEvent, rule)
						So(len(metaAlarmEvents), ShouldEqual, 1)
						So(metaAlarmEvents[0].EventType, ShouldEqual, "metaalarm")
					})
				})
				mapPattern = bson.M{
					"list": []bson.M{
						{
							"v": bson.M{
								"connector_name": bson.M{
									"regex_match": "tesk.$",
								},
							},
						},
					},
				}
				bsonPattern, err = bson.Marshal(mapPattern)
				So(err, ShouldBeNil)

				var w2 alarmPatternListWrapper
				err = bson.Unmarshal(bsonPattern, &w2)
				So(err, ShouldBeNil)

				rule.Config.AlarmPatterns = w2.PatternList

				So(fmt.Sprintf("%+v", rule.Config.AlarmPatterns.AsMongoDriverQuery()["$or"]), ShouldEqual, "[map[v.connector_name:map[$regex:tesk.$]]]")

				Convey("mismatch event complex rule", func() {
					now := types.CpsTime{Time: time.Now()}
					testEvent = types.Event{
						Component:     "testComponent",
						Resource:      "testResource2",
						EventType:     types.EventTypeCheck,
						SourceType:    types.SourceTypeResource,
						Connector:     "test",
						ConnectorName: "test",
						State:         types.AlarmStateCritical,
						Output:        "",
						Timestamp:     now,
					}

					entity2 := types.NewEntity("testResource2/testComponent", "testResource", "resource", nil, nil, nil)
					err = entityAdapter.Insert(ctx, entity2)
					So(err, ShouldBeNil)

					alarm2, err := types.NewAlarm(testEvent, c)
					So(err, ShouldBeNil)

					err = alarmAdapter.Insert(ctx, alarm2)
					So(err, ShouldBeNil)

					testEvent.Alarm = &alarm2
					testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

					Convey("Test RuleMatched", func() {
						So(testEvent.Alarm, ShouldNotBeNil)
						metaAlarmEvents, _ := applicator.Apply(ctx, &testEvent, rule)
						So(len(metaAlarmEvents), ShouldEqual, 0)
					})
				})

				So(rule.Config.TimeInterval, ShouldEqual, 3)
			})
		})

		Convey("timebased rule match", func() {
			Convey("mismatch event complex rule", func() {
				rule := metaalarm.Rule{
					ID:   "timebased-test",
					Type: "timebased",
					Config: metaalarm.RuleConfig{
						TimeInterval: 3,
					},
				}

				testEvent1 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource3",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity3 := types.NewEntity("testResource3/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(ctx, entity3)
				So(err, ShouldBeNil)

				alarm3, err := types.NewAlarm(testEvent1, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(ctx, alarm3)
				So(err, ShouldBeNil)

				testEvent1.Alarm = &alarm3
				testEvent1.Entity = &types.Entity{ID: testEvent1.GetEID()}

				testEvent4 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource4",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity4 := types.NewEntity("testResource4/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(ctx, entity4)
				So(err, ShouldBeNil)

				alarm4, err := types.NewAlarm(testEvent4, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(ctx, alarm4)
				So(err, ShouldBeNil)

				testEvent4.Alarm = &alarm4
				testEvent4.Entity = &types.Entity{ID: testEvent4.GetEID()}

				metaAlarmEvents, _ := applicator.Apply(ctx, &testEvent1, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				metaAlarmEvents, _ = applicator.Apply(ctx, &testEvent4, rule)
				So(len(metaAlarmEvents), ShouldEqual, 1)
				So(metaAlarmEvents[0].MetaAlarmRuleID, ShouldEqual, rule.ID)

				time.Sleep(time.Second * 3)

				err = alarmAdapter.Insert(ctx, types.Alarm{
					ID:       "",
					Time:     types.CpsTime{},
					EntityID: "",
					Value: types.AlarmValue{
						Meta: rule.ID,
						Children: []string{
							testEvent1.GetEID(),
							testEvent4.GetEID(),
						},
					},
				})
				So(err, ShouldBeNil)

				testEvent1.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(ctx, &testEvent1, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				testEvent4.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(ctx, &testEvent4, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0) //New metaalarm shouldn't be created, since there is an existing one

				time.Sleep(time.Second * 3)

				testEvent5 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource5",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity5 := types.NewEntity("testResource5/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(ctx, entity5)
				So(err, ShouldBeNil)

				alarm5, err := types.NewAlarm(testEvent5, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(ctx, alarm5)
				So(err, ShouldBeNil)

				testEvent5.Alarm = &alarm5
				testEvent5.Entity = &types.Entity{ID: testEvent5.GetEID()}

				testEvent6 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource6",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity6 := types.NewEntity("testResource6/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(ctx, entity6)
				So(err, ShouldBeNil)

				alarm6, err := types.NewAlarm(testEvent6, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(ctx, alarm6)
				So(err, ShouldBeNil)

				testEvent6.Alarm = &alarm6
				testEvent6.Entity = &types.Entity{ID: testEvent6.GetEID()}

				testEvent5.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(ctx, &testEvent5, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				testEvent6.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(ctx, &testEvent6, rule)
				So(len(metaAlarmEvents), ShouldEqual, 1)
				So(metaAlarmEvents[0].MetaAlarmRuleID, ShouldEqual, rule.ID)
			})
		})
	})

}
