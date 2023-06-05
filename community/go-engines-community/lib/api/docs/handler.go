package docs

import (
	"net/http"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func GetHandler(generatedSchemasContent []byte, contents ...[]byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		var mergedContent map[string]interface{}

		for _, content := range contents {
			replacedSchemasContent := strings.ReplaceAll(string(content), "schemas_swagger.yaml#/definitions/", "#/components/schemas/")
			var parsed map[string]interface{}
			err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
			if err != nil {
				panic(err)
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

		replacedSchemasContent := strings.ReplaceAll(string(generatedSchemasContent), "#/definitions/", "#/components/schemas/")
		var parsed map[string]interface{}
		err := yaml.Unmarshal([]byte(replacedSchemasContent), &parsed)
		if err != nil {
			panic(err)
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

		if info, ok := mergedContent["info"].(map[string]interface{}); ok {
			buildInfo := canopsis.GetBuildInfo()

			info["version"] = buildInfo.Version
			mergedContent["info"] = info
		}

		c.YAML(http.StatusOK, mergedContent)
	}
}
