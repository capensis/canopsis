package oldpattern_test

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type infoPatternWrapper struct {
	Pattern oldpattern.InfoPattern `bson:"pattern"`
}

func TestInfoPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an info query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"name":        "test",
				"description": "test info",
				"value":       bson.M{"regex_match": "abc-.*-def"},
			},
		}
		mongoFilter := bson.M{
			"name": bson.M{
				"$eq": "test",
			},
			"description": bson.M{
				"$eq": "test info",
			},
			"value": bson.M{
				"$regex": "abc-.*-def",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInfoPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil info query", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M(nil)
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w infoPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.InfoPattern
	}{
		{
			TestName: "test for full info pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"name":        "test",
					"description": "test info",
					"value":       bson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: oldpattern.InfoPattern{
				InfoFields: oldpattern.InfoFields{
					Name: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test",
							},
						},
					},
					Description: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test info",
							},
						},
					},
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
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
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"description": "test info",
					"value":       bson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: oldpattern.InfoPattern{
				InfoFields: oldpattern.InfoFields{
					Description: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test info",
							},
						},
					},
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
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
			Pattern: oldpattern.InfoPattern{
				ShouldNotBeSet: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.InfoPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w infoPatternWrapper
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
