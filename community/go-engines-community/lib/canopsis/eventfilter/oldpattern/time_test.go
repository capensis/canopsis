package oldpattern_test

import (
	"reflect"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type timePatternWrapper struct {
	Pattern oldpattern.TimePattern `bson:"pattern"`
}

type timeRefPatternWrapper struct {
	Pattern oldpattern.TimeRefPattern `bson:"pattern"`
}

func timestamp(value int64) types.CpsTime {
	return types.CpsTime{Time: time.Unix(int64(value), 0)}
}

func TestValidEqualTimePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": 7,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w timePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				So(p.Matches(timestamp(7)), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(timestamp(-1)), ShouldBeFalse)
				So(p.Matches(timestamp(0)), ShouldBeFalse)
				So(p.Matches(timestamp(3)), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualNilTimeRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with nil", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsTime
			var w timeRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				value = timestamp(-1)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(0)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(3)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(7)
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidEqualTimeRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing equality with an integer", t, func() {
		mapPattern := bson.M{
			"pattern": 7,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsTime
			var w timeRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				value = timestamp(7)
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The pattern should reject any other value", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = timestamp(-1)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(0)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(3)
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesTimePatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var w timePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				So(p.Matches(timestamp(1)), ShouldBeTrue)
				So(p.Matches(timestamp(2)), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(timestamp(-1)), ShouldBeFalse)
				So(p.Matches(timestamp(0)), ShouldBeFalse)
				So(p.Matches(timestamp(3)), ShouldBeFalse)
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
			var w timePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				So(p.Matches(timestamp(0)), ShouldBeTrue)
				So(p.Matches(timestamp(1)), ShouldBeTrue)
				So(p.Matches(timestamp(2)), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(timestamp(-1)), ShouldBeFalse)
				So(p.Matches(timestamp(3)), ShouldBeFalse)
				So(p.Matches(timestamp(4)), ShouldBeFalse)
			})
		})
	})
}

func TestValidInequalitiesTimeRefPatternMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern representing an integer range", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				">": 0,
				"<": 3,
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var value types.CpsTime
			var w timeRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				value = timestamp(1)
				So(p.Matches(&value), ShouldBeTrue)
				value = timestamp(2)
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = timestamp(-1)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(0)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(3)
				So(p.Matches(&value), ShouldBeFalse)
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
			var value types.CpsTime
			var w timeRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
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
				value = timestamp(0)
				So(p.Matches(&value), ShouldBeTrue)
				value = timestamp(1)
				So(p.Matches(&value), ShouldBeTrue)
				value = timestamp(2)
				So(p.Matches(&value), ShouldBeTrue)
			})

			Convey("The filter should reject values that are out of the range", func() {
				So(p.Matches(nil), ShouldBeFalse)
				value = timestamp(-1)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(3)
				So(p.Matches(&value), ShouldBeFalse)
				value = timestamp(4)
				So(p.Matches(&value), ShouldBeFalse)
			})
		})
	})
}

func TestInvalidTimePatternMongoDriver(t *testing.T) {
	Convey("Given a nil pattern", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w timePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := bson.M{
			"pattern": "string",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w timePatternWrapper
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
			var w timePatternWrapper
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
			var w timePatternWrapper
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
			var w timePatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestInvalidTimeRefPatternMongoDriver(t *testing.T) {
	Convey("Given a BSON pattern with a string", t, func() {
		mapPattern := bson.M{
			"pattern": "string",
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var w timeRefPatternWrapper
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
			var w timeRefPatternWrapper
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
			var w timeRefPatternWrapper
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
			var w timeRefPatternWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldNotBeNil)
		})
	})
}

func TestTimePatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.TimePattern
	}{
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: bson.M{
				"pattern": int64(7),
			},
			Pattern: oldpattern.TimePattern{
				IntegerPattern: oldpattern.IntegerPattern{
					IntegerConditions: oldpattern.IntegerConditions{
						Equal: types.OptionalInt64{
							Set:   true,
							Value: 7,
						},
					},
				},
			},
		},
		{
			TestName: "test for closed interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: oldpattern.TimePattern{
				IntegerPattern: oldpattern.IntegerPattern{
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
		},
		{
			TestName: "test for open interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">": int64(0),
					"<": int64(3),
				},
			},
			Pattern: oldpattern.TimePattern{
				IntegerPattern: oldpattern.IntegerPattern{
					IntegerConditions: oldpattern.IntegerConditions{
						Gt: types.OptionalInt64{
							Set:   true,
							Value: 0,
						},
						Lt: types.OptionalInt64{
							Set:   true,
							Value: 3,
						},
					},
				},
			},
		},
		{
			TestName: "test for half-open interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">=": int64(2),
					"<":  int64(3),
				},
			},
			Pattern: oldpattern.TimePattern{
				IntegerPattern: oldpattern.IntegerPattern{
					IntegerConditions: oldpattern.IntegerConditions{
						Gte: types.OptionalInt64{
							Set:   true,
							Value: 2,
						},
						Lt: types.OptionalInt64{
							Set:   true,
							Value: 3,
						},
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.TimePattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w timePatternWrapper
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

func TestTimeRefPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.TimeRefPattern
	}{
		{
			TestName: "test for equal",
			ExpectedUnmarshalled: bson.M{
				"pattern": int64(7),
			},
			Pattern: oldpattern.TimeRefPattern{
				IntegerRefPattern: oldpattern.IntegerRefPattern{
					IntegerPattern: oldpattern.IntegerPattern{
						IntegerConditions: oldpattern.IntegerConditions{
							Equal: types.OptionalInt64{
								Set:   true,
								Value: 7,
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for closed interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">=": int64(0),
					"<=": int64(2),
				},
			},
			Pattern: oldpattern.TimeRefPattern{
				IntegerRefPattern: oldpattern.IntegerRefPattern{
					IntegerPattern: oldpattern.IntegerPattern{
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
			},
		},
		{
			TestName: "test for open interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">": int64(0),
					"<": int64(3),
				},
			},
			Pattern: oldpattern.TimeRefPattern{
				IntegerRefPattern: oldpattern.IntegerRefPattern{
					IntegerPattern: oldpattern.IntegerPattern{
						IntegerConditions: oldpattern.IntegerConditions{
							Gt: types.OptionalInt64{
								Set:   true,
								Value: 0,
							},
							Lt: types.OptionalInt64{
								Set:   true,
								Value: 3,
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for half-open interval",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					">=": int64(2),
					"<":  int64(3),
				},
			},
			Pattern: oldpattern.TimeRefPattern{
				IntegerRefPattern: oldpattern.IntegerRefPattern{
					IntegerPattern: oldpattern.IntegerPattern{
						IntegerConditions: oldpattern.IntegerConditions{
							Gte: types.OptionalInt64{
								Set:   true,
								Value: 2,
							},
							Lt: types.OptionalInt64{
								Set:   true,
								Value: 3,
							},
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
			Pattern: oldpattern.TimeRefPattern{
				IntegerRefPattern: oldpattern.IntegerRefPattern{
					EqualNil: true,
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.TimeRefPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w timeRefPatternWrapper
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
