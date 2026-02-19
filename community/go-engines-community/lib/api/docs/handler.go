package docs

import (
	"maps"
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/gin-gonic/gin"
	"go.yaml.in/yaml/v3"
)

func GetHandler(errorResponder httperror.Responder, generatedSchemasContent []byte, contents ...[]byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		var mergedContent map[string]any

		for _, content := range contents {
			replacedSchemasContent := strings.ReplaceAll(string(content), "schemas_swagger.yaml#/definitions/", "#/components/schemas/")
			var parsed map[string]any
			err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
			if err != nil {
				errorResponder.Respond(c, err)

				return
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

		replacedSchemasContent := strings.ReplaceAll(string(generatedSchemasContent), "#/definitions/", "#/components/schemas/")
		var parsed map[string]any
		err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
		if err != nil {
			errorResponder.Respond(c, err)

			return
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

		if info, ok := mergedContent["info"].(map[string]any); ok {
			buildInfo := canopsis.GetBuildInfo()

			info["version"] = buildInfo.Version
			mergedContent["info"] = info
		}

		c.YAML(http.StatusOK, mergedContent)
	}
}
