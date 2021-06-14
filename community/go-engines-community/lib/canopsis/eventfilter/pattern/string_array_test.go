package pattern_test

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"testing"
)

type stringArrayPatternWrapper struct {
	Pattern pattern.StringArrayPattern `bson:"pattern"`
}

func TestValidHasStringArrayPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON patterns representing has string array pattern with single value", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{"test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeTrue)
				So(p.HasEvery.Value, ShouldContain, "test2")
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has string array pattern with several values", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{"test1", "test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeTrue)
				So(p.HasEvery.Value, ShouldContain, "test1")
				So(p.HasEvery.Value, ShouldContain, "test2")
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test3"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_not string array pattern with single value", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_not": []string{"test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeTrue)
				So(p.HasNot.Value, ShouldContain, "test2")
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_not string array pattern with several value", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_not": []string{"test1", "test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeTrue)
				So(p.HasNot.Value, ShouldContain, "test1")
				So(p.HasNot.Value, ShouldContain, "test2")
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_one_of string array pattern with single value", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_one_of": []string{"test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasOneOf.Set, ShouldBeTrue)
				So(p.HasOneOf.Value, ShouldContain, "test2")
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_one_of string array pattern with several values", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_one_of": []string{"test1", "test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasOneOf.Set, ShouldBeTrue)
				So(p.HasOneOf.Value, ShouldContain, "test1")
				So(p.HasOneOf.Value, ShouldContain, "test2")
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has and has_not string array patterns", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{"test3"},
				"has_not":   []string{"test1", "test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeTrue)
				So(p.HasEvery.Value, ShouldContain, "test3")
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeTrue)
				So(p.HasNot.Value, ShouldContain, "test1")
				So(p.HasNot.Value, ShouldContain, "test2")
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_every, has_one_of and has_not string array patterns", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_one_of": []string{"test1", "test2"},
				"has_not":    []string{"test1"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeTrue)
				So(p.HasOneOf.Value, ShouldContain, "test1")
				So(p.HasOneOf.Value, ShouldContain, "test2")
				So(p.HasNot.Set, ShouldBeTrue)
				So(p.HasNot.Value, ShouldContain, "test1")
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test1"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has and has_not string array patterns, which intersect each other, which should cause always false result", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{"test2"},
				"has_not":   []string{"test1", "test2"},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeTrue)
				So(p.HasEvery.Value, ShouldContain, "test2")
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeTrue)
				So(p.HasNot.Value, ShouldContain, "test1")
				So(p.HasNot.Value, ShouldContain, "test2")
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test1", "test2"}), ShouldBeFalse)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeFalse)
				So(p.Matches([]string{"test3"}), ShouldBeFalse)
				So(p.Matches([]string{}), ShouldBeFalse)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has string array empty pattern", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing has_not string array empty patterns", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_not": []string{},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing both has and has_not string array empty patterns", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{},
				"has_not":   []string{},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing both has and has_not string array empty patterns", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every": []string{},
				"has_not":   []string{},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing pattern with empty document", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
			p := w.Pattern

			Convey("The decoded values should be correct", func() {
				So(p.HasEvery.Set, ShouldBeFalse)
				So(p.HasEvery.Value, ShouldBeEmpty)
				So(p.HasNot.Set, ShouldBeFalse)
				So(p.HasNot.Value, ShouldBeEmpty)
				So(p.HasOneOf.Set, ShouldBeFalse)
				So(p.HasOneOf.Value, ShouldBeEmpty)
			})

			Convey("The pattern should accept the value of the string", func() {
				So(p.Matches([]string{"test1", "test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test2", "test3"}), ShouldBeTrue)
				So(p.Matches([]string{"test3"}), ShouldBeTrue)
				So(p.Matches([]string{}), ShouldBeTrue)
			})
		})
	})

	Convey("Given a valid BSON patterns representing pattern with extra fields", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_not": []string{"test2"},
				"extra":   "asd",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w stringArrayPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestStringArrayPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"has_every":  []string{"test1"},
				"has_one_of": []string{"test2"},
				"has_not":    []string{"test3"},
			},
		}
		mongoFilter := bson.M{
			"$all": []string{"test1"},
			"$in":  []string{"test2"},
			"$nin": []string{"test3"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w stringArrayPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestStringArrayPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              pattern.StringArrayPattern
	}{
		{
			TestName: "test for single has_every value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_every": bson.A{"test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasEvery: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test2"},
					},
				},
			},
		},
		{
			TestName: "test for several has_every values",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_every": bson.A{"test1", "test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasEvery: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1", "test2"},
					},
				},
			},
		},
		{
			TestName: "test for single has_not value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_not": bson.A{"test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasNot: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test2"},
					},
				},
			},
		},
		{
			TestName: "test for several has_not value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_not": bson.A{"test1", "test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasNot: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1", "test2"},
					},
				},
			},
		},
		{
			TestName: "test for several has_not value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_not": bson.A{"test1", "test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasNot: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1", "test2"},
					},
				},
			},
		},
		{
			TestName: "test for single has_one_of value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_one_of": bson.A{"test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasOneOf: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test2"},
					},
				},
			},
		},
		{
			TestName: "test for several has_one_of value",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_one_of": bson.A{"test1", "test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasOneOf: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1", "test2"},
					},
				},
			},
		},
		{
			TestName: "test for several has_not and has_every values",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"has_every": bson.A{"test3"},
					"has_not":   bson.A{"test1", "test2"},
				},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{
					HasEvery: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test3"},
					},
					HasNot: utils.OptionalStringArray{
						Set:   true,
						Value: []string{"test1", "test2"},
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.StringArrayPattern{
				StringArrayConditions: pattern.StringArrayConditions{},
			},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w stringArrayPatternWrapper
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
