package builder

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

func TestBuild_StringOperators(t *testing.T) {
	regexpCond, err := pattern.NewRegexpCondition(pattern.ConditionRegexp, `^alarm-\d+$`)
	if err != nil {
		t.Fatalf("NewRegexpCondition: %v", err)
	}

	tests := []struct {
		name             string
		input            Input
		expected         [][]pattern.FieldCondition
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given eq operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("alarm-1")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alarm-1"),
			}}},
		},
		{
			name: "given neq operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionNotEqual,
						Value: ValueSpec{String: new("alarm-1")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionNotEqual, "alarm-1"),
			}}},
		},
		{
			name: "given contain operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionContain,
						Value: ValueSpec{String: new("disk")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionContain, "disk"),
			}}},
		},
		{
			name: "given not_contain operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionNotContain,
						Value: ValueSpec{String: new("disk")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionNotContain, "disk"),
			}}},
		},
		{
			name: "given begin_with operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionBeginWith,
						Value: ValueSpec{String: new("srv-")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionBeginWith, "srv-"),
			}}},
		},
		{
			name: "given not_begin_with operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionNotBeginWith,
						Value: ValueSpec{String: new("srv-")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionNotBeginWith, "srv-"),
			}}},
		},
		{
			name: "given end_with operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionEndWith,
						Value: ValueSpec{String: new(".local")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionEndWith, ".local"),
			}}},
		},
		{
			name: "given not_end_with operator with string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionNotEndWith,
						Value: ValueSpec{String: new(".local")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionNotEndWith, ".local"),
			}}},
		},
		{
			name: "given regexp operator with valid pattern should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionRegexp,
						Value: ValueSpec{String: new(`^alarm-\d+$`)},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: regexpCond,
			}}},
		},
		{
			name: "given eq operator with empty string value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("")},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, ""),
			}}},
		},
		{
			name: "given is_one_of operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsOneOf,
						Value: ValueSpec{Strings: []string{"a", "b"}},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionIsOneOf, []string{"a", "b"}),
			}}},
		},
		{
			name: "given is_not_one_of operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsNotOneOf,
						Value: ValueSpec{Strings: []string{"a", "b"}},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionIsNotOneOf, []string{"a", "b"}),
			}}},
		},
		{
			name: "given exist operator with true should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionExist,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}},
		},
		{
			name: "given exist operator with false should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionExist,
						Value: ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.display_name",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
			}}},
		},
		{
			name: "given eq operator and missing string value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionEqual,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a string value`,
		},
		{
			name: "given contain operator and strings list instead of string should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionContain,
						Value: ValueSpec{Strings: []string{"a"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a string value`,
		},
		{
			name: "given neq operator and boolean instead of string should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionNotEqual,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a string value`,
		},
		{
			name: "given regexp operator and missing string value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionRegexp,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a string value`,
		},
		{
			name: "given regexp operator and empty string value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionRegexp,
						Value: ValueSpec{String: new("")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty regexp`,
		},
		{
			name: "given regexp operator and invalid pattern should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionRegexp,
						Value: ValueSpec{String: new("[invalid")},
					}}},
				}},
			},
			expectedContains: "error parsing regexp",
		},
		{
			name: "given is_one_of operator and empty string list should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsOneOf,
						Value: ValueSpec{Strings: nil},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty string list`,
		},
		{
			name: "given is_not_one_of operator and empty string list should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsNotOneOf,
						Value: ValueSpec{Strings: []string{}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty string list`,
		},
		{
			name: "given is_one_of operator and list containing empty string should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsOneOf,
						Value: ValueSpec{Strings: []string{"ok", ""}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `empty strings`,
		},
		{
			name: "given is_one_of operator and single string instead of list should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionIsOneOf,
						Value: ValueSpec{String: new("a")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty string list`,
		},
		{
			name: "given exist operator and missing boolean value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionExist,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given exist operator and string instead of boolean should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionExist,
						Value: ValueSpec{String: new("true")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given string field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: pattern.ConditionGT,
						Value: ValueSpec{Integer: new(int64(1))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionGT,
		},
		{
			name: "given unknown operator name should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.display_name", Operator: "like",
						Value: ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: "like",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Build(tt.input, "widget_filter_AlarmsList")
			if tt.expectedErr != nil || tt.expectedContains != "" {
				if err == nil {
					t.Fatal("Build() expected error, got nil")
				}
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("error %v does not wrap %v", err, tt.expectedErr)
				}
				if tt.expectedContains != "" && !strings.Contains(err.Error(), tt.expectedContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.expectedContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("Build() unexpected error: %v", err)
			}

			for i := range tt.expected {
				if !reflect.DeepEqual(result.AlarmPattern[i], tt.expected[i]) {
					t.Errorf("expected field condition %v, got %v.", tt.expected[i], result.AlarmPattern[i])
				}
			}
		})
	}
}

func TestBuild_IntOperators(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         [][]pattern.FieldCondition
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given eq operator with integer value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(3))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionEqual, 3),
			}}},
		},
		{
			name: "given neq operator with integer value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionNotEqual,
						Value: ValueSpec{Integer: new(int64(1))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionNotEqual, 1),
			}}},
		},
		{
			name: "given gt operator with integer value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionGT,
						Value: ValueSpec{Integer: new(int64(2))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionGT, 2),
			}}},
		},
		{
			name: "given lt operator with integer value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionLT,
						Value: ValueSpec{Integer: new(int64(4))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionLT, 4),
			}}},
		},
		{
			name: "given eq operator with negative integer should return error for state enum",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(-1))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `is not a valid value for "v.state.val"`,
		},
		{
			name: "given eq operator with out-of-range state should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(9))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `is not a valid value for "v.state.val"`,
		},
		{
			name: "given lt 0 on state should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionLT,
						Value: ValueSpec{Integer: new(int64(0))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `lt 0 matches no values`,
		},
		{
			name: "given gt 3 on state should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionGT,
						Value: ValueSpec{Integer: new(int64(3))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `gt 3 matches no values`,
		},
		{
			name: "given gt -1 on state should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionGT,
						Value: ValueSpec{Integer: new(int64(-1))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionGT, -1),
			}}},
		},
		{
			name: "given lt 1 on state should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionLT,
						Value: ValueSpec{Integer: new(int64(1))},
					}}},
				}},
			},
			expected: [][]pattern.FieldCondition{{{
				Field:     "v.state.val",
				Condition: pattern.NewIntCondition(pattern.ConditionLT, 1),
			}}},
		},
		{
			name: "given int operator and missing integer value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires an integer value`,
		},
		{
			name: "given int operator and string instead of integer should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("3")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires an integer value`,
		},
		{
			name: "given int operator and boolean instead of integer should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionLT,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires an integer value`,
		},
		{
			name: "given int field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.state.val", Operator: pattern.ConditionContain,
						Value: ValueSpec{String: new("3")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionContain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Build(tt.input, "widget_filter_AlarmsList")
			if tt.expectedErr != nil || tt.expectedContains != "" {
				if err == nil {
					t.Fatal("Build() expected error, got nil")
				}
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("error %v does not wrap %v", err, tt.expectedErr)
				}
				if tt.expectedContains != "" && !strings.Contains(err.Error(), tt.expectedContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.expectedContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("Build() unexpected error: %v", err)
			}

			for i := range tt.expected {
				if !reflect.DeepEqual(result.AlarmPattern[i], tt.expected[i]) {
					t.Errorf("expected field condition %v, got %v.", tt.expected[i], result.AlarmPattern[i])
				}
			}
		})
	}
}

