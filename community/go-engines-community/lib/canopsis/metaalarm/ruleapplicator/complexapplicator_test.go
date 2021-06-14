package ruleapplicator_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/entity"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm/ruleapplicator"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/metaalarm/service"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/log"
	"git.canopsis.net/canopsis/go-engines/lib/mongo"
	"git.canopsis.net/canopsis/go-engines/lib/redis"
	"git.canopsis.net/canopsis/go-engines/lib/testutils"
	"github.com/bsm/redislock"
	"github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
)

func testNewComplexApplicator() (*ruleapplicator.ComplexApplicator, alarm.Adapter, entity.Adapter, error) {
	logger := log.NewLogger(true)

	redisClient, err := redis.NewSession(redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}
	redisClient.FlushDB()

	redisClient2, err := redis.NewSession(redis.CorrelationLockStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}
	redisClient2.FlushDB()

	redisLockClient := redislock.New(redisClient2)

	mongo, err := mongo.NewSession(mongo.Timeout)
	if err != nil {
		return nil, nil, nil, err
	}

	alarmCollection := alarm.DefaultCollection(mongo)
	_, err = alarmCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, err
	}

	entityCollection := entity.DefaultCollection(mongo)
	_, err = entityCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, err
	}

	alarmAdapter := alarm.NewAdapterLegacy(alarmCollection, entity.EntityCollectionName, alarmCollection.NewBulk(1000))

	entityAdapter := entity.NewAdapter(entityCollection)

	rulesCollection := metaalarm.DefaultRulesCollection(mongo)
	_, err = rulesCollection.RemoveAll()
	if err != nil {
		return nil, nil, nil, err
	}
	rulesAdapter := metaalarm.NewRuleAdapter(rulesCollection)

	metaAlarmService := service.NewMetaAlarmService(alarmAdapter, entityAdapter, rulesAdapter, testutils.GetTestConf(), log.NewLogger(true))

	redisClient3, err := redis.NewSession(redis.RuleTotalEntitiesStorage, logger, 0, 0)
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
	c := testutils.GetTestConf()
	applicator, alarmAdapter, entityAdapter, err := testNewComplexApplicator()

	Convey("Create RuleMatch test", t, func() {
		So(err, ShouldBeNil)

		testEvent := types.Event{}
		rule := metaalarm.Rule{}

		Convey("empty event and rule", func() {
			metaAlarmEvents, _ := applicator.Apply(&testEvent, rule)
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
						bson.M{
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

				So(fmt.Sprintf("%+v", rule.Config.AlarmPatterns.AsMongoQuery()["$or"]), ShouldEqual, "[map[v.connector_name:map[$eq:test899]]]")

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
					err = entityAdapter.Insert(entity1)
					So(err, ShouldBeNil)

					alarm1, err := types.NewAlarm(testEvent, c)
					So(err, ShouldBeNil)

					err = alarmAdapter.Insert(alarm1)
					So(err, ShouldBeNil)

					testEvent.Alarm = &alarm1
					testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

					Convey("Test RuleMatched", func() {
						metaAlarmEvents, _ := applicator.Apply(&testEvent, rule)
						So(len(metaAlarmEvents), ShouldEqual, 1)
						So(metaAlarmEvents[0].EventType, ShouldEqual, "metaalarm")
					})
				})
				mapPattern = bson.M{
					"list": []bson.M{
						bson.M{
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

				err = bson.Unmarshal(bsonPattern, &w)
				So(err, ShouldBeNil)

				rule.Config.AlarmPatterns = w.PatternList

				So(fmt.Sprintf("%+v", rule.Config.AlarmPatterns.AsMongoQuery()["$or"]), ShouldEqual, "[map[v.connector_name:map[$regex:tesk.$]]]")

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
					err = entityAdapter.Insert(entity2)
					So(err, ShouldBeNil)

					alarm2, err := types.NewAlarm(testEvent, c)
					So(err, ShouldBeNil)

					err = alarmAdapter.Insert(alarm2)
					So(err, ShouldBeNil)

					testEvent.Alarm = &alarm2
					testEvent.Entity = &types.Entity{ID: testEvent.GetEID()}

					Convey("Test RuleMatched", func() {
						So(testEvent.Alarm, ShouldNotBeNil)
						metaAlarmEvents, _ := applicator.Apply(&testEvent, rule)
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
				err = entityAdapter.Insert(entity3)
				So(err, ShouldBeNil)

				alarm3, err := types.NewAlarm(testEvent1, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(alarm3)
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
				err = entityAdapter.Insert(entity4)
				So(err, ShouldBeNil)

				alarm4, err := types.NewAlarm(testEvent4, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(alarm4)
				So(err, ShouldBeNil)

				testEvent4.Alarm = &alarm4
				testEvent4.Entity = &types.Entity{ID: testEvent4.GetEID()}

				metaAlarmEvents, _ := applicator.Apply(&testEvent1, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				metaAlarmEvents, _ = applicator.Apply(&testEvent4, rule)
				So(len(metaAlarmEvents), ShouldEqual, 1)
				So(metaAlarmEvents[0].MetaAlarmRuleID, ShouldEqual, rule.ID)

				time.Sleep(time.Second * 3)

				err = alarmAdapter.Insert(types.Alarm{
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
				metaAlarmEvents, _ = applicator.Apply(&testEvent1, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				testEvent4.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(&testEvent4, rule)
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
				err = entityAdapter.Insert(entity5)
				So(err, ShouldBeNil)

				alarm5, err := types.NewAlarm(testEvent5, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(alarm5)
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
				err = entityAdapter.Insert(entity6)
				So(err, ShouldBeNil)

				alarm6, err := types.NewAlarm(testEvent6, c)
				So(err, ShouldBeNil)

				err = alarmAdapter.Insert(alarm6)
				So(err, ShouldBeNil)

				testEvent6.Alarm = &alarm6
				testEvent6.Entity = &types.Entity{ID: testEvent6.GetEID()}

				testEvent5.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(&testEvent5, rule)
				So(len(metaAlarmEvents), ShouldEqual, 0)
				testEvent6.Alarm.Value.LastUpdateDate = types.CpsTime{Time: time.Now()}
				metaAlarmEvents, _ = applicator.Apply(&testEvent6, rule)
				So(len(metaAlarmEvents), ShouldEqual, 1)
				So(metaAlarmEvents[0].MetaAlarmRuleID, ShouldEqual, rule.ID)
			})
		})
	})

}
