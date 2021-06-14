package engineinfo

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/engine"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetRunInfo(ctx context.Context, manager engine.RunInfoManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		infos, err := manager.GetAll(ctx)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, infos)
	}
}
