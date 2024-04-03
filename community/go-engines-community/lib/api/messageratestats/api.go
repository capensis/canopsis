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

// List
// @Success 200 {object} StatsListResponse
func (a *api) List(c *gin.Context) {
	var r = ListRequest{}

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	stats, err := a.store.Find(c, r)
	if err != nil {
		panic(err)
	}

	response := StatsListResponse{
		Data: stats,
	}

	if r.Interval == IntervalHour {
		response.Meta.DeletedBefore, err = a.store.GetDeletedBeforeForHours(c)
		if err != nil {
			panic(err)
		}
	}

	c.JSON(http.StatusOK, response)
}
