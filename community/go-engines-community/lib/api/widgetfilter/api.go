package widgetfilter

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type API interface {
	crud.API
	UpdatePositions(c *gin.Context)
}

type api struct {
	store          Store
	enforcer       security.Enforcer
	errorResponder httperror.Responder
}

func NewApi(
	store Store,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		enforcer:       enforcer,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, _, err := a.checkAccessByWidget(c, r.Widget, userID, model.PermissionRead)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	users, err := a.store.Find(c, r, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, users)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	id := c.Param("id")
	ok, _, err := a.checkAccess(c, id, userID, model.PermissionRead)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	filter, err := a.store.GetOneBy(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if filter == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, filter)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := CreateRequest{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	var granted bool
	perm := model.PermissionUpdate
	if *request.IsUserPreference {
		perm = model.PermissionRead
	}

	granted, request.IsPrivate, err = a.checkAccessByWidget(c, request.Widget, userID, perm)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !granted {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	if !*request.IsUserPreference && !request.IsPrivate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	filter, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, filter)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	var granted bool
	perm := model.PermissionUpdate
	if *request.IsUserPreference {
		perm = model.PermissionRead
	}

	granted, request.IsPrivate, err = a.checkAccess(c, request.ID, userID, perm)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !granted {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	if !*request.IsUserPreference && !request.IsPrivate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	filter, err := a.store.GetOneBy(c, request.ID, request.Author)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if filter == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if filter.IsUserPreference != *request.IsUserPreference {
		err = validation.NewError(
			validator.ValidationErrors{validation.NewFieldError("unchangeable", "IsUserPreference", "IsUserPreference")},
			request,
		)
		a.errorResponder.Respond(c, err)

		return
	}

	filter, err = a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if filter == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, filter)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	id := c.Param("id")

	filter, err := a.store.GetOneBy(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if filter == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	perm := model.PermissionUpdate
	if filter.IsUserPreference {
		perm = model.PermissionRead
	}

	granted, isPrivate, err := a.checkAccess(c, id, userID, perm)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !granted {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	if !filter.IsUserPreference && !isPrivate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	ok, err := a.store.Delete(c, id, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) UpdatePositions(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := EditPositionRequest{}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if len(request.Items) == 0 {
		err = validation.NewError(
			validator.ValidationErrors{validation.NewFieldError("required", "Items", "Items")},
			request,
		)
		a.errorResponder.Respond(c, err)

		return
	}

	firstItem := request.Items[0]
	firstFilter, err := a.store.GetOneBy(c, firstItem, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if firstFilter == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	isUserPreference := firstFilter.IsUserPreference
	perm := model.PermissionUpdate
	if isUserPreference {
		perm = model.PermissionRead
	}

	granted, isPrivate, err := a.checkAccess(c, firstItem, userID, perm)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !granted {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	if !isUserPreference && !isPrivate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.ObjView, model.PermissionUpdate)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	ok, err := a.store.UpdatePositions(c, request.Items, firstFilter.Widget, userID, isUserPreference)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.Status(http.StatusNoContent)
}

func (a *api) checkAccess(ctx context.Context, id string, userID, perm string) (bool, bool, error) {
	viewId, author, isPrivate, err := a.store.FindViewId(ctx, id)
	if err != nil || viewId == "" {
		return false, false, err
	}

	if isPrivate {
		return author == userID, true, nil
	}

	granted, err := a.enforcer.Enforce(userID, viewId, perm)

	return granted, isPrivate, err
}

func (a *api) checkAccessByWidget(ctx context.Context, id, userID, perm string) (bool, bool, error) {
	viewId, author, isPrivate, err := a.store.FindViewIdByWidget(ctx, id)
	if err != nil || viewId == "" {
		return false, false, err
	}

	if isPrivate {
		return author == userID, true, nil
	}

	granted, err := a.enforcer.Enforce(userID, viewId, perm)

	return granted, isPrivate, err
}
