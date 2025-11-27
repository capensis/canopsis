package pbehaviorreason

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/rpc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewApi(
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) crud.API {
	return &api{
		store:          store,
		computeChan:    computeChan,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

type api struct {
	store          Store
	computeChan    chan<- rpc.PbehaviorRecomputeEvent
	errorResponder httperror.Responder
	logger         zerolog.Logger
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Insert(c, request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Get
// @Param id path string true "reason id"
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	reason, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if reason == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, reason)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Update(c, request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isLinked, err := a.store.IsLinkedToPbehavior(c, res.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if isLinked {
		a.sendComputeTask()
	}
	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(authctx.UserKey).(string))
	if err != nil {
		var valErr ValidationError
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}
