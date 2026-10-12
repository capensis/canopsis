package messageratestats

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder response.ErrorResponder
}

func NewApi(store Store, errorResponder response.ErrorResponder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} StatsListResponse
func (a *api) List(c *gin.Context) {
	var r = ListRequest{}

	if err := binding.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	stats, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response := StatsListResponse{
		Data: stats,
	}

	if r.Interval == IntervalHour {
		response.Meta.DeletedBefore, err = a.store.GetDeletedBeforeForHours(c)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}
	}

	c.JSON(http.StatusOK, response)
}
