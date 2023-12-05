package pattern_test

import (
	"errors"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEntity_Match(t *testing.T) {
	dataSets := getEntityMatchDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := data.pattern.Match(data.entity)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}

			ok, _, err = data.pattern.MatchWithRegexMatches(data.entity)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

func TestEntity_ToMongoQuery(t *testing.T) {
	dataSets := getEntityToMongoQueryDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			query, err := data.pattern.ToMongoQuery("entity")
			if !errors.Is(err, data.mongoQueryErr) {
				t.Errorf("expected error %v but got %v", data.mongoQueryErr, err)
			}
			if diff := pretty.Compare(query, data.mongoQueryResult); diff != "" {
				t.Errorf("unexpected result %s", diff)
			}
		})
	}
}

func TestEntity_ToSql(t *testing.T) {
	dataSets := getEntityToSqlDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			sql, err := data.pattern.ToSql("entity")
			if !errors.Is(err, data.sqlErr) {
				t.Errorf("expected error %v but got %v", data.sqlErr, err)
			}
			if sql != data.sqlResult {
				t.Errorf("expected\n%s\nbut got\n%s", sql, data.sqlResult)
			}
		})
	}
}

func BenchmarkEntity_Match_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name 2"),
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityMatch(b, cond, entity)
}

func BenchmarkEntity_Match_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "name",
		Condition: regexpCondition,
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityMatch(b, cond, entity)
}

func BenchmarkEntity_Match_HasOneOf(b *testing.B) {
	const condValueSize = 100
	const valueSize = 1000
	condValue := make([]string, condValueSize)
	for i := 0; i < condValueSize; i++ {
		condValue[i] = fmt.Sprintf("test-cond-val-%d", i)
	}
	value := make([]string, valueSize)
	for i := 0; i < valueSize; i++ {
		value[i] = fmt.Sprintf("test-val-%d", i)
	}

	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeStringArray,
		Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, condValue),
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:        "test",
				Description: "test",
				Value:       value,
			},
		},
	}

	benchmarkEntityMatch(b, cond, entity)
}

func BenchmarkEntity_Match_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test 2"),
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkEntityMatch(b, cond, entity)
}

func BenchmarkEntity_Match_Infos_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: regexpCondition,
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkEntityMatch(b, cond, entity)
}

func BenchmarkEntity_UnmarshalBsonAndMatch_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test name 2",
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityUnmarshalBsonAndMatch(b, cond, entity)
}

func BenchmarkEntity_UnmarshalBsonAndMatch_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	entity := types.Entity{
		Name: "test name",
	}

	benchmarkEntityUnmarshalBsonAndMatch(b, cond, entity)
}

func BenchmarkEntity_UnmarshalBsonAndMatch_HasOneOf(b *testing.B) {
	const condValueSize = 100
	const valueSize = 1000
	condValue := make([]string, condValueSize)
	for i := 0; i < condValueSize; i++ {
		condValue[i] = fmt.Sprintf("test-cond-val-%d", i)
	}
	value := make([]string, valueSize)
	for i := 0; i < valueSize; i++ {
		value[i] = fmt.Sprintf("test-val-%d", i)
	}

	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeStringArray,
		Condition: pattern.Condition{
			Type:  pattern.ConditionHasOneOf,
			Value: condValue,
		},
	}
	entity := types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:        "test",
				Description: "test",
				Value:       value,
			},
		},
	}

	benchmarkEntityUnmarshalBsonAndMatch(b, cond, entity)
}

func BenchmarkEntity_UnmarshalBsonAndMatch_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test 2",
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

	benchmarkEntityUnmarshalBsonAndMatch(b, cond, entity)
}

func BenchmarkEntity_UnmarshalBsonAndMatch_Infos_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
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

	benchmarkEntityUnmarshalBsonAndMatch(b, cond, entity)
}

