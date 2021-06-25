package entityservice

import (
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/entityservice"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
)

type API interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetDependencies(c *gin.Context)
	GetImpacts(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type api struct {
	store                 Store
	serviceChangeListener chan<- entityservice.ChangeEntityMessage
	actionLogger          logger.ActionLogger
	logger                zerolog.Logger
}

func NewApi(
	store Store,
	serviceChangeListener chan<- entityservice.ChangeEntityMessage,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:                 store,
		serviceChangeListener: serviceChangeListener,
		actionLogger:          actionLogger,
		logger:                logger,
	}
}

// Get entity service by id
// @Summary Get entity service by id
// @Description Get entity service by id
// @Tags entityservices
// @ID entityservices-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "entity service id"
// @Success 200 {object} Response
// @Failure 404 {object} common.ErrorResponse
// @Router /entityservices/{id} [get]
func (a *api) Get(c *gin.Context) {
	service, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if service == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, service)
}

// Get entity service's dependencies by id
// @Summary Get entity service's dependencies by id
// @Description Get entity service's dependencies by id
// @Tags entityservices
// @ID entityservices-get-dependencies-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id query string true "entity service id"
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Success 200 {object} common.PaginatedListResponse{data=[]AlarmWithEntity}
// @Failure 404 {object} common.ErrorResponse
// @Router /entityservice-dependencies [get]
func (a *api) GetDependencies(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.GetDependencies(c.Request.Context(), r.ID, r.Query)
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

// Get entity's impacted services by id
// @Summary Get entity's impacted services by id
// @Description Get entity's impacted services by id
// @Tags entityservices
// @ID entityservices-get-impacts-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id query string true "entity id"
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Success 200 {object} common.PaginatedListResponse{data=[]AlarmWithEntity}
// @Failure 404 {object} common.ErrorResponse
// @Router /entityservice-impacts [get]
func (a *api) GetImpacts(c *gin.Context) {
	var r ContextGraphRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.GetImpacts(c.Request.Context(), r.ID, r.Query)
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

// Create entity service
// @Summary Create entity service
// @Description Create entity service
// @Tags entityservices
// @ID entityservices-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /entityservices [post]
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	service, err := a.store.Create(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if service.Enabled {
		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      service.ID,
			IsService:               true,
			IsServicePatternChanged: true,
		})
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   service.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, service)
}

// Update entity service by id
// @Summary Update entity service by id
// @Description Update entity service by id
// @Tags entityservices
// @ID entityservices-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "entity service id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /entityservices/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	service, serviceChanges, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if service == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if service.Enabled || serviceChanges.IsToggled {
		a.sendChangeMsg(entityservice.ChangeEntityMessage{
			ID:                      service.ID,
			IsService:               true,
			IsServicePatternChanged: serviceChanges.IsPatternChanged,
			IsToggled:               serviceChanges.IsToggled,
		})
	}

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   service.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, service)
}

// Delete entity service by id
// @Summary Delete entity service by id
// @Description Delete entity service by id
// @Tags entityservices
// @ID entityservices-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "entity service id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /entityservices/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, alarm, err := a.store.Delete(c.Request.Context(), id)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendChangeMsg(entityservice.ChangeEntityMessage{
		ID:                      id,
		IsService:               true,
		IsServicePatternChanged: true,
		ServiceAlarm:            alarm,
	})

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEntityService,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

func (a *api) sendChangeMsg(msg entityservice.ChangeEntityMessage) {
	select {
	case a.serviceChangeListener <- msg:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("service_id", msg.ID).
			Msg("fail to send change message")
	}
}
