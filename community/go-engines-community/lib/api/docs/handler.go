package docs

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetHandler(errorResponder httperror.Responder, generatedSchemasContents [][]byte, contents [][]byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		mergedContent, err := mergeContent(contents)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		mergedContent, err = mergeSchemas(generatedSchemasContents, mergedContent)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		if info, ok := mergedContent["info"].(map[string]interface{}); ok {
			buildInfo := canopsis.GetBuildInfo()

			info["version"] = buildInfo.Version
			mergedContent["info"] = info
		}

		c.YAML(http.StatusOK, mergedContent)
	}
}

func mergeContent(contents [][]byte) (map[string]any, error) {
	var mergedContent map[string]any
	for _, content := range contents {
		replacedSchemasContent := strings.ReplaceAll(string(content), "schemas_swagger.yaml#/definitions/", "#/components/schemas/")
		var parsed map[string]interface{}
		err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
		if err != nil {
			return nil, err
		}

		if mergedContent == nil {
			mergedContent = parsed
			continue
		}

		if paths, ok := parsed["paths"].(map[string]interface{}); ok {
			if mergedPaths, ok := mergedContent["paths"].(map[string]interface{}); ok {
				for k, v := range paths {
					mergedPaths[k] = v
				}
			} else {
				mergedContent["paths"] = paths
			}
		}

		if components, ok := parsed["components"].(map[string]interface{}); ok {
			if schemas, ok := components["schemas"].(map[string]interface{}); ok {
				if mergedComponents, ok := mergedContent["components"].(map[string]interface{}); ok {
					if mergedSchemas, ok := mergedComponents["schemas"].(map[string]interface{}); ok {
						for k, v := range schemas {
							mergedSchemas[k] = v
						}
					} else {
						mergedComponents["schemas"] = schemas
					}
				}
			}
		}
	}

	return mergedContent, nil
}

func mergeSchemas(generatedSchemasContents [][]byte, mergedContent map[string]any) (map[string]any, error) {
	for _, generatedSchemasContent := range generatedSchemasContents {
		replacedSchemasContent := strings.ReplaceAll(string(generatedSchemasContent), "#/definitions/", "#/components/schemas/")
		var parsed map[string]interface{}
		err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
		if err != nil {
			return nil, err
		}

		if schemas, ok := parsed["definitions"].(map[string]interface{}); ok {
			if mergedComponents, ok := mergedContent["components"].(map[string]interface{}); ok {
				if mergedSchemas, ok := mergedComponents["schemas"].(map[string]interface{}); ok {
					for k, v := range schemas {
						mergedSchemas[k] = v
					}
				} else {
					mergedComponents["schemas"] = schemas
				}
			}
		}
	}

	return mergedContent, nil
}
