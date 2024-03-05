package entityinfodictionary

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	ListKeys(c *gin.Context)
	ListValues(c *gin.Context)
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

// List info dictionary keys
// @Success 200 {object} common.PaginatedListResponse{data=[]Result}
func (a *api) ListKeys(c *gin.Context) {
	var request ListKeysRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	findResult, err := a.store.FindKeys(c.Request.Context(), request)
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

// List info dictionary values
// @Success 200 {object} common.PaginatedListResponse{data=[]Result}
func (a *api) ListValues(c *gin.Context) {
	var request ListValuesRequest
	request.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	findResult, err := a.store.FindValues(c.Request.Context(), request)
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
