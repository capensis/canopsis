package widget

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Copy(c *gin.Context)
	UpdateGridPositions(c *gin.Context)
}

type api struct {
	store       Store
	enforcer    security.Enforcer
	transformer *RequestTransformer

	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	transformer *RequestTransformer,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		enforcer:     enforcer,
		transformer:  transformer,
		actionLogger: actionLogger,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	widget, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		panic(err)
	}
	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.checkAccessByTab(c, widget.Tab, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	c.JSON(http.StatusOK, widget)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{}
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.checkAccessByTab(c, request.Tab, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	err = a.transformer.Transform(c, &request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	widget, err := a.store.Insert(c, request)
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

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	ok, err := a.checkAccess(c, []string{request.ID}, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.checkAccessByTab(c, request.Tab, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	err = a.transformer.Transform(c, &request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	widget, err := a.store.Update(c, request)
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

func (a *api) Delete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")

	ok, err := a.checkAccess(c, []string{id}, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.Delete(c, id)
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

// Copy
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Copy(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")
	request := EditRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	widget, err := a.store.GetOneBy(c, id)
	if err != nil {
		panic(err)
	}
	if widget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	ok, err := a.checkAccessByTab(c, widget.Tab, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.checkAccessByTab(c, request.Tab, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	newWidget, err := a.store.Copy(c, *widget, request)
	if err != nil {
		panic(err)
	}

	if newWidget == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeWidget,
		ValueID:   newWidget.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, newWidget)
}

// UpdateGridPositions
// @Param body body []EditGridPositionItemRequest true "body"
func (a *api) UpdateGridPositions(c *gin.Context) {
	request := EditGridPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	userId := c.MustGet(auth.UserKey).(string)
	ids := make([]string, len(request.Items))
	for i, item := range request.Items {
		ids[i] = item.ID
	}
	ok, err := a.checkAccess(c, ids, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.UpdateGridPositions(c, request.Items)
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
	viewId, err := a.store.FindViewIdByTab(ctx, tabId)
	if err != nil || viewId == "" {
		return false, err
	}

	return a.enforcer.Enforce(userId, viewId, perm)
}
