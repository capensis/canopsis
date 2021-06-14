package pattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type alarmPatternWrapper struct {
	Pattern pattern.AlarmPattern `bson:"pattern"`
}

func TestAlarmPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_id": bson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": bson.M{
					"ack": bson.M{
						"_t": "test1",
					},
					"canceled": bson.M{
						"_t": "test2",
					},
					"done": bson.M{
						"_t": "test3",
					},
					"snooze": bson.M{
						"_t": "test4",
					},
					"state": bson.M{
						"_t": "test5",
					},
					"status": bson.M{
						"_t": "test6",
					},
					"ticket": bson.M{
						"_t": "test7",
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  bson.M{">": 10},
					"extra": bson.M{
						"key": "test11",
					},
					"hard_limit":                        bson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  bson.M{"<=": 12},
					"last_event_date":                   bson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": bson.M{">=": 14},
					"total_state_changes":               bson.M{">=": 15},
				},
			},
		}
		mongoFilter := bson.M{
			"_id": bson.M{"$regex": "abc"},
			"t":   bson.M{"$eq": int64(9)},
			"d":   bson.M{"$eq": "def"},

			"v.ack._t":           bson.M{"$eq": "test1"},
			"v.canceled._t":      bson.M{"$eq": "test2"},
			"v.done._t":          bson.M{"$eq": "test3"},
			"v.snooze._t":        bson.M{"$eq": "test4"},
			"v.state._t":         bson.M{"$eq": "test5"},
			"v.status._t":        bson.M{"$eq": "test6"},
			"v.ticket._t":        bson.M{"$eq": "test7"},
			"v.component":        bson.M{"$eq": "test8"},
			"v.connector":        bson.M{"$eq": "test9"},
			"v.connector_name":   bson.M{"$eq": "test10"},
			"v.creation_date":    bson.M{"$gt": int64(10)},
			"v.extra.key":        bson.M{"$eq": "test11"},
			"v.hard_limit":       bson.M{"$lt": int64(11)},
			"v.initial_output":   bson.M{"$eq": "test12"},
			"v.last_update_date": bson.M{"$lte": int64(12)},
			"v.last_event_date":  bson.M{"$gte": int64(13)},
			"v.resource":         bson.M{"$eq": "test13"},
			"v.resolved":         bson.M(nil),

			"v.state_changes_since_status_update": bson.M{"$gte": int64(14)},
			"v.total_state_changes":               bson.M{"$gte": int64(15)},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_id": bson.M{"regex_match": "abc"},
				"t":   9,
				"d":   "def",
				"v": bson.M{
					"ack": bson.M{
						"_t": "test1",
					},
					"canceled": bson.M{
						"_t": "test2",
					},
					"done": bson.M{
						"_t": "test3",
					},
					"snooze": bson.M{
						"_t": "test4",
					},
					"state": bson.M{
						"_t": "test5",
					},
					"status": bson.M{
						"_t": "test6",
					},
					"ticket": bson.M{
						"_t": "test7",
						"data": bson.M{
							"priority_id": bson.M{
								"regex_match": "^(?!1 - Critical).*$",
							},
						},
					},
					"component":      "test8",
					"connector":      "test9",
					"connector_name": "test10",
					"creation_date":  bson.M{">": 10},
					"extra": bson.M{
						"key": "test11",
					},
					"hard_limit":                        bson.M{"<": 11},
					"initial_output":                    "test12",
					"last_update_date":                  bson.M{"<=": 12},
					"last_event_date":                   bson.M{">=": 13},
					"resource":                          "test13",
					"resolved":                          nil,
					"state_changes_since_status_update": bson.M{">=": 14},
					"total_state_changes":               bson.M{">=": 15},
				},
			},
		}

		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		ExpectedUnmarshalled bson.M
		Pattern              pattern.AlarmPattern
	}{
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"_id": bson.M{"regex_match": "abc-.*-def"},
					"t":   int64(9),
					"v": bson.M{
						"ack": bson.M{
							"_t": "test1",
						},
						"canceled": bson.M{
							"_t": "test2",
						},
						"done": bson.M{
							"_t": "test3",
						},
						"snooze": bson.M{
							"_t": "test4",
						},
						"state": bson.M{
							"_t": "test5",
						},
						"status": bson.M{
							"_t": "test6",
						},
						"ticket": bson.M{
							"_t": "test7",
						},
						"component":      "test8",
						"connector":      "test9",
						"connector_name": "test10",
						"creation_date":  bson.M{">": int64(10)},
						"extra": bson.M{
							"key": "test11",
						},
						"hard_limit":                        bson.M{"<": int64(11)},
						"initial_output":                    "test12",
						"last_update_date":                  bson.M{"<=": int64(12)},
						"last_event_date":                   bson.M{">=": int64(13)},
						"resource":                          "test13",
						"resolved":                          nil,
						"state_changes_since_status_update": bson.M{">=": int64(14)},
						"total_state_changes":               bson.M{">=": int64(15)},
					},
				},
			},
			Pattern: pattern.AlarmPattern{
				ShouldNotBeNil: true,
				ShouldBeNil:    false,
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
												Set:   true,
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
												Set:   true,
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
												Set:   true,
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
												Set:   true,
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
												Set:   true,
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
												Set:   true,
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
										Set:   true,
										Value: "test8",
									},
								},
							},
							Connector: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
										Value: "test9",
									},
								},
							},
							ConnectorName: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
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
											Set:   true,
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
										Set:   true,
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

			resultBson, err := bson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled bson.M
			err = bson.Unmarshal(resultBson, &unmarshalled)
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
// AlarmPatternList because bson.Unmarshal does not work when called with an
// array.
type alarmPatternListWrapper struct {
	PatternList pattern.AlarmPatternList `bson:"list"`
}

func TestAlarmPatternListMatchMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"d": "id1",
				},
				bson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

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
		mapPattern := bson.M{}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

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
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"_id": "id1",
				},
				bson.M{
					"_id": bson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestAlarmPatternListToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"d": "id1",
				},
				bson.M{
					"d": "id2",
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := bson.M{
					"$or": []bson.M{
						bson.M{
							"d": bson.M{
								"$eq": "id1",
							},
						},
						bson.M{
							"d": bson.M{
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
		mapPattern := bson.M{}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := bson.M{}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := bson.M{
			"list": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w alarmPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := bson.M{}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})
}

func TestAlarmPatternListMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              pattern.AlarmPatternList
	}{
		{
			TestName: "test for pattern list",
			ExpectedUnmarshalled: bson.M{
				"list": bson.A{
					bson.M{
						"d": "id1",
					},
					bson.M{
						"d": "id2",
					},
				},
			},
			Pattern: pattern.AlarmPatternList{
				Set:   true,
				Valid: true,
				Patterns: []pattern.AlarmPattern{
					{
						AlarmFields: pattern.AlarmFields{
							EntityID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
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
										Set:   true,
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
			ExpectedUnmarshalled: bson.M{
				"list": nil,
			},
			Pattern: pattern.AlarmPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmPatternListWrapper
			w.PatternList = dataset.Pattern

			resultBson, err := bson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled bson.M
			err = bson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}
