package ruleapplicator

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/storage"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/alarm"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/correlation/service"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entity"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/log"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/redis"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

func testNewValueApplicator() (*ValueApplicator, alarm.Adapter, entity.Adapter, error) {
	logger := log.NewLogger(true)
	ctx := context.Background()

	redisClient, err := redis.NewSession(ctx, redis.AlarmGroupStorage, logger, 0, 0)
	if err != nil {
		panic(err)
	}

	dbClient, err := mongo.NewClient(0, 0)
	if err != nil {
		panic(err)
	}

	alarmAdapter := alarm.NewAdapter(dbClient)

	entityAdapter := entity.NewAdapter(dbClient)

	metaAlarmService := service.NewMetaAlarmService(alarmAdapter, config.NewAlarmConfigProvider(config.CanopsisConf{}, logger), log.NewTestLogger())

	redisClient3, err := redis.NewSession(ctx, redis.RuleTotalEntitiesStorage, logger, 0, 0)
	if err != nil {
		return nil, nil, nil, err
	}

	valueGroupEntityCounter := correlation.NewValueGroupEntityCounter(dbClient, redisClient3, logger)

	applicator := NewValueGroupApplicator(alarmAdapter, metaAlarmService, storage.NewRedisGroupingStorage(), redisClient, valueGroupEntityCounter, logger)
	return &applicator, alarmAdapter, entityAdapter, nil
}

