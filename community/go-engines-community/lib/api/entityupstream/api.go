package entityupstream

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

type API interface {
	GetDownstreams(c *gin.Context)
	GetUpstream(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

func NewApi(store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// GetDownstreams
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) GetDownstreams(c *gin.Context) {
	var r DownstreamsRequest
	r.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.GetDownstreams(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if aggregationResult == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// GetUpstream
// @Success 200 {object} Response
func (a *api) GetUpstream(c *gin.Context) {
	var r UpstreamRequest
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, entityExists, err := a.store.GetUpstream(c, r.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !entityExists {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if res == nil {
		c.JSON(http.StatusOK, map[string]string{})

		return
	}

	c.JSON(http.StatusOK, res)
}
