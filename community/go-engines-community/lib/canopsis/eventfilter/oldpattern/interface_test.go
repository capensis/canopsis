package oldpattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type interfacePatternWrapper struct {
	Pattern oldpattern.InterfacePattern `bson:"pattern"`
}

func TestIntRangeInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				">=": 0,
				"<=": 2,
			},
		}
		mongoFilter := bson.M{
			"$gte": int64(0),
			"$lte": int64(2),
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringRegexInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := bson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringArrayInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every":  []string{"test1"},
				"has_one_of": []string{"test2"},
				"has_not":    []string{"test3"},
				"is_empty":   false,
			},
		}
		mongoFilter := bson.M{
			"$all":    []string{"test1"},
			"$in":     []string{"test2"},
			"$nin":    []string{"test3"},
			"$exists": true, "$ne": bson.A{},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilInterfacePatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a integer range query", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M(nil)
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w interfacePatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidEqualStringInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := bson.M{
			"pattern": "value",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidRegexInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidEqualNilInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidEqualInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": 7,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidInequalitiesInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var m oldpattern.RegexMatches

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
		mapPattern := bson.M{
			"pattern": bson.M{
				">=": 0,
				"<=": 2,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestInvalidInterfacePatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with a string", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{">": "string"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an inequality with nil", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{">": nil},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is not a string", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": 3},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a regexp that is nil", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": nil},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-(.*-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with conditions on integers and strings", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"regex_match": "abc-(.*)-def",
				">":           3,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w interfacePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
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
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.InterfacePattern
	}{
		{
			TestName: "test for equal string",
			ExpectedUnmarshalled: bson.M{
				"pattern": "value",
			},
			Pattern: oldpattern.InterfacePattern{
				StringConditions: oldpattern.StringConditions{
					Equal: types.OptionalString{
						Set:   true,
						Value: "value",
					},
				},
			},
		},
		{
			TestName: "test for regexp",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{"regex_match": "abc-.*-def"},
			},
			Pattern: oldpattern.InterfacePattern{
				StringConditions: oldpattern.StringConditions{
					RegexMatch: types.OptionalRegexp{
						Set:   true,
						Value: testRegexp,
					},
				},
			},
		},
		{
			TestName: "test for string array",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_every":  bson.A{"test1"},
					"has_one_of": bson.A{"test2"},
					"has_not":    bson.A{"test3"},
					"is_empty":   false,
				},
			},
			Pattern: oldpattern.InterfacePattern{
				StringArrayConditions: oldpattern.StringArrayConditions{
					HasEvery: types.OptionalStringArray{
						Set:   true,
						Value: []string{"test1"},
					},
					HasOneOf: types.OptionalStringArray{
						Set:   true,
						Value: []string{"test2"},
					},
					HasNot: types.OptionalStringArray{
						Set:   true,
						Value: []string{"test3"},
					},
					IsEmpty: types.OptionalBool{
						Set:   true,
						Value: false,
					},
				},
			},
		},
		{
			TestName: "test for equal int value",
			ExpectedUnmarshalled: bson.M{
				"pattern": int64(5),
			},
			Pattern: oldpattern.InterfacePattern{
				IntegerConditions: oldpattern.IntegerConditions{
					Equal: types.OptionalInt64{
						Set:   true,
						Value: 5,
					},
				},
			},
		},
		{
			TestName: "test for range int values",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: oldpattern.InterfacePattern{
				IntegerConditions: oldpattern.IntegerConditions{
					Gte: types.OptionalInt64{
						Set:   true,
						Value: 0,
					},
					Lte: types.OptionalInt64{
						Set:   true,
						Value: 2,
					},
				},
			},
		},
		{
			TestName: "test for nil value",
			ExpectedUnmarshalled: bson.M{
				"pattern": nil,
			},
			Pattern: oldpattern.InterfacePattern{
				EqualNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.InterfacePattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w interfacePatternWrapper
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
