package pattern_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	mgobson "github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	mongobson "go.mongodb.org/mongo-driver/bson"
)

type integerPatternWrapper struct {
	Pattern pattern.IntegerPattern `bson:"pattern"`
}

type integerRefPatternWrapper struct {
	Pattern pattern.IntegerRefPattern `bson:"pattern"`
}

func TestRangeIntegerPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		mongoFilter := mgobson.M{
			"$gte": int64(0),
			"$lte": int64(2),
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRangeIntegerPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		mongoFilter := mongobson.M{
			"$gte": int64(0),
			"$lte": int64(2),
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualNilIntegerRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M(nil)
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualNilIntegerRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M(nil)
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualIntegerRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 5,
		}
		mongoFilter := mgobson.M{
			"$eq": int64(5),
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualIntegerRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 5,
		}
		mongoFilter := mongobson.M{
			"$eq": int64(5),
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w integerRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidEqualIntegerPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, 7)

				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(7), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(0), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualIntegerPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, 7)

				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(7), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(0), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualNilIntegerRefPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with nil", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeTrue)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(nil), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
				value = 7
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualNilIntegerRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with nil", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeTrue)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(nil), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
				value = 7
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualIntegerRefPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, 7)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				value = 7
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualIntegerRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, 7)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the integer", func() {
				value = 7
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gt.Set, ShouldBeTrue)
				So(p.Gt.Value, ShouldEqual, 0)
				So(p.Lt.Set, ShouldBeTrue)
				So(p.Lt.Value, ShouldEqual, 3)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(1), ShouldBeTrue)
				So(p.Matches(2), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(0), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON filter representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gte.Set, ShouldBeTrue)
				So(p.Gte.Value, ShouldEqual, 0)
				So(p.Lte.Set, ShouldBeTrue)
				So(p.Lte.Value, ShouldEqual, 2)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(0), ShouldBeTrue)
				So(p.Matches(1), ShouldBeTrue)
				So(p.Matches(2), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
				So(p.Matches(4), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gt.Set, ShouldBeTrue)
				So(p.Gt.Value, ShouldEqual, 0)
				So(p.Lt.Set, ShouldBeTrue)
				So(p.Lt.Value, ShouldEqual, 3)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(1), ShouldBeTrue)
				So(p.Matches(2), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(0), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON filter representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gte.Set, ShouldBeTrue)
				So(p.Gte.Value, ShouldEqual, 0)
				So(p.Lte.Set, ShouldBeTrue)
				So(p.Lte.Value, ShouldEqual, 2)

				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(0), ShouldBeTrue)
				So(p.Matches(1), ShouldBeTrue)
				So(p.Matches(2), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1), ShouldBeFalse)
				So(p.Matches(3), ShouldBeFalse)
				So(p.Matches(4), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesRefPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gt.Set, ShouldBeTrue)
				So(p.Gt.Value, ShouldEqual, 0)
				So(p.Lt.Set, ShouldBeTrue)
				So(p.Lt.Value, ShouldEqual, 3)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				value = 1
				So(p.Matches(&value), ShouldBeTrue)
				value = 2
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON filter representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gte.Set, ShouldBeTrue)
				So(p.Gte.Value, ShouldEqual, 0)
				So(p.Lte.Set, ShouldBeTrue)
				So(p.Lte.Value, ShouldEqual, 2)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				value = 0
				So(p.Matches(&value), ShouldBeTrue)
				value = 1
				So(p.Matches(&value), ShouldBeTrue)
				value = 2
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
				value = 4
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gt.Set, ShouldBeTrue)
				So(p.Gt.Value, ShouldEqual, 0)
				So(p.Lt.Set, ShouldBeTrue)
				So(p.Lt.Value, ShouldEqual, 3)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gte.Set, ShouldBeFalse)
				So(p.Lte.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				value = 1
				So(p.Matches(&value), ShouldBeTrue)
				value = 2
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 0
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON filter representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">=": 0,
				"<=": 2,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsNumber
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.Gte.Set, ShouldBeTrue)
				So(p.Gte.Value, ShouldEqual, 0)
				So(p.Lte.Set, ShouldBeTrue)
				So(p.Lte.Value, ShouldEqual, 2)

				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.Gt.Set, ShouldBeFalse)
				So(p.Lt.Set, ShouldBeFalse)
			})

			Convey("The filter should accept values that are in the range", func() {
				value = 0
				So(p.Matches(&value), ShouldBeTrue)
				value = 1
				So(p.Matches(&value), ShouldBeTrue)
				value = 2
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = -1
				So(p.Matches(&value), ShouldBeFalse)
				value = 3
				So(p.Matches(&value), ShouldBeFalse)
				value = 4
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestInvalidIntegerPatternMgoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": "string",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{">": "string"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with nil", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{">": nil},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidIntegerPatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": "string",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{">": "string"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with nil", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{">": nil},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidIntegerRefPatternMgoDriver(t *testing.T) {
	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": "string",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{">": "string"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with nil", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{">": nil},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidIntegerRefPatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": "string",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{">": "string"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with nil", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{">": nil},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w integerRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestIntegerPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.IntegerPattern
	} {
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": int64(7),
			},
			Pattern: pattern.IntegerPattern{
				IntegerConditions: pattern.IntegerConditions{
					Equal: utils.OptionalInt64{
						Set: true,
						Value: 7,
					},
				},
			},
		},
		{
			TestName: "test for closed interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: pattern.IntegerPattern{
				IntegerConditions: pattern.IntegerConditions{
					Gte: utils.OptionalInt64{
						Set: true,
						Value: 0,
					},
					Lte: utils.OptionalInt64{
						Set: true,
						Value: 2,
					},
				},
			},
		},
		{
			TestName: "test for open interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">": int64(0),
					"<": int64(3),
				},
			},
			Pattern: pattern.IntegerPattern{
				IntegerConditions: pattern.IntegerConditions{
					Gt: utils.OptionalInt64{
						Set: true,
						Value: 0,
					},
					Lt: utils.OptionalInt64{
						Set: true,
						Value: 3,
					},
				},
			},
		},
		{
			TestName: "test for half-open interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">=": int64(2),
					"<": int64(3),
				},
			},
			Pattern: pattern.IntegerPattern{
				IntegerConditions: pattern.IntegerConditions{
					Gte: utils.OptionalInt64{
						Set: true,
						Value: 2,
					},
					Lt: utils.OptionalInt64{
						Set: true,
						Value: 3,
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.IntegerPattern{
				IntegerConditions: pattern.IntegerConditions{},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w integerPatternWrapper
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

func TestIntegerRefPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.IntegerRefPattern
	} {
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": int64(7),
			},
			Pattern: pattern.IntegerRefPattern{
				IntegerPattern: pattern.IntegerPattern{
					IntegerConditions: pattern.IntegerConditions{
						Equal: utils.OptionalInt64{
							Set: true,
							Value: 7,
						},
					},
				},
			},
		},
		{
			TestName: "test for closed interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: pattern.IntegerRefPattern{
				IntegerPattern: pattern.IntegerPattern{
					IntegerConditions: pattern.IntegerConditions{
						Gte: utils.OptionalInt64{
							Set: true,
							Value: 0,
						},
						Lte: utils.OptionalInt64{
							Set: true,
							Value: 2,
						},
					},
				},
			},
		},
		{
			TestName: "test for open interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">": int64(0),
					"<": int64(3),
				},
			},
			Pattern: pattern.IntegerRefPattern{
				IntegerPattern: pattern.IntegerPattern{
					IntegerConditions: pattern.IntegerConditions{
						Gt: utils.OptionalInt64{
							Set: true,
							Value: 0,
						},
						Lt: utils.OptionalInt64{
							Set: true,
							Value: 3,
						},
					},
				},
			},
		},
		{
			TestName: "test for half-open interval",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">=": int64(2),
					"<": int64(3),
				},
			},
			Pattern: pattern.IntegerRefPattern{
				IntegerPattern: pattern.IntegerPattern{
					IntegerConditions: pattern.IntegerConditions{
						Gte: utils.OptionalInt64{
							Set: true,
							Value: 2,
						},
						Lt: utils.OptionalInt64{
							Set: true,
							Value: 3,
						},
					},
				},
			},
		},
		{
			TestName: "test for ref nil",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": nil,
			},
			Pattern: pattern.IntegerRefPattern{
				EqualNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.IntegerRefPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w integerRefPatternWrapper
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
