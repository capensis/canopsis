package pbehaviorexception

import (
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
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
	logger zerolog.Logger,
) API {
	return &api{
		store:       store,
		computeChan: computeChan,
		maxFileSize: maxFileSize,
		logger:      logger,
	}
}

type api struct {
	store       Store
	computeChan chan<- rpc.PbehaviorRecomputeEvent
	maxFileSize uint64
	logger      zerolog.Logger
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.Find(c, r)
	if err != nil {
		panic(err)
	}

	res := pagination.NewResponse(r.Query, aggregationResult)
	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
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

		panic(err)
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

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
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

		panic(err)
	}

	if res == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isLinked, err := a.store.IsLinked(c, res.ID)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	if exception == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(authctx.UserKey).(string))

	if err != nil {
		if errors.Is(err, ErrLinkedException) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		panic(err)
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

	exception, err := a.store.Import(c, name, pbhType, c.MustGet(authctx.UserKey).(string), f, fh)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	c.JSON(http.StatusOK, exception)
}

func (a *api) sendComputeTask() {
	a.computeChan <- rpc.PbehaviorRecomputeEvent{}
}
