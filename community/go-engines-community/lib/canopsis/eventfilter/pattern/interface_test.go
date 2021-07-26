package pattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	mgobson "github.com/globalsign/mgo/bson"
	. "github.com/smartystreets/goconvey/convey"
	mongobson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type interfacePatternWrapper struct {
	Pattern pattern.InterfacePattern `bson:"pattern"`
}

func TestIntRangeInterfacePatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
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

		var w interfacePatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestIntRangeInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
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

		var w interfacePatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringRegexInterfacePatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mgobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringRegexInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mongobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringArrayInterfacePatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"has_every":  []string{"test1"},
				"has_one_of": []string{"test2"},
				"has_not":    []string{"test3"},
			},
		}
		mongoFilter := mgobson.M{
			"$all": []string{"test1"},
			"$in":  []string{"test2"},
			"$nin": []string{"test3"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringArrayInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"has_every":  []string{"test1"},
				"has_one_of": []string{"test2"},
				"has_not":    []string{"test3"},
				"is_empty":   false,
			},
		}
		mongoFilter := mongobson.M{
			"$all":    []string{"test1"},
			"$in":     []string{"test2"},
			"$nin":    []string{"test3"},
			"$exists": true, "$ne": mongobson.A{},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInterfacePatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M(nil)
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M(nil)
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidEqualStringInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches("value", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualStringInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches("value", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidRegexInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept values that match the regex", func() {
				So(p.Matches("abc-bla-def", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches("abc-ok-!-def", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject non string values", func() {
				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a valid BSON pattern representing a regular expression with subexpressions", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept values that match the regex and return the values of the subexpressions", func() {
				So(p.Matches("abc-bla-def", &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "bla")

				So(p.Matches("abc-ok-!-def", &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "ok-!")
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject non string values", func() {
				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidRegexInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept values that match the regex", func() {
				So(p.Matches("abc-bla-def", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches("abc-ok-!-def", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject non string values", func() {
				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})

	Convey("Given a valid BSON pattern representing a regular expression with subexpressions", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept values that match the regex and return the values of the subexpressions", func() {
				So(p.Matches("abc-bla-def", &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "bla")

				So(p.Matches("abc-ok-!-def", &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "ok-!")
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject non string values", func() {
				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualNilInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept nil", func() {
				So(p.Matches(nil, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("value", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualNilInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept nil", func() {
				So(p.Matches(nil, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("value", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(12, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(5, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(7, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(0, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("value", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 7,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The pattern should accept the value of the integer", func() {
				So(p.Matches(7, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(0, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("value", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(nil, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidInequalitiesInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var m pattern.RegexMatches

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(1, &m), ShouldBeTrue)
				So(p.Matches(2, &m), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(0, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
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
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(0, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches(1, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches(2, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(4, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidInequalitiesInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var m pattern.RegexMatches

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(1, &m), ShouldBeTrue)
				So(p.Matches(2, &m), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(0, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
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
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The filter should accept values that are in the range", func() {
				So(p.Matches(0, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches(1, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				So(p.Matches(2, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(-1, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(4, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches(3.14, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestInvalidInterfacePatternMgoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
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
			var w interfacePatternWrapper
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
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is not a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": 3},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is nil", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": nil},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-(.*-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with conditions on integers and strings", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"regex_match": "abc-(.*)-def",
				">":           3,
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
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
			var w interfacePatternWrapper
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
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is not a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": 3},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is nil", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": nil},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-(.*-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with conditions on integers and strings", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"regex_match": "abc-(.*)-def",
				">":           3,
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInterfacePatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.InterfacePattern
	}{
		{
			TestName: "test for equal string",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": "value",
			},
			Pattern: pattern.InterfacePattern{
				StringConditions: pattern.StringConditions{
					Equal: utils.OptionalString{
						Set:   true,
						Value: "value",
					},
				},
			},
		},
		{
			TestName: "test for regexp",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{"regex_match": "abc-.*-def"},
			},
			Pattern: pattern.InterfacePattern{
				StringConditions: pattern.StringConditions{
					RegexMatch: utils.OptionalRegexp{
						Set:   true,
						Value: testRegexp,
					},
				},
			},
		},
		{
			TestName: "test for string array",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"has_every":  mongobson.A{"test1"},
					"has_one_of": mongobson.A{"test2"},
					"has_not":    mongobson.A{"test3"},
					"is_empty":   false,
				},
			},
			Pattern: pattern.InterfacePattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasEvery: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1"},
					},
					HasOneOf: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test2"},
					},
					HasNot: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test3"},
					},
					IsEmpty: utils.OptionalBool{
						Set:   true,
						Value: false,
					},
				},
			},
		},
		{
			TestName: "test for equal int value",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": int64(5),
			},
			Pattern: pattern.InterfacePattern{
				IntegerConditions: pattern.IntegerConditions{
					Equal: utils.OptionalInt64{
						Set:   true,
						Value: 5,
					},
				},
			},
		},
		{
			TestName: "test for range int values",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: pattern.InterfacePattern{
				IntegerConditions: pattern.IntegerConditions{
					Gte: utils.OptionalInt64{
						Set:   true,
						Value: 0,
					},
					Lte: utils.OptionalInt64{
						Set:   true,
						Value: 2,
					},
				},
			},
		},
		{
			TestName: "test for nil value",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": nil,
			},
			Pattern: pattern.InterfacePattern{
				EqualNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.InterfacePattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w interfacePatternWrapper
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
