package engineinfo

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/engine"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRunInfo(manager engine.RunInfoManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		infos, err := manager.GetAll()
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, infos)
	}
}
