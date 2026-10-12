package pattern_test

import (
	"slices"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

func TestEntity_GetInfosNames(t *testing.T) {
	dataSet := []struct {
		testName string
		pattern  pattern.Entity
		expected []string
	}{
		{
			testName: "empty pattern returns no keys",
			pattern:  pattern.Entity{},
			expected: nil,
		},
		{
			testName: "pattern without infos returns no keys",
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			expected: nil,
		},
		{
			testName: "single infos key in single group",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
			},
			expected: []string{"team"},
		},
		{
			testName: "multiple infos keys in single group",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
					{
						Field:     "infos.owner",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
					},
				},
			},
			expected: []string{"team", "owner"},
		},
		{
			testName: "infos keys across multiple groups",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
				{
					{
						Field:     "infos.zone",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
					},
				},
			},
			expected: []string{"team", "zone"},
		},
		{
			testName: "duplicate infos keys are deduplicated",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
					},
					{
						Field:     "infos.owner",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "admin"),
					},
				},
			},
			expected: []string{"team", "owner"},
		},
		{
			testName: "component_infos keys are ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
				},
			},
			expected: []string{"team"},
		},
		{
			testName: "non-infos fields are ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "comp-1"),
					},
				},
			},
			expected: []string{"team"},
		},
		{
			testName: "infos prefix without key is ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "infos",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
					},
				},
			},
			expected: nil,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result := data.pattern.GetInfosNames()
			slices.Sort(result)
			expected := slices.Clone(data.expected)
			slices.Sort(expected)

			if slices.Compare(result, expected) != 0 {
				t.Errorf("expected %v but got %v", expected, result)
			}
		})
	}
}

func TestEntity_GetComponentInfosNames(t *testing.T) {
	dataSet := []struct {
		testName string
		pattern  pattern.Entity
		expected []string
	}{
		{
			testName: "empty pattern returns no keys",
			pattern:  pattern.Entity{},
			expected: nil,
		},
		{
			testName: "pattern without component_infos returns no keys",
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
				},
			},
			expected: nil,
		},
		{
			testName: "single component_infos key in single group",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
				},
			},
			expected: []string{"env"},
		},
		{
			testName: "multiple component_infos keys in single group",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
					{
						Field:     "component_infos.zone",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu"),
					},
				},
			},
			expected: []string{"env", "zone"},
		},
		{
			testName: "component_infos keys across multiple groups",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
				},
				{
					{
						Field:     "component_infos.region",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu-west"),
					},
				},
			},
			expected: []string{"env", "region"},
		},
		{
			testName: "duplicate component_infos keys are deduplicated",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
				},
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "dev"),
					},
					{
						Field:     "component_infos.region",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "eu-west"),
					},
				},
			},
			expected: []string{"env", "region"},
		},
		{
			testName: "infos keys are ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
					{
						Field:     "infos.team",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "core"),
					},
				},
			},
			expected: []string{"env"},
		},
		{
			testName: "non-component_infos fields are ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "name",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "test"),
					},
					{
						Field:     "component_infos.env",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "prod"),
					},
					{
						Field:     "component",
						Condition: pattern.NewStringCondition(pattern.ConditionEqual, "comp-1"),
					},
				},
			},
			expected: []string{"env"},
		},
		{
			testName: "component_infos prefix without key is ignored",
			pattern: pattern.Entity{
				{
					{
						Field:     "component_infos",
						Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
					},
				},
			},
			expected: nil,
		},
	}

	for _, data := range dataSet {
		t.Run(data.testName, func(t *testing.T) {
			result := data.pattern.GetComponentInfosNames()
			slices.Sort(result)
			expected := slices.Clone(data.expected)
			slices.Sort(expected)

			if slices.Compare(result, expected) != 0 {
				t.Errorf("expected %v but got %v", expected, result)
			}
		})
	}
}
