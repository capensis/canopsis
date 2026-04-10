package pattern

import (
	"reflect"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

func TestGetLiteralToFieldCombinations(t *testing.T) {
	testCases := []struct {
		desc        string
		sets        map[string][]LiteralFieldStats
		takenFields map[string]bool
		expected    []map[string]string
	}{
		{
			desc:        "Given empty sets - should return empty slice",
			sets:        map[string][]LiteralFieldStats{},
			takenFields: map[string]bool{},
			expected:    []map[string]string{},
		},
		{
			desc: "Given single key with single value - should return one combination",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "field1"},
			},
		},
		{
			desc: "Given single key with two values - should return both combinations progressively",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
					{FieldName: "field2", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "field1"},
				{"literal1": "field2"},
			},
		},
		{
			desc: "Given two keys with single values - should return one combination",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
				},
				"literal2": {
					{FieldName: "field2", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "field1", "literal2": "field2"},
			},
		},
		{
			desc: "Given two keys with multiple values - should generate combinations progressively by max index",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
					{FieldName: "field2", LiteralFoundTimes: 1},
				},
				"literal2": {
					{FieldName: "field3", LiteralFoundTimes: 1},
					{FieldName: "field4", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "field1", "literal2": "field3"},
				{"literal1": "field1", "literal2": "field4"},
				{"literal1": "field2", "literal2": "field3"},
				{"literal1": "field2", "literal2": "field4"},
			},
		},
		{
			desc: "Given banned field - should filter out combinations with banned fields",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
					{FieldName: "banned", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{"banned": true},
			expected: []map[string]string{
				{"literal1": "field1"},
			},
		},
		{
			desc: "Given banned field in multi-key combination - should generate partial combinations",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
					{FieldName: "banned", LiteralFoundTimes: 1},
				},
				"literal2": {
					{FieldName: "field2", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{"banned": true},
			expected: []map[string]string{
				{"literal1": "field1", "literal2": "field2"},
				{"literal2": "field2"},
			},
		},
		{
			desc: "Given more combinations than maxCalculatedSuggestions - should limit to maxCalculatedSuggestions",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "f1", LiteralFoundTimes: 1},
					{FieldName: "f2", LiteralFoundTimes: 1},
					{FieldName: "f3", LiteralFoundTimes: 1},
					{FieldName: "f4", LiteralFoundTimes: 1},
				},
				"literal2": {
					{FieldName: "f5", LiteralFoundTimes: 1},
					{FieldName: "f6", LiteralFoundTimes: 1},
					{FieldName: "f7", LiteralFoundTimes: 1},
				},
				"literal3": {
					{FieldName: "f8", LiteralFoundTimes: 1},
					{FieldName: "f9", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "f1", "literal2": "f5", "literal3": "f8"},
				{"literal1": "f1", "literal2": "f5", "literal3": "f9"},
				{"literal1": "f1", "literal2": "f6", "literal3": "f8"},
				{"literal1": "f1", "literal2": "f6", "literal3": "f9"},
				{"literal1": "f2", "literal2": "f5", "literal3": "f8"},
				{"literal1": "f2", "literal2": "f5", "literal3": "f9"},
				{"literal1": "f2", "literal2": "f6", "literal3": "f8"},
				{"literal1": "f2", "literal2": "f6", "literal3": "f9"},
				{"literal1": "f1", "literal2": "f7", "literal3": "f8"},
				{"literal1": "f1", "literal2": "f7", "literal3": "f9"},
				{"literal1": "f2", "literal2": "f7", "literal3": "f8"},
				{"literal1": "f2", "literal2": "f7", "literal3": "f9"},
				{"literal1": "f3", "literal2": "f5", "literal3": "f8"},
				{"literal1": "f3", "literal2": "f5", "literal3": "f9"},
				{"literal1": "f3", "literal2": "f6", "literal3": "f8"},
				{"literal1": "f3", "literal2": "f6", "literal3": "f9"},
				{"literal1": "f3", "literal2": "f7", "literal3": "f8"},
				{"literal1": "f3", "literal2": "f7", "literal3": "f9"},
				{"literal1": "f4", "literal2": "f5", "literal3": "f8"},
				{"literal1": "f4", "literal2": "f5", "literal3": "f9"},
			},
		},
		{
			desc: "Given three keys with different sizes - should handle asymmetric sets",
			sets: map[string][]LiteralFieldStats{
				"literal1": {
					{FieldName: "field1", LiteralFoundTimes: 1},
				},
				"literal2": {
					{FieldName: "field2", LiteralFoundTimes: 1},
					{FieldName: "field3", LiteralFoundTimes: 1},
					{FieldName: "field4", LiteralFoundTimes: 1},
				},
				"literal3": {
					{FieldName: "field5", LiteralFoundTimes: 1},
					{FieldName: "field6", LiteralFoundTimes: 1},
				},
			},
			takenFields: map[string]bool{},
			expected: []map[string]string{
				{"literal1": "field1", "literal2": "field2", "literal3": "field5"},
				{"literal1": "field1", "literal2": "field2", "literal3": "field6"},
				{"literal1": "field1", "literal2": "field3", "literal3": "field5"},
				{"literal1": "field1", "literal2": "field3", "literal3": "field6"},
				{"literal1": "field1", "literal2": "field4", "literal3": "field5"},
				{"literal1": "field1", "literal2": "field4", "literal3": "field6"},
			},
		},
		{
			desc: "Given literal with banned field - should include only non-banned literals",
			sets: map[string][]LiteralFieldStats{
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
			takenFields: map[string]bool{
				"info": true,
			},
			expected: []map[string]string{
				{"literal_1": "name"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := getLiteralToFieldCombinations(tc.sets, tc.takenFields)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestGetPatternsCombinations(t *testing.T) {
	testCases := []struct {
		desc              string
		suggestedPatterns [][][]pattern.FieldCondition
		expected          []pattern.Entity
	}{
		{
			desc:              "Given empty patterns - should return empty result",
			suggestedPatterns: [][][]pattern.FieldCondition{},
			expected:          []pattern.Entity{},
		},
		{
			desc: "Given single pattern group with single option - should return one combination",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
				},
			},
		},
		{
			desc: "Given single pattern group with two options - should return both progressively",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
			},
		},
		{
			desc: "Given two pattern groups with single options - should return one combination",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
			},
		},
		{
			desc: "Given two pattern groups with multiple options - should generate combinations progressively",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
				{
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
				},
			},
		},
		{
			desc: "Given pattern groups with asymmetric sizes - should handle different sizes",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
				},
			},
		},
		{
			desc: "Given more combinations than maxCalculatedSuggestions - should limit to maxCalculatedSuggestions",
			suggestedPatterns: [][][]pattern.FieldCondition{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
				},
				{
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
				},
				{
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
			},
			expected: []pattern.Entity{
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field1", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value1")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field4", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value4")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
				{
					{
						{Field: "field2", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value2")},
					},
					{
						{Field: "field5", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value5")},
					},
					{
						{Field: "field8", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value8")},
					},
				},
				{
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field6", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value6")},
					},
				},
				{
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field3", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value3")},
					},
					{
						{Field: "field7", Condition: pattern.NewStringCondition(pattern.ConditionEqual, "value7")},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			result := getPatternsCombinations(tc.suggestedPatterns)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d combinations, got %d", len(tc.expected), len(result))
				return
			}

			for i := range tc.expected {
				if !reflect.DeepEqual(result[i], tc.expected[i]) {
					t.Errorf("combination %d: expected %+v, got %+v", i, tc.expected[i], result[i])
				}
			}
		})
	}
}
