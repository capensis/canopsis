package docs

import (
	"embed"
	"maps"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

//go:embed static/*
var UIFiles embed.FS

// generatedDefinitionsRefRegexp matches refs to generated files, e.g. schemas_swagger.yaml#/definitions/ or fields_swagger.yaml#/definitions/.
var generatedDefinitionsRefRegexp = regexp.MustCompile(`[A-Za-z0-9_]+_swagger\.yaml#/definitions/`)

// generatedParametersRefRegexp matches refs to parameters of generated files, e.g. fields_swagger.yaml#/parameters/.
var generatedParametersRefRegexp = regexp.MustCompile(`[A-Za-z0-9_]+_swagger\.yaml#/parameters/`)

func MergeContent(contents [][]byte) (map[string]any, error) {
	var mergedContent map[string]any
	for _, content := range contents {
		replacedSchemasContent := generatedDefinitionsRefRegexp.ReplaceAll(content, []byte(`#/components/schemas/`))
		replacedSchemasContent = generatedParametersRefRegexp.ReplaceAll(replacedSchemasContent, []byte(`#/components/parameters/`))
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

		if parameters, ok := parsed["parameters"].(map[string]any); ok {
			if mergedComponents, ok := mergedContent["components"].(map[string]any); ok {
				if mergedParameters, ok := mergedComponents["parameters"].(map[string]any); ok {
					maps.Copy(mergedParameters, parameters)
				} else {
					mergedComponents["parameters"] = parameters
				}
			}
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
