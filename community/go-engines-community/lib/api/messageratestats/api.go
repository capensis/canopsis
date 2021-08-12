package messageratestats

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @Success 200 {object} common.PaginatedListResponse{data=[]StatsListResponse}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /message-rate-stats [get]
func (a *api) List(c *gin.Context) {
	var r = ListRequest{}
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	stats, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, stats)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	data, ok := res.Data.([]StatsResponse)
	if !ok {
		panic(fmt.Errorf("paginated response data should be []StatsResponse"))
	}

	response := StatsListResponse{}
	response.Data = data
	response.Meta.PaginatedMeta = res.Meta

	if r.Interval == IntervalHour {
		response.Meta.DeletedBefore, err = a.store.GetDeletedBeforeForHours(c.Request.Context())
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, response)
}
