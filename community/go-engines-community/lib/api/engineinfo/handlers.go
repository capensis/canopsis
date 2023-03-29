package engineinfo

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/gin-gonic/gin"
)

func GetRunInfo(manager engine.RunInfoManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		infos, err := manager.GetEngineQueues(c)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, infos)
	}
}
