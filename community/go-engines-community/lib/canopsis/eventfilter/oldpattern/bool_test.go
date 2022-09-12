package oldpattern_test

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type boolPatternWrapper struct {
	Pattern oldpattern.BoolPattern `bson:"pattern"`
}

func TestBoolPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is false", t, func() {
		mapPattern := bson.M{
			"pattern": false,
		}
		mongoFilter := bson.M{
			"$eq": false,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w boolPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidBoolPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is true", t, func() {
		mapPattern := bson.M{
			"pattern": true,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Set, ShouldBeTrue)
				So(p.Value, ShouldBeTrue)
			})

			Convey("The pattern should accept true", func() {
				So(p.Matches(true), ShouldBeTrue)
			})

			Convey("The pattern should reject false", func() {
				So(p.Matches(false), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON pattern checking that a value is false", t, func() {
		mapPattern := bson.M{
			"pattern": false,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Set, ShouldBeTrue)
				So(p.Value, ShouldBeFalse)
			})

			Convey("The pattern should accept false", func() {
				So(p.Matches(false), ShouldBeTrue)
			})

			Convey("The pattern should reject true", func() {
				So(p.Matches(true), ShouldBeFalse)
			})
		})
	})

	Convey("Given an empty pattern", t, func() {
		mapPattern := bson.M{}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept true and false", func() {
				So(p.Matches(true), ShouldBeTrue)
				So(p.Matches(false), ShouldBeTrue)
			})
		})
	})
}

func TestBoolPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.BoolPattern
	}{
		{
			TestName: "test for true",
			ExpectedUnmarshalled: bson.M{
				"pattern": true,
			},
			Pattern: oldpattern.BoolPattern{
				OptionalBool: types.OptionalBool{
					Set:   true,
					Value: true,
				},
			},
		},
		{
			TestName: "test for false",
			ExpectedUnmarshalled: bson.M{
				"pattern": false,
			},
			Pattern: oldpattern.BoolPattern{
				OptionalBool: types.OptionalBool{
					Set:   true,
					Value: false,
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.BoolPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w boolPatternWrapper
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
