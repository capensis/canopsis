package oldpattern_test

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type entityPatternWrapper struct {
	Pattern oldpattern.EntityPattern `bson:"pattern"`
}

func TestEntityPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing an entity query", t, func() {
		mapPattern := bson.M{
			"pattern": bson.M{
				"_id":     "123456",
				"name":    "test",
				"enabled": true,
				"infos": bson.M{
					"toto": bson.M{
						"description": "test",
					},
					"tata": bson.M{
						"description": "test",
						"value":       "testtest",
					},
				},
				"component_infos": bson.M{
					"toto": bson.M{
						"description": "test",
					},
					"tata": bson.M{
						"description": "test",
						"value":       "testtest",
					},
				},
				"type":      bson.M{"regex_match": "abc-.*-def"},
				"component": bson.M{"regex_match": "abc-.*-def"},
				"category":  bson.M{"regex_match": "abc-.*-def"},
			},
		}
		mongoFilter := bson.M{
			"_id": bson.M{
				"$eq": "123456",
			},
			"name": bson.M{
				"$eq": "test",
			},
			"enabled": bson.M{
				"$eq": true,
			},
			"infos.toto.description": bson.M{
				"$eq": "test",
			},
			"infos.tata.description": bson.M{
				"$eq": "test",
			},
			"infos.tata.value": bson.M{
				"$eq": "testtest",
			},
			"component_infos.toto.description": bson.M{
				"$eq": "test",
			},
			"component_infos.tata.description": bson.M{
				"$eq": "test",
			},
			"component_infos.tata.value": bson.M{
				"$eq": "testtest",
			},
			"type": bson.M{
				"$regex": "abc-.*-def",
			},
			"component": bson.M{
				"$regex": "abc-.*-def",
			},
			"category": bson.M{
				"$regex": "abc-.*-def",
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

func TestNilEntityPatternToMongoDriverQuery(t *testing.T) {
	Convey("Given a valid BSON pattern representing a nil entity query", t, func() {
		mapPattern := bson.M{
			"pattern": nil,
		}
		mongoFilter := bson.M(nil)
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		var w entityPatternWrapper
		So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)
		p := w.Pattern

		Convey("The pattern should be converted to BSON as mongo query without error", func() {
			So(p.AsMongoDriverQuery(), ShouldResemble, mongoFilter)
		})
	})
}

// entityPatternListWrapper is a type that wraps an EntityPatternList into a
// struct. It is required to test the unmarshalling of an array into an
// EntityPatternList because bson.Unmarshal does not work when called with an
// array.
type entityPatternListWrapper struct {
	PatternList oldpattern.EntityPatternList `bson:"list"`
}

func TestValidEntityPatternListMongoDriver(t *testing.T) {
	Convey("Given a valid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				{
					"_id": "id1",
				},
				{
					"_id": "id2",
				},
				{
					"infos": bson.M{
						"toto": bson.M{
							"value": "tata",
						},
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

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
		mapPattern := bson.M{}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

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

func TestInvalidEntityPatternListMongoDriver(t *testing.T) {
	Convey("Given an invalid BSON pattern list", t, func() {
		mapPattern := bson.M{
			"list": []bson.M{
				bson.M{
					"_id": "id1",
				},
				bson.M{
					"_id": bson.M{
						">=": 3,
					},
				},
			},
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeTrue)
			So(p.IsValid(), ShouldBeFalse)
		})
	})

	Convey("Given a nil pattern list", t, func() {
		mapPattern := bson.M{
			"list": nil,
		}
		bsonPattern, err := bson.Marshal(mapPattern)
		So(err, ShouldBeNil)

		Convey("The pattern list should be decoded without errors", func() {
			var w entityPatternListWrapper
			So(bson.Unmarshal(bsonPattern, &w), ShouldBeNil)

			p := w.PatternList
			So(p.IsSet(), ShouldBeFalse)
			So(p.IsValid(), ShouldBeTrue)

			Convey("The pattern list should be converted to the right mongo query", func() {
				mongoFilter := bson.M{}
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
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.EntityPattern
	}{
		{
			TestName: "test for pattern",
			ExpectedUnmarshalled: bson.M{
				"pattern": bson.M{
					"_id":     "123456",
					"enabled": true,
					"infos": bson.M{
						"toto": bson.M{
							"description": "test",
						},
						"tata": bson.M{
							"description": "test",
							"value":       "testtest",
						},
					},
					"type": bson.M{"regex_match": "abc-.*-def"},
				},
			},
			Pattern: oldpattern.EntityPattern{
				ShouldNotBeNil: true,
				ShouldBeNil:    false,
				EntityFields: oldpattern.EntityFields{
					ID: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "123456",
							},
						},
					},
					Enabled: oldpattern.BoolPattern{
						OptionalBool: types.OptionalBool{
							Set:   true,
							Value: true,
						},
					},
					Infos: map[string]oldpattern.InfoPattern{
						"toto": {
							InfoFields: oldpattern.InfoFields{
								Description: oldpattern.StringPattern{
									StringConditions: oldpattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "test",
										},
									},
								},
							},
						},
						"tata": {
							InfoFields: oldpattern.InfoFields{
								Description: oldpattern.StringPattern{
									StringConditions: oldpattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "test",
										},
									},
								},
								Value: oldpattern.StringPattern{
									StringConditions: oldpattern.StringConditions{
										Equal: types.OptionalString{
											Set:   true,
											Value: "testtest",
										},
									},
								},
							},
						},
					},
					Type: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
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
			ExpectedUnmarshalled: bson.M{
				"pattern": nil,
			},
			Pattern: oldpattern.EntityPattern{
				ShouldBeNil: true,
			},
		},
		{
			TestName: "test for undefined",
			ExpectedUnmarshalled: bson.M{
				"pattern": primitive.Undefined{},
			},
			Pattern: oldpattern.EntityPattern{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w entityPatternWrapper
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

func TestEntityListPatternMarshalBSON(t *testing.T) {
	datasets := []struct {
		TestName             string
		ExpectedUnmarshalled bson.M
		Pattern              oldpattern.EntityPatternList
	}{
		{
			TestName: "test for list",
			ExpectedUnmarshalled: bson.M{
				"list": bson.A{
					bson.M{
						"_id": "id1",
					},
					bson.M{
						"_id": "id2",
					},
					bson.M{
						"infos": bson.M{
							"toto": bson.M{
								"value": "tata",
							},
						},
					},
				},
			},
			Pattern: oldpattern.EntityPatternList{
				Set:   true,
				Valid: true,
				Patterns: []oldpattern.EntityPattern{
					{
						ShouldNotBeNil: true,
						ShouldBeNil:    false,
						EntityFields: oldpattern.EntityFields{
							ID: oldpattern.StringPattern{
								StringConditions: oldpattern.StringConditions{
									Equal: types.OptionalString{
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
						EntityFields: oldpattern.EntityFields{
							ID: oldpattern.StringPattern{
								StringConditions: oldpattern.StringConditions{
									Equal: types.OptionalString{
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
						EntityFields: oldpattern.EntityFields{
							Infos: map[string]oldpattern.InfoPattern{
								"toto": {
									InfoFields: oldpattern.InfoFields{
										Value: oldpattern.StringPattern{
											StringConditions: oldpattern.StringConditions{
												Equal: types.OptionalString{
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
			ExpectedUnmarshalled: bson.M{
				"list": nil,
			},
			Pattern: oldpattern.EntityPatternList{},
		},
	}

	for _, dataset := range datasets {
		t.Run(dataset.TestName, func(t *testing.T) {
			var w entityPatternListWrapper
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

			if !reflect.DeepEqual(dataset.ExpectedUnmarshalled, unmarshalled) {
				t.Errorf("expected unmarshalled value = %v, got %v", dataset.ExpectedUnmarshalled["pattern"], unmarshalled["pattern"])
			}
		})
	}
}

func BenchmarkEntityPatternList_Matches_Equal(b *testing.B) {
	cond := oldpattern.EntityFields{
		Name: oldpattern.StringPattern{
			StringConditions: oldpattern.StringConditions{
				Equal: types.OptionalString{
					Set:   true,
					Value: "test name2",
				},
			},
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityPatternListMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_Matches_Regexp(b *testing.B) {
	testRegexp, err := utils.NewRegexExpression("^test .+name$")
	if err != nil {
		b.Fatalf("err is not expected: %s", err)
	}
	cond := oldpattern.EntityFields{
		Name: oldpattern.StringPattern{
			StringConditions: oldpattern.StringConditions{
				RegexMatch: types.OptionalRegexp{
					Set:   true,
					Value: testRegexp,
				},
			},
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityPatternListMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_Matches_Infos_Equal(b *testing.B) {
	cond := oldpattern.EntityFields{
		Infos: map[string]oldpattern.InfoPattern{
			"test": {
				InfoFields: oldpattern.InfoFields{
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test 2",
							},
						},
					},
				},
			},
		},
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkEntityPatternListMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_Matches_Infos_Regexp(b *testing.B) {
	testRegexp, err := utils.NewRegexExpression("^test .+val")
	if err != nil {
		b.Fatalf("err is not expected: %s", err)
	}
	cond := oldpattern.EntityFields{
		Infos: map[string]oldpattern.InfoPattern{
			"test": {
				InfoFields: oldpattern.InfoFields{
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
				},
			},
		},
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test val",
			},
		},
	}

	benchmarkEntityPatternListMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_UnmarshalBsonAndMatches_Equal(b *testing.B) {
	cond := oldpattern.EntityFields{
		Name: oldpattern.StringPattern{
			StringConditions: oldpattern.StringConditions{
				Equal: types.OptionalString{
					Set:   true,
					Value: "test name2",
				},
			},
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityPatternListUnmarshalBsonAndMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_UnmarshalBsonAndMatches_Regexp(b *testing.B) {
	testRegexp, err := utils.NewRegexExpression("^test .+name$")
	if err != nil {
		b.Fatalf("err is not expected: %s", err)
	}
	cond := oldpattern.EntityFields{
		Name: oldpattern.StringPattern{
			StringConditions: oldpattern.StringConditions{
				RegexMatch: types.OptionalRegexp{
					Set:   true,
					Value: testRegexp,
				},
			},
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityPatternListUnmarshalBsonAndMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_UnmarshalBsonAndMatches_Infos_Equal(b *testing.B) {
	cond := oldpattern.EntityFields{
		Infos: map[string]oldpattern.InfoPattern{
			"test": {
				InfoFields: oldpattern.InfoFields{
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							Equal: types.OptionalString{
								Set:   true,
								Value: "test 2",
							},
						},
					},
				},
			},
		},
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkEntityPatternListUnmarshalBsonAndMatches(b, cond, entity)
}

func BenchmarkEntityPatternList_UnmarshalBsonAndMatches_Infos_Regexp(b *testing.B) {
	testRegexp, err := utils.NewRegexExpression("^test .+val")
	if err != nil {
		b.Fatalf("err is not expected: %s", err)
	}
	cond := oldpattern.EntityFields{
		Infos: map[string]oldpattern.InfoPattern{
			"test": {
				InfoFields: oldpattern.InfoFields{
					Value: oldpattern.StringPattern{
						StringConditions: oldpattern.StringConditions{
							RegexMatch: types.OptionalRegexp{
								Set:   true,
								Value: testRegexp,
							},
						},
					},
				},
			},
		},
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test val",
			},
		},
	}

	benchmarkEntityPatternListUnmarshalBsonAndMatches(b, cond, entity)
}

func benchmarkEntityPatternListMatches(b *testing.B, cond oldpattern.EntityFields, entity types.Entity) {
	size := 100
	patterns := make([]oldpattern.EntityPattern, size)
	for i := 0; i < size; i++ {
		patterns[i] = oldpattern.EntityPattern{
			ShouldNotBeNil: true,
			EntityFields:   cond,
		}
	}

	p := oldpattern.EntityPatternList{
		Patterns: patterns,
		Set:      true,
		Valid:    true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Matches(&entity)
	}
}

func benchmarkEntityPatternListUnmarshalBsonAndMatches(b *testing.B, cond oldpattern.EntityFields, entity types.Entity) {
	size := 100
	patterns := make([]oldpattern.EntityPattern, size)
	for i := 0; i < size; i++ {
		patterns[i] = oldpattern.EntityPattern{
			ShouldNotBeNil: true,
			EntityFields:   cond,
		}
	}

	type wrapper struct {
		Pattern oldpattern.EntityPatternList `bson:"pattern"`
	}
	p := oldpattern.EntityPatternList{
		Patterns: patterns,
		Set:      true,
		Valid:    true,
	}

	bytes, err := bson.Marshal(wrapper{Pattern: p})
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var w wrapper
		err := bson.Unmarshal(bytes, &w)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}

		_ = w.Pattern.Matches(&entity)
	}
}
