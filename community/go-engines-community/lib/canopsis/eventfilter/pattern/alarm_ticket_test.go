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

type alarmTicketRefPatternWrapper struct {
	Pattern pattern.AlarmTicketRefPattern `bson:"pattern"`
}

func TestAlarmTicketRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_t":  "test1",
				"t":   mgobson.M{">": 10},
				"a":   "test2",
				"m":   mgobson.M{"regex_match": "abc-.*-def"},
				"val": "a12",
				"data": mongobson.M{
					"priority_id": mongobson.M{
						"regex_match": "^(?!1 - Critical).*$",
					},
				},
			},
		}
		mongoFilter := mgobson.M{
			"ticket._t": mgobson.M{
				"$eq": "test1",
			},
			"ticket.t": mgobson.M{
				"$gt": int64(10),
			},
			"ticket.a": mgobson.M{
				"$eq": "test2",
			},
			"ticket.m": mgobson.M{
				"$regex": "abc-.*-def",
			},
			"ticket.val": mgobson.M{
				"$eq": "a12",
			},
			"ticket.data.priority_id": mgobson.M{
				"$regex": "^(?!1 - Critical).*$",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mgobson.M{}
			p.AsMongoQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmTicketRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_t":  "test1",
				"t":   mongobson.M{">": 10},
				"a":   "test2",
				"m":   mongobson.M{"regex_match": "abc-.*-def"},
				"val": "a12",
				"data": mongobson.M{
					"priority_id": mongobson.M{
						"regex_match": "^(?!1 - Critical).*$",
					},
				},
			},
		}
		mongoFilter := mongobson.M{
			"ticket._t": mongobson.M{
				"$eq": "test1",
			},
			"ticket.t": mongobson.M{
				"$gt": int64(10),
			},
			"ticket.a": mongobson.M{
				"$eq": "test2",
			},
			"ticket.m": mongobson.M{
				"$regex": "abc-.*-def",
			},
			"ticket.val": mongobson.M{
				"$eq": "a12",
			},
			"ticket.data.priority_id": mongobson.M{
				"$regex": "^(?!1 - Critical).*$",
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mongobson.M{}
			p.AsMongoDriverQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmTicketRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M{"ticket": nil}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mgobson.M{}
			p.AsMongoQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmTicketRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M{"ticket": nil}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mongobson.M{}
			p.AsMongoDriverQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmTicketRefPatternMatchesMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_t":  "test1",
				"t":   mgobson.M{">": 10},
				"a":   "test2",
				"m":   mgobson.M{"regex_match": "abc-.*-def"},
				"val": "a12",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern matches the right alarm steps", func() {
			matches := pattern.NewAlarmTicketRegexMatches()

			alarmTicket1 := types.AlarmTicket{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Author:    "test2",
				Message:   "abc-1-def",
				Value:     "a12",
			}
			So(p.Matches(&alarmTicket1, &matches), ShouldBeTrue)

			alarmTicket2 := types.AlarmTicket{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Author:    "test2",
				Message:   "abc-1-def",
			}
			So(p.Matches(&alarmTicket2, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmTicketRefPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"data": mongobson.M{
					"priority_id": mongobson.M{
						"regex_match": "^(?!1 - Critical).*$",
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern matches the right alarm steps", func() {
			matches := pattern.NewAlarmTicketRegexMatches()

			alarmTicket1 := types.AlarmTicket{
				Data: map[string]string{
					"priority_id": "2 - Critical",
				},
			}
			So(p.Matches(&alarmTicket1, &matches), ShouldBeTrue)

			alarmTicket2 := types.AlarmTicket{
				Data: map[string]string{
					"priority_id": "1 - Critical",
				},
			}
			So(p.Matches(&alarmTicket2, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmTicketRefPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.AlarmTicketRefPattern
	} {
		{
			TestName: "test for full pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"_t":  "test1",
					"t":   mongobson.M{">": int64(10)},
					"a":   "test2",
					"m":   mongobson.M{"regex_match": "abc-.*-def"},
					"val": "a12",
				},
			},
			Pattern: pattern.AlarmTicketRefPattern{
				AlarmTicketFields: pattern.AlarmTicketFields{
					Type: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set: true,
								Value: "test1",
							},
						},
					},
					Timestamp: pattern.TimePattern{
						IntegerPattern: pattern.IntegerPattern{
							IntegerConditions: pattern.IntegerConditions{
								Gt: utils.OptionalInt64{
									Set:   true,
									Value: 10,
								},
							},
						},
					},
					Author: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set: true,
								Value: "test2",
							},
						},
					},
					Message: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set: true,
								Value: testRegexp,
							},
						},
					},
					Value: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set: true,
								Value: "a12",
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for partial pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"_t":  "test1",
					"m":   mongobson.M{"regex_match": "abc-.*-def"},
					"val": "a12",
				},
			},
			Pattern: pattern.AlarmTicketRefPattern{
				AlarmTicketFields: pattern.AlarmTicketFields{
					Type: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set: true,
								Value: "test1",
							},
						},
					},
					Message: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set: true,
								Value: testRegexp,
							},
						},
					},
					Value: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set: true,
								Value: "a12",
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": nil,
			},
			Pattern: pattern.AlarmTicketRefPattern{
				ShouldBeNil: true,
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmTicketRefPatternWrapper
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
