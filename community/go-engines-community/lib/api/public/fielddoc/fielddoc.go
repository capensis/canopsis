// Package fielddoc describes the fields accepted by sort and search_by query parameters
// so that OpenAPI schemas can be generated from the same lists which are used for validation.
package fielddoc

import (
	"regexp"
	"slices"
	"strings"
)

type Fields struct {
	// Static contains exact field names.
	Static []string
	// Patterns contains unanchored regular expressions of dynamic field names (e.g. infos).
	Patterns []string
	// Labels contains human-readable allowed values (e.g. wildcard patterns) shown in descriptions.
	// Static is used when Labels is empty.
	Labels []string
}

// Regexp returns an unanchored regular expression which matches one field.
func (f Fields) Regexp() string {
	alts := make([]string, 0, len(f.Static)+len(f.Patterns))
	for _, s := range f.Static {
		alts = append(alts, regexp.QuoteMeta(s))
	}

	alts = append(alts, f.Patterns...)

	return "(?:" + strings.Join(alts, "|") + ")"
}

// SortRegexp returns an anchored regular expression which matches a comma separated list of field[:asc|desc].
func (f Fields) SortRegexp() string {
	item := f.Regexp() + "(?::(?:asc|desc))?"

	return "^" + item + "(?:," + item + ")*$"
}

// FieldSchema returns an OpenAPI 3 schema of one field.
func (f Fields) FieldSchema(description string) map[string]any {
	schema := map[string]any{
		"type":        "string",
		"description": description + "\n\n" + f.AvailableValues(),
	}
	if len(f.Patterns) == 0 {
		schema["enum"] = slices.Clone(f.Static)

		return schema
	}

	anyOf := make([]any, 0, len(f.Patterns)+1)
	if len(f.Static) > 0 {
		anyOf = append(anyOf, map[string]any{"type": "string", "enum": slices.Clone(f.Static)})
	}

	for _, p := range f.Patterns {
		anyOf = append(anyOf, map[string]any{"type": "string", "pattern": "^" + p + "$"})
	}

	schema["anyOf"] = anyOf

	return schema
}

// SortSchema returns an OpenAPI 3 schema of the sort query parameter.
func (f Fields) SortSchema(description string) map[string]any {
	return map[string]any{
		"type":        "string",
		"pattern":     f.SortRegexp(),
		"description": description + "\n\n" + f.AvailableValues(),
	}
}

// AvailableValues returns a Markdown line listing the human-readable allowed values of the field.
func (f Fields) AvailableValues() string {
	labels := f.Labels
	if len(labels) == 0 {
		labels = f.Static
	}

	quoted := make([]string, 0, len(labels))
	for _, label := range labels {
		quoted = append(quoted, "`"+label+"`")
	}

	return "<i>Available values</i> : " + strings.Join(quoted, ", ")
}

// SortParameter returns a sort query parameter which references the schema definition.
func (f Fields) SortParameter(definition, description, example string) map[string]any {
	param := map[string]any{
		"name":        "sort",
		"in":          "query",
		"description": description + "\n\n" + f.AvailableValues(),
		"schema":      map[string]any{"$ref": "#/definitions/" + definition},
	}
	if example != "" {
		param["example"] = example
	}

	return param
}

// SearchByParameter returns a search_by query parameter whose items reference the schema definition.
func (f Fields) SearchByParameter(definition, description string) map[string]any {
	return map[string]any{
		"name":        "search_by",
		"in":          "query",
		"style":       "form",
		"explode":     true,
		"description": description + "\n\n" + f.AvailableValues(),
		"schema": map[string]any{
			"type":  "array",
			"items": map[string]any{"$ref": "#/definitions/" + definition},
		},
	}
}
