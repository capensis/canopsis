package oldpattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/google/go-cmp/cmp"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
)

type eventPatternWrapper struct {
	Pattern oldpattern.EventPattern `bson:"pattern"`
}

func TestPatternUnmarshalMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern", t, func() {
		mapPattern := bson.M{
			"state": bson.M{
				">": 0,
				"<": 3,
			},
			"component": "component",
			"resource":  bson.M{"regex_match": "service-(?P<id>\\d+)"},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern should be decoded without errors", func() {
			var p oldpattern.EventPattern
			So(bson.Unmarshal(bsonPattern, &p), ShouldBeNil)

			Convey("The values of the pattern should be correct", func() {
				So(p.State.Gt.Set, ShouldBeTrue)
				So(p.State.Gt.Value, ShouldEqual, 0)
				So(p.State.Lt.Set, ShouldBeTrue)
				So(p.State.Lt.Value, ShouldEqual, 3)

				So(p.State.Equal.Set, ShouldBeFalse)
				So(p.State.Gte.Set, ShouldBeFalse)
				So(p.State.Lte.Set, ShouldBeFalse)

				So(p.Component.Equal.Set, ShouldBeTrue)
				So(p.Component.Equal.Value, ShouldEqual, "component")

				So(p.Component.RegexMatch.Set, ShouldBeFalse)

				So(p.Resource.RegexMatch.Set, ShouldBeTrue)
				So(p.Resource.RegexMatch.Value.String(), ShouldEqual, "service-(?P<id>\\d+)")

				So(p.Resource.Equal.Set, ShouldBeFalse)
			})
		})
	})

	Convey("Given a BSON pattern with an invalid regexp", t, func() {
		mapPattern := bson.M{
			"component": bson.M{
				"regex_match": "abc-(.*-def",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p oldpattern.EventPattern
			So(bson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected entity field", t, func() {
		mapPattern := bson.M{
			"current_entity": bson.M{
				"unexpected_field": "",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p oldpattern.EventPattern
			So(bson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with an unexpected info field", t, func() {
		mapPattern := bson.M{
			"current_entity": bson.M{
				"infos": bson.M{
					"info_name": bson.M{
						"unexpected_field": "",
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p oldpattern.EventPattern
			So(bson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})

	Convey("Given a BSON pattern with a mistyped info field", t, func() {
		mapPattern := bson.M{
			"current_entity": bson.M{
				"infos": bson.M{
					"info_name": bson.M{
						"value": 3,
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("Decoding the pattern returns an error", func() {
			var p oldpattern.EventPattern
			So(bson.Unmarshal(bsonPattern, &p), ShouldNotBeNil)
		})
	})
}

func TestEventPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("service-(?P<id>\\d+)")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.EventPattern
	}{
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"state": bson.M{
						">": int64(0),
						"<": int64(3),
					},
					"component": "component",
					"resource":  bson.M{"regex_match": "service-(?P<id>\\d+)"},
					"current_entity": bson.M{
						"infos": bson.M{
							"info_name": bson.M{
								"value": "info_value",
							},
						},
					},
					"extra_info": bson.M{
						"has_every": bson.A{"test1", "test2"},
						"has_not":   bson.A{"test3"},
					},
					"extra_info_2": bson.M{
						"is_empty": false,
					},
				},
			},
			Pattern: oldpattern.EventPattern{
				State: oldpattern.IntegerPattern{
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
				Component: oldpattern.StringPattern{
					StringConditions: oldpattern.StringConditions{
						Equal: types.OptionalString{
							Set:   true,
							Value: "component",
						},
					},
				},
				Resource: oldpattern.StringPattern{
					StringConditions: oldpattern.StringConditions{
						RegexMatch: types.OptionalRegexp{
							Set:   true,
							Value: testRegexp,
						},
					},
				},
				Entity: oldpattern.EntityPattern{
					EntityFields: oldpattern.EntityFields{
						Infos: map[string]oldpattern.InfoPattern{
							"info_name": {
								InfoFields: oldpattern.InfoFields{
									Value: oldpattern.StringPattern{
										StringConditions: oldpattern.StringConditions{
											Equal: types.OptionalString{
												Set:   true,
												Value: "info_value",
											},
										},
									},
								},
							},
						},
					},
				},
				ExtraInfos: map[string]oldpattern.InterfacePattern{
					"extra_info": {
						StringArrayConditions: oldpattern.StringArrayConditions{
							HasEvery: types.OptionalStringArray{
								Set:   true,
								Value: []string{"test1", "test2"},
							},
							HasNot: types.OptionalStringArray{
								Set:   true,
								Value: []string{"test3"},
							},
						},
					},
					"extra_info_2": {
						StringArrayConditions: oldpattern.StringArrayConditions{
							IsEmpty: types.OptionalBool{
								Set:   true,
								Value: false,
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{},
			},
			Pattern: oldpattern.EventPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w eventPatternWrapper
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

// eventPatternListWrapper is a type that wraps an EventPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// EventPatternList because bson.Unmarshal does not work when called with an
// array.
type eventPatternListWrapper struct {
	PatternList oldpattern.EventPatternList `bson:"list"`
}

func TestValidEventPatternListMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				{
					"component": "component1",
				},
				{
					"component": "component2",
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := bson.M{}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding events", func() {
				event1 := types.Event{
					Component: "component1",
				}
				So(p.Matches(event1), ShouldBeTrue)

				event2 := types.Event{
					Component: "component2",
				}
				So(p.Matches(event2), ShouldBeTrue)

				event3 := types.Event{
					Component: "component3",
				}
				So(p.Matches(event3), ShouldBeTrue)
			})
		})
	})
}

func TestInvalidEventPatternListMongoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				{
					"component": "component1",
				},
				{
					"component": bson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w eventPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})
}

func TestEventPatternListMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("service-(?P<id>\\d+)")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.EventPatternList
	}{
		{
			TestName: "test for pattern list",
			ExpectedUnmarshalled: bson.M{
				"list": bson.A{
					bson.M{
						"state": bson.M{
							">": int64(0),
							"<": int64(3),
						},
						"component": "component",
						"resource":  bson.M{"regex_match": "service-(?P<id>\\d+)"},
						"current_entity": bson.M{
							"infos": bson.M{
								"info_name": bson.M{
									"value": "info_value",
								},
							},
						},
						"extra_info": bson.M{
							"has_every": bson.A{"test1", "test2"},
							"has_not":   bson.A{"test3"},
						},
					},
					bson.M{
						"state": bson.M{
							">": int64(0),
							"<": int64(3),
						},
						"component": "component",
						"resource":  bson.M{"regex_match": "service-(?P<id>\\d+)"},
						"current_entity": bson.M{
							"infos": bson.M{
								"info_name": bson.M{
									"value": "info_value",
								},
							},
						},
						"extra_info": bson.M{
							"has_every": bson.A{"test1", "test2"},
							"has_not":   bson.A{"test3"},
						},
						"extra_info_2": "test_info",
						"extra_info_3": bson.M{
							"is_empty": true,
						},
					},
				},
			},
			Pattern: oldpattern.EventPatternList{
				Set:   true,
				Valid: true,
				Patterns: []oldpattern.EventPattern{
					{
						State: oldpattern.IntegerPattern{
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
						Component: oldpattern.StringPattern{
							StringConditions: oldpattern.StringConditions{
								Equal: types.OptionalString{
									Set:   true,
									Value: "component",
								},
							},
						},
						Resource: oldpattern.StringPattern{
							StringConditions: oldpattern.StringConditions{
								RegexMatch: types.OptionalRegexp{
									Set:   true,
									Value: testRegexp,
								},
							},
						},
						Entity: oldpattern.EntityPattern{
							EntityFields: oldpattern.EntityFields{
								Infos: map[string]oldpattern.InfoPattern{
									"info_name": {
										InfoFields: oldpattern.InfoFields{
											Value: oldpattern.StringPattern{
												StringConditions: oldpattern.StringConditions{
													Equal: types.OptionalString{
														Set:   true,
														Value: "info_value",
													},
												},
											},
										},
									},
								},
							},
						},
						ExtraInfos: map[string]oldpattern.InterfacePattern{
							"extra_info": {
								StringArrayConditions: oldpattern.StringArrayConditions{
									HasEvery: types.OptionalStringArray{
										Set:   true,
										Value: []string{"test1", "test2"},
									},
									HasNot: types.OptionalStringArray{
										Set:   true,
										Value: []string{"test3"},
									},
								},
							},
						},
					},
					{
						State: oldpattern.IntegerPattern{
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
						Component: oldpattern.StringPattern{
							StringConditions: oldpattern.StringConditions{
								Equal: types.OptionalString{
									Set:   true,
									Value: "component",
								},
							},
						},
						Resource: oldpattern.StringPattern{
							StringConditions: oldpattern.StringConditions{
								RegexMatch: types.OptionalRegexp{
									Set:   true,
									Value: testRegexp,
								},
							},
						},
						Entity: oldpattern.EntityPattern{
							EntityFields: oldpattern.EntityFields{
								Infos: map[string]oldpattern.InfoPattern{
									"info_name": {
										InfoFields: oldpattern.InfoFields{
											Value: oldpattern.StringPattern{
												StringConditions: oldpattern.StringConditions{
													Equal: types.OptionalString{
														Set:   true,
														Value: "info_value",
													},
												},
											},
										},
									},
								},
							},
						},
						ExtraInfos: map[string]oldpattern.InterfacePattern{
							"extra_info": {
								StringArrayConditions: oldpattern.StringArrayConditions{
									HasEvery: types.OptionalStringArray{
										Set:   true,
										Value: []string{"test1", "test2"},
									},
									HasNot: types.OptionalStringArray{
										Set:   true,
										Value: []string{"test3"},
									},
								},
							},
							"extra_info_2": {
								StringConditions: oldpattern.StringConditions{
									Equal: types.OptionalString{
										Set:   true,
										Value: "test_info",
									},
								},
							},
							"extra_info_3": {
								StringArrayConditions: oldpattern.StringArrayConditions{
									IsEmpty: types.OptionalBool{
										Set:   true,
										Value: true,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: bson.M{
				"list": nil,
			},
			Pattern: oldpattern.EventPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w eventPatternListWrapper
			w.PatternList = dataset.Pattern

			resultBson, err := bson.Marshal(w)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			var unmarshalled bson.M
			err = bson.Unmarshal(resultBson, &unmarshalled)
			if err != nil {
				t.Fatalf("err is not expected: %s", err)
			}

			// after unmarshall json integer becomes float64
			// it would fail in equality comparision so need to specify compare func
			compareIntFloat := cmp.Comparer(func(x, y interface{}) bool {
				return reflect.ValueOf(x).Int() == int64(reflect.ValueOf(y).Float())
			})

			// This option handles slices and maps of any type.
			opt := cmp.FilterValues(func(x, y interface{}) bool {
				vx, vy := reflect.ValueOf(x), reflect.ValueOf(y)
				return vx.IsValid() && vy.IsValid() && vx.Kind() == reflect.Int64 && vy.Kind() == reflect.Float64
			}, compareIntFloat)

			if !cmp.Equal(dataset.ExpectedUnmarshalled, unmarshalled, opt) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled, unmarshalled)
			}
		})
	}
}
