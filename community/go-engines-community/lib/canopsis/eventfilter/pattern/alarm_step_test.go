package pattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type alarmStepRefPatternWrapper struct {
	Pattern pattern.AlarmStepRefPattern `bson:"pattern"`
}

func TestAlarmStepRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_t":  "test1",
				"t":   bson.M{">": 10},
				"a":   "test2",
				"m":   bson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		mongoFilter := bson.M{
			"ack._t": bson.M{
				"$eq": "test1",
			},
			"ack.t": bson.M{
				"$gt": int64(10),
			},
			"ack.a": bson.M{
				"$eq": "test2",
			},
			"ack.m": bson.M{
				"$regex": "abc-.*-def",
			},
			"ack.val": bson.M{
				"$eq": int64(3),
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := bson.M{}
			p.AsMongoDriverQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestNilAlarmStepRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil alarm step query", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M{"ack": nil}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			query := bson.M{}
			p.AsMongoDriverQuery("ack", query)
			So(query, ShouldResemble, mongoFilter)
		})
	})
}

func TestAlarmStepRefPatternMatchesMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an alarm step", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_t":  "test1",
				"t":   bson.M{">": 10},
				"a":   "test2",
				"m":   bson.M{"regex_match": "abc-.*-def"},
				"val": 3,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w alarmStepRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		ExpectedUnmarshalled bson.M
		Pattern              pattern.AlarmStepRefPattern
	}{
		{
			TestName: "test for full pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"_t":        "test1",
					"t":         bson.M{">": int64(10)},
					"a":         "test2",
					"m":         bson.M{"regex_match": "abc-.*-def"},
					"val":       int64(3),
					"initiator": bson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: pattern.AlarmStepRefPattern{
				AlarmStepFields: pattern.AlarmStepFields{
					Type: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
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
								Set:   true,
								Value: "test2",
							},
						},
					},
					Message: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
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
					Initiator: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
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
					"_t":        "test1",
					"m":         bson.M{"regex_match": "abc-.*-def"},
					"val":       int64(3),
					"initiator": bson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: pattern.AlarmStepRefPattern{
				AlarmStepFields: pattern.AlarmStepFields{
					Type: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
								Value: "test1",
							},
						},
					},
					Message: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
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
					Initiator: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
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
			Pattern: pattern.AlarmStepRefPattern{
				ShouldBeNil: true,
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w alarmStepRefPatternWrapper
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