func TestBuild_ForbiddenFields(t *testing.T) {
	tests := []struct {
		name             string
		promptContext    string
		input            Input
		expectedErr      error
		expectedContains string
	}{
		{
			name:          "given idle_rule context and forbidden alarm field v.last_event_date should return error",
			promptContext: "idle_rule",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_event_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 1, To: 2}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "v.last_event_date" is forbidden`,
		},
		{
			name:          "given idle_rule context and relative_time on absolute-only field should return error",
			promptContext: "idle_rule",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.creation_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 1, Unit: "h"},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "v.creation_date" is forbidden`,
		},
		{
			name:          "given widget_filter context and otherwise forbidden field should succeed",
			promptContext: "widget_filter_AlarmsList",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_event_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 1, To: 2}},
					}}},
				}},
			},
		},
		{
			name:          "given idle_rule context and forbidden entity field last_event_date should return error",
			promptContext: "idle_rule",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "last_event_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 1, To: 2}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "last_event_date" is forbidden`,
		},
		{
			name:          "given state_settings context and forbidden component_infos field should return error",
			promptContext: "state_settings",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "component_infos.team" is forbidden`,
		},
		{
			name:          "given state_settings context and forbidden type field should return error",
			promptContext: "state_settings",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("resource")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "type" is forbidden`,
		},
		{
			name:          "given state_settings context and forbidden component field should return error",
			promptContext: "state_settings",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("comp-1")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "component" is forbidden`,
		},
		{
			name:          "given state_settings_service_inherited context and allowed component_infos field should succeed",
			promptContext: "state_settings_service_inherited",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
		},
		{
			name:          "given state_settings_service_inherited context and allowed type and component fields should succeed",
			promptContext: "state_settings_service_inherited",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{
						{
							{Field: "type", Operator: pattern.ConditionEqual, Value: ValueSpec{String: new("resource")}},
							{Field: "component", Operator: pattern.ConditionEqual, Value: ValueSpec{String: new("comp-1")}},
						},
					},
				}},
			},
		},
		{
			name:          "given state_settings_service_inherited context and forbidden last_event_date field should return error",
			promptContext: "state_settings_service_inherited",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "last_event_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 1, To: 2}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "last_event_date" is forbidden`,
		},
		{
			name:          "given entity_service context and allowed component_infos field should succeed",
			promptContext: "entity_service",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
		},
		{
			name:          "given widget_filter_ServiceWeather context and allowed component_infos field should succeed",
			promptContext: "widget_filter_ServiceWeather",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
		},
		{
			name:          "given entity_service context and forbidden connector field should return error",
			promptContext: "entity_service",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "connector", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("nagios")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "connector" is forbidden`,
		},
		{
			name:          "given dynamic_infos context and forbidden v.infos field should return error",
			promptContext: "dynamic_infos",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "owner", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("alice")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrForbiddenField,
			expectedContains: `field "v.infos.owner" is forbidden`,
		},
		{
			name:          "given unknown prompt context should return error",
			promptContext: "not_a_real_context",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "name", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `unknown prompt context "not_a_real_context"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Build(tt.input, tt.promptContext)
			if tt.expectedErr != nil || tt.expectedContains != "" {
				if err == nil {
					t.Fatal("Build() expected error, got nil")
				}
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("error %v does not wrap %v", err, tt.expectedErr)
				}
				if tt.expectedContains != "" && !strings.Contains(err.Error(), tt.expectedContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.expectedContains)
				}
				return
			}
			if err != nil {
				t.Fatalf("Build() unexpected error: %v", err)
			}
		})
	}
}

