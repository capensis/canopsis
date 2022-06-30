package entityinfodictionary

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type API interface {
	List(c *gin.Context)
}

type api struct {
	store  Store
	logger zerolog.Logger
}

func NewApi(
	store Store,
	logger zerolog.Logger,
) API {
	return &api{
		store:  store,
		logger: logger,
	}
}

// List info dictionary
// @Success 200 {object} common.PaginatedListResponse{data=[]Result}
func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	findResult, err := a.store.Find(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(request.Query, findResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}
