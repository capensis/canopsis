package utils

import "testing"

func TestToSnakeCase(t *testing.T) {
	dataSets := []struct {
		input    string
		expected string
	}{
		{
			input:    "",
			expected: "",
		},
		{
			input:    "TestCase",
			expected: "test_case",
		},
		{
			input:    "testCase",
			expected: "test_case",
		},
		{
			input:    "TESTCase",
			expected: "test_case",
		},
		{
			input:    "TestCASE",
			expected: "test_case",
		},
		{
			input:    "Test0Case",
			expected: "test0_case",
		},
		{
			input:    "TEST0Case",
			expected: "test0_case",
		},
		{
			input:    "test_case",
			expected: "test_case",
		},
	}

	for _, data := range dataSets {
		output := ToSnakeCase(data.input)

		if output != data.expected {
			t.Errorf("%v: expected %v but got %v", data.input, data.expected, output)
		}
	}
}