func TestApply(t *testing.T) {
	displayNameScheme, err := config.CreateDisplayNameTpl(config.AlarmDefaultNameScheme)
	if err != nil {
		panic(err)
	}
	cfg := config.AlarmConfig{
		FlappingFreqLimit:    0,
		FlappingInterval:     0,
		StealthyInterval:     0,
		BaggotTime:           time.Second,
		EnableLastEventDate:  true,
		CancelAutosolveDelay: time.Hour,
		DisplayNameScheme:    displayNameScheme,
		OutputLength:         10,
	}

	ctx := context.Background()

	Convey("Test valuegroup rule match", t, func() {
		applicator, alarmAdapter, entityAdapter, err := testNewValueApplicator()
		So(err, ShouldBeNil)

		Convey("check simple valuegroup rule", func() {
			threshold := int64(2)
			rule := correlation.Rule{
				ID:   "valuegroup-test",
				Type: "valuegroup",
				Config: correlation.RuleConfig{
					TimeInterval: 300,
					ValuePaths: []string{
						"entity.infos.customer.value",
						"entity.infos.location.value",
					},
					ThresholdCount: &threshold,
				},
			}

			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource1",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity1 := types.NewEntity("testResource1/testComponent", "testResource", "resource", nil, nil, nil)
			err := entityAdapter.Insert(entity1)
			So(err, ShouldBeNil)

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)

			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &types.Entity{}
			testEvent1.Entity.Infos = make(map[string]types.Info)
			testEvent1.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent1.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource2",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity2 := types.NewEntity("testResource2/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity2)
			So(err, ShouldBeNil)

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &types.Entity{}
			testEvent2.Entity.Infos = make(map[string]types.Info)
			testEvent2.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent2.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource3",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity3 := types.NewEntity("testResource3/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity3)
			So(err, ShouldBeNil)

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &types.Entity{}
			testEvent3.Entity.Infos = make(map[string]types.Info)
			testEvent3.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent3.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

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

			alarm4, err := types.NewAlarm(testEvent4, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm4)
			So(err, ShouldBeNil)

			testEvent4.Alarm = &alarm4
			testEvent4.Entity = &types.Entity{}
			testEvent4.Entity.Infos = make(map[string]types.Info)
			testEvent4.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent4.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			testEvent5 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource5",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity5 := types.NewEntity("testResource5/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity5)
			So(err, ShouldBeNil)

			alarm5, err := types.NewAlarm(testEvent5, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm5)
			So(err, ShouldBeNil)

			testEvent5.Alarm = &alarm5
			testEvent5.Entity = &types.Entity{}
			testEvent5.Entity.Infos = make(map[string]types.Info)
			testEvent5.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent5.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

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

			alarm6, err := types.NewAlarm(testEvent6, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm6)
			So(err, ShouldBeNil)

			testEvent6.Alarm = &alarm6
			testEvent6.Entity = &types.Entity{}
			testEvent6.Entity.Infos = make(map[string]types.Info)
			testEvent6.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent6.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			testEvent7 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource7",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity7 := types.NewEntity("testResource7/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity7)
			So(err, ShouldBeNil)

			alarm7, err := types.NewAlarm(testEvent7, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm7)
			So(err, ShouldBeNil)

			testEvent7.Alarm = &alarm7
			testEvent7.Entity = &types.Entity{}
			testEvent7.Entity.Infos = make(map[string]types.Info)
			testEvent7.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent7.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent8 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource8",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity8 := types.NewEntity("testResource8/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity8)
			So(err, ShouldBeNil)

			alarm8, err := types.NewAlarm(testEvent8, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm8)
			So(err, ShouldBeNil)

			testEvent8.Alarm = &alarm8
			testEvent8.Entity = &types.Entity{}
			testEvent8.Entity.Infos = make(map[string]types.Info)
			testEvent8.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent8.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent := metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-1.location-1")
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent4, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent = metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-1.location-2")
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent5, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent6, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent7, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent = metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-2.location-1")
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent8, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent = metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-2.location-2")

			Convey("it shouldn't add to a metaalarm if some valuepath is empty", func() {
				threshold := int64(2)
				rule := correlation.Rule{
					ID:   "valuegroup-test-2",
					Type: "valuegroup",
					Config: correlation.RuleConfig{
						TimeInterval: 300,
						ValuePaths: []string{
							"entity.infos.customer.value",
							"entity.infos.location.value",
						},
						ThresholdCount: &threshold,
					},
				}

				testEvent1 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource10",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,

					Output:    "",
					Timestamp: types.CpsTime{Time: time.Now()},
				}

				entity10 := types.NewEntity("testResource10/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(entity10)
				So(err, ShouldBeNil)

				alarm1, err := types.NewAlarm(testEvent1, cfg)
				So(err, ShouldBeNil)
				err = alarmAdapter.Insert(alarm1)
				So(err, ShouldBeNil)

				testEvent1.Alarm = &alarm1
				testEvent1.Entity = &types.Entity{}
				testEvent1.Entity.Infos = make(map[string]types.Info)
				testEvent1.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

				testEvent2 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource20",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity20 := types.NewEntity("testResource20/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(entity20)
				So(err, ShouldBeNil)

				alarm2, err := types.NewAlarm(testEvent2, cfg)
				So(err, ShouldBeNil)
				err = alarmAdapter.Insert(alarm2)
				So(err, ShouldBeNil)

				testEvent2.Alarm = &alarm2
				testEvent2.Entity = &types.Entity{}
				testEvent2.Entity.Infos = make(map[string]types.Info)
				testEvent2.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")

				testEvent3 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource30",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,

					Output:    "",
					Timestamp: types.CpsTime{Time: time.Now()},
				}

				entity30 := types.NewEntity("testResource30/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(entity30)
				So(err, ShouldBeNil)

				alarm3, err := types.NewAlarm(testEvent3, cfg)
				So(err, ShouldBeNil)
				err = alarmAdapter.Insert(alarm3)
				So(err, ShouldBeNil)

				testEvent3.Alarm = &alarm3
				testEvent3.Entity = &types.Entity{}
				testEvent3.Entity.Infos = make(map[string]types.Info)
				testEvent3.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")

				testEvent4 := types.Event{
					Component:     "testComponent",
					Resource:      "testResource40",
					EventType:     types.EventTypeCheck,
					SourceType:    types.SourceTypeResource,
					Connector:     "test",
					ConnectorName: "test",
					State:         types.AlarmStateCritical,
					Output:        "",
					Timestamp:     types.CpsTime{Time: time.Now()},
				}

				entity40 := types.NewEntity("testResource40/testComponent", "testResource", "resource", nil, nil, nil)
				err = entityAdapter.Insert(entity40)
				So(err, ShouldBeNil)

				alarm4, err := types.NewAlarm(testEvent4, cfg)
				So(err, ShouldBeNil)
				err = alarmAdapter.Insert(alarm4)
				So(err, ShouldBeNil)

				testEvent4.Alarm = &alarm4
				testEvent4.Entity = &types.Entity{}
				testEvent4.Entity.Infos = make(map[string]types.Info)
				testEvent4.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

				metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule)
				So(metaAlarmEventArray, ShouldBeNil)
				metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule)
				So(metaAlarmEventArray, ShouldBeNil)
				metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule)
				So(metaAlarmEventArray, ShouldBeNil)
				metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent4, rule)
				So(metaAlarmEventArray, ShouldBeNil)
			})
		})

		Convey("check valuegroup rule with single item in paths", func() {
			threshold := int64(2)
			rule := correlation.Rule{
				ID:   "valuegroup-test-3",
				Type: "valuegroup",
				Config: correlation.RuleConfig{
					TimeInterval: 300,
					ValuePaths: []string{
						"entity.infos.customer.value",
					},
					ThresholdCount: &threshold,
				},
			}

			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource100",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity100 := types.NewEntity("testResource100/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity100)
			So(err, ShouldBeNil)

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &types.Entity{}
			testEvent1.Entity.Infos = make(map[string]types.Info)
			testEvent1.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent1.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource200",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity200 := types.NewEntity("testResource200/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity200)
			So(err, ShouldBeNil)

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &types.Entity{}
			testEvent2.Entity.Infos = make(map[string]types.Info)
			testEvent2.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent2.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource300",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity300 := types.NewEntity("testResource300/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity300)
			So(err, ShouldBeNil)

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &types.Entity{}
			testEvent3.Entity.Infos = make(map[string]types.Info)
			testEvent3.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent3.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent4 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource400",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity400 := types.NewEntity("testResource400/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity400)
			So(err, ShouldBeNil)

			alarm4, err := types.NewAlarm(testEvent4, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm4)
			So(err, ShouldBeNil)

			testEvent4.Alarm = &alarm4
			testEvent4.Entity = &types.Entity{}
			testEvent4.Entity.Infos = make(map[string]types.Info)
			testEvent4.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-2")
			testEvent4.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent := metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test-3")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-1")
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent4, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent = metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test-3")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-2")
		})

		Convey("shouldn't work without valuegroups", func() {
			threshold := int64(2)
			rule := correlation.Rule{
				ID:   "valuegroup-test-empty",
				Type: "valuegroup",
				Config: correlation.RuleConfig{
					TimeInterval:   300,
					ThresholdCount: &threshold,
				},
			}

			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource1000",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity1000 := types.NewEntity("testResource1000/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity1000)
			So(err, ShouldBeNil)

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &types.Entity{}
			testEvent1.Entity.Infos = make(map[string]types.Info)
			testEvent1.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent1.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource2000",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity2000 := types.NewEntity("testResource2000/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity2000)
			So(err, ShouldBeNil)

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &types.Entity{}
			testEvent2.Entity.Infos = make(map[string]types.Info)
			testEvent2.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent2.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource3000",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity3000 := types.NewEntity("testResource3000/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity3000)
			So(err, ShouldBeNil)

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &types.Entity{}
			testEvent3.Entity.Infos = make(map[string]types.Info)
			testEvent3.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent3.Entity.Infos["location"] = types.NewInfo("location", "location", "location-1")

			testEvent4 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource4000",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity4000 := types.NewEntity("testResource4000/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity4000)
			So(err, ShouldBeNil)

			alarm4, err := types.NewAlarm(testEvent4, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm4)
			So(err, ShouldBeNil)

			testEvent4.Alarm = &alarm4
			testEvent4.Entity = &types.Entity{}
			testEvent4.Entity.Infos = make(map[string]types.Info)
			testEvent4.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-1")
			testEvent4.Entity.Infos["location"] = types.NewInfo("location", "location", "location-2")

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent4, rule)
			So(metaAlarmEventArray, ShouldBeNil)
		})

		Convey("should calculate group length properly if the group contains resolved alarms", func() {
			threshold := int64(2)
			rule := correlation.Rule{
				ID:   "valuegroup-test",
				Type: "valuegroup",
				Config: correlation.RuleConfig{
					TimeInterval: 300,
					ValuePaths: []string{
						"entity.infos.customer.value",
						"entity.infos.location.value",
					},
					ThresholdCount: &threshold,
				},
			}

			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource01",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity1 := types.NewEntity("testResource01/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity1)
			So(err, ShouldBeNil)

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &types.Entity{}
			testEvent1.Entity.Infos = make(map[string]types.Info)
			testEvent1.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-3")
			testEvent1.Entity.Infos["location"] = types.NewInfo("location", "location", "location-3")

			//resolve the alarm1 so it shouldn't count in group len
			resolveTime := types.NewCpsTime(time.Now().Unix())
			alarm1.Resolve(&resolveTime)

			err = alarmAdapter.Update(alarm1)
			So(err, ShouldBeNil)

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource02",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			entity2 := types.NewEntity("testResource02/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity2)
			So(err, ShouldBeNil)

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &types.Entity{}
			testEvent2.Entity.Infos = make(map[string]types.Info)
			testEvent2.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-3")
			testEvent2.Entity.Infos["location"] = types.NewInfo("location", "location", "location-3")

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource03",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			entity3 := types.NewEntity("testResource03/testComponent", "testResource", "resource", nil, nil, nil)
			err = entityAdapter.Insert(entity3)
			So(err, ShouldBeNil)

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &types.Entity{}
			testEvent3.Entity.Infos = make(map[string]types.Info)
			testEvent3.Entity.Infos["customer"] = types.NewInfo("customer", "customer", "customer-3")
			testEvent3.Entity.Infos["location"] = types.NewInfo("location", "location", "location-3")

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent := metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-3.location-3")
		})
	})
}

