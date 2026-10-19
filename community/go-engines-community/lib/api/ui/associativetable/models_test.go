package associativetable

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestContentUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected any
	}{
		"object": {
			input:    `{"a":1,"b":"x"}`,
			expected: map[string]any{"a": float64(1), "b": "x"},
		},
		"array": {
			input:    `[1,"x",true]`,
			expected: []any{float64(1), "x", true},
		},
		"string": {
			input:    `"x"`,
			expected: "x",
		},
		"null": {
			input:    `null`,
			expected: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var c Content
			if err := json.Unmarshal([]byte(tc.input), &c); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(c.value, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, c.value)
			}
		})
	}
}

func TestContentUnmarshalJSONInvalidKeepsPreviousValue(t *testing.T) {
	c := Content{value: "previous"}
	if err := c.UnmarshalJSON([]byte(`{"a":`)); err == nil {
		t.Fatal("expected error")
	}
	if c.value != "previous" {
		t.Fatalf("expected previous value to be kept, got %#v", c.value)
	}
}

func TestContentUnmarshalJSONReplacesPreviousValue(t *testing.T) {
	c := Content{value: map[string]any{"stale": true}}
	if err := c.UnmarshalJSON([]byte(`{"fresh":1}`)); err != nil {
		t.Fatal(err)
	}
	expected := map[string]any{"fresh": float64(1)}
	if !reflect.DeepEqual(c.value, expected) {
		t.Fatalf("expected %#v, got %#v", expected, c.value)
	}
}

var benchContentInput = []byte(`{"name":"table","rows":[{"id":1,"value":"a"},{"id":2,"value":"b"}],"enabled":true}`)

func BenchmarkContentUnmarshalJSON(b *testing.B) {
	var c Content
	b.ReportAllocs()

	for b.Loop() {
		if err := c.UnmarshalJSON(benchContentInput); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAssociativeTableUnmarshalJSON(b *testing.B) {
	input := []byte(`{"name":"t","content":{"name":"table","rows":[{"id":1,"value":"a"},{"id":2,"value":"b"}],"enabled":true}}`)
	b.ReportAllocs()

	for b.Loop() {
		var table AssociativeTable
		if err := json.Unmarshal(input, &table); err != nil {
			b.Fatal(err)
		}
	}
}
