package widgetfilter

import (
	"context"
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	common.CrudAPI
	UpdatePositions(c *gin.Context)
}

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
) API {
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

	ok, _, err := a.checkAccessByWidget(c, r.Widget, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	users, err := a.store.Find(c, r, userId)
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
	ok, _, err := a.checkAccess(c, id, userId, model.PermissionRead)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	filter, err := a.store.GetOneBy(c, id, userId)
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
	request := CreateRequest{}

	if err := c.ShouldBind(&request); err != nil {
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

	if request.WidgetPrivate != nil {
		if !*request.WidgetPrivate {
			ok, err := a.enforcer.Enforce(userId, apisecurity.ObjView, model.PermissionUpdate)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}

		var granted bool
		if *request.WidgetPrivate {
			granted, request.IsPrivate, err = a.checkAccessByWidget(c, request.Widget, userId, model.PermissionRead)
		} else {
			granted, request.IsPrivate, err = a.checkAccessByWidget(c, request.Widget, userId, model.PermissionUpdate)
		}
		if err != nil {
			panic(err)
		}

		if !granted {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	filter, err := a.store.Insert(c, request)
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
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := c.ShouldBind(&request); err != nil {
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

	if request.WidgetPrivate != nil {
		if !*request.WidgetPrivate {
			ok, err := a.enforcer.Enforce(userId, apisecurity.ObjView, model.PermissionUpdate)
			if err != nil {
				panic(err)
			}

			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
				return
			}
		}

		var granted bool
		if *request.WidgetPrivate {
			granted, request.IsPrivate, err = a.checkAccess(c, request.ID, userId, model.PermissionRead)
		} else {
			granted, request.IsPrivate, err = a.checkAccess(c, request.ID, userId, model.PermissionUpdate)
		}
		if err != nil {
			panic(err)
		}

		if !granted {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	filter, err := a.store.GetOneBy(c, request.ID, request.Author)
	if err != nil {
		panic(err)
	}
	if filter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if filter.WidgetPrivate != *request.WidgetPrivate {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"is_user_preference": "WidgetPrivate cannot be changed"}})
		return
	}

	filter, err = a.store.Update(c, request)
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

	filter, err := a.store.GetOneBy(c, id, userId)
	if err != nil {
		panic(err)
	}
	if filter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if !filter.WidgetPrivate {
		ok, err := a.enforcer.Enforce(userId, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	var granted bool
	if filter.WidgetPrivate {
		granted, _, err = a.checkAccess(c, id, userId, model.PermissionRead)
	} else {
		granted, _, err = a.checkAccess(c, id, userId, model.PermissionUpdate)
	}
	if err != nil {
		panic(err)
	}

	if !granted {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err := a.store.Delete(c, id, userId)
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

func (a *api) UpdatePositions(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditPositionRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	if len(request.Items) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"items": "Cannot be blank"}})
		return
	}

	firstItem := request.Items[0]
	firstFilter, err := a.store.GetOneBy(c, firstItem, userId)
	if err != nil {
		panic(err)
	}
	if firstFilter == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	isPrivate := firstFilter.WidgetPrivate
	if !isPrivate {
		ok, err := a.enforcer.Enforce(userId, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			panic(err)
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
			return
		}
	}

	perm := model.PermissionUpdate
	if isPrivate {
		perm = model.PermissionRead
	}

	granted, _, err := a.checkAccess(c, firstItem, userId, perm)
	if err != nil {
		panic(err)
	}

	if !granted {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	ok, err := a.store.UpdatePositions(c, request.Items, firstFilter.Widget, userId, isPrivate)
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

func (a *api) checkAccess(ctx context.Context, id string, userId, perm string) (bool, bool, error) {
	viewId, author, isPrivate, err := a.store.FindViewId(ctx, id)
	if err != nil || viewId == "" {
		return false, false, err
	}

	if isPrivate {
		return author == userId, true, nil
	}

	granted, err := a.enforcer.Enforce(userId, viewId, perm)

	return granted, isPrivate, err
}

func (a *api) checkAccessByWidget(ctx context.Context, id, userId, perm string) (bool, bool, error) {
	viewId, author, isPrivate, err := a.store.FindViewIdByWidget(ctx, id)
	if err != nil || viewId == "" {
		return false, false, err
	}

	if isPrivate {
		return author == userId, true, nil
	}

	granted, err := a.enforcer.Enforce(userId, viewId, perm)

	return granted, isPrivate, err
}

func (a *api) transformEditRequest(ctx context.Context, request *EditRequest) error {
	var err error
	request.AlarmPatternFieldsRequest, err = a.transformer.TransformAlarmPatternFieldsRequest(
		ctx,
		request.AlarmPatternFieldsRequest,
		*request.WidgetPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}
	request.EntityPatternFieldsRequest, err = a.transformer.TransformEntityPatternFieldsRequest(
		ctx,
		request.EntityPatternFieldsRequest,
		*request.WidgetPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}
	request.PbehaviorPatternFieldsRequest, err = a.transformer.TransformPbehaviorPatternFieldsRequest(
		ctx,
		request.PbehaviorPatternFieldsRequest,
		*request.WidgetPrivate,
		request.Author,
	)
	if err != nil {
		return err
	}

	return nil
}
