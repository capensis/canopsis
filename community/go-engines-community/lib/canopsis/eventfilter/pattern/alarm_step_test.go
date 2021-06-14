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

type alarmStepRefPatternWrapper struct {
	Pattern pattern.AlarmStepRefPattern `bson:"pattern"`
}

func TestAlarmStepRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_t":  "test1",
				"t":   mgobson.M{">": 10},
				"a":   "test2",
				"m":   mgobson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		mongoFilter := mgobson.M{
			"ack._t": mgobson.M{
				"$eq": "test1",
			},
			"ack.t": mgobson.M{
				"$gt": int64(10),
			},
			"ack.a": mgobson.M{
				"$eq": "test2",
			},
			"ack.m": mgobson.M{
				"$regex": "abc-.*-def",
			},
			"ack.val": mgobson.M{
				"$eq": int64(3),
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mgobson.M{}
			p.AsMongoQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmStepRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_t":  "test1",
				"t":   mongobson.M{">": 10},
				"a":   "test2",
				"m":   mongobson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		mongoFilter := mongobson.M{
			"ack._t": mongobson.M{
				"$eq": "test1",
			},
			"ack.t": mongobson.M{
				"$gt": int64(10),
			},
			"ack.a": mongobson.M{
				"$eq": "test2",
			},
			"ack.m": mongobson.M{
				"$regex": "abc-.*-def",
			},
			"ack.val": mongobson.M{
				"$eq": int64(3),
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mongobson.M{}
			p.AsMongoDriverQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmStepRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M{"ack": nil}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mgobson.M{}
			p.AsMongoQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmStepRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M{"ack": nil}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := mongobson.M{}
			p.AsMongoDriverQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmStepRefPatternMatchesMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_t":  "test1",
				"t":   mgobson.M{">": 10},
				"a":   "test2",
				"m":   mgobson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern matches the right alarm steps", func() {
			matches := pattern.AlarmStepRegexMatches{}

			alarmStep1 := types.AlarmStep{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Author:    "test2",
				Message:   "abc-1-def",
				Value:     3,
			}
			So(p.Matches(&alarmStep1, &matches), ShouldBeTrue)

			alarmStep2 := types.AlarmStep{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Message:   "abc-1-def",
				Value:     3,
			}
			So(p.Matches(&alarmStep2, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmStepRefPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_t":  "test1",
				"t":   mongobson.M{">": 10},
				"a":   "test2",
				"m":   mongobson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern matches the right alarm steps", func() {
			matches := pattern.AlarmStepRegexMatches{}

			alarmStep1 := types.AlarmStep{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Author:    "test2",
				Message:   "abc-1-def",
				Value:     3,
			}
			So(p.Matches(&alarmStep1, &matches), ShouldBeTrue)

			alarmStep2 := types.AlarmStep{
				Type:      "test1",
				Timestamp: types.NewCpsTime(12),
				Message:   "abc-1-def",
				Value:     3,
			}
			So(p.Matches(&alarmStep2, &matches), ShouldBeFalse)
		})
	})
}

func TestAlarmStepRefPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.AlarmStepRefPattern
	} {
		{
			TestName: "test for full pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"_t":  "test1",
					"t":   mongobson.M{">": int64(10)},
					"a":   "test2",
					"m":   mongobson.M{"regex_match": "abc-.*-def"},
					"val": int64(3),
				},
			},
			Pattern: pattern.AlarmStepRefPattern{
				AlarmStepFields: pattern.AlarmStepFields{
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
					Value: pattern.IntegerPattern{
						IntegerConditions: pattern.IntegerConditions{
							Equal: utils.OptionalInt64{
								Set:   true,
								Value: 3,
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
					"val": int64(3),
				},
			},
			Pattern: pattern.AlarmStepRefPattern{
				AlarmStepFields: pattern.AlarmStepFields{
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
					Value: pattern.IntegerPattern{
						IntegerConditions: pattern.IntegerConditions{
							Equal: utils.OptionalInt64{
								Set:   true,
								Value: 3,
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
			Pattern: pattern.AlarmStepRefPattern{
				ShouldBeNil: true,
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmStepRefPatternWrapper
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
