package docs

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/docs"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis"
	"github.com/gin-gonic/gin"
)

func GetHandler(errorResponder response.ErrorResponder, generatedSchemasContents [][]byte, contents [][]byte) func(c *gin.Context) {
	return func(c *gin.Context) {
		mergedContent, err := docs.MergeContent(contents)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		mergedContent, err = docs.MergeSchemas(generatedSchemasContents, mergedContent)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		if info, ok := mergedContent["info"].(map[string]any); ok {
			buildInfo := canopsis.GetBuildInfo()

			info["version"] = buildInfo.Version
			mergedContent["info"] = info
		}

		c.YAML(http.StatusOK, mergedContent)
	}
}