func createEntity(id string, customer string, location string) types.Entity {
	infos := make(map[string]types.Info)
	infos["customer"] = types.NewInfo("customer", "customer", customer)
	infos["location"] = types.NewInfo("location", "location", location)

	return types.Entity{
		ID:    id,
		Infos: infos,
	}
}

func TestApplyWithRate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	displayNameScheme, err := config.CreateDisplayNameTpl(config.AlarmDefaultNameScheme)
	if err != nil {
		panic(err)
	}
	cfg := config.AlarmConfig{
		FlappingFreqLimit:    0,
		FlappingInterval:     0,
		StealthyInterval:     0,
		BaggotTime:           time.Second,
		EnableLastEventDate:  true,
		CancelAutosolveDelay: time.Hour,
		DisplayNameScheme:    displayNameScheme,
		OutputLength:         10,
	}

	Convey("Init data", t, func() {
		applicator, alarmAdapter, _, err := testNewValueApplicator()
		So(err, ShouldBeNil)

		dbClient, err := mongo.NewClient(0, 0)
		if err != nil {
			panic(err)
		}

		entityCollection := dbClient.Collection(mongo.EntityMongoCollection)
		_, err = entityCollection.DeleteMany(ctx, bson.M{})
		if err != nil {
			panic(err)
		}

		entity1 := createEntity("testResource1/testComponent", "customer-1", "location-1")
		_, err = entityCollection.InsertOne(ctx, entity1)
		So(err, ShouldBeNil)
		entity2 := createEntity("testResource2/testComponent", "customer-1", "location-1")
		_, err = entityCollection.InsertOne(ctx, entity2)
		So(err, ShouldBeNil)
		entity3 := createEntity("testResource3/testComponent", "customer-1", "location-1")
		_, err = entityCollection.InsertOne(ctx, entity3)
		So(err, ShouldBeNil)
		entity4 := createEntity("testResource4/testComponent", "customer-1", "location-1")
		_, err = entityCollection.InsertOne(ctx, entity4)
		So(err, ShouldBeNil)
		entity5 := createEntity("testResource5/testComponent", "customer-1", "location-2")
		_, err = entityCollection.InsertOne(ctx, entity5)
		So(err, ShouldBeNil)
		entity6 := createEntity("testResource6/testComponent", "customer-1", "location-2")
		_, err = entityCollection.InsertOne(ctx, entity6)
		So(err, ShouldBeNil)
		entity7 := createEntity("testResource7/testComponent", "customer-1", "location-2")
		_, err = entityCollection.InsertOne(ctx, entity7)
		So(err, ShouldBeNil)
		entity8 := createEntity("testResource8/testComponent", "customer-1", "location-2")
		_, err = entityCollection.InsertOne(ctx, entity8)
		So(err, ShouldBeNil)

		thresholdRate1 := 0.7
		rule1 := correlation.Rule{
			ID:   "valuegroup-test-rate-1",
			Type: "valuegroup",
			Config: correlation.RuleConfig{
				TimeInterval: 300,
				ValuePaths: []string{
					"entity.infos.customer.value",
					"entity.infos.location.value",
				},
				ThresholdRate: &thresholdRate1,
			},
		}

		thresholdRate2 := 0.8
		rule2 := correlation.Rule{
			ID:   "valuegroup-test-rate-2",
			Type: "valuegroup",
			Config: correlation.RuleConfig{
				TimeInterval: 300,
				ValuePaths: []string{
					"entity.infos.customer.value",
					"entity.infos.location.value",
				},
				ThresholdRate: &thresholdRate2,
			},
		}

		logger := log.NewLogger(true)
		redisClient, err := redis.NewSession(ctx, redis.RuleTotalEntitiesStorage, logger, 0, 0)
		if err != nil {
			So(err, ShouldBeNil)
		}

		redisClient.FlushAll(ctx)

		valueGroupEntityCounter := correlation.NewValueGroupEntityCounter(dbClient, redisClient, logger)
		err = valueGroupEntityCounter.CountTotalEntitiesAmount(ctx, rule1)
		So(err, ShouldBeNil)
		total, err := valueGroupEntityCounter.GetTotalEntitiesAmount(ctx, rule1.ID, "customer-1.location-1")
		So(err, ShouldBeNil)
		So(total, ShouldEqual, 4)
		total, err = valueGroupEntityCounter.GetTotalEntitiesAmount(ctx, rule1.ID, "customer-1.location-2")
		So(err, ShouldBeNil)
		So(total, ShouldEqual, 4)
		err = valueGroupEntityCounter.CountTotalEntitiesAmount(ctx, rule2)
		So(err, ShouldBeNil)
		total, err = valueGroupEntityCounter.GetTotalEntitiesAmount(ctx, rule2.ID, "customer-1.location-1")
		So(err, ShouldBeNil)
		So(total, ShouldEqual, 4)
		total, err = valueGroupEntityCounter.GetTotalEntitiesAmount(ctx, rule2.ID, "customer-1.location-2")
		So(err, ShouldBeNil)
		So(total, ShouldEqual, 4)

		Convey("check simple valuegroup rule1, rate should be reached - 3/4 > 0.7", func() {
			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource1",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &entity1

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource2",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &entity2

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource3",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &entity3

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

			alarm4, err := types.NewAlarm(testEvent4, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm4)
			So(err, ShouldBeNil)

			testEvent4.Alarm = &alarm4
			testEvent4.Entity = &entity4

			testEvent5 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource5",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm5, err := types.NewAlarm(testEvent5, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm5)
			So(err, ShouldBeNil)

			testEvent5.Alarm = &alarm5
			testEvent5.Entity = &entity5

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

			alarm6, err := types.NewAlarm(testEvent6, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm6)
			So(err, ShouldBeNil)

			testEvent6.Alarm = &alarm6
			testEvent6.Entity = &entity6

			testEvent7 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource7",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm7, err := types.NewAlarm(testEvent7, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm7)
			So(err, ShouldBeNil)

			testEvent7.Alarm = &alarm7
			testEvent7.Entity = &entity7

			testEvent8 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource8",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			alarm8, err := types.NewAlarm(testEvent8, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm8)
			So(err, ShouldBeNil)

			testEvent8.Alarm = &alarm8
			testEvent8.Entity = &entity8

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule1)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule1)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule1)
			So(metaAlarmEventArray, ShouldNotBeNil)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent := metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test-rate-1")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-1.location-1")
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent5, rule1)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent6, rule1)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent7, rule1)
			So(len(metaAlarmEventArray), ShouldEqual, 1)
			metaAlarmEvent = metaAlarmEventArray[0]
			So(metaAlarmEvent.MetaAlarmRuleID, ShouldEqual, "valuegroup-test-rate-1")
			So(metaAlarmEvent.MetaAlarmValuePath, ShouldEqual, "customer-1.location-2")
		})
		Convey("check simple valuegroup rule1, rate shouldn't be reached - 3/4 < 0.8", func() {
			testEvent1 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource1",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm1, err := types.NewAlarm(testEvent1, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm1)
			So(err, ShouldBeNil)

			testEvent1.Alarm = &alarm1
			testEvent1.Entity = &entity1

			testEvent2 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource2",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			alarm2, err := types.NewAlarm(testEvent2, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm2)
			So(err, ShouldBeNil)

			testEvent2.Alarm = &alarm2
			testEvent2.Entity = &entity2

			testEvent3 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource3",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm3, err := types.NewAlarm(testEvent3, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm3)
			So(err, ShouldBeNil)

			testEvent3.Alarm = &alarm3
			testEvent3.Entity = &entity3

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

			alarm4, err := types.NewAlarm(testEvent4, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm4)
			So(err, ShouldBeNil)

			testEvent4.Alarm = &alarm4
			testEvent4.Entity = &entity4

			testEvent5 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource5",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm5, err := types.NewAlarm(testEvent5, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm5)
			So(err, ShouldBeNil)

			testEvent5.Alarm = &alarm5
			testEvent5.Entity = &entity5

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

			alarm6, err := types.NewAlarm(testEvent6, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm6)
			So(err, ShouldBeNil)

			testEvent6.Alarm = &alarm6
			testEvent6.Entity = &entity6

			testEvent7 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource7",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,

				Output:    "",
				Timestamp: types.CpsTime{Time: time.Now()},
			}

			alarm7, err := types.NewAlarm(testEvent7, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm7)
			So(err, ShouldBeNil)

			testEvent7.Alarm = &alarm7
			testEvent7.Entity = &entity7

			testEvent8 := types.Event{
				Component:     "testComponent",
				Resource:      "testResource8",
				EventType:     types.EventTypeCheck,
				SourceType:    types.SourceTypeResource,
				Connector:     "test",
				ConnectorName: "test",
				State:         types.AlarmStateCritical,
				Output:        "",
				Timestamp:     types.CpsTime{Time: time.Now()},
			}

			alarm8, err := types.NewAlarm(testEvent8, cfg)
			So(err, ShouldBeNil)
			err = alarmAdapter.Insert(alarm8)
			So(err, ShouldBeNil)

			testEvent8.Alarm = &alarm8
			testEvent8.Entity = &entity8

			metaAlarmEventArray, _ := applicator.Apply(ctx, testEvent1, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent2, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent3, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent5, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent6, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
			metaAlarmEventArray, _ = applicator.Apply(ctx, testEvent7, rule2)
			So(metaAlarmEventArray, ShouldBeNil)
		})
	})
}
