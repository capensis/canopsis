package view

import (
	"context"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	common.BulkCrudAPI
	UpdatePositions(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]viewgroup.View}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	authorizedIds, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		r.Ids = authorizedIds.([]string)
	}

	views := &AggregationResult{}
	var err error

	if len(r.Ids) > 0 {
		views, err = a.store.Find(c.Request.Context(), r)
	}

	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, views)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} viewgroup.View
func (a *api) Get(c *gin.Context) {
	view, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if view == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, view)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} viewgroup.View
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)

	views, err := a.store.Insert(c.Request.Context(), userID, []EditRequest{request})
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(context.Background(), userID, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeView,
		ValueID:   views[0].ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, views[0])
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} viewgroup.View
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	views, err := a.store.Update(c.Request.Context(), []BulkUpdateRequestItem{{
		ID:              request.ID,
		BaseEditRequest: request.BaseEditRequest,
	}})
	if err != nil {
		panic(err)
	}

	if len(views) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeView,
		ValueID:   views[0].ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, views[0])
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	ok, err := a.store.Delete(c.Request.Context(), []string{id})

	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeView,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// UpdatePositions
// @Param body body []EditPositionItemRequest true "body"
func (a *api) UpdatePositions(c *gin.Context) {
	request := EditPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	v, ok := c.Get(middleware.AuthorizedIds)
	var authorizedIds []string
	if ok {
		authorizedIds = v.([]string)
	}

	canUpdate := make(map[string]bool, len(authorizedIds))
	for _, id := range authorizedIds {
		canUpdate[id] = true
	}

	for _, item := range request.Items {
		for _, view := range item.Views {
			if !canUpdate[view] {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}
	}

	ok, err := a.store.UpdatePositions(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) BulkCreate(c *gin.Context) {
	var request BulkCreateRequest
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)

	views, err := a.store.Insert(c.Request.Context(), userID, request.Items)
	if err != nil {
		panic(err)
	}

	for _, view := range views {
		err = a.actionLogger.Action(context.Background(), userID, logger.LogEntry{
			Action:    logger.ActionCreate,
			ValueType: logger.ValueTypeView,
			ValueID:   view.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusCreated, views)
}

func (a *api) BulkUpdate(c *gin.Context) {
	request := BulkUpdateRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	v, ok := c.Get(middleware.AuthorizedIds)
	var authorizedIds []string
	if ok {
		authorizedIds = v.([]string)
	}

	canUpdate := make(map[string]bool, len(authorizedIds))
	for _, id := range authorizedIds {
		canUpdate[id] = true
	}

	for _, item := range request.Items {
		if !canUpdate[item.ID] {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	views, err := a.store.Update(c.Request.Context(), request.Items)
	if err != nil {
		panic(err)
	}

	if len(views) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	for _, view := range views {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionUpdate,
			ValueType: logger.ValueTypeView,
			ValueID:   view.ID,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.JSON(http.StatusOK, views)
}

func (a *api) BulkDelete(c *gin.Context) {
	request := BulkDeleteRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	v, ok := c.Get(middleware.AuthorizedIds)
	var authorizedIds []string
	if ok {
		authorizedIds = v.([]string)
	}

	canUpdate := make(map[string]bool, len(authorizedIds))
	for _, id := range authorizedIds {
		canUpdate[id] = true
	}

	for _, id := range request.IDs {
		if !canUpdate[id] {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	ok, err := a.store.Delete(c.Request.Context(), request.IDs)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	for _, id := range request.IDs {
		err = a.actionLogger.Action(context.Background(), c.MustGet(auth.UserKey).(string), logger.LogEntry{
			Action:    logger.ActionDelete,
			ValueType: logger.ValueTypeView,
			ValueID:   id,
		})
		if err != nil {
			a.actionLogger.Err(err, "failed to log action")
		}
	}

	c.Status(http.StatusNoContent)
}
