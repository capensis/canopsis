package eventfilter

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	common.BulkCrudAPI
	ListFailures(c *gin.Context)
	ReadFailures(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
	transformer  common.PatternFieldsTransformer
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
	transformer common.PatternFieldsTransformer,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		logger:       logger,
		transformer:  transformer,
	}
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	var err error

	if err = c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err = a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	eventfilter, err := a.store.Insert(c, request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   eventfilter.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, eventfilter)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var query FilteredQuery
	query.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, query))
		return
	}

	aggregationResult, err := a.store.Find(c, query)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(query.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	evf, err := a.store.GetById(c, c.Param("id"))

	if errors.Is(err, mongodriver.ErrNoDocuments) || evf == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, evf)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	err := a.transformEditRequest(c, &request.EditRequest)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	eventfilter, err := a.store.Update(c, request)
	if err != nil {
		panic(err)
	}

	if eventfilter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   eventfilter.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, eventfilter)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeEventFilter,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request CreateRequest) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		eventfilter, err := a.store.Insert(c, request)
		if err != nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   eventfilter.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return eventfilter.ID, nil
	}, a.logger)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		err := a.transformEditRequest(c, &request.EditRequest)
		if err != nil {
			return "", err
		}

		eventfilter, err := a.store.Update(c, UpdateRequest(request))
		if err != nil || eventfilter == nil {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   eventfilter.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return eventfilter.ID, nil
	}, a.logger)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID)
		if err != nil || !ok {
			return "", err
		}

		err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeEventFilter,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return request.ID, nil
	}, a.logger)
}

// ListFailures
// @Success 200 {object} common.PaginatedListResponse{data=[]FailureResponse}
func (a *api) ListFailures(c *gin.Context) {
	r := FailureRequest{}
	r.Query = pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	aggregationResult, err := a.store.FindFailures(c, c.Param("id"), r)
	if err != nil {
		panic(err)
	}

	if aggregationResult == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	res, err := common.NewPaginatedResponse(r.Query, aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) ReadFailures(c *gin.Context) {
	exists, err := a.store.ReadFailures(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !exists {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error

	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(ctx, request.EntityPatternFieldsRequest)
	return err
}
