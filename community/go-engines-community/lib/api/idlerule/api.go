package idlerule

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/scenario"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	common.CrudAPI
	CountPatterns(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	conf         config.UserInterfaceConfigProvider
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	conf config.UserInterfaceConfigProvider,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		conf:         conf,
	}
}

// Find all idle rules
// @Summary Find idle rules
// @Description Get paginated list of idle rules
// @Tags idle rules
// @ID idlerules-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]Rule}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /idle-rules [get]
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	rules, err := a.store.Find(c.Request.Context(), query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, rules)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get idle rule by id
// @Summary Get idle rule by id
// @Description Get idle rule by id
// @Tags idlerules
// @ID idlerules-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "rule id"
// @Success 200 {object} Rule
// @Failure 404 {object} common.ErrorResponse
// @Router /idle-rules/{id} [get]
func (a *api) Get(c *gin.Context) {
	rule, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if rule == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, rule)
}

// Create idle rule
// @Summary Create idle rule
// @Description Create idle rule
// @Tags idlerules
// @ID idlerules-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Rule
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /idle-rules [post]
func (a *api) Create(c *gin.Context) {
	var request EditRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	setOperationParameterAuthor(&request, request.Author)
	rule, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeIdleRule,
		ValueID:   rule.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, rule)
}

// Update idle rule by id
// @Summary Update idle rule by id
// @Description Update idle rule by id
// @Tags idlerules
// @ID idlerules-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "rule id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Rule
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /idle-rules/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	setOperationParameterAuthor(&request, request.Author)
	rule, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if rule == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeIdleRule,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, rule)
}

// Delete idle rule by id
// @Summary Delete idle rule by id
// @Description Delete idle rule by id
// @Tags idlerules
// @ID idlerules-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "rule id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /idle-rules/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeIdleRule,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// Count entities and alarm matching patterns
// @Summary Count entities and alarm matching patterns
// @Description Count entities and alarm matching patterns
// @Tags idlerules
// @ID idlerules-countpatterns
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body CountByPatternRequest true "body"
// @Success 200 {object} CountByPatternResult
// @Failure 400 {object} common.ErrorResponse
// @Failure 408 {object} common.ErrorResponse
// @Router /idle-rules/count [post]
func (a *api) CountPatterns(c *gin.Context) {
	var request CountByPatternRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	data, err := a.store.CountByPatterns(c.Request.Context(), request, a.conf.Get().CheckCountRequestTimeout, a.conf.Get().MaxMatchedItems)
	if errors.Is(err, context.DeadlineExceeded) {
		c.AbortWithStatusJSON(http.StatusRequestTimeout, common.ErrTimeoutResponse)
		return
	} else if err != nil {
		panic(err)
	}

	if int(data.TotalCountEntities) > a.conf.Get().MaxMatchedItems || int(data.TotalCountAlarms) > a.conf.Get().MaxMatchedItems {
		data.OverLimit = true
	}

	c.JSON(http.StatusOK, data)
}

func setOperationParameterAuthor(request *EditRequest, value string) {
	if request.Operation == nil {
		return
	}

	switch v := request.Operation.Parameters.(type) {
	case scenario.SnoozeParametersRequest:
		v.Author = value
		request.Operation.Parameters = v
	case scenario.ChangeStateParametersRequest:
		v.Author = value
		request.Operation.Parameters = v
	case scenario.AssocTicketParametersRequest:
		v.Author = value
		request.Operation.Parameters = v
	case scenario.PbehaviorParametersRequest:
		v.Author = value
		request.Operation.Parameters = v
	case scenario.ParametersRequest:
		v.Author = value
		request.Operation.Parameters = v
	}
}