func TestBuild_BoolOperators(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given eq operator with boolean value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "weather_service_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "is_grey", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{WeatherServicePattern: [][]pattern.FieldCondition{{{
				Field:     "is_grey",
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
			}}}},
		},
		{
			name: "given bool operator and missing boolean value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "weather_service_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "is_grey", Operator: pattern.ConditionEqual,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given bool field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "weather_service_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "is_grey", Operator: pattern.ConditionNotEqual,
						Value: ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionNotEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_RefOperators(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given exist operator with boolean value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.ack", Operator: pattern.ConditionExist,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.ack",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}}},
		},
		{
			name: "given exist operator and missing boolean value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.ack", Operator: pattern.ConditionExist,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given ref field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.ack", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_DurationOperators(t *testing.T) {
	gtCond := mustDurationCondition(t, pattern.ConditionGT, datetime.DurationWithUnit{Value: 2, Unit: datetime.DurationUnitDay})
	ltCond := mustDurationCondition(t, pattern.ConditionLT, datetime.DurationWithUnit{Value: 1, Unit: datetime.DurationUnitHour})

	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given gt operator with duration value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionGT,
						Value: ValueSpec{Duration: &DurationSpec{Value: 2, Unit: datetime.DurationUnitDay}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.duration",
				Condition: gtCond,
			}}}},
		},
		{
			name: "given lt operator with duration value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionLT,
						Value: ValueSpec{Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.duration",
				Condition: ltCond,
			}}}},
		},
		{
			name: "given duration operator and missing duration value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionGT,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a duration value`,
		},
		{
			name: "given duration operator and non-positive duration value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionGT,
						Value: ValueSpec{Duration: &DurationSpec{Value: 0, Unit: datetime.DurationUnitHour}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must be a positive integer`,
		},
		{
			name: "given duration operator and unknown duration unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionLT,
						Value: ValueSpec{Duration: &DurationSpec{Value: 1, Unit: "ms"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unknown unit`,
		},
		{
			name: "given duration operator with month unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionGT,
						Value: ValueSpec{Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitMonth}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unknown unit`,
		},
		{
			name: "given duration field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.duration", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_TagsOperators(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given has_every operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasEvery,
						Value: ValueSpec{Strings: []string{"prod", "critical"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"prod", "critical"}),
			}}}},
		},
		{
			name: "given has_one_of operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasOneOf,
						Value: ValueSpec{Strings: []string{"prod"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"prod"}),
			}}}},
		},
		{
			name: "given has_not operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasNot,
						Value: ValueSpec{Strings: []string{"dev"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"dev"}),
			}}}},
		},
		{
			name: "given has_labels operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasLabels,
						Value: ValueSpec{Strings: []string{"env"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasLabels, []string{"env"}),
			}}}},
		},
		{
			name: "given has_not_labels operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasNotLabels,
						Value: ValueSpec{Strings: []string{"env"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNotLabels, []string{"env"}),
			}}}},
		},
		{
			name: "given is_empty operator with boolean value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionIsEmpty,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "tags",
				Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, true),
			}}}},
		},
		{
			name: "given is_empty operator and missing boolean value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionIsEmpty,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given has_every operator and empty string list should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionHasEvery,
						Value: ValueSpec{Strings: nil},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty string list`,
		},
		{
			name: "given tags field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "tags", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("prod")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionEqual,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_StringArrayOperators(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given has_every operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasEvery,
						Value:    ValueSpec{Strings: []string{"admin", "ops"}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"admin", "ops"}),
			}}}},
		},
		{
			name: "given has_one_of operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasOneOf,
						Value:    ValueSpec{Strings: []string{"admin"}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"admin"}),
			}}}},
		},
		{
			name: "given has_not operator with string list should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasNot,
						Value:    ValueSpec{Strings: []string{"guest"}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasNot, []string{"guest"}),
			}}}},
		},
		{
			name: "given is_empty operator with boolean value should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionIsEmpty,
						Value:    ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			}}}},
		},
		{
			name: "given is_empty operator and missing boolean value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionIsEmpty,
						Value:    ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a boolean value`,
		},
		{
			name: "given has_one_of operator and empty string list should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasOneOf,
						Value:    ValueSpec{Strings: []string{}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `non-empty string list`,
		},
		{
			name: "given string_array field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasLabels,
						Value:    ValueSpec{Strings: []string{"env"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionHasLabels,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_TimeOperators(t *testing.T) {
	newerThan := mustDurationCondition(t, pattern.ConditionTimeRelative, datetime.DurationWithUnit{Value: 2, Unit: datetime.DurationUnitHour})
	olderThan := mustDurationCondition(t, pattern.ConditionTimeRelative, datetime.DurationWithUnit{}, datetime.DurationWithUnit{Value: 3, Unit: datetime.DurationUnitDay})
	between := mustDurationCondition(t, pattern.ConditionTimeRelative,
		datetime.DurationWithUnit{Value: 2, Unit: datetime.DurationUnitHour},
		datetime.DurationWithUnit{Value: 1, Unit: datetime.DurationUnitHour},
	)

	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given relative_time newer_than with duration should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.last_update_date",
				Condition: newerThan,
			}}}},
		},
		{
			name: "given relative_time older_than with duration should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "older_than",
							Duration: &DurationSpec{Value: 3, Unit: datetime.DurationUnitDay},
						}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.last_update_date",
				Condition: olderThan,
			}}}},
		},
		{
			name: "given relative_time between with from and to should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.last_update_date",
				Condition: between,
			}}}},
		},
		{
			name: "given absolute_time with from and to should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 100, To: 200}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.last_update_date",
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 100, 200),
			}}}},
		},
		{
			name: "given relative_time operator and missing relative_time value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires a relative_time value`,
		},
		{
			name: "given absolute_time operator and missing absolute_time value should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires an absolute_time value`,
		},
		{
			name: "given time field and unsupported operator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(1))},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionEqual,
		},
		{
			name: "given relative_time newer_than and missing duration should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{Kind: "newer_than"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `kind "newer_than" requires duration`,
		},
		{
			name: "given relative_time older_than and missing duration should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{Kind: "older_than"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `kind "older_than" requires duration`,
		},
		{
			name: "given relative_time older_than with from/to should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "older_than",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitDay},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitDay},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must use duration only`,
		},
		{
			name: "given relative_time between with duration should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "between",
							Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitDay},
							From:     &DurationSpec{Value: 2, Unit: datetime.DurationUnitDay},
							To:       &DurationSpec{Value: 1, Unit: datetime.DurationUnitDay},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must use from and to only`,
		},
		{
			name: "given relative_time between and missing from or to should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `kind "between" requires from and to`,
		},
		{
			name: "given relative_time newer_than and non-positive duration should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 0, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must be a positive integer`,
		},
		{
			name: "given relative_time older_than and unknown duration unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "older_than",
							Duration: &DurationSpec{Value: 1, Unit: "ms"},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unknown unit`,
		},
		{
			name: "given relative_time between and non-positive from duration should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: -1, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"between" from`,
		},
		{
			name: "given relative_time between and unknown to duration unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 1, Unit: "ms"},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"between" to`,
		},
		{
			name: "given relative_time newer_than with month unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitMonth},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unknown unit`,
		},
		{
			name: "given relative_time older_than with year unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "older_than",
							Duration: &DurationSpec{Value: 1, Unit: datetime.DurationUnitYear},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unknown unit`,
		},
		{
			name: "given relative_time between with month from unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitMonth},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `relative_time "between" from`,
		},
		{
			name: "given relative_time between with year to unit should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitYear},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `relative_time "between" to`,
		},
		{
			name: "given relative_time with unsupported kind should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{Kind: "around"}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `unsupported relative_time kind`,
		},
		{
			name: "given absolute_time and non-positive timestamps should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 0, To: 10}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must be positive Unix timestamps in seconds`,
		},
		{
			name: "given absolute_time in milliseconds should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 1710000000000, To: 1710003600000}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `not milliseconds`,
		},
		{
			name: "given absolute_time and to less than or equal to from should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 100, To: 100}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must be greater than from`,
		},
		{
			name: "given absolute_time and reversed range should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 200, To: 100}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must be greater than from`,
		},
		{
			name: "given activation_date with relative_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.activation_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.activation_date",
				Condition: newerThan,
			}}}},
		},
		{
			name: "given activation_date with absolute_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.activation_date", Operator: pattern.ConditionTimeAbsolute,
						Value: ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 100, To: 200}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.activation_date",
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 100, 200),
			}}}},
		},
		{
			name: "given relative_time between and equal from and to should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `must not be equal`,
		},
		{
			name: "given relative_time between and from less than to should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.last_update_date", Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind: "between",
							From: &DurationSpec{Value: 1, Unit: datetime.DurationUnitHour},
							To:   &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `requires from > to`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_ResolveField(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given pattern group with no conditions should return error",
			input: Input{
				Patterns: []Spec{{
					Type:   "alarm_pattern",
					Groups: [][]ConditionSpec{{}},
				}},
			},
			expectedErr: pattern.ErrEmptyGroup,
		},
		{
			name: "given unknown field name should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.unknown", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedField,
			expectedContains: "v.unknown",
		},
		{
			name: "given dual-type field and exist operator should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.activation_date", Operator: pattern.ConditionExist,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.activation_date",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}}},
		},
		{
			name: "given dual-type field and operator not allowed for that field should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.activation_date", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: pattern.ConditionEqual,
		},
		{
			name: "given dynamic extra field with field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "extra.*", Key: "site", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expected: Result{EventPattern: [][]pattern.FieldCondition{{{
				Field:     "extra.site",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "paris"),
			}}}},
		},
		{
			name: "given dynamic extra field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "extra.*", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expectedContains: `requires key`,
		},
		{
			name: "given dynamic extra field exist without field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "extra.*", Key: "site",
						Operator: pattern.ConditionExist,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{EventPattern: [][]pattern.FieldCondition{{{
				Field:     "extra.site",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}}},
		},
		{
			name: "given dynamic extra field non-exist without field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "extra.*", Key: "site",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expectedContains: `field_type is required`,
		},
		{
			name: "given dynamic extra field unsupported field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "extra.*", Key: "site", FieldType: "float",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedFieldType,
			expectedContains: "float",
		},
		{
			name: "given ticket_data dynamic field without field_type should succeed as string",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.ticket.ticket_data.*", Key: "number",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("42")},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.ticket.ticket_data.number",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "42"),
			}}}},
		},
		{
			name: "given ticket_data dynamic field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field:    "v.ticket.ticket_data.*",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("42")},
					}}},
				}},
			},
			expectedContains: `requires key`,
		},
		{
			name: "given pbehavior_pattern with canonical_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "pbehavior_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "pbehavior_info.canonical_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("maintenance")},
					}}},
				}},
			},
			expected: Result{PbehaviorPattern: [][]pattern.FieldCondition{{{
				Field:     "pbehavior_info.canonical_type",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "maintenance"),
			}}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_EntityInfos(t *testing.T) {
	relativeNewer := mustDurationCondition(t, pattern.ConditionTimeRelative, datetime.NewDurationWithUnit(1, datetime.DurationUnitHour))

	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given infos string field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "city", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.city",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "paris"),
			}}}},
		},
		{
			name: "given infos int field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "level", FieldType: pattern.FieldTypeInt,
						Operator: pattern.ConditionGT,
						Value:    ValueSpec{Integer: new(int64(1))},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.level",
				FieldType: pattern.FieldTypeInt,
				Condition: pattern.NewIntCondition(pattern.ConditionGT, 1),
			}}}},
		},
		{
			name: "given infos bool field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "active", FieldType: pattern.FieldTypeBool,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.active",
				FieldType: pattern.FieldTypeBool,
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, false),
			}}}},
		},
		{
			name: "given infos string_array field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasOneOf,
						Value:    ValueSpec{Strings: []string{"a", "b"}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"a", "b"}),
			}}}},
		},
		{
			name: "given infos timestamp field_type with absolute_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeAbsolute,
						Value:    ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 10, To: 20}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 10, 20),
			}}}},
		},
		{
			name: "given infos timestamp field_type with relative_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 1, Unit: "h"},
						}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: relativeNewer,
			}}}},
		},
		{
			name: "given infos exist operator without field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "city",
						Operator: pattern.ConditionExist,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "infos.city",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}}},
		},
		{
			name: "given infos field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `requires key`,
		},
		{
			name: "given infos non-exist operator without field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "city",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("paris")},
					}}},
				}},
			},
			expectedContains: `field_type is required`,
		},
		{
			name: "given infos unsupported field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "infos.*", Key: "city", FieldType: "float",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedFieldType,
			expectedContains: "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_ComponentInfos(t *testing.T) {
	relativeOlder := mustDurationCondition(t, pattern.ConditionTimeRelative, datetime.DurationWithUnit{}, datetime.DurationWithUnit{Value: 2, Unit: datetime.DurationUnitDay})

	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given component_infos string field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.team",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "ops"),
			}}}},
		},
		{
			name: "given component_infos int field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "rank", FieldType: pattern.FieldTypeInt,
						Operator: pattern.ConditionGT,
						Value:    ValueSpec{Integer: new(int64(3))},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.rank",
				FieldType: pattern.FieldTypeInt,
				Condition: pattern.NewIntCondition(pattern.ConditionGT, 3),
			}}}},
		},
		{
			name: "given component_infos bool field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "critical", FieldType: pattern.FieldTypeBool,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.critical",
				FieldType: pattern.FieldTypeBool,
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
			}}}},
		},
		{
			name: "given component_infos string_array field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasOneOf,
						Value:    ValueSpec{Strings: []string{"a", "b"}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasOneOf, []string{"a", "b"}),
			}}}},
		},
		{
			name: "given component_infos timestamp field_type with absolute_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeAbsolute,
						Value:    ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 10, To: 20}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 10, 20),
			}}}},
		},
		{
			name: "given component_infos timestamp field_type with relative_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "older_than",
							Duration: &DurationSpec{Value: 2, Unit: "d"},
						}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: relativeOlder,
			}}}},
		},
		{
			name: "given component_infos exist operator without field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team",
						Operator: pattern.ConditionExist,
						Value:    ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Field:     "component_infos.team",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
			}}}},
		},
		{
			name: "given component_infos field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `requires key`,
		},
		{
			name: "given component_infos non-exist operator without field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("ops")},
					}}},
				}},
			},
			expectedContains: `field_type is required`,
		},
		{
			name: "given component_infos unsupported field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "component_infos.*", Key: "team", FieldType: "float",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedFieldType,
			expectedContains: "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_AlarmInfos(t *testing.T) {
	relativeNewer := mustDurationCondition(t, pattern.ConditionTimeRelative, datetime.DurationWithUnit{Value: 2, Unit: datetime.DurationUnitHour})

	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given v.infos string field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "owner", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("alice")},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.owner",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alice"),
			}}}},
		},
		{
			name: "given v.infos int field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "priority", FieldType: pattern.FieldTypeInt,
						Operator: pattern.ConditionLT,
						Value:    ValueSpec{Integer: new(int64(5))},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.priority",
				FieldType: pattern.FieldTypeInt,
				Condition: pattern.NewIntCondition(pattern.ConditionLT, 5),
			}}}},
		},
		{
			name: "given v.infos bool field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "escalated", FieldType: pattern.FieldTypeBool,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.escalated",
				FieldType: pattern.FieldTypeBool,
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
			}}}},
		},
		{
			name: "given v.infos string_array field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "labels", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionHasEvery,
						Value:    ValueSpec{Strings: []string{"prod", "critical"}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.labels",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewStringArrayCondition(pattern.ConditionHasEvery, []string{"prod", "critical"}),
			}}}},
		},
		{
			name: "given v.infos timestamp field_type with absolute_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeAbsolute,
						Value:    ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 100, To: 200}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 100, 200),
			}}}},
		},
		{
			name: "given v.infos timestamp field_type with relative_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeRelative,
						Value: ValueSpec{RelativeTime: &RelativeTimeSpec{
							Kind:     "newer_than",
							Duration: &DurationSpec{Value: 2, Unit: datetime.DurationUnitHour},
						}},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: relativeNewer,
			}}}},
		},
		{
			name: "given v.infos exist operator without field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "owner",
						Operator: pattern.ConditionExist,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{AlarmPattern: [][]pattern.FieldCondition{{{
				Field:     "v.infos.owner",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, true),
			}}}},
		},
		{
			name: "given v.infos field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `requires key`,
		},
		{
			name: "given v.infos non-exist operator without field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "owner",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("alice")},
					}}},
				}},
			},
			expectedContains: `field_type is required`,
		},
		{
			name: "given v.infos unsupported field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.infos.*", Key: "owner", FieldType: "float",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedFieldType,
			expectedContains: "float",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuild_Alias(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given alias string field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "city", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionContain,
						Value:    ValueSpec{String: new("Paris")},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "city",
				FieldType: pattern.FieldTypeString,
				Condition: pattern.NewStringCondition(pattern.ConditionContain, "Paris"),
			}}}},
		},
		{
			name: "given alias int field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "level", FieldType: pattern.FieldTypeInt,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{Integer: new(int64(2))},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "level",
				FieldType: pattern.FieldTypeInt,
				Condition: pattern.NewIntCondition(pattern.ConditionEqual, 2),
			}}}},
		},
		{
			name: "given alias bool field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "active", FieldType: pattern.FieldTypeBool,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{Boolean: new(true)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "active",
				FieldType: pattern.FieldTypeBool,
				Condition: pattern.NewBoolCondition(pattern.ConditionEqual, true),
			}}}},
		},
		{
			name: "given alias string_array field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "roles", FieldType: pattern.FieldTypeStringArray,
						Operator: pattern.ConditionIsEmpty,
						Value:    ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "roles",
				FieldType: pattern.FieldTypeStringArray,
				Condition: pattern.NewBoolCondition(pattern.ConditionIsEmpty, false),
			}}}},
		},
		{
			name: "given alias timestamp field_type with absolute_time should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "seen_at", FieldType: pattern.FieldTypeTimestamp,
						Operator: pattern.ConditionTimeAbsolute,
						Value:    ValueSpec{AbsoluteTime: &AbsoluteTimeSpec{From: 10, To: 20}},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "seen_at",
				FieldType: pattern.FieldTypeTimestamp,
				Condition: pattern.NewTimeIntervalCondition(pattern.ConditionTimeAbsolute, 10, 20),
			}}}},
		},
		{
			name: "given alias exist operator without field_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "city",
						Operator: pattern.ConditionExist,
						Value:    ValueSpec{Boolean: new(false)},
					}}},
				}},
			},
			expected: Result{EntityPattern: [][]pattern.FieldCondition{{{
				Alias:     "city",
				Condition: pattern.NewBoolCondition(pattern.ConditionExist, false),
			}}}},
		},
		{
			name: "given alias field without key should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", FieldType: pattern.FieldTypeString,
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `field "alias" requires key`,
		},
		{
			name: "given alias non-exist operator without field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "city",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedContains: `field_type is required`,
		},
		{
			name: "given alias unsupported field_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "alias", Key: "city", FieldType: "invalid",
						Operator: pattern.ConditionEqual,
						Value:    ValueSpec{String: new("x")},
					}}},
				}},
			},
			expectedErr:      pattern.ErrUnsupportedFieldType,
			expectedContains: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, "widget_filter_AlarmsList", tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func assertBuildResult(t *testing.T, input Input, promptContext string, expected Result, expectedErr error, expectedContains string) {
	t.Helper()

	result, err := Build(input, promptContext)
	if expectedErr != nil || expectedContains != "" {
		if err == nil {
			t.Fatal("Build() expected error, got nil")
		}
		if expectedErr != nil && !errors.Is(err, expectedErr) {
			t.Errorf("error %v does not wrap %v", err, expectedErr)
		}
		if expectedContains != "" && !strings.Contains(err.Error(), expectedContains) {
			t.Errorf("error %q does not contain %q", err.Error(), expectedContains)
		}
		return
	}
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func mustDurationCondition(t *testing.T, operator string, d ...datetime.DurationWithUnit) pattern.Condition {
	t.Helper()

	cond, err := pattern.NewDurationCondition(operator, d...)
	if err != nil {
		t.Fatalf("NewDurationCondition: %v", err)
	}

	return cond
}

func TestBuild_UnknownPatternType(t *testing.T) {
	assertBuildResult(t, Input{
		Patterns: []Spec{{
			Type: "unknown_pattern",
			Groups: [][]ConditionSpec{{{
				Field: "name", Operator: pattern.ConditionEqual,
				Value: ValueSpec{String: new("x")},
			}}},
		}},
	}, "widget_filter_AlarmsList", Result{}, nil, `unknown pattern type "unknown_pattern"`)
}

func TestBuild_MultiPatternComposition(t *testing.T) {
	assertBuildResult(t, Input{
		Patterns: []Spec{
			{
				Type: "alarm_pattern",
				Groups: [][]ConditionSpec{{{
					Field: "v.display_name", Operator: pattern.ConditionEqual,
					Value: ValueSpec{String: new("alarm-1")},
				}}},
			},
			{
				Type: "entity_pattern",
				Groups: [][]ConditionSpec{{{
					Field: "name", Operator: pattern.ConditionEqual,
					Value: ValueSpec{String: new("entity-1")},
				}}},
			},
		},
	}, "widget_filter_AlarmsList", Result{
		AlarmPattern: [][]pattern.FieldCondition{{{
			Field:     "v.display_name",
			Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alarm-1"),
		}}},
		EntityPattern: [][]pattern.FieldCondition{{{
			Field:     "name",
			Condition: pattern.NewStringCondition(pattern.ConditionEqual, "entity-1"),
		}}},
	}, nil, "")
}

func TestBuild_SameTypePatternsAppend(t *testing.T) {
	assertBuildResult(t, Input{
		Patterns: []Spec{
			{
				Type: "alarm_pattern",
				Groups: [][]ConditionSpec{{{
					Field: "v.display_name", Operator: pattern.ConditionEqual,
					Value: ValueSpec{String: new("alarm-1")},
				}}},
			},
			{
				Type: "alarm_pattern",
				Groups: [][]ConditionSpec{{{
					Field: "v.display_name", Operator: pattern.ConditionEqual,
					Value: ValueSpec{String: new("alarm-2")},
				}}},
			},
		},
	}, "widget_filter_AlarmsList", Result{
		AlarmPattern: [][]pattern.FieldCondition{
			{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alarm-1"),
			}},
			{{
				Field:     "v.display_name",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "alarm-2"),
			}},
		},
	}, nil, "")
}

