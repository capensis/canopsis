package pattern

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

func TestSuggestConditions(t *testing.T) {
	testCases := []struct {
		desc               string
		counts             map[string][]LiteralFieldStats
		groups             [][]string
		takenFields        map[string]bool
		expectedConditions [][]pattern.FieldCondition
		expectedError      bool
	}{
		{
			desc:   "",
			counts: map[string][]LiteralFieldStats{},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields:        make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 10,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal"),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 10,
					},
				},
			},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields: map[string]bool{
				"name": true,
			},
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 10,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal",
				},
			},
			takenFields: map[string]bool{
				"name": true,
			},
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal"),
					},
				},
			},
		},

		{
			desc:   "",
			counts: map[string][]LiteralFieldStats{},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields:        make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_2"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_1"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_1"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
			},
			takenFields: make(map[string]bool),
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: map[string]bool{
				"name": true,
			},
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_2"),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: map[string]bool{
				"name": true,
				"info": true,
			},
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_2"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_1"),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1",
				},
				{
					"literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_2"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_1"),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_1"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_2"),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},

		{
			desc:   "",
			counts: map[string][]LiteralFieldStats{},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields:        make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: map[string]bool{
				"info": true,
			},
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: map[string]bool{
				"info": true,
				"name": true,
			},
			expectedConditions: [][]pattern.FieldCondition{},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},
		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
			},
		},

		{
			desc: "",
			counts: map[string][]LiteralFieldStats{
				"literal_1": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_2": {
					{
						FieldName:         "name",
						LiteralFoundTimes: 1,
					},
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
				"literal_3": {
					{
						FieldName:         "info",
						LiteralFoundTimes: 1,
					},
				},
			},
			groups: [][]string{
				{
					"literal_1", "literal_2",
				},
				{
					"literal_3",
				},
			},
			takenFields: make(map[string]bool),
			expectedConditions: [][]pattern.FieldCondition{
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "literal_3"),
					},
					{
						Field:     "name",
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2"}),
					},
				},
				{
					{
						Field:     "infos.info",
						FieldType: pattern.FieldTypeString,
						Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"literal_1", "literal_2", "literal_3"}),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			patterns := suggestConditions(tc.counts, tc.groups, tc.takenFields)
			if len(patterns) != len(tc.expectedConditions) {
				t.Errorf("expected %d, got %d.", len(tc.expectedConditions), len(patterns))
				return
			}

			for i := range tc.expectedConditions {
				if !reflect.DeepEqual(patterns[i], tc.expectedConditions[i]) {
					t.Errorf("expected field condition %v, got %v.", tc.expectedConditions[i], patterns[i])
				}
			}
		})
	}
}
