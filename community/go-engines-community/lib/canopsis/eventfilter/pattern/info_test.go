package pattern_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	mgobson "github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

type infoPatternWrapper struct {
	Pattern pattern.InfoPattern `bson:"pattern"`
}

func TestInfoPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an info query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"name":        "test",
				"description": "test info",
				"value":       mgobson.M{"regex_match": "abc-.*-def"},
			},
		}
		mongoFilter := mgobson.M{
			"name": mgobson.M{
				"$eq": "test",
			},
			"description": mgobson.M{
				"$eq": "test info",
			},
			"value": mgobson.M{
				"$regex": "abc-.*-def",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestInfoPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an info query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"name":        "test",
				"description": "test info",
				"value":       mongobson.M{"regex_match": "abc-.*-def"},
			},
		}
		mongoFilter := mongobson.M{
			"name": mongobson.M{
				"$eq": "test",
			},
			"description": mongobson.M{
				"$eq": "test info",
			},
			"value": mongobson.M{
				"$regex": "abc-.*-def",
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInfoPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil info query", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M(nil)
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInfoPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil info query", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M(nil)
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestInfoPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.InfoPattern
	} {
		{
			TestName: "test for full info pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"name":        "test",
					"description": "test info",
					"value":       mongobson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: pattern.InfoPattern{
				InfoFields: pattern.InfoFields{
					Name: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
								Value: "test",
							},
						},
					},
					Description: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
								Value: "test info",
							},
						},
					},
					Value: pattern.StringPattern{
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
			TestName: "test for partial info pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"description": "test info",
					"value":       mongobson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: pattern.InfoPattern{
				InfoFields: pattern.InfoFields{
					Description: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
								Value: "test info",
							},
						},
					},
					Value: pattern.StringPattern{
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
			ExpectedUnmarshalled: mongobson.M{
				"pattern": nil,
			},
			Pattern: pattern.InfoPattern{
				ShouldNotBeSet: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.InfoPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w infoPatternWrapper
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