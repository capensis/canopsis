package pbehaviortimespan

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

// GetTimeSpans gets all pbehavior timespans
// @Summary Get all pbehavior timespans
// @Description Get time spans of calendar event within view span; {by_date: false} adds exception spans with types, {by_date: true} merges adjacent spans if gap between sequential ones less than 24 hours
// @Tags pbehavior-timespans
// @ID pbehavior-timespans-get-all
// @Accept json
// @Produce json
// @Param body body TimespansRequest true "body"
// @Success 200 {array} ItemResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-timespans [post]
func GetTimeSpans(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request TimespansRequest

		if err := c.ShouldBind(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
			return
		}

		res, err := s.GetTimespans(c.Request.Context(), request)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, res)
	}
}
