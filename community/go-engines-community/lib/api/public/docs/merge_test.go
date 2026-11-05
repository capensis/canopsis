package docs

import (
	"os"
	"regexp"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"go.yaml.in/yaml/v3"
)

func TestMerge_GivenGeneratedFieldsDefinitions_ShouldResolveAllRefs(t *testing.T) {
	content := mustRead(t, "swagger.yaml")
	merged, err := docs.MergeContent([][]byte{content})
	if err != nil {
		t.Fatal(err)
	}

	merged, err = docs.MergeSchemas([][]byte{mustRead(t, "schemas_swagger.yaml"), mustRead(t, "fields_swagger.yaml")}, merged)
	if err != nil {
		t.Fatal(err)
	}

	b, err := yaml.Marshal(merged)
	if err != nil {
		t.Fatal(err)
	}

	if out := os.Getenv("MERGED_SWAGGER_OUT"); out != "" {
		if err := os.WriteFile(out, b, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	if regexp.MustCompile(`_swagger\.yaml#`).Match(b) {
		t.Error("merged swagger contains unresolved file refs")
	}

	schemas := merged["components"].(map[string]any)["schemas"].(map[string]any)
	for _, m := range regexp.MustCompile(`#/components/schemas/([^'"\s]+)`).FindAllSubmatch(b, -1) {
		if _, ok := schemas[string(m[1])]; !ok {
			t.Errorf("unresolved ref %s", m[0])
		}
	}

	for _, name := range []string{"alarm.Sort", "alarm.SearchByField", "alarm.StepSort", "entity.Sort", "entity.SearchByField"} {
		if _, ok := schemas[name]; !ok {
			t.Errorf("missing schema %s", name)
		}
	}

	paths := merged["paths"].(map[string]any)
	for path, sortSchema := range map[string]string{
		"/alarms/search":   "alarm.Sort",
		"/alarm-exports":   "alarm.Sort",
		"/entities/search": "entity.Sort",
		"/entity-exports":  "entity.Sort",
	} {
		if ref := postBodySortRef(paths, path); ref != "#/components/schemas/"+sortSchema {
			t.Errorf("expected POST %s body sort to reference %s but got %q", path, sortSchema, ref)
		}
	}
}

func postBodySortRef(paths map[string]any, path string) string {
	defer func() { _ = recover() }()

	schema := paths[path].(map[string]any)["post"].(map[string]any)["requestBody"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
	for _, item := range schema["allOf"].([]any) {
		if props, ok := item.(map[string]any)["properties"].(map[string]any); ok {
			return props["sort"].(map[string]any)["$ref"].(string)
		}
	}

	return ""
}

func mustRead(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	return b
}
