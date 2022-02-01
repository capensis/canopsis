package pbehaviortype

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type api struct {
	store        Store
	transformer  ModelTransformer
	computeChan  chan<- pbehavior.ComputeTask
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

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

// Find all pbehavior types
// @Summary Find pbehavior types
// @Description Get paginated list of behavior types
// @Tags pbehavior-types
// @ID pbehavior-types-find-all
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
// @Success 200 {object} common.PaginatedListResponse{data=[]pbehavior.Type}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-types [get]
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	types, err := a.store.Find(c.Request.Context(), r)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, types)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get pbehavior type by id
// @Summary Get pbehavior type by id
// @Description Get pbehavior type by id
// @Tags pbehavior-types
// @ID pbehavior-types-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "type id"
// @Success 200 {object} pbehavior.Type
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-types/{id} [get]
func (a *api) Get(c *gin.Context) {
	pt, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if pt == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pt)
}

// Create pbehavior type
// @Summary Create pbehavior type
// @Description Create pbehavior type
// @Tags pbehavior-types
// @ID pbehavior-types-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} pbehavior.Type
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-types [post]
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pt := a.transformer.TransformCreateRequestToModel(request)
	if err := a.store.Insert(c.Request.Context(), pt); err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		panic(err)
	}

	err := a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehaviorType,
		ValueID:   pt.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.computeChan <- pbehavior.ComputeTask{}
	c.JSON(http.StatusCreated, pt)
}

// Update behavior type by id
// @Summary Update behavior type by id
// @Description Update behavior type by id
// @Tags pbehavior-types
// @ID pbehavior-types-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "type id"
// @Param body body EditRequest true "body"
// @Success 200 {object} pbehavior.Type
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-types/{id} [put]
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pt := a.transformer.TransformUpdateRequestToModel(request)
	ok, err := a.store.Update(c.Request.Context(), c.Param("id"), pt)
	if err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		panic(err)
	}

	if !ok {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehaviorType,
		ValueID:   pt.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	a.sendComputeTask(pt.ID)
	c.JSON(http.StatusOK, pt)
}

// Delete pbehavior type by id
// @Summary Delete pbehavior type by id
// @Description Delete pbehavior type by id
// @Tags pbehavior-types
// @ID pbehavior-types-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "type id"
// @Success 204
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-types/{id} [delete]
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
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehaviorType,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

func (a *api) sendComputeTask(typeID string) {
	task := pbehavior.ComputeTask{}

	select {
	case a.computeChan <- task:
	default:
		a.logger.Err(errors.New("channel is full")).
			Str("type", typeID).
			Msg("fail to start linked pbehaviors recompute on type update")
	}
}
