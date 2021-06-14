package watcherweather

import (
	"net/http"

	"git.canopsis.net/canopsis/go-engines/lib/api/auth"
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"github.com/gin-gonic/gin"
)

func NewApi(store Store) API {
	return &api{
		store: store,
	}
}

type API interface {
	List(c *gin.Context)
	EntityList(c *gin.Context)
}

type api struct {
	store Store
}

// Find all watchers
// @Summary Find all watchers
// @Description Get paginated list of watchers
// @Tags weather-watchers
// @ID weather-watchers-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param filter query string false "filter query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Watcher}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /weather-watchers [get]
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Find all entity by watcher id
// @Summary Find all entity by watcher id
// @Description Get paginated list of entities
// @Tags weather-watchers
// @ID weather-watchers-find-all-entities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "watcher id"
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Entity}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /weather-watchers/{id} [get]
func (a *api) EntityList(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	id := c.Param("id")
	apiKey := ""
	if key, ok := c.Get(auth.ApiKey); ok {
		apiKey = key.(string)
	}
	aggregationResult, err := a.store.FindEntities(id, apiKey, query)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}
