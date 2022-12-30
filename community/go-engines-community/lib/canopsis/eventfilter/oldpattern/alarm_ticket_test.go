package oldpattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type alarmTicketRefPatternWrapper struct {
	Pattern oldpattern.AlarmTicketRefPattern `bson:"pattern"`
}

func TestAlarmTicketRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_t":  "test1",
				"t":   bson.M{">": 10},
				"a":   "test2",
				"m":   bson.M{"regex_match": "abc-.*-def"},
				"val": "a12",
				"data": bson.M{
					"priority_id": bson.M{
						"regex_match": "^(?!1 - Critical).*$",
					},
				},
			},
		}
		mongoFilter := bson.M{
			"ticket._t": bson.M{
				"$eq": "test1",
			},
			"ticket.t": bson.M{
				"$gt": int64(10),
			},
			"ticket.a": bson.M{
				"$eq": "test2",
			},
			"ticket.m": bson.M{
				"$regex": "abc-.*-def",
			},
			"ticket.ticket": bson.M{
				"$eq": "a12",
			},
			"ticket.ticket_data.priority_id": bson.M{
				"$regex": "^(?!1 - Critical).*$",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := bson.M{}
			p.AsMongoDriverQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmTicketRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M{"ticket": nil}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := bson.M{}
			p.AsMongoDriverQuery("ticket", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmTicketRefPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"data": bson.M{
					"priority_id": bson.M{
						"regex_match": "^(?!1 - Critical).*$",
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmTicketRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern matches the right alarm steps", func() {
			matches := oldpattern.NewAlarmTicketRegexMatches()

			alarmTicket1 := &types.AlarmStep{
				TicketInfo: types.TicketInfo{
					TicketData: map[string]string{
						"priority_id": "2 - Critical",
					},
				},
			}
			So(p.Matches(alarmTicket1, &matches), ShouldBeTrue)

			alarmTicket2 := &types.AlarmStep{
				TicketInfo: types.TicketInfo{
					TicketData: map[string]string{
						"priority_id": "1 - Critical",
					},
				},
			}
			So(p.Matches(alarmTicket2, &matches), ShouldBeFalse)
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
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.AlarmTicketRefPattern
	}{
		{
			TestName: "test for full pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"_t":  "test1",
					"t":   bson.M{">": int64(10)},
					"a":   "test2",
					"m":   bson.M{"regex_match": "abc-.*-def"},
					"val": "a12",
				},
			},
			Pattern: oldpattern.AlarmTicketRefPattern{
				AlarmTicketFields: oldpattern.AlarmTicketFields{
					Type: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test1",
							},
						},
					},
					Timestamp: oldpattern.TimePattern{
						IntegerPattern: oldpattern.IntegerPattern{
							IntegerConditions: oldpattern.IntegerConditions{
								Gt: types.OptionalInt64{
									Set:   true,
									Value: 10,
								},
							},
						},
					},
					Author: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test2",
							},
						},
					},
					Message: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "a12",
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for partial pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"_t":  "test1",
					"m":   bson.M{"regex_match": "abc-.*-def"},
					"val": "a12",
				},
			},
			Pattern: oldpattern.AlarmTicketRefPattern{
				AlarmTicketFields: oldpattern.AlarmTicketFields{
					Type: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test1",
							},
						},
					},
					Message: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "a12",
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: bson.M{
				"pattern": nil,
			},
			Pattern: oldpattern.AlarmTicketRefPattern{
				ShouldBeNil: true,
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmTicketRefPatternWrapper
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
