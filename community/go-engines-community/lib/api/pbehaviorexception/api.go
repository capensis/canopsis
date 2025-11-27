package pbehaviorexception

import (
	"errors"
	"fmt"
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

type API interface {
	crud.API
	Import(c *gin.Context)
}

func NewApi(
	store Store,
	computeChan chan<- rpc.PbehaviorRecomputeEvent,
	maxFileSize uint64,
	errorResponder httperror.Responder,
	logger zerolog.Logger,
) API {
	return &api{
		store:          store,
		computeChan:    computeChan,
		maxFileSize:    maxFileSize,
		errorResponder: errorResponder,
		logger:         logger,
	}
}

type api struct {
	store          Store
	computeChan    chan<- rpc.PbehaviorRecomputeEvent
	maxFileSize    uint64
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
	var request CreateRequest
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

		if errors.Is(err, ErrTypeNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
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

		if errors.Is(err, ErrTypeNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		a.errorResponder.Respond(c, err)

		return
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isLinked, err := a.store.IsLinked(c, res.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if isLinked {
		a.sendComputeTask()
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	exception, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if exception == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, c.Param("id"), userID)

	if err != nil {
		if errors.Is(err, ErrLinkedException) {
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

// Import
// @Success 200 {object} Response
func (a *api) Import(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationError("file", "File is missing.").ValidationErrorResponse())
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: "request has invalid structure"})
		return
	}
	defer f.Close()

	name := c.Request.FormValue("name")
	pbhType := c.Request.FormValue("type")
	valErrors := make(map[string]string)
	if a.maxFileSize > 0 && uint64(fh.Size) > a.maxFileSize {
		valErrors["file"] = fmt.Sprintf("File size %d exceeds limit %d", fh.Size, a.maxFileSize)
	}

	if name == "" {
		valErrors["name"] = "Name is missing."
	}
	if pbhType == "" {
		valErrors["type"] = "Type is missing."
	}

	if len(valErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrors(valErrors).ValidationErrorResponse())

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	exception, err := a.store.Import(c, name, pbhType, userID, f, fh)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}
