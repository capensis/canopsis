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

type stringPatternWrapper struct {
	Pattern pattern.StringPattern `bson:"pattern"`
}

type stringRefPatternWrapper struct {
	Pattern pattern.StringRefPattern `bson:"pattern"`
}

func TestEqualStringPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := mgobson.M{
			"pattern": "toto",
		}
		mongoFilter := mgobson.M{
			"$eq": "toto",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualStringPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := mongobson.M{
			"pattern": "toto",
		}
		mongoFilter := mongobson.M{
			"$eq": "toto",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mgobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mongobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualNilStringRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M(nil)
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualNilStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M(nil)
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualStringRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := mgobson.M{
			"pattern": "toto",
		}
		mongoFilter := mgobson.M{
			"$eq": "toto",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := mongobson.M{
			"pattern": "toto",
		}
		mongoFilter := mongobson.M{
			"$eq": "toto",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringRefPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mgobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := mongobson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidEqualStringPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, "value")
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches("value", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualStringPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, "value")
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches("value", &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches("test", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				So(p.Matches("", &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidRegexPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-.*-def")
			})

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
		})
	})

	Convey("Given a valid BSON pattern representing a regular expression with subexpressions", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-(?P<sub>.*)-def")
			})

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
		})
	})
}

func TestValidRegexPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-.*-def")
			})

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
		})
	})

	Convey("Given a valid BSON pattern representing a regular expression with subexpressions", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-(?P<sub>.*)-def")
			})

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
		})
	})
}

func TestValidEqualNilStringRefPatternMgoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeTrue)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept nil", func() {
				So(p.Matches(nil, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				value := "value"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

		})
	})
}

func TestValidEqualNilStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeTrue)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept nil", func() {
				So(p.Matches(nil, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				value := "value"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})

		})
	})
}

func TestValidEqualStringRefPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mgobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, "value")
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the string", func() {
				value := "value"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidEqualStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := mongobson.M{
			"pattern": "value",
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeTrue)
				So(p.Equal.Value, ShouldEqual, "value")
				So(p.RegexMatch.Set, ShouldBeFalse)
			})

			Convey("The pattern should accept the value of the string", func() {
				value := "value"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject any other value", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidRegexRefPatternMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-.*-def")
			})

			Convey("The pattern should accept values that match the regex", func() {
				value := "abc-bla-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				value = "abc-ok-!-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
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
			var w stringRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-(?P<sub>.*)-def")
			})

			Convey("The pattern should accept values that match the regex and return the values of the subexpressions", func() {
				value := "abc-bla-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "bla")

				value = "abc-ok-!-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "ok-!")
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestValidRegexRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.EqualNil, ShouldBeFalse)
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-.*-def")
			})

			Convey("The pattern should accept values that match the regex", func() {
				value := "abc-bla-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)

				value = "abc-ok-!-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m, ShouldBeEmpty)
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)

				value = ""
				So(p.Matches(&value, &m), ShouldBeFalse)
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
			var w stringRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m pattern.RegexMatches

			Convey("The decoded values should be correct", func() {
				So(p.Equal.Set, ShouldBeFalse)
				So(p.RegexMatch.Set, ShouldBeTrue)
				So(p.RegexMatch.Value.String(), ShouldEqual, "abc-(?P<sub>.*)-def")
			})

			Convey("The pattern should accept values that match the regex and return the values of the subexpressions", func() {
				value := "abc-bla-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "bla")

				value = "abc-ok-!-def"
				So(p.Matches(&value, &m), ShouldBeTrue)
				So(m["sub"], ShouldEqual, "ok-!")
			})

			Convey("The pattern should reject values that do not match the regex", func() {
				value := "test"
				So(p.Matches(&value, &m), ShouldBeFalse)
				So(m, ShouldBeEmpty)
			})
		})
	})
}

func TestInvalidStringPatternMgoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 3,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidStringPatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 3,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidStringRefPatternMgoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := mgobson.M{
			"pattern": 3,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := mongobson.M{
			"pattern": 3,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestStringPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.StringPattern
	} {
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": "test",
			},
			Pattern: pattern.StringPattern{
				StringConditions: pattern.StringConditions{
					Equal: utils.OptionalString{
						Set: true,
						Value: "test",
					},
				},
			},
		},
		{
			TestName: "test for regexp",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"regex_match": "abc-.*-def",
				},
			},
			Pattern: pattern.StringPattern{
				StringConditions: pattern.StringConditions{
					RegexMatch: utils.OptionalRegexp{
						Set: true,
						Value: testRegexp,
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.StringPattern{
				StringConditions: pattern.StringConditions{},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w stringPatternWrapper
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

func TestStringRefPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.StringRefPattern
	} {
		{
			TestName: "test for ref equal",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": "test",
			},
			Pattern: pattern.StringRefPattern{
				EqualNil:      false,
				StringPattern: pattern.StringPattern{
					StringConditions: pattern.StringConditions{
						Equal: utils.OptionalString{
							Set: true,
							Value: "test",
						},
					},
				},
			},
		},
		{
			TestName: "test for ref regexp",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"regex_match": "abc-.*-def",
				},
			},
			Pattern: pattern.StringRefPattern{
				EqualNil:      false,
				StringPattern: pattern.StringPattern{
					StringConditions: pattern.StringConditions{
						RegexMatch: utils.OptionalRegexp{
							Set: true,
							Value: testRegexp,
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
			Pattern: pattern.StringRefPattern{
				EqualNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.StringRefPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w stringRefPatternWrapper
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
