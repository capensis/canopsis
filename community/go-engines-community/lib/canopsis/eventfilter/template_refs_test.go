package eventfilter

import (
	"reflect"
	"testing"
	"text/template"

	libtemplate "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
)

func TestTemplateExternalDataRefs(t *testing.T) {
	refs := []string{"missing", "matched", "missing_extra", "ref.with.dot"}
	for _, test := range []struct {
		name     string
		text     string
		expected []string
	}{
		{
			name:     "direct field",
			text:     `{{ .ExternalData.missing.status }}`,
			expected: []string{"missing"},
		},
		{
			name:     "reference without child field",
			text:     `{{ .ExternalData.missing }}`,
			expected: []string{"missing"},
		},
		{
			name:     "index",
			text:     `{{ index .ExternalData "missing" }}`,
			expected: []string{"missing"},
		},
		{
			name:     "nested index",
			text:     `{{ index .ExternalData "missing" "status" }}`,
			expected: []string{"missing"},
		},
		{
			name:     "index followed by field",
			text:     `{{ (index .ExternalData "missing").status }}`,
			expected: []string{"missing"},
		},
		{
			name:     "numeric index within reference",
			text:     `{{ index .ExternalData.missing.items 0 }}`,
			expected: []string{"missing"},
		},
		{
			name:     "reference name containing dots",
			text:     `{{ index .ExternalData "ref.with.dot" }}`,
			expected: []string{"ref.with.dot"},
		},
		{
			name:     "similar reference names remain distinct",
			text:     `{{ .ExternalData.missing_extra.status }}`,
			expected: []string{"missing_extra"},
		},
		{
			name:     "pipeline",
			text:     `{{ .ExternalData.missing.status | printf "%s" }}`,
			expected: []string{"missing"},
		},
		{
			name:     "function arguments",
			text:     `{{ printf "%s %s" .ExternalData.matched.status .ExternalData.missing.status }}`,
			expected: []string{"missing", "matched"},
		},
		{
			name:     "result follows reference order and includes each reference once",
			text:     `{{ .ExternalData.matched.status }}{{ .ExternalData.missing.status }}{{ .ExternalData.matched.name }}`,
			expected: []string{"missing", "matched"},
		},
		{
			name:     "with changes dot to external data",
			text:     `{{ with .ExternalData }}{{ .missing.status }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "with only accesses another reference",
			text:     `{{ with .ExternalData }}{{ .matched.status }}{{ end }}`,
			expected: []string{"matched"},
		},
		{
			name:     "nested with",
			text:     `{{ with .ExternalData }}{{ with .missing }}{{ .status }}{{ end }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "with else retains outer dot",
			text:     `{{ with .Event }}{{ .Output }}{{ else }}{{ .ExternalData.missing.status }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "dot is restored after with",
			text:     `{{ with .ExternalData }}{{ .missing.status }}{{ end }}{{ .ExternalData.matched.status }}`,
			expected: []string{"missing", "matched"},
		},
		{
			name:     "root variable inside with",
			text:     `{{ with .Event }}{{ $.ExternalData.missing.status }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "external data variable",
			text:     `{{ $data := .ExternalData }}{{ index $data "missing" }}`,
			expected: []string{"missing"},
		},
		{
			name:     "reference variable",
			text:     `{{ $row := index .ExternalData "missing" }}{{ $row.status }}`,
			expected: []string{"missing"},
		},
		{
			name:     "with variable declaration",
			text:     `{{ with $data := .ExternalData }}{{ $data.missing.status }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "variable shadowing stays inside with",
			text:     `{{ $data := .ExternalData }}{{ with .Event }}{{ $data := .ExtraInfos }}{{ $data.missing }}{{ end }}{{ $data.matched.status }}`,
			expected: []string{"matched"},
		},
		{
			name:     "if condition and both branches",
			text:     `{{ if .ExternalData.missing.enabled }}{{ .ExternalData.matched.status }}{{ else }}{{ .ExternalData.missing_extra.status }}{{ end }}`,
			expected: []string{"missing", "matched", "missing_extra"},
		},
		{
			name:     "if variable declaration keeps outer dot",
			text:     `{{ if $data := .ExternalData }}{{ $data.missing.status }}{{ .Event.Output }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "if branch variable shadowing does not affect else",
			text:     `{{ $data := .ExternalData }}{{ if .Event.State }}{{ $data := .Event }}{{ $data.missing }}{{ else }}{{ $data.matched.status }}{{ end }}`,
			expected: []string{"matched"},
		},
		{
			name:     "range within one reference",
			text:     `{{ range .ExternalData.missing.items }}{{ .name }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "range value variable",
			text:     `{{ range $row := .ExternalData }}{{ $row.status }}{{ end }}`,
			expected: refs,
		},
		{
			name: "range key is not a reference value",
			text: `{{ range $key, $row := .ExternalData }}{{ $key.ExternalData.missing }}{{ end }}`,
		},
		{
			name:     "range key and value variables",
			text:     `{{ range $key, $row := .ExternalData }}{{ $key }}{{ $row.status }}{{ end }}`,
			expected: refs,
		},
		{
			name:     "range else keeps outer dot",
			text:     `{{ range .Event.Tags }}{{ . }}{{ else }}{{ .ExternalData.missing.status }}{{ end }}`,
			expected: []string{"missing"},
		},
		{
			name:     "named template receives external data as dot",
			text:     `{{ define "info" }}{{ .missing.status }}{{ end }}{{ template "info" .ExternalData }}`,
			expected: []string{"missing"},
		},
		{
			name:     "named template receives a reference as root",
			text:     `{{ define "info" }}{{ $.status }}{{ end }}{{ template "info" .ExternalData.missing }}`,
			expected: []string{"missing"},
		},
		{
			name:     "named template called with different scopes",
			text:     `{{ define "info" }}{{ .missing.status }}{{ end }}{{ template "info" .Event }}{{ template "info" .ExternalData }}`,
			expected: []string{"missing"},
		},
		{
			name:     "recursive template terminates dependency traversal",
			text:     `{{ define "info" }}{{ .missing.status }}{{ template "info" . }}{{ end }}{{ template "info" .ExternalData }}`,
			expected: []string{"missing"},
		},
		{
			name: "unused named template",
			text: `{{ define "unused" }}{{ .ExternalData.missing.status }}{{ end }}literal`,
		},
		{
			name: "named template without data has no root reference",
			text: `{{ define "info" }}{{ .ExternalData.missing.status }}{{ end }}{{ template "info" }}`,
		},
		{
			name: "undefined named template does not panic",
			text: `{{ template "undefined" .ExternalData }}`,
		},
		{
			name: "other data",
			text: `{{ .Event.Output }}`,
		},
		{
			name: "external data field inside another object",
			text: `{{ with .Event }}{{ .ExternalData.missing.status }}{{ end }}`,
		},
		{
			name: "unknown reference",
			text: `{{ .ExternalData.unknown.status }}`,
		},
		{
			name: "comment",
			text: `{{/* .ExternalData.missing.status */}}literal`,
		},
		{
			name: "quoted text",
			text: `{{ printf "%s" ".ExternalData.missing.status" }}`,
		},
		{
			name: "literal text",
			text: `.ExternalData.missing.status`,
		},
		{
			name: "empty template",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			tpl, err := template.New("test").Parse(test.text)
			if err != nil {
				t.Fatal(err)
			}
			actual := templateExternalDataRefs(libtemplate.ParsedTemplate{Tpl: tpl}, refs)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("expected references %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestTemplateExternalDataRefsWithoutTree(t *testing.T) {
	for _, tpl := range []libtemplate.ParsedTemplate{
		{},
		{Tpl: template.New("test")},
	} {
		if actual := templateExternalDataRefs(tpl, []string{"missing"}); actual != nil {
			t.Fatalf("expected nil references, got %v", actual)
		}
	}
}

func TestTemplateExternalDataRefsWithoutReferences(t *testing.T) {
	tpl, err := template.New("test").Parse(`{{ .ExternalData.missing.status }}`)
	if err != nil {
		t.Fatal(err)
	}
	for _, refs := range [][]string{nil, {}} {
		if actual := templateExternalDataRefs(libtemplate.ParsedTemplate{Tpl: tpl}, refs); actual != nil {
			t.Fatalf("expected nil references, got %v", actual)
		}
	}
}
