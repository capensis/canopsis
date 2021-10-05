package messageratestats

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(c *gin.Context)
}

type api struct {
	store Store
}

func NewApi(
	store Store,
) API {
	return &api{
		store: store,
	}
}

// Find message rate stats
// @Summary Find message rate stats
// @Description Get paginated list of stats
// @Tags message-rate-stats
// @ID message-rate-stats-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query ListRequest true "request"
// @Success 200 {object} StatsListResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /message-rate-stats [get]
func (a *api) List(c *gin.Context) {
	var r = ListRequest{}

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	stats, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	response := StatsListResponse{
		Data: stats,
	}

	if r.Interval == IntervalHour {
		response.Meta.DeletedBefore, err = a.store.GetDeletedBeforeForHours(c.Request.Context())
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, response)
}
