package pbehaviorexception

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type API interface {
	common.CrudAPI
	Import(c *gin.Context)
}

func NewApi(
	transformer ModelTransformer,
	store Store,
	computeChan chan<- []string,
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

type api struct {
	transformer  ModelTransformer
	store        Store
	computeChan  chan<- []string
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Exception}
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

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create
// @Param body body Request true "body"
// @Success 201 {object} Exception
func (a *api) Create(c *gin.Context) {
	var request CreateRequest

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	exception, err := a.transformer.TransformCreateRequestToModel(c, request)
	if err != nil {
		if errors.Is(err, ErrTypeNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		panic(err)
	}

	err = a.store.Insert(c, exception)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypePbehaviorException,
		ValueID:   exception.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, exception)
}

// Update
// @Param body body Request true "body"
// @Success 200 {object} Exception
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	exception, err := a.transformer.TransformUpdateRequestToModel(c, request)
	if err != nil {
		if errors.Is(err, ErrTypeNotExists) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}
		panic(err)
	}

	ok, err := a.store.Update(c, exception)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isLinked, err := a.store.IsLinked(c, exception.ID)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypePbehaviorException,
		ValueID:   exception.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	if isLinked {
		a.sendComputeTask()
	}

	c.JSON(http.StatusOK, exception)
}

// Get
// @Success 200 {object} Exception
func (a *api) Get(c *gin.Context) {
	exception, err := a.store.GetOneById(c, c.Param("id"))
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
	ok, err := a.store.Delete(c, c.Param("id"))

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

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypePbehaviorException,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}

// Import
// @Success 200 {object} Exception
func (a *api) Import(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{
				"file": "File is missing.",
			}})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: "request has invalid structure"})
		return
	}
	defer f.Close()

	name := c.Request.FormValue("name")
	pbhType := c.Request.FormValue("type")
	valErrors := make(map[string]string)
	if name == "" {
		valErrors["name"] = "Name is missing."
	}
	if pbhType == "" {
		valErrors["type"] = "Type is missing."
	}

	if len(valErrors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: valErrors})
		return
	}

	exception, err := a.store.Import(c, name, pbhType, f, fh)
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
	a.computeChan <- []string{}
}