func benchmarkEntityMatch(b *testing.B, fieldCond pattern.FieldCondition, entity types.Entity) {
	const size = 100
	p := make(pattern.Entity, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := p.Match(entity)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func benchmarkEntityUnmarshalBsonAndMatch(b *testing.B, fieldCond pattern.FieldCondition, entity types.Entity) {
	const size = 100
	p := make(pattern.Entity, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	type wrapper struct {
		Pattern pattern.Entity `bson:"pattern"`
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

		_, err = w.Pattern.Match(entity)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func getEntityMatchDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given empty pattern should match": {
			pattern: pattern.Entity{},
			entity: types.Entity{
				Name: "test name",
			},
			matchResult: true,
		},
		"given string field condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				Name: "test name",
			},
			matchResult: true,
		},
		"given string field condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				Name: "test another name",
			},
			matchResult: false,
		},
		"given string field condition and not string field should return error": {
			pattern: pattern.Entity{
				{
					{
						Field:     "impact_level",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity:   types.Entity{},
			matchErr: pattern.ErrWrongConditionValue,
		},
		"given string field condition and unknown field should return error": {
			pattern: pattern.Entity{
				{
					{
						Field:     "created",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity:   types.Entity{},
			matchErr: pattern.ErrUnsupportedField,
		},
		"given string info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given string info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test another name",
					},
				},
			},
			matchResult: false,
		},
		"given string info condition and not string info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       2,
					},
				},
			},
			matchErr:    pattern.ErrWrongConditionValue,
			matchResult: false,
		},
		"given string info condition and unknown info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity:      types.Entity{},
			matchResult: false,
		},
		"given exist info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given exist info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
		"given not exist info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given not exist info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
		"given string component info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given string component info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test another name",
					},
				},
			},
			matchResult: false,
		},
		"given string component info condition and not string info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       2,
					},
				},
			},
			matchErr:    pattern.ErrWrongConditionValue,
			matchResult: false,
		},
		"given string component info condition and unknown info should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			entity:      types.Entity{},
			matchResult: false,
		},
		"given exist component info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given exist component info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
		"given not exist component info condition should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_another_name": {
						Name:        "info_another_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: true,
		},
		"given not exist component info condition should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"info_name": {
						Name:        "info_name",
						Description: "test description",
						Value:       "test name",
					},
				},
			},
			matchResult: false,
		},
	}
}

func getEntityToMongoQueryDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given one condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
				}},
			}},
		},
		"given multiple conditions": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
					{"entity.category": bson.M{"$eq": "test category"}},
				}},
			}},
		},
		"given multiple groups": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.name": bson.M{"$eq": "test name"}},
				}},
				{"$and": []bson.M{
					{"entity.category": bson.M{"$eq": "test category"}},
				}},
			}},
		},
		"given invalid condition type": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
					},
				},
			},
			mongoQueryErr: pattern.ErrUnsupportedConditionType,
		},
		"given invalid condition value": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 2),
					},
				},
			},
			mongoQueryErr: pattern.ErrWrongConditionValue,
		},
		"given infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.infos.info_name.value": bson.M{"$eq": 3}},
				}},
			}},
		},
		"given component infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			mongoQueryResult: bson.M{"$or": []bson.M{
				{"$and": []bson.M{
					{"entity.component_infos.info_name.value": bson.M{"$eq": 3}},
				}},
			}},
		},
	}
}

func getEntityToSqlDataSets() map[string]entityDataSet {
	return map[string]entityDataSet{
		"given one condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name')",
		},
		"given multiple conditions": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name' AND entity.category = 'test category')",
		},
		"given multiple groups": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name"),
					},
				},
				{
					{
						Field:     "category",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test category"),
					},
				},
			},
			sqlResult: "(entity.name = 'test name') OR (entity.category = 'test category')",
		},
		"given invalid condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionIsNotOneOf, "test name"),
					},
				},
			},
			sqlErr: pattern.ErrWrongConditionValue,
		},
		"given infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(infos->'info_name') = 'number' AND (CASE WHEN jsonb_typeof(infos->'info_name') = 'number' THEN (infos->'info_name')::numeric END) = 3))`,
		},
		"given component infos condition": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.info_name",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
					},
				},
			},
			sqlResult: `((jsonb_typeof(component_infos->'info_name') = 'number' AND (CASE WHEN jsonb_typeof(component_infos->'info_name') = 'number' THEN (component_infos->'info_name')::numeric END) = 3))`,
		},
	}
}

type entityDataSet struct {
	pattern pattern.Entity
	entity  types.Entity

	matchErr    error
	matchResult bool

	mongoQueryErr    error
	mongoQueryResult bson.M

	sqlErr    error
	sqlResult string
}
