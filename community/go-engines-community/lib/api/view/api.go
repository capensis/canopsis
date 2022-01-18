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

// Find all views
// @Summary Find views
// @Description Get paginated list of views
// @Tags views
// @ID views-find-all
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Param search query string false "search query"
// @Param sort query string false "sort query"
// @Param sort_by query string false "sort query"
// @Success 200 {object} common.PaginatedListResponse{data=[]viewgroup.View}
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /views [get]
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

// Get view by id
// @Summary Get view by id
// @Description Get view by id
// @Tags views
// @ID views-get-by-id
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view id"
// @Success 200 {object} viewgroup.View
// @Failure 404 {object} common.ErrorResponse
// @Router /views/{id} [get]
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

// Create view
// @Summary Create view
// @Description Create view
// @Tags views
// @ID views-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} viewgroup.View
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /views [post]
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

// Update view by id
// @Summary Update view by id
// @Description Update view by id
// @Tags views
// @ID views-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view id"
// @Param body body EditRequest true "body"
// @Success 200 {object} viewgroup.View
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /views/{id} [put]
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

// Delete view by id
// @Summary Delete view by id
// @Description Delete view by id
// @Tags views
// @ID views-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "view id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /views/{id} [delete]
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

// Update views positions
// @Summary Update views positions
// @Description Update views positions
// @Tags views
// @ID views-update-positions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body []EditPositionItemRequest true "body"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /view-positions [put]
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

// Bulk create views
// @Summary Bulk create views
// @Description Bulk create views
// @Tags views
// @ID views-bulk-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body []EditRequest true "body"
// @Success 201 {array} viewgroup.View
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /bulk/views [post]
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

// Bulk update views by id
// @Summary Bulk update views by id
// @Description Bulk update views by id
// @Tags views
// @ID views-bulk-update-by-id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body []BulkUpdateRequestItem true "body"
// @Success 200 {array} viewgroup.View
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /bulk/views [put]
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

// Bulk delete views by id
// @Summary Bulk delete views by id
// @Description Bulk delete views by id
// @Tags views
// @ID views-bulk-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param request query BulkDeleteRequest true "request"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /bulk/views [delete]
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
