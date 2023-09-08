package colortheme

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
)

type API interface {
	common.CrudAPI
	BulkDelete(c *gin.Context)
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
	logger zerolog.Logger,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
		logger:       logger,
	}
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
	logger       zerolog.Logger
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Theme
func (a api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	theme, err := a.store.Insert(c, request)
	if err != nil {
		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeColorTheme,
		ValueID:   theme.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, theme)
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Theme}
func (a api) List(c *gin.Context) {
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
// @Success 200 {object} Theme
func (a api) Get(c *gin.Context) {
	theme, err := a.store.GetById(c, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if theme == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, theme)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Theme
func (a api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	theme, err := a.store.Update(c, request)
	if err != nil {
		if errors.Is(err, ErrDefaultTheme) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
			return
		}

		validationErr := common.ValidationError{}
		if errors.As(err, &validationErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationErr.ValidationErrorResponse())
			return
		}

		panic(err)
	}

	if theme == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeColorTheme,
		ValueID:   theme.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, theme)
}

func (a api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
	if err != nil {
		if errors.Is(err, ErrDefaultTheme) {
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
		ValueType: logger.ValueTypeColorTheme,
		ValueID:   c.Param("id"),
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusNoContent, nil)
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
			ValueType: logger.ValueTypeColorTheme,
			ValueID:   request.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}

		return request.ID, nil
	}, a.logger)
}
