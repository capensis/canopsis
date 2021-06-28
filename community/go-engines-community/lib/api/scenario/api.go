package scenario

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"git.canopsis.net/canopsis/go-engines/lib/api/logger"
	"git.canopsis.net/canopsis/go-engines/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	store Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) common.CrudAPI {
	return &api{
		store: store,
		actionLogger: actionLogger,
	}
}

// Find all scenarios
// @Summary Find scenarios
// @Description Get paginated list of scenarios
// @Tags scenarios
// @ID scenarios-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Scenario}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /scenarios [get]
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	scenarios, err := a.store.Find(query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, scenarios)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get scenario by id
// @Summary Get scenario by id
// @Description Get scenario by id
// @Tags scenarios
// @ID scenarios-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Success 200 {object} Scenario
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [get]
func (a *api) Get(c *gin.Context) {
	scenario, err := a.store.GetOneBy(c.Param("id"))
	if err != nil {
		panic(err)
	}
	if scenario == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, scenario)
}

// Create scenario
// @Summary Create scenario
// @Description Create scenario
// @Tags scenarios
// @ID scenarios-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Scenario
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /scenarios [post]
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	scenario, err := a.store.Insert(request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeScenario,
		ValueID:   scenario.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, scenario)
}

// Update scenario by id
// @Summary Update scenario by id
// @Description Update scenario by id
// @Tags scenarios
// @ID scenarios-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Scenario
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	scenario, err := a.store.Update(request)
	if err != nil {
		panic(err)
	}

	if scenario == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeScenario,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, scenario)
}

// Delete scenario by id
// @Summary Delete scenario by id
// @Description Delete scenario by id
// @Tags scenarios
// @ID scenarios-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "scenario id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /scenarios/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeScenario,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}
