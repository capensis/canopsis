package docs

import (
	"regexp"
	"strings"
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/public/alarm"
	"go.yaml.in/yaml/v3"
)

func TestMerge_GivenGeneratedParameters_ShouldResolveAllParameterRefs(t *testing.T) {
	merged, err := docs.MergeContent([][]byte{mustRead(t, "swagger.yaml")})
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

	parameters := merged["components"].(map[string]any)["parameters"].(map[string]any)
	for _, name := range []string{"offset", "alarm.Sort", "alarm.SearchBy", "alarm.StepSort", "entity.Sort", "entity.SearchBy"} {
		if _, ok := parameters[name]; !ok {
			t.Errorf("missing parameter %s", name)
		}
	}

	for _, m := range regexp.MustCompile(`#/components/parameters/([^'"\s]+)`).FindAllSubmatch(b, -1) {
		if _, ok := parameters[string(m[1])]; !ok {
			t.Errorf("unresolved ref %s", m[0])
		}
	}

	paths := merged["paths"].(map[string]any)
	for path, want := range map[string][]string{
		"/alarms":            {"alarm.Sort", "alarm.SearchBy"},
		"/alarms/{id}/steps": {"alarm.StepSort"},
		"/entities":          {"entity.Sort", "entity.SearchBy"},
	} {
		params := paths[path].(map[string]any)["get"].(map[string]any)["parameters"].([]any)
		for _, name := range want {
			if !hasParameterRef(params, "#/components/parameters/"+name) {
				t.Errorf("expected GET %s to reference parameter %s", path, name)
			}
		}
	}

	sortLabels := strings.Fields(alarm.SortErrParam())
	for i, label := range sortLabels {
		sortLabels[i] = "`" + label + "`"
	}
	wantLabels := "<i>Available values</i> : " + strings.Join(sortLabels, ", ")
	if desc, _ := parameters["alarm.Sort"].(map[string]any)["description"].(string); !strings.Contains(desc, wantLabels) {
		t.Errorf("expected alarm.Sort description to contain %q, got %q", wantLabels, desc)
	} else {
		t.Logf("labels %q", wantLabels)
	}
}

func hasParameterRef(params []any, ref string) bool {
	for _, p := range params {
		if m, ok := p.(map[string]any); ok && m["$ref"] == ref {
			return true
		}
	}

	return false
}
