package pbehaviortimespan

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

// GetTimeSpans
// @Param body body TimespansRequest true "body"
// @Success 200 {array} ItemResponse
func GetTimeSpans(s Service, errorResponder httperror.Responder) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TimespansRequest

		if err := validation.Bind(c, &request); err != nil {
			errorResponder.Respond(c, err)

			return
		}

		res, err := s.GetTimespans(c, request)
		if err != nil {
			errorResponder.Respond(c, err)

			return
		}

		c.JSON(http.StatusOK, res)
	}
}
