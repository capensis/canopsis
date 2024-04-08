package pbehaviortype

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	common.CrudAPI
	GetNextPriority(c *gin.Context)
}

type api struct {
	store        Store
	transformer  ModelTransformer
	computeChan  chan<- rpc.PbehaviorRecomputeEvent
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

func NewApi(
	transformer ModelTransformer,
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		transformer:  transformer,
		store:        store,
		computeChan:  computeChan,
		actionLogger: actionLogger,
		logger:       logger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]pbehavior.Type}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	types, err := a.store.Find(c, r)
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

// Get
// @Success 200 {object} pbehavior.Type
func (a *api) Get(c *gin.Context) {
	pt, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}
	if pt == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, pt)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} pbehavior.Type
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pt := a.transformer.TransformCreateRequestToModel(request)
	if err := a.store.Insert(c, pt); err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		var fieldValErr common.ValidationError
		if errors.As(err, &fieldValErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, fieldValErr.ValidationErrorResponse())
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

	a.sendComputeTask()
	c.JSON(http.StatusCreated, pt)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} pbehavior.Type
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	pt := a.transformer.TransformUpdateRequestToModel(request)
	ok, err := a.store.Update(c, c.Param("id"), pt)
	if err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		var fieldValErr common.ValidationError
		if errors.As(err, &fieldValErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, fieldValErr.ValidationErrorResponse())
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

	a.sendComputeTask()
	c.JSON(http.StatusOK, pt)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
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

// GetNextPriority
// @Success 200 {object} PriorityResponse
func (a *api) GetNextPriority(c *gin.Context) {
	priority, err := a.store.GetNextPriority(c)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, PriorityResponse{Priority: priority})
}

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}