func TestBuild_ForbiddenAlias(t *testing.T) {
	// No production context currently forbids infos (which also forbids aliases).
	// Temporarily enable it on idle_rule so Build exercises the alias error label.
	forbidden := patternfields.GetForbiddenFieldsInEntityPattern(mongo.IdleRuleMongoCollection)
	forbidden["infos"] = true
	t.Cleanup(func() { delete(forbidden, "infos") })

	assertBuildResult(t, Input{
		Patterns: []Spec{{
			Type: "entity_pattern",
			Groups: [][]ConditionSpec{{{
				Field: "alias", Key: "city", FieldType: pattern.FieldTypeString,
				Operator: pattern.ConditionEqual,
				Value:    ValueSpec{String: new("paris")},
			}}},
		}},
	}, "idle_rule", Result{}, pattern.ErrForbiddenField, `field "alias" is forbidden`)
}

func TestBuild_StringEnumValues(t *testing.T) {
	tests := []struct {
		name             string
		input            Input
		promptContext    string
		expected         Result
		expectedErr      error
		expectedContains string
	}{
		{
			name: "given valid source_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "source_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("resource")},
					}}},
				}},
			},
			promptContext: "eventfilter",
			expected: Result{EventPattern: [][]pattern.FieldCondition{{{
				Field:     "source_type",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "resource"),
			}}}},
		},
		{
			name: "given invalid source_type network should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "source_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("network")},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"network" is not a valid value for "source_type"`,
		},
		{
			name: "given invalid source_type in is_one_of should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "source_type", Operator: pattern.ConditionIsOneOf,
						Value: ValueSpec{Strings: []string{"resource", "network"}},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"network" is not a valid value for "source_type"`,
		},
		{
			name: "given contain operator on source_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "source_type", Operator: pattern.ConditionContain,
						Value: ValueSpec{String: new("res")},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrUnsupportedConditionType,
			expectedContains: `field "source_type" only supports`,
		},
		{
			name: "given invalid event_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "event_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("ping")},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"ping" is not a valid value for "event_type"`,
		},
		{
			name: "given valid event_type should succeed",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "event_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("check")},
					}}},
				}},
			},
			promptContext: "eventfilter",
			expected: Result{EventPattern: [][]pattern.FieldCondition{{{
				Field:     "event_type",
				Condition: pattern.NewStringCondition(pattern.ConditionEqual, "check"),
			}}}},
		},
		{
			name: "given invalid entity type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "entity_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("host")},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"host" is not a valid value for "type"`,
		},
		{
			name: "given invalid initiator should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "initiator", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("admin")},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"admin" is not a valid value for "initiator"`,
		},
		{
			name: "given invalid canonical_type should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "pbehavior_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "pbehavior_info.canonical_type", Operator: pattern.ConditionEqual,
						Value: ValueSpec{String: new("down")},
					}}},
				}},
			},
			promptContext:    "pbehavior",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `"down" is not a valid value for "pbehavior_info.canonical_type"`,
		},
		{
			name: "given invalid status should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "alarm_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "v.status.val", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(99))},
					}}},
				}},
			},
			promptContext:    "widget_filter_AlarmsList",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `is not a valid value for "v.status.val"`,
		},
		{
			name: "given invalid event state should return error",
			input: Input{
				Patterns: []Spec{{
					Type: "event_pattern",
					Groups: [][]ConditionSpec{{{
						Field: "state", Operator: pattern.ConditionEqual,
						Value: ValueSpec{Integer: new(int64(5))},
					}}},
				}},
			},
			promptContext:    "eventfilter",
			expectedErr:      pattern.ErrWrongConditionValue,
			expectedContains: `is not a valid value for "state"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertBuildResult(t, tt.input, tt.promptContext, tt.expected, tt.expectedErr, tt.expectedContains)
		})
	}
}

func TestBuildCondition_UnsupportedType(t *testing.T) {
	_, err := buildCondition("not_a_type", pattern.ConditionEqual, ValueSpec{String: new("x")})
	if err == nil {
		t.Fatal("expected unsupported field type error, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported field type: not_a_type") {
		t.Errorf("error %q does not mention unsupported field type", err.Error())
	}
}

func TestResult_IsZero(t *testing.T) {
	if !(Result{}).IsZero() {
		t.Fatal("empty Result should be zero")
	}

	nonZero := []Result{
		{EventPattern: [][]pattern.FieldCondition{{{Field: "event_type"}}}},
		{AlarmPattern: [][]pattern.FieldCondition{{{Field: "v.display_name"}}}},
		{EntityPattern: [][]pattern.FieldCondition{{{Field: "name"}}}},
		{PbehaviorPattern: [][]pattern.FieldCondition{{{Field: "pbehavior_info.canonical_type"}}}},
		{WeatherServicePattern: [][]pattern.FieldCondition{{{Field: "state.val"}}}},
	}
	for i, r := range nonZero {
		if r.IsZero() {
			t.Errorf("case %d: expected non-zero Result", i)
		}
	}
}
