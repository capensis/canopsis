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

type entityPatternWrapper struct {
	Pattern pattern.EntityPattern `bson:"pattern"`
}

func TestEntityPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an entity query", t, func() {
		mapPattern := mgobson.M{
			"pattern": mgobson.M{
				"_id":     "123456",
				"name":    "test",
				"enabled": true,
				"infos": mgobson.M{
					"toto": mgobson.M{
						"description": "test",
					},
					"tata": mgobson.M{
						"description": "test",
						"value":       "testtest",
					},
				},
				"type":      mgobson.M{"regex_match": "abc-.*-def"},
				"component": mgobson.M{"regex_match": ".*MYSERVER.*"},
			},
		}
		mongoFilter := mgobson.M{
			"_id": mgobson.M{
				"$eq": "123456",
			},
			"name": mgobson.M{
				"$eq": "test",
			},
			"enabled": mgobson.M{
				"$eq": true,
			},
			"infos.toto.description": mgobson.M{
				"$eq": "test",
			},
			"infos.tata.description": mgobson.M{
				"$eq": "test",
			},
			"infos.tata.value": mgobson.M{
				"$eq": "testtest",
			},
			"type": mgobson.M{
				"$regex": "abc-.*-def",
			},
			"component": mgobson.M{
				"$regex": ".*MYSERVER.*",
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestEntityPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an entity query", t, func() {
		mapPattern := mongobson.M{
			"pattern": mongobson.M{
				"_id":     "123456",
				"name":    "test",
				"enabled": true,
				"infos": mongobson.M{
					"toto": mongobson.M{
						"description": "test",
					},
					"tata": mongobson.M{
						"description": "test",
						"value":       "testtest",
					},
				},
				"component_infos": mongobson.M{
					"toto": mongobson.M{
						"description": "test",
					},
					"tata": mongobson.M{
						"description": "test",
						"value":       "testtest",
					},
				},
				"type":      mongobson.M{"regex_match": "abc-.*-def"},
				"component": mongobson.M{"regex_match": "abc-.*-def"},
			},
		}
		mongoFilter := mongobson.M{
			"_id": mongobson.M{
				"$eq": "123456",
			},
			"name": mongobson.M{
				"$eq": "test",
			},
			"enabled": mongobson.M{
				"$eq": true,
			},
			"infos.toto.description": mongobson.M{
				"$eq": "test",
			},
			"infos.tata.description": mongobson.M{
				"$eq": "test",
			},
			"infos.tata.value": mongobson.M{
				"$eq": "testtest",
			},
			"component_infos.toto.description": mongobson.M{
				"$eq": "test",
			},
			"component_infos.tata.description": mongobson.M{
				"$eq": "test",
			},
			"component_infos.tata.value": mongobson.M{
				"$eq": "testtest",
			},
			"type": mongobson.M{
				"$regex": "abc-.*-def",
			},
			"component": mongobson.M{
				"$regex": "abc-.*-def",
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilEntityPatternToMongoQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil entity query", t, func() {
		mapPattern := mgobson.M{
			"pattern": nil,
		}
		mongoFilter := mgobson.M(nil)
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilEntityPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil entity query", t, func() {
		mapPattern := mongobson.M{
			"pattern": nil,
		}
		mongoFilter := mongobson.M(nil)
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

// entityPatternListWrapper is a type that wraps an EntityPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// EntityPatternList because mgobson.Unmarshal does not work when called with an
// array.
type entityPatternListWrapper struct {
	PatternList pattern.EntityPatternList `bson:"list"`
}

func TestValidEntityPatternListMgoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"_id": "id1",
				},
				mgobson.M{
					"_id": "id2",
				},
				mgobson.M{
					"infos": mgobson.M{
						"toto": mgobson.M{
							"value": "tata",
						},
					},
				},
				mgobson.M{
					"component": mgobson.M{
						"regex_match": ".*MYSERVER.*",
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding entities", func() {
				entity1 := types.Entity{
					ID: "id1",
				}
				So(p.Matches(&entity1), ShouldBeTrue)

				entity2 := types.Entity{
					ID: "id2",
				}
				So(p.Matches(&entity2), ShouldBeTrue)

				testInfo := types.Info{Value: "tata"}
				entityInfos := types.Entity{
					Infos: map[string]types.Info{"toto": testInfo},
				}
				So(p.Matches(&entityInfos), ShouldBeTrue)

				entity3 := types.Entity{
					ID:        "id3",
					Component: "empty",
				}
				So(p.Matches(&entity3), ShouldBeFalse)

				entity4 := types.Entity{
					Component: "zzzMYSERVERzzz",
				}
				So(p.Matches(&entity4), ShouldBeTrue)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mgobson.M{}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding entities", func() {
				entity1 := types.Entity{
					ID: "id1",
				}
				So(p.Matches(&entity1), ShouldBeTrue)

				entity2 := types.Entity{
					ID: "id2",
				}
				So(p.Matches(&entity2), ShouldBeTrue)

				entity3 := types.Entity{
					ID: "id3",
				}
				So(p.Matches(&entity3), ShouldBeTrue)
			})
		})
	})
}

func TestValidEntityPatternListMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				{
					"_id": "id1",
				},
				{
					"_id": "id2",
				},
				{
					"infos": mongobson.M{
						"toto": mongobson.M{
							"value": "tata",
						},
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding entities", func() {
				entity1 := types.Entity{
					ID: "id1",
				}
				So(p.Matches(&entity1), ShouldBeTrue)

				entity2 := types.Entity{
					ID: "id2",
				}
				So(p.Matches(&entity2), ShouldBeTrue)

				testInfo := types.Info{Value: "tata"}
				entityInfos := types.Entity{
					Infos: map[string]types.Info{"toto": testInfo},
				}
				So(p.Matches(&entityInfos), ShouldBeTrue)

				entity3 := types.Entity{
					ID: "id3",
				}
				So(p.Matches(&entity3), ShouldBeFalse)
			})
		})
	})

	Convey("Given an unset pattern list", t, func() {
		mapPattern := mongobson.M{}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should match the corresponding entities", func() {
				entity1 := types.Entity{
					ID: "id1",
				}
				So(p.Matches(&entity1), ShouldBeTrue)

				entity2 := types.Entity{
					ID: "id2",
				}
				So(p.Matches(&entity2), ShouldBeTrue)

				entity3 := types.Entity{
					ID: "id3",
				}
				So(p.Matches(&entity3), ShouldBeTrue)
			})
		})
	})
}

func TestInvalidEntityPatternListMgoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": []mgobson.M{
				mgobson.M{
					"_id": "id1",
				},
				mgobson.M{
					"_id": mgobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := mgobson.M{
			"list": nil,
		}
		bsonPattern, err := mgobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mgobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mgobson.M{}
				So(p.AsMongoQuery(), ShouldResemble, mongoFilter)
			})
		})
	})
}

func TestInvalidEntityPatternListMongoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": []mongobson.M{
				mongobson.M{
					"_id": "id1",
				},
				mongobson.M{
					"_id": mongobson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := mongobson.M{
			"list": nil,
		}
		bsonPattern, err := mongobson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(mongobson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := mongobson.M{}
				So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
			})
		})
	})
}

func TestEntityPatternMarshalBSON(t *testing.T) {
	testRegexp, err := utils.NewRegexExpression("abc-.*-def")
	if err != nil {
		t.Fatalf("err is not expected: %s", err)
	}

	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.EntityPattern
	}{
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": mongobson.M{
					"_id":     "123456",
					"enabled": true,
					"infos": mongobson.M{
						"toto": mongobson.M{
							"description": "test",
						},
						"tata": mongobson.M{
							"description": "test",
							"value":       "testtest",
						},
					},
					"type": mongobson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: pattern.EntityPattern{
				ShouldNotBeNil: true,
				ShouldBeNil:    false,
				EntityFields: pattern.EntityFields{
					ID: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							Equal: utils.OptionalString{
								Set:   true,
								Value: "123456",
							},
						},
					},
					Enabled: pattern.BoolPattern{
						OptionalBool: utils.OptionalBool{
							Set:   true,
							Value: true,
						},
					},
					Infos: map[string]pattern.InfoPattern{
						"toto": {
							InfoFields: pattern.InfoFields{
								Description: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: utils.OptionalString{
											Set:   true,
											Value: "test",
										},
									},
								},
							},
						},
						"tata": {
							InfoFields: pattern.InfoFields{
								Description: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: utils.OptionalString{
											Set:   true,
											Value: "test",
										},
									},
								},
								Value: pattern.StringPattern{
									StringConditions: pattern.StringConditions{
										Equal: utils.OptionalString{
											Set:   true,
											Value: "testtest",
										},
									},
								},
							},
						},
					},
					Type: pattern.StringPattern{
						StringConditions: pattern.StringConditions{
							RegexMatch: utils.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
				},
			},
		},
		{
			TestName: "test for nil",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": nil,
			},
			Pattern: pattern.EntityPattern{
				ShouldBeNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: mongobson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: pattern.EntityPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w entityPatternWrapper
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

func TestEntityListPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled mongobson.M
		Pattern              pattern.EntityPatternList
	}{
		{
			TestName: "test for list",
			ExpectedUnmarshalled: mongobson.M{
				"list": mongobson.A{
					mongobson.M{
						"_id": "id1",
					},
					mongobson.M{
						"_id": "id2",
					},
					mongobson.M{
						"infos": mongobson.M{
							"toto": mongobson.M{
								"value": "tata",
							},
						},
					},
				},
			},
			Pattern: pattern.EntityPatternList{
				Set:   true,
				Valid: true,
				Patterns: []pattern.EntityPattern{
					{
						ShouldNotBeNil: true,
						ShouldBeNil:    false,
						EntityFields: pattern.EntityFields{
							ID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
										Value: "id1",
									},
								},
							},
						},
					},
					{
						ShouldNotBeNil: true,
						ShouldBeNil:    false,
						EntityFields: pattern.EntityFields{
							ID: pattern.StringPattern{
								StringConditions: pattern.StringConditions{
									Equal: utils.OptionalString{
										Set:   true,
										Value: "id2",
									},
								},
							},
						},
					},
					{
						ShouldNotBeNil: true,
						ShouldBeNil:    false,
						EntityFields: pattern.EntityFields{
							Infos: map[string]pattern.InfoPattern{
								"toto": {
									InfoFields: pattern.InfoFields{
										Value: pattern.StringPattern{
											StringConditions: pattern.StringConditions{
												Equal: utils.OptionalString{
													Set:   true,
													Value: "tata",
												},
											},
										},
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
			ExpectedUnmarshalled: mongobson.M{
				"list": nil,
			},
			Pattern: pattern.EntityPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w entityPatternListWrapper
			w.PatternList = dataset.Pattern

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
