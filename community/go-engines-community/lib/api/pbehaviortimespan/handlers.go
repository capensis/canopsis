package pbehaviortimespan

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

// GetTimeSpans
// @Param body body TimespansRequest true "body"
// @Success 200 {array} ItemResponse
func GetTimeSpans(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TimespansRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		res, err := s.GetTimespans(c, request)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, res)
	}
}
