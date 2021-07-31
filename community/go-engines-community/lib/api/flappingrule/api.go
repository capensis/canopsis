package flappingrule

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

// Create flappingrule
// @Summary Create flappingrule
// @Description Create flappingrule
// @Tags flappingrules
// @ID flappingrules-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body CreateRequest true "body"
// @Success 201 {object} CreateRequest
// @Failure 400 {object} common.ErrorResponse
// @Router /baggot-rules [post]
func (a api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	rule, err := a.store.Insert(c.Request.Context(), &request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeFlappingRule,
		ValueID:   request.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, rule)
}

// Find all flappingrule
// @Summary Find all flappingrule
// @Description Get paginated list of flappingrule
// @Tags flappingrules
// @ID flappingrules-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Success 200 {object} common.PaginatedListResponse{data=[]CreateRequest}
// @Failure 400 {object} common.ErrorResponse
// @Router /baggot-rules [get]
func (a api) List(c *gin.Context) {
	var query FilteredQuery
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

// Get flappingrule by id
// @Summary Get flappingrule by id
// @Description Get flappingrule by id
// @Tags flappingrules
// @ID flappingrules-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "flappingrule id"
// @Success 200 {object} CreateRequest
// @Failure 404 {object} common.ErrorResponse
// @Router /baggot-rules/{id} [get]
func (a api) Get(c *gin.Context) {
	rule, err := a.store.GetById(c.Request.Context(), c.Param("id"))

	if err != nil {
		panic(err)
	}

	if rule == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, rule)
}

// Update flappingrule by id
// @Summary Update flappingrule by id
// @Description Update flappingrule by id
// @Tags flappingrules
// @ID flappingrules-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "flappingrule id"
// @Param body body CreateRequest true "body"
// @Success 200 {object} CreateRequest
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /baggot-rules/{id} [put]
func (a api) Update(c *gin.Context) {
	var request UpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	var data = request.CreateRequest
	data.ID = c.Param("id")
	rule, err := a.store.Update(c.Request.Context(), &data)

	if err != nil {
		panic(err)
	}

	if rule == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeFlappingRule,
		ValueID:   data.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, rule)
}

// Delete flappingrule by id
// @Summary Delete flappingrule by id
// @Description Delete flappingrule by id
// @Tags flappingrules
// @ID flappingrules-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth`
// @Param id path string true "flappingrule id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /baggot-rules/{id} [delete]
func (a api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeFlappingRule,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) common.CrudAPI {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		logger:       logger,
	}
}
