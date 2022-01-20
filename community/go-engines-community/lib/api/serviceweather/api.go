package serviceweather

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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

// Find all services
// @Summary Find all services
// @Description Get paginated list of services
// @Tags weather-services
// @ID weather-services-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param filter query string false "filter query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Param category query string false "category"
// @Success 200 {object} common.PaginatedListResponse{data=[]Service}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /weather-services [get]
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), query)
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

// Find all entity by service id
// @Summary Find all entity by service id
// @Description Get paginated list of entities
// @Tags weather-services
// @ID weather-services-find-all-entities
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "service id"
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Param with_instructions query boolean false "show assigned instructions and execution flags"
// @Success 200 {object} common.PaginatedListResponse{data=[]Entity}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /weather-services/{id} [get]
func (a *api) EntityList(c *gin.Context) {
	var query EntitiesListRequest
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
	aggregationResult, err := a.store.FindEntities(c.Request.Context(), id, apiKey, query)
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
