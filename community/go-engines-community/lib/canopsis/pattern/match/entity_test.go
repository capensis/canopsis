package match_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"go.mongodb.org/mongo-driver/v2/bson"
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

func TestValidateEntityPattern(t *testing.T) {
	dataSets := map[string]struct {
		pattern         pattern.Entity
		forbiddenFields map[string]bool
		expectedResult  bool
	}{
		"given empty pattern should be valid": {
			pattern:        pattern.Entity{},
			expectedResult: true,
		},
		"given empty group should be invalid": {
			pattern:        pattern.Entity{{}},
			expectedResult: false,
		},
		"given valid string field condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			expectedResult: true,
		},
		"given invalid string field condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "name",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 10),
					},
				},
			},
			expectedResult: false,
		},
		"given valid int field condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "impact_level",
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 2),
					},
				},
			},
			expectedResult: true,
		},
		"given invalid int field condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "impact_level",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not-an-int"),
					},
				},
			},
			expectedResult: false,
		},
		"given valid time field condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "last_event_date",
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			expectedResult: true,
		},
		"given invalid time field condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "last_event_date",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "not-a-time"),
					},
				},
			},
			expectedResult: false,
		},
		"given unsupported field should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "created",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			expectedResult: false,
		},
		"given exact forbidden field should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			forbiddenFields: map[string]bool{"name": true},
			expectedResult:  false,
		},
		"given forbidden infos field but not infos in pattern should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			forbiddenFields: map[string]bool{"infos": true},
			expectedResult:  true,
		},
		"given forbidden component infos field but not component infos in pattern should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test-name"),
					},
				},
			},
			forbiddenFields: map[string]bool{"component_infos": true},
			expectedResult:  true,
		},
		"given forbidden infos should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "infos.team",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
			},
			forbiddenFields: map[string]bool{"infos": true},
			expectedResult:  false,
		},
		"given forbidden component infos should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "component_infos.team",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
			},
			forbiddenFields: map[string]bool{"component_infos": true},
			expectedResult:  false,
		},
		"given valid alias condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alias-value"),
						Alias:     "display_name",
					},
				},
			},
			expectedResult: true,
		},
		"given valid alias condition and forbidden infos should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alias-value"),
						Alias:     "display_name",
					},
				},
			},
			forbiddenFields: map[string]bool{"infos": true},
			expectedResult:  false,
		},
		"given invalid alias condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						FieldType: "invalid",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alias-value"),
						Alias:     "display_name",
					},
				},
			},
			expectedResult: false,
		},
		"given valid info condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "infos.team",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
			},
			expectedResult: true,
		},
		"given invalid info condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "infos.team",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 10),
					},
				},
			},
			expectedResult: false,
		},
		"given valid component info condition should be valid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "component_infos.rank",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewIntCondition(pattern.ConditionEqual, 10),
					},
				},
			},
			expectedResult: true,
		},
		"given invalid component info condition should be invalid": {
			pattern: pattern.Entity{
				{
					pattern.FieldCondition{
						Field:     "component_infos.rank",
						FieldType: pattern.FieldTypeInt,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "wrong-type"),
					},
				},
			},
			expectedResult: false,
		},
	}

	for name, test := range dataSets {
		t.Run(name, func(t *testing.T) {
			actual := match.ValidateEntityPattern(test.pattern, test.forbiddenFields)
			if actual != test.expectedResult {
				t.Fatalf("expected %v, got %v", test.expectedResult, actual)
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
		"given timestamp info condition with absolute time should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       1625097600,
					},
				},
			},
			matchResult: true,
		},
		"given timestamp info condition with absolute time should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       1672531200,
					},
				},
			},
			matchResult: false,
		},
		"given timestamp info condition with relative time should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       time.Now().Add(-30 * time.Minute).Unix(),
					},
				},
			},
			matchResult: true,
		},
		"given timestamp info condition with relative time should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       time.Now().Add(-2 * time.Hour).Unix(),
					},
				},
			},
			matchResult: false,
		},
		"given timestamp info condition and non-timestamp info should return error": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				Infos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       "not a timestamp",
					},
				},
			},
			matchErr:    pattern.ErrWrongConditionValue,
			matchResult: false,
		},

		"given timestamp component info condition with absolute time should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       1625097600,
					},
				},
			},
			matchResult: true,
		},
		"given timestamp component info condition with absolute time should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       1672531200,
					},
				},
			},
			matchResult: false,
		},
		"given timestamp component info condition with relative time should match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       time.Now().Add(-30 * time.Minute).Unix(),
					},
				},
			},
			matchResult: true,
		},
		"given timestamp component info condition with relative time should not match": {
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: func() pattern.Condition {
							condition, _ := pattern.NewDurationCondition(pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(3600, datetime.DurationUnitSecond))
							return condition
						}(),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       time.Now().Add(-2 * time.Hour).Unix(),
					},
				},
			},
			matchResult: false,
		},
		"given timestamp component info condition and non-timestamp info should return error": {
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.timestamp_info",
						FieldType: pattern.FieldTypeTimestamp,
						Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 1609459200, 1640995200),
					},
				},
			},
			entity: types.Entity{
				ComponentInfos: map[string]types.Info{
					"timestamp_info": {
						Name:        "timestamp_info",
						Description: "test timestamp",
						Value:       "not a timestamp",
					},
				},
			},
			matchErr:    pattern.ErrWrongConditionValue,
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
