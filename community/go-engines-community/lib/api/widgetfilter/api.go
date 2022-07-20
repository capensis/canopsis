package widgetfilter

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type api struct {
	store        Store
	enforcer     security.Enforcer
	transformer  PatternFieldsTransformer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	transformer PatternFieldsTransformer,
	actionLogger logger.ActionLogger,
) common.CrudAPI {
	return &api{
		store:        store,
		enforcer:     enforcer,
		transformer:  transformer,
		actionLogger: actionLogger,
	}
}

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	ok, err := a.checkAccessByWidget(c.Request.Context(), r.Widget, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	users, err := a.store.Find(c.Request.Context(), r, userId)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r.Query, users)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")
	ok, err := a.checkAccess(c.Request.Context(), id, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	filter, err := a.store.GetOneBy(c.Request.Context(), id, userId)
	if err != nil {
		panic(err)
	}
	if filter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, filter)
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

	err := a.transformEditRequest(c.Request.Context(), &request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	ok, err := a.checkAccessByWidget(c.Request.Context(), request.Widget, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	filter, err := a.store.Insert(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionCreate,
		ValueType: logger.ValueTypeWidgetFilter,
		ValueID:   filter.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusCreated, filter)
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

	err := a.transformEditRequest(c.Request.Context(), &request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	ok, err := a.checkAccess(c.Request.Context(), request.ID, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	filter, err := a.store.GetOneBy(c.Request.Context(), request.ID, request.Author)
	if err != nil {
		panic(err)
	}
	if filter == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if filter.Widget != request.Widget {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"widget": "Widget cannot be changed"}})
		return
	}

	if *filter.IsPrivate != *request.IsPrivate {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"is_private": "IsPrivate cannot be changed"}})
		return
	}

	filter, err = a.store.Update(c.Request.Context(), request)
	if err != nil {
		panic(err)
	}

	if filter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionUpdate,
		ValueType: logger.ValueTypeWidgetFilter,
		ValueID:   filter.ID,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, filter)
}

func (a *api) Delete(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	id := c.Param("id")

	ok, err := a.checkAccess(c.Request.Context(), id, userId, model.PermissionUpdate)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err = a.store.Delete(c.Request.Context(), id, userId)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	err = a.actionLogger.Action(c, userId, logger.LogEntry{
		Action:    logger.ActionDelete,
		ValueType: logger.ValueTypeWidgetFilter,
		ValueID:   id,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.Status(http.StatusNoContent)
}

func (a *api) checkAccess(ctx context.Context, id string, userId, perm string) (bool, error) {
	viewId, err := a.store.FindViewId(ctx, id)
	if err != nil || viewId == "" {
		return false, err
	}

	return a.enforcer.Enforce(userId, viewId, perm)
}

func (a *api) checkAccessByWidget(ctx context.Context, id string, userId, perm string) (bool, error) {
	viewId, err := a.store.FindViewIdByWidget(ctx, id)
	if err != nil || viewId == "" {
		return false, err
	}

	return a.enforcer.Enforce(userId, viewId, perm)
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.AlarmPatternFieldsRequest, err = a.transformer.TransformAlarmPatternFieldsRequest(
		ctx,
		request.AlarmPatternFieldsRequest,
		*request.IsPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(
		ctx,
		request.EntityPatternFieldsRequest,
		*request.IsPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}
	request.PbehaviorPatternFieldsRequest, err = a.transformer.TransformPbehaviorPatternFieldsRequest(
		ctx,
		request.PbehaviorPatternFieldsRequest,
		*request.IsPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}

	return nil
}
