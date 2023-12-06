package match_test

import (
	"errors"
	"fmt"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMatchEntityPattern(t *testing.T) {
	dataSets := getMatchEntityPatternDataSets()

	for name, data := range dataSets {
		t.Run(name, func(t *testing.T) {
			ok, err := match.MatchEntityPattern(data.pattern, &data.entity)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}

			ok, _, err = match.MatchEntityPatternWithRegexMatches(data.pattern, &data.entity)
			if !errors.Is(err, data.matchErr) {
				t.Errorf("expected error %v but got %v", data.matchErr, err)
			}
			if ok != data.matchResult {
				t.Errorf("expected result %v but got %v", data.matchResult, ok)
			}
		})
	}
}

type entityDataSet struct {
	pattern     pattern.Entity
	entity      types.Entity
	matchErr    error
	matchResult bool
}

func getMatchEntityPatternDataSets() map[string]entityDataSet {
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

func BenchmarkMatchEntityPattern_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "name",
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test name 2"),
	}
	entity := &types.Entity{
		Name: "test name",
	}

	benchmarkMatchEntityPattern(b, cond, entity)
}

func BenchmarkMatchEntityPattern_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "name",
		Condition: regexpCondition,
	}
	entity := &types.Entity{
		Name: "test name",
	}

	benchmarkMatchEntityPattern(b, cond, entity)
}

func BenchmarkMatchEntityPattern_HasOneOf(b *testing.B) {
	const condValueSize = 10
	const valueSize = 10
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
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:        "test",
				Description: "test",
				Value:       value,
			},
		},
	}

	benchmarkMatchEntityPattern(b, cond, entity)
}

func BenchmarkMatchEntityPattern_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test 2"),
	}
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkMatchEntityPattern(b, cond, entity)
}

func BenchmarkMatchEntityPattern_Infos_Regexp(b *testing.B) {
	regexpCondition, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, "^test .+name$")
	if err != nil {
		b.Fatalf("unexpected error %v", err)
	}
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: regexpCondition,
	}
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkMatchEntityPattern(b, cond, entity)
}

func BenchmarkMatchEntityPattern_UnmarshalBson_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test name 2",
		},
	}
	entity := &types.Entity{
		Name: "test name",
	}

	benchmarkMatchEntityPatternUnmarshalBson(b, cond, entity)
}

func BenchmarkMatchEntityPattern_UnmarshalBson_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field: "name",
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	entity := &types.Entity{
		Name: "test name",
	}

	benchmarkMatchEntityPatternUnmarshalBson(b, cond, entity)
}

func BenchmarkMatchEntityPattern_UnmarshalBson_HasOneOf(b *testing.B) {
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
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:        "test",
				Description: "test",
				Value:       value,
			},
		},
	}

	benchmarkMatchEntityPatternUnmarshalBson(b, cond, entity)
}

func BenchmarkMatchEntityPattern_UnmarshalBson_Infos_Equal(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionEqual,
			Value: "test 2",
		},
	}
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkMatchEntityPatternUnmarshalBson(b, cond, entity)
}

func BenchmarkMatchEntityPattern_UnmarshalBson_Infos_Regexp(b *testing.B) {
	cond := pattern.FieldCondition{
		Field:     "infos.test",
		FieldType: pattern.FieldTypeString,
		Condition: pattern.Condition{
			Type:  pattern.ConditionRegexp,
			Value: "^test .+name$",
		},
	}
	entity := &types.Entity{
		Infos: map[string]types.Info{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
	}

	benchmarkMatchEntityPatternUnmarshalBson(b, cond, entity)
}

func benchmarkMatchEntityPattern(b *testing.B, fieldCond pattern.FieldCondition, entity *types.Entity) {
	const size = 100
	p := make(pattern.Entity, size)
	for i := 0; i < size; i++ {
		p[i] = []pattern.FieldCondition{fieldCond}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := match.MatchEntityPattern(p, entity)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}

func benchmarkMatchEntityPatternUnmarshalBson(b *testing.B, fieldCond pattern.FieldCondition, entity *types.Entity) {
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

		_, err = match.MatchEntityPattern(w.Pattern, entity)
		if err != nil {
			b.Fatalf("unexpected error %v", err)
		}
	}
}
