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

type boolPatternWrapper struct {
	Pattern pattern.BoolPattern `bson:"pattern"`
}

func TestBoolPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is false", t, func() {
		mapPattern := mgobson.M{
			"pattern": false,
		}
		mongoFilter := mgobson.M{
			"$eq": false,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w boolPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestBoolPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is false", t, func() {
		mapPattern := mongobson.M{
			"pattern": false,
		}
		mongoFilter := mongobson.M{
			"$eq": false,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w boolPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidBoolPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is true", t, func() {
		mapPattern := mgobson.M{
			"pattern": true,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		mapPattern := mgobson.M{
			"pattern": false,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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

func TestValidBoolPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern checking that a value is true", t, func() {
		mapPattern := mongobson.M{
			"pattern": true,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		mapPattern := mongobson.M{
			"pattern": false,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w boolPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.BoolPattern
	} {
		{
			TestName: "test for true",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": true,
			},
			Pattern: pattern.BoolPattern{
				OptionalBool: utils.OptionalBool{
					Set:   true,
					Value: true,
				},
			},
		},
		{
			TestName: "test for false",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": false,
			},
			Pattern: pattern.BoolPattern{
				OptionalBool: utils.OptionalBool{
					Set:   true,
					Value: false,
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.BoolPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w boolPatternWrapper
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