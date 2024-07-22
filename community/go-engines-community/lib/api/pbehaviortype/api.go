package pbehaviortype

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
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
	store       Store
	computeChan chan<- rpc.PbehaviorRecomputeEvent
	logger      zerolog.Logger
}

func NewApi(
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	logger zerolog.Logger,
) API {
	return &api{
		store:       store,
		computeChan: computeChan,
		logger:      logger,
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
	pt, err := a.store.GetByID(c, c.Param("id"))
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

	res, err := a.store.Insert(c, request)
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

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendComputeTask()
	c.JSON(http.StatusCreated, res)
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

	res, err := a.store.Update(c, request)
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

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.sendComputeTask()
	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(auth.UserKey).(string))
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
