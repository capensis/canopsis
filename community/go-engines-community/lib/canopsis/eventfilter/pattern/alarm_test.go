package pattern_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

type alarmPatternWrapper struct {
	Pattern pattern.AlarmPattern `bson:"pattern"`
}

func TestAlarmPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_id": mgobson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": mgobson.M{
					"ack": mgobson.M{
						"_t": "test1",
					},
					"canceled": mgobson.M{
						"_t": "test2",
					},
					"done": mgobson.M{
						"_t": "test3",
					},
					"snooze": mgobson.M{
						"_t": "test4",
					},
					"state": mgobson.M{
						"_t": "test5",
					},
					"status": mgobson.M{
						"_t": "test6",
					},
					"ticket": mgobson.M{
						"_t": "test7",
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  mgobson.M{">": 10},
					"extra": mgobson.M{
						"key": "test11",
					},
					"hard_limit":                        mgobson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  mgobson.M{"<=": 12},
					"last_event_date":                   mgobson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": mgobson.M{">=": 14},
					"total_state_changes":               mgobson.M{">=": 15},
				},
			},
		}
		mongoFilter := mgobson.M{
			"_id": mgobson.M{"$regex": "abc"},
			"t":   mgobson.M{"$eq": int64(9)},
			"d":   mgobson.M{"$eq": "def"},

			"v.ack._t":           mgobson.M{"$eq": "test1"},
			"v.canceled._t":      mgobson.M{"$eq": "test2"},
			"v.done._t":          mgobson.M{"$eq": "test3"},
			"v.snooze._t":        mgobson.M{"$eq": "test4"},
			"v.state._t":         mgobson.M{"$eq": "test5"},
			"v.status._t":        mgobson.M{"$eq": "test6"},
			"v.ticket._t":        mgobson.M{"$eq": "test7"},
			"v.component":        mgobson.M{"$eq": "test8"},
			"v.connector":        mgobson.M{"$eq": "test9"},
			"v.connector_name":   mgobson.M{"$eq": "test10"},
			"v.creation_date":    mgobson.M{"$gt": int64(10)},
			"v.extra.key":        mgobson.M{"$eq": "test11"},
			"v.hard_limit":       mgobson.M{"$lt": int64(11)},
			"v.initial_output":   mgobson.M{"$eq": "test12"},
			"v.last_update_date": mgobson.M{"$lte": int64(12)},
			"v.last_event_date":  mgobson.M{"$gte": int64(13)},
			"v.resource":         mgobson.M{"$eq": "test13"},
			"v.resolved":         mgobson.M(nil),

			"v.state_changes_since_status_update": mgobson.M{"$gte": int64(14)},
			"v.total_state_changes":               mgobson.M{"$gte": int64(15)},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_id": mongobson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": mongobson.M{
					"ack": mongobson.M{
						"_t": "test1",
					},
					"canceled": mongobson.M{
						"_t": "test2",
					},
					"done": mongobson.M{
						"_t": "test3",
					},
					"snooze": mongobson.M{
						"_t": "test4",
					},
					"state": mongobson.M{
						"_t": "test5",
					},
					"status": mongobson.M{
						"_t": "test6",
					},
					"ticket": mongobson.M{
						"_t": "test7",
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  mongobson.M{">": 10},
					"extra": mongobson.M{
						"key": "test11",
					},
					"hard_limit":                        mongobson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  mongobson.M{"<=": 12},
					"last_event_date":                   mongobson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": mongobson.M{">=": 14},
					"total_state_changes":               mongobson.M{">=": 15},
				},
			},
		}
		mongoFilter := mongobson.M{
			"_id": mongobson.M{"$regex": "abc"},
			"t":   mongobson.M{"$eq": int64(9)},
			"d":   mongobson.M{"$eq": "def"},

			"v.ack._t":           mongobson.M{"$eq": "test1"},
			"v.canceled._t":      mongobson.M{"$eq": "test2"},
			"v.done._t":          mongobson.M{"$eq": "test3"},
			"v.snooze._t":        mongobson.M{"$eq": "test4"},
			"v.state._t":         mongobson.M{"$eq": "test5"},
			"v.status._t":        mongobson.M{"$eq": "test6"},
			"v.ticket._t":        mongobson.M{"$eq": "test7"},
			"v.component":        mongobson.M{"$eq": "test8"},
			"v.connector":        mongobson.M{"$eq": "test9"},
			"v.connector_name":   mongobson.M{"$eq": "test10"},
			"v.creation_date":    mongobson.M{"$gt": int64(10)},
			"v.extra.key":        mongobson.M{"$eq": "test11"},
			"v.hard_limit":       mongobson.M{"$lt": int64(11)},
			"v.initial_output":   mongobson.M{"$eq": "test12"},
			"v.last_update_date": mongobson.M{"$lte": int64(12)},
			"v.last_event_date":  mongobson.M{"$gte": int64(13)},
			"v.resource":         mongobson.M{"$eq": "test13"},
			"v.resolved":         mongobson.M(nil),

			"v.state_changes_since_status_update": mongobson.M{"$gte": int64(14)},
			"v.total_state_changes":               mongobson.M{"$gte": int64(15)},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmPatternMatchesMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_id": mgobson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": mgobson.M{
					"ack": mgobson.M{
						"_t": "test1",
					},
					"canceled": mgobson.M{
						"_t": "test2",
					},
					"done": mgobson.M{
						"_t": "test3",
					},
					"snooze": mgobson.M{
						"_t": "test4",
					},
					"state": mgobson.M{
						"_t": "test5",
					},
					"status": mgobson.M{
						"_t": "test6",
					},
					"ticket": mgobson.M{
						"_t": "test7",
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  mgobson.M{">": 10},
					"extra": mgobson.M{
						"key": "test11",
					},
					"hard_limit":                        mgobson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  mgobson.M{"<=": 12},
					"last_event_date":                   mgobson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": mgobson.M{">=": 14},
					"total_state_changes":               mgobson.M{">=": 15},
				},
			},
		}

		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should match the right alarm values", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test1",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
					},
					Component:     "test8",
					Connector:     "test9",
					ConnectorName: "test10",
					CreationDate:  types.NewCpsTime(12),
					Extra: map[string]interface{}{
						"key": "test11",
					},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}
			So(p.Matches(&alarm, &matches), ShouldBeTrue)
		})

		Convey("The pattern should not match an alarm value with the wrong ack", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test2",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
					},
					Component:     "test8",
					Connector:     "test9",
					ConnectorName: "test10",
					CreationDate:  types.NewCpsTime(12),
					Extra: map[string]interface{}{
						"key": "test11",
					},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}

			So(p.Matches(&alarm, &matches), ShouldBeFalse)
		})

		Convey("The pattern should not match an alarm value with the wrong extra", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test2",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
					},
					Component:                     "test8",
					Connector:                     "test9",
					ConnectorName:                 "test10",
					CreationDate:                  types.NewCpsTime(12),
					Extra:                         map[string]interface{}{},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}
			So(p.Matches(&alarm, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_id": mongobson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": mongobson.M{
					"ack": mongobson.M{
						"_t": "test1",
					},
					"canceled": mongobson.M{
						"_t": "test2",
					},
					"done": mongobson.M{
						"_t": "test3",
					},
					"snooze": mongobson.M{
						"_t": "test4",
					},
					"state": mongobson.M{
						"_t": "test5",
					},
					"status": mongobson.M{
						"_t": "test6",
					},
					"ticket": mongobson.M{
						"_t": "test7",
						"data": mongobson.M{
							"priority_id": mongobson.M{
								"regex_match": "^(?!1 - Critical).*$",
							},
						},
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  mongobson.M{">": 10},
					"extra": mongobson.M{
						"key": "test11",
					},
					"hard_limit":                        mongobson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  mongobson.M{"<=": 12},
					"last_event_date":                   mongobson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": mongobson.M{">=": 14},
					"total_state_changes":               mongobson.M{">=": 15},
				},
			},
		}

		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should match the right alarm values", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test1",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
						Data: map[string]string{
							"priority_id": "2 - Critical",
						},
					},
					Component:     "test8",
					Connector:     "test9",
					ConnectorName: "test10",
					CreationDate:  types.NewCpsTime(12),
					Extra: map[string]interface{}{
						"key": "test11",
					},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}
			So(p.Matches(&alarm, &matches), ShouldBeTrue)
		})

		Convey("The pattern should not match an alarm value with the wrong ack", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test2",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
					},
					Component:     "test8",
					Connector:     "test9",
					ConnectorName: "test10",
					CreationDate:  types.NewCpsTime(12),
					Extra: map[string]interface{}{
						"key": "test11",
					},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}

			So(p.Matches(&alarm, &matches), ShouldBeFalse)
		})

		Convey("The pattern should not match an alarm value with the wrong extra", func() {
			matches := pattern.NewAlarmRegexMatches()
			hardLimit := types.CpsNumber(10)
			alarm := types.Alarm{
				ID:       "abc",
				Time:     types.NewCpsTime(9),
				EntityID: "def",
				Value: types.AlarmValue{
					ACK: &types.AlarmStep{
						Type: "test2",
					},
					Canceled: &types.AlarmStep{
						Type: "test2",
					},
					Done: &types.AlarmStep{
						Type: "test3",
					},
					Snooze: &types.AlarmStep{
						Type: "test4",
					},
					State: &types.AlarmStep{
						Type: "test5",
					},
					Status: &types.AlarmStep{
						Type: "test6",
					},
					Ticket: &types.AlarmTicket{
						Type: "test7",
					},
					Component:                     "test8",
					Connector:                     "test9",
					ConnectorName:                 "test10",
					CreationDate:                  types.NewCpsTime(12),
					Extra:                         map[string]interface{}{},
					HardLimit:                     &hardLimit,
					InitialOutput:                 "test12",
					LastUpdateDate:                types.NewCpsTime(12),
					LastEventDate:                 types.NewCpsTime(13),
					Resource:                      "test13",
					Resolved:                      nil,
					StateChangesSinceStatusUpdate: types.CpsNumber(14),
					TotalStateChanges:             types.CpsNumber(15),
				},
			}
			So(p.Matches(&alarm, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.AlarmPattern
	} {
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"_id": mongobson.M{"regex_match": "abc-.*-def"},
					"t": int64(9),
					"v": mongobson.M{
						"ack": mongobson.M{
							"_t": "test1",
						},
						"canceled": mongobson.M{
							"_t": "test2",
						},
						"done": mongobson.M{
							"_t": "test3",
						},
						"snooze": mongobson.M{
							"_t": "test4",
						},
						"state": mongobson.M{
							"_t": "test5",
						},
						"status": mongobson.M{
							"_t": "test6",
						},
						"ticket": mongobson.M{
							"_t": "test7",
						},
						"component":      "test8",
						"connector":      "test9",
						"connector_name": "test10",
						"creation_date":  mongobson.M{">": int64(10)},
						"extra": mongobson.M{
							"key": "test11",
						},
						"hard_limit":                        mongobson.M{"<": int64(11)},
						"initial_output":                    "test12",
						"last_update_date":                  mongobson.M{"<=": int64(12)},
						"last_event_date":                   mongobson.M{">=": int64(13)},
						"resource":                          "test13",
						"resolved":                          nil,
						"state_changes_since_status_update": mongobson.M{">=": int64(14)},
						"total_state_changes":               mongobson.M{">=": int64(15)},
					},
				},
			},
			Pattern: pattern.AlarmPattern{
				ShouldNotBeNil: true,
				ShouldBeNil: false,
				AlarmFields: pattern.AlarmFields{
					ID: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
					Time: pattern.TimePattern{
						IntegerPattern: pattern.IntegerPattern{
							IntegerConditions: pattern.IntegerConditions{
								Equal: utils.OptionalInt64{
									Set:   true,
									Value: 9,
								},
							},
						},
					},
					Value: pattern.AlarmValuePattern{
						AlarmValueFields: pattern.AlarmValueFields{
							ACK: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test1",
											},
										},
									},
								},
							},
							Canceled: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test2",
											},
										},
									},
								},
							},
							Done: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test3",
											},
										},
									},
								},
							},
							Snooze: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test4",
											},
										},
									},
								},
							},
							State: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test5",
											},
										},
									},
								},
							},
							Status: pattern.AlarmStepRefPattern{
								AlarmStepFields: pattern.AlarmStepFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set: true,
												Value: "test6",
											},
										},
									},
								},
							},
							Ticket: pattern.AlarmTicketRefPattern{
								AlarmTicketFields: pattern.AlarmTicketFields{
									Type: pattern.StringPattern{
										StringConditions: pattern.StringConditions{
											Equal: utils.OptionalString{
												Set:   true,
												Value: "test7",
											},
										},
									},
								},
							},
							Component: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "test8",
									},
								},
							},
							Connector: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "test9",
									},
								},
							},
							ConnectorName: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "test10",
									},
								},
							},
							CreationDate: pattern.TimePattern{
								IntegerPattern: pattern.IntegerPattern{
									IntegerConditions: pattern.IntegerConditions{
										Gt: utils.OptionalInt64{
											Set:   true,
											Value: 10,
										},
									},
								},
							},
							Extra: map[string]pattern.InterfacePattern{
								"key": {
									StringConditions: pattern.StringConditions{
										Equal: utils.OptionalString{
											Set:   true,
											Value: "test11",
										},
									},
								},
							},
							HardLimit: pattern.IntegerRefPattern{
								IntegerPattern: pattern.IntegerPattern{
									IntegerConditions: pattern.IntegerConditions{
										Lt: utils.OptionalInt64{
											Set: true,
											Value: 11,
										},
									},
								},
							},
							InitialOutput: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
										Value: "test12",
									},
								},
							},
							LastUpdateDate: pattern.TimePattern{
								IntegerPattern: pattern.IntegerPattern{
									IntegerConditions: pattern.IntegerConditions{
										Lte: utils.OptionalInt64{
											Set:   true,
											Value: 12,
										},
									},
								},
							},
							LastEventDate: pattern.TimePattern{
								IntegerPattern: pattern.IntegerPattern{
									IntegerConditions: pattern.IntegerConditions{
										Gte: utils.OptionalInt64{
											Set:   true,
											Value: 13,
										},
									},
								},
							},
							Resource: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "test13",
									},
								},
							},
							Resolved: pattern.TimeRefPattern{
								IntegerRefPattern: pattern.IntegerRefPattern{
									EqualNil: true,
								},
							},
							StateChangesSinceStatusUpdate: pattern.IntegerPattern{
								IntegerConditions: pattern.IntegerConditions{
									Gte: utils.OptionalInt64{
										Set:   true,
										Value: 14,
									},
								},
							},
							TotalStateChanges: pattern.IntegerPattern{
								IntegerConditions: pattern.IntegerConditions{
									Gte: utils.OptionalInt64{
										Set:   true,
										Value: 15,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmPatternWrapper
			w.Pattern = dataset.Pattern

			resultBson, err := mongobson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled mongobson.M
			err = mongobson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}

// alarmPatternListWrapper is a type that wraps an AlarmPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// AlarmPatternList because mgobson.Unmarshal does not work when called with an
// array.
type alarmPatternListWrapper struct {
	PatternList pattern.AlarmPatternList `bson:"list"`
}

func TestAlarmPatternListMatchMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"d": "id1",
				},
				mgobson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding alarms", func() {
				alarm1 := types.Alarm{
					EntityID: "id1",
				}
				So(p.Matches(&alarm1), ShouldBeTrue)

				alarm2 := types.Alarm{
					EntityID: "id2",
				}
				So(p.Matches(&alarm2), ShouldBeTrue)

				alarm3 := types.Alarm{
					EntityID: "id3",
				}
				So(p.Matches(&alarm3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding alarms", func() {
				alarm1 := types.Alarm{
					EntityID: "id1",
				}
				So(p.Matches(&alarm1), ShouldBeTrue)

				alarm2 := types.Alarm{
					EntityID: "id2",
				}
				So(p.Matches(&alarm2), ShouldBeTrue)

				alarm3 := types.Alarm{
					EntityID: "id3",
				}
				So(p.Matches(&alarm3), ShouldBeTrue)
			})
		})
	})

	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"_id": "id1",
				},
				mgobson.M{
					"_id": mgobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestAlarmPatternListMatchMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"d": "id1",
				},
				mongobson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding alarms", func() {
				alarm1 := types.Alarm{
					EntityID: "id1",
				}
				So(p.Matches(&alarm1), ShouldBeTrue)

				alarm2 := types.Alarm{
					EntityID: "id2",
				}
				So(p.Matches(&alarm2), ShouldBeTrue)

				alarm3 := types.Alarm{
					EntityID: "id3",
				}
				So(p.Matches(&alarm3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding alarms", func() {
				alarm1 := types.Alarm{
					EntityID: "id1",
				}
				So(p.Matches(&alarm1), ShouldBeTrue)

				alarm2 := types.Alarm{
					EntityID: "id2",
				}
				So(p.Matches(&alarm2), ShouldBeTrue)

				alarm3 := types.Alarm{
					EntityID: "id3",
				}
				So(p.Matches(&alarm3), ShouldBeTrue)
			})
		})
	})

	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"_id": "id1",
				},
				mongobson.M{
					"_id": mongobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestAlarmPatternListToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"d": "id1",
				},
				mgobson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mgobson.M{
					"$or": []mgobson.M{
						mgobson.M{
							"d": mgobson.M{
								"$eq": "id1",
							},
						},
						mgobson.M{
							"d": mgobson.M{
								"$eq": "id2",
							},
						},
					},
				}
				So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mgobson.M{}
				So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
			})
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mgobson.M{}
				So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
			})
		})
	})
}

func TestAlarmPatternListToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"d": "id1",
				},
				mongobson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mongobson.M{
					"$or": []mongobson.M{
						mongobson.M{
							"d": mongobson.M{
								"$eq": "id1",
							},
						},
						mongobson.M{
							"d": mongobson.M{
								"$eq": "id2",
							},
						},
					},
				}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mongobson.M{}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mongobson.M{}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})
}

func TestAlarmPatternListMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.AlarmPatternList
	} {
		{
			TestName: "test for pattern list",
			ExpectedUnmarshalled: mongobson.M{
				"list": mongobson.A{
					mongobson.M{
						"d": "id1",
					},
					mongobson.M{
						"d": "id2",
					},
				},
			},
			Pattern: pattern.AlarmPatternList{
				Set:      true,
				Valid:    true,
				Patterns: []pattern.AlarmPattern{
					{
						AlarmFields: pattern.AlarmFields{
							EntityID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "id1",
									},
								},
							},
						},
					},
					{
						AlarmFields: pattern.AlarmFields{
							EntityID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set: true,
										Value: "id2",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: mongobson.M{
				"list": nil,
			},
			Pattern: pattern.AlarmPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmPatternListWrapper
			w.PatternList = dataset.Pattern

			resultBson, err := mongobson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled mongobson.M
			err = mongobson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}