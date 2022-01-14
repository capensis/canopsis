package widget

import (
	"context"
	"errors"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	UpdatePositions(c *gin.Context)
}

type api struct {
	store        Store
	tabStore     viewtab.Store
	enforcer     security.Enforcer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	tabStore viewtab.Store,
	enforcer security.Enforcer,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		tabStore:     tabStore,
		enforcer:     enforcer,
		actionLogger: actionLogger,
	}
}

// Get widget by id
// @Summary Get widget by id
// @Description Get widget by id
// @Tags widgets
// @ID widgets-get-by-id
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "widget id"
// @Success 200 {object} view.Widget
// @Failure 404 {object} common.ErrorResponse
// @Router /widgets/{id} [get]
func (a *api) Get(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	widget, err := a.store.GetOneBy(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}
	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.checkAccessByTab(c.Request.Context(), widget.Tab, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	c.JSON(http.StatusOK, widget)
}

// Create widget
// @Summary Create widget
// @Description Create widget
// @Tags widgets
// @ID widgets-create
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 201 {object} view.Widget
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /widgets [post]
func (a *api) Create(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.checkAccessByTab(c.Request.Context(), request.Tab, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	widget, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeWidget,
		ValueID:   widget.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, widget)
}

// Update widget by id
// @Summary Update widget by id
// @Description Update widget by id
// @Tags widgets
// @ID widgets-update-by-id
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "widget id"
// @Param body body EditRequest true "body"
// @Success 200 {object} view.Widget
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /widgets/{id} [put]
func (a *api) Update(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.checkAccess(c.Request.Context(), []string{request.ID}, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.checkAccessByTab(c.Request.Context(), request.Tab, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	widget, err := a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeWidget,
		ValueID:   widget.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, widget)
}

// Delete widget by id
// @Summary Delete widget by id
// @Description Delete widget by id
// @Tags widgets
// @ID widgets-delete-by-id
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "widget id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /widgets/{id} [delete]
func (a *api) Delete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")

	ok, err := a.checkAccess(c.Request.Context(), []string{id}, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.Delete(c.Request.Context(), id)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeWidget,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

// Update widgets positions
// @Summary Update widgets positions
// @Description Update widgets positions
// @Tags widgets
// @ID widgets-update-positions
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body []string true "body"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /widget-positions [put]
func (a *api) UpdatePositions(c *gin.Context) {
	request := EditPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)
	ok, err := a.checkAccess(c.Request.Context(), request.Items, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.UpdatePositions(c.Request.Context(), request.Items)
	if err != nil {
		valErr := ValidationErr{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) checkAccess(ctx context.Context, ids []string, userId, perm string) (bool, error) {
	viewIds, err := a.store.FindViewIds(ctx, ids)
	if err != nil || len(viewIds) != len(ids) {
		return false, err
	}

	for _, viewId := range viewIds {
		ok, err := a.enforcer.Enforce(userId, viewId, perm)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}

func (a *api) checkAccessByTab(ctx context.Context, tabId string, userId, perm string) (bool, error) {
	tab, err := a.tabStore.GetOneBy(ctx, tabId)
	if err != nil || tab == nil {
		return false, err
	}

	return a.enforcer.Enforce(userId, tab.View, perm)
}
