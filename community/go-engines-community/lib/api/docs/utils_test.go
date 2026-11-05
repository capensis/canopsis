package docs

import "testing"

func TestMergeContent_GivenParameterRefs_ShouldRewriteToComponents(t *testing.T) {
	content := []byte(`paths:
  /alarms:
    get:
      parameters:
        - $ref: 'fields_swagger.yaml#/parameters/alarm.Sort'
        - $ref: 'fields_swagger.yaml#/definitions/alarm.Sort'
`)

	merged, err := MergeContent([][]byte{content})
	if err != nil {
		t.Fatal(err)
	}

	params := merged["paths"].(map[string]any)["/alarms"].(map[string]any)["get"].(map[string]any)["parameters"].([]any)
	if ref := params[0].(map[string]any)["$ref"]; ref != "#/components/parameters/alarm.Sort" {
		t.Errorf("unexpected parameter ref %v", ref)
	}

	if ref := params[1].(map[string]any)["$ref"]; ref != "#/components/schemas/alarm.Sort" {
		t.Errorf("unexpected schema ref %v", ref)
	}
}

func TestMergeSchemas_GivenGeneratedParameters_ShouldMergeIntoComponents(t *testing.T) {
	base := map[string]any{
		"components": map[string]any{
			"parameters": map[string]any{"offset": map[string]any{"name": "offset"}},
		},
	}
	generated := []byte(`definitions:
  alarm.Sort:
    type: string
parameters:
  alarm.Sort:
    name: sort
    in: query
    schema:
      $ref: '#/definitions/alarm.Sort'
`)

	merged, err := MergeSchemas([][]byte{generated}, base)
	if err != nil {
		t.Fatal(err)
	}

	components := merged["components"].(map[string]any)
	params := components["parameters"].(map[string]any)
	if _, ok := params["offset"]; !ok {
		t.Error("expected existing parameter offset to be kept")
	}

	param, ok := params["alarm.Sort"].(map[string]any)
	if !ok {
		t.Fatalf("expected parameter alarm.Sort, got %v", params)
	}

	if ref := param["schema"].(map[string]any)["$ref"]; ref != "#/components/schemas/alarm.Sort" {
		t.Errorf("unexpected schema ref %v", ref)
	}

	if _, ok := components["schemas"].(map[string]any)["alarm.Sort"]; !ok {
		t.Error("expected schema alarm.Sort to be merged")
	}
}
