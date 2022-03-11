package pbehaviorreason

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func NewApi(
	transformer ModelTransformer,
	store Store,
	computeChan chan<- pbehavior.ComputeTask,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) common.CrudAPI {
	return &api{
		transformer:  transformer,
		store:        store,
		computeChan:  computeChan,
		actionLogger: actionLogger,
		logger:       logger,
	}
}

type api struct {
	transformer  ModelTransformer
	store        Store
	computeChan  chan<- pbehavior.ComputeTask
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

// Find all pbehavior reasons
// @Summary Find pbehavior reasons
// @Description Get paginated list of behavior reasons
// @Tags pbehavior-reasons
// @ID pbehavior-reasons-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Param with_flags query bool false "with flags"
// @Success 200 {object} common.PaginatedListResponse{data=[]Reason}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-reasons [get]
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create pbehavior reason
// @Summary Create pbehavior reason
// @Description Create pbehavior reason
// @Tags pbehavior-reasons
// @ID pbehavior-reasons-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body Request true "body"
// @Success 201 {object} Reason
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-reasons [post]
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	reason := a.transformer.TransformCreateRequestToModel(request)
	err := a.store.Insert(c.Request.Context(), reason)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehaviorReason,
		ValueID:   reason.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, reason)
}

// Get pbehavior reason by id
// @Summary Get pbehavior reason by id
// @Description Get pbehavior reason by id
// @Tags pbehavior-reasons
// @ID pbehavior-reasons-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "reason id"
// @Success 200 {object} Reason
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-reasons/{id} [get]
func (a *api) Get(c *gin.Context) {
	reason, err := a.store.GetOneBy(c.Request.Context(), bson.M{"_id": c.Param("id")})
	if err != nil {
		panic(err)
	}

	if reason == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, reason)
}

// Update behavior reason by id
// @Summary Update behavior reason by id
// @Description Update behavior reason by id
// @Tags pbehavior-reasons
// @ID pbehavior-reasons-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "reason id"
// @Param body body Request true "body"
// @Success 200 {object} Reason
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-reasons/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	reason := a.transformer.TransformUpdateRequestToModel(request)
	ok, err := a.store.Update(c.Request.Context(), reason)

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isLinked, err := a.store.IsLinkedToPbehavior(c.Request.Context(), reason.ID)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehaviorReason,
		ValueID:   reason.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	if isLinked {
		a.sendComputeTask(reason.ID)
	}
	c.JSON(http.StatusOK, reason)
}

// Delete pbehavior reason by id
// @Summary Delete pbehavior reason by id
// @Description Delete pbehavior reason by id
// @Tags pbehavior-reasons
// @ID pbehavior-reasons-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "reason id"
// @Success 204
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-reasons/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehaviorReason,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}

func (a *api) sendComputeTask(reasonID string) {
	task := pbehavior.ComputeTask{}

	select {
	case a.computeChan <- task:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("reason", reasonID).
			Msg("fail to start linked pbehaviors recompute on reason update")
	}
}
