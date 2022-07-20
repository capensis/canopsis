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

type stringPatternWrapper struct {
	Pattern oldpattern.StringPattern `bson:"pattern"`
}

type stringRefPatternWrapper struct {
	Pattern oldpattern.StringRefPattern `bson:"pattern"`
}

func TestEqualStringPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := bson.M{
			"pattern": "toto",
		}
		mongoFilter := bson.M{
			"$eq": "toto",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := bson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualNilStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M(nil)
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEqualStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an equality", t, func() {
		mapPattern := bson.M{
			"pattern": "toto",
		}
		mongoFilter := bson.M{
			"$eq": "toto",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestRegexStringRefPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regex", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		mongoFilter := bson.M{
			"$regex": "abc-.*-def",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringRefPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestValidEqualStringPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := bson.M{
			"pattern": "value",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidRegexPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidEqualNilStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidEqualStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with a string", t, func() {
		mapPattern := bson.M{
			"pattern": "value",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestValidRegexRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing a regular expression", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-.*-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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
		mapPattern := bson.M{
			"pattern": bson.M{"regex_match": "abc-(?P<sub>.*)-def"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			var m oldpattern.RegexMatches

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

func TestInvalidStringPatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": 3,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
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
			var w stringPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidStringRefPatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": 3,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected field", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{"unexpected_field": "test1"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
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
			var w stringRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
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
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.StringPattern
	}{
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: bson.M{
				"pattern": "test",
			},
			Pattern: oldpattern.StringPattern{
				StringConditions: oldpattern.StringConditions{
					Equal: types.OptionalString{
						Set:   true,
						Value: "test",
					},
				},
			},
		},
		{
			TestName: "test for regexp",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"regex_match": "abc-.*-def",
				},
			},
			Pattern: oldpattern.StringPattern{
				StringConditions: oldpattern.StringConditions{
					RegexMatch: types.OptionalRegexp{
						Set:   true,
						Value: testRegexp,
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.StringPattern{
				StringConditions: oldpattern.StringConditions{},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w stringPatternWrapper
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

func TestStringRefPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.StringRefPattern
	}{
		{
			TestName: "test for ref equal",
			ExpectedUnmarshalled: bson.M{
				"pattern": "test",
			},
			Pattern: oldpattern.StringRefPattern{
				EqualNil: false,
				StringPattern: oldpattern.StringPattern{
					StringConditions: oldpattern.StringConditions{
						Equal: types.OptionalString{
							Set:   true,
							Value: "test",
						},
					},
				},
			},
		},
		{
			TestName: "test for ref regexp",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"regex_match": "abc-.*-def",
				},
			},
			Pattern: oldpattern.StringRefPattern{
				EqualNil: false,
				StringPattern: oldpattern.StringPattern{
					StringConditions: oldpattern.StringConditions{
						RegexMatch: types.OptionalRegexp{
							Set:   true,
							Value: testRegexp,
						},
					},
				},
			},
		},
		{
			TestName: "test for ref nil",
			ExpectedUnmarshalled: bson.M{
				"pattern": nil,
			},
			Pattern: oldpattern.StringRefPattern{
				EqualNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.StringRefPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w stringRefPatternWrapper
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
