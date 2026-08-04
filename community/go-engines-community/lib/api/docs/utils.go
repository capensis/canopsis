package docs

import (
	"bytes"
	"embed"
	"maps"
	"strings"

	"go.yaml.in/yaml/v3"
)

//go:embed static/*
var UIFiles embed.FS

func MergeContent(contents [][]byte) (map[string]any, error) {
	var mergedContent map[string]any
	for _, content := range contents {
		replacedSchemasContent := bytes.ReplaceAll(
			content,
			[]byte(`schemas_swagger.yaml#/definitions/`),
			[]byte(`#/components/schemas/`),
		)
		var parsed map[string]any
		err := yaml.Unmarshal(replacedSchemasContent, &parsed)
		if err != nil {
			return nil, err
		}

		if mergedContent == nil {
			mergedContent = parsed
			continue
		}

		if paths, ok := parsed["paths"].(map[string]any); ok {
			if mergedPaths, ok := mergedContent["paths"].(map[string]any); ok {
				maps.Copy(mergedPaths, paths)
			} else {
				mergedContent["paths"] = paths
			}
		}

		if components, ok := parsed["components"].(map[string]any); ok {
			if schemas, ok := components["schemas"].(map[string]any); ok {
				if mergedComponents, ok := mergedContent["components"].(map[string]any); ok {
					if mergedSchemas, ok := mergedComponents["schemas"].(map[string]any); ok {
						maps.Copy(mergedSchemas, schemas)
					} else {
						mergedComponents["schemas"] = schemas
					}
				}
			}
		}
	}

	return mergedContent, nil
}

func MergeSchemas(generatedSchemasContents [][]byte, mergedContent map[string]any) (map[string]any, error) {
	for _, generatedSchemasContent := range generatedSchemasContents {
		replacedSchemasContent := strings.ReplaceAll(string(generatedSchemasContent), "#/definitions/", "#/components/schemas/")
		var parsed map[string]any
		err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
		if err != nil {
			return nil, err
		}

		if schemas, ok := parsed["definitions"].(map[string]any); ok {
			if mergedComponents, ok := mergedContent["components"].(map[string]any); ok {
				if mergedSchemas, ok := mergedComponents["schemas"].(map[string]any); ok {
					maps.Copy(mergedSchemas, schemas)
				} else {
					mergedComponents["schemas"] = schemas
				}
			}
		}
	}

	return mergedContent, nil
}
