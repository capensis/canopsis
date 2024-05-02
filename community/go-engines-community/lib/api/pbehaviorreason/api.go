package pbehaviorreason

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
	"go.mongodb.org/mongo-driver/bson"
)

func NewApi(
	transformer ModelTransformer,
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
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
	computeChan  chan<- rpc.PbehaviorRecomputeEvent
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Reason}
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

// Create
// @Param body body Request true "body"
// @Success 201 {object} Reason
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

// Get
// @Param id path string true "reason id"
// @Success 200 {object} Reason
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

// Update
// @Param body body Request true "body"
// @Success 200 {object} Reason
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
		a.sendComputeTask()
	}
	c.JSON(http.StatusOK, reason)
}

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

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}
