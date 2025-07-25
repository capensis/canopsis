package entityupstream

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

type API interface {
	GetDownstreams(c *gin.Context)
	GetUpstream(c *gin.Context)
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

// GetDownstreams
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) GetDownstreams(c *gin.Context) {
	var r DownstreamsRequest
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	aggregationResult, err := a.store.GetDownstreams(c, r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))

		return
	}

	c.JSON(http.StatusOK, res)
}

// GetUpstream
// @Success 200 {object} Response
func (a *api) GetUpstream(c *gin.Context) {
	var r UpstreamRequest
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))

		return
	}

	res, entityExists, err := a.store.GetUpstream(c, r.ID)
	if err != nil {
		panic(err)
	}

	if !entityExists {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)

		return
	}

	if res == nil {
		c.JSON(http.StatusOK, map[string]string{})

		return
	}

	c.JSON(http.StatusOK, res)
}
