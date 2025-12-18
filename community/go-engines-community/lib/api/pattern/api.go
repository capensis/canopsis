package pattern

import (
	"context"
	"net/http"
	"time"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.API
	BulkDelete(c *gin.Context)
	CountAlarms(c *gin.Context)
	CountEntities(c *gin.Context)
}

type api struct {
	store          Store
	enforcer       security.Enforcer
	errorResponder httperror.Responder

	configProvider config.UserInterfaceConfigProvider
}

func NewApi(
	store Store,
	configProvider config.UserInterfaceConfigProvider,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		enforcer:       enforcer,
		errorResponder: errorResponder,

		configProvider: configProvider,
	}
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if *request.IsCorporate {
		userID, err := authctx.GetUserKey(c)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		ok, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	pattern, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, pattern)
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var request ListRequest
	request.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	aggregationResult, err := a.store.Find(c, request, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(request.Query, aggregationResult)
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

	pattern, err := a.store.GetByID(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, pattern)
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
	request := EditRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	pattern, err := a.store.GetByID(c, request.ID, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if pattern.Type != request.Type {
		err = validation.NewSingleError("unchangeable", "Type", "Type", request)
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern.IsCorporate != *request.IsCorporate {
		err = validation.NewSingleError("unchangeable", "IsCorporate", "IsCorporate", request)
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern.IsCorporate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	pattern, err = a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, pattern)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	pattern, err := a.store.GetByID(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if pattern == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if pattern.IsCorporate {
		ok, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	ok, err := a.store.Delete(c, *pattern, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	canDeleteCorporate, err := a.enforcer.Enforce(userID, apisecurity.PermCorporatePattern, model.PermissionCan)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		pattern, err := a.store.GetByID(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if pattern == nil {
			return "", httperror.ErrNotFound
		}

		if pattern.IsCorporate && !canDeleteCorporate {
			return "", httperror.NewForbiddenError("")
		}

		ok, err := a.store.Delete(c, *pattern, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		return pattern.ID, nil
	}, a.errorResponder)
}

// CountAlarms
// @Param body body CountRequest true "body"
// @Success 200 {object} CountAlarmsResponse
func (a *api) CountAlarms(c *gin.Context) {
	request := CountRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	conf := a.configProvider.Get()
	ctx, cancel := context.WithTimeout(c, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
	defer cancel()

	res, err := a.store.CountAlarms(ctx, request, int64(conf.MaxMatchedItems))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

// CountEntities
// @Param body body CountRequest true "body"
// @Success 200 {object} CountEntitiesResponse
func (a *api) CountEntities(c *gin.Context) {
	request := CountRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	conf := a.configProvider.Get()
	ctx, cancel := context.WithTimeout(c, time.Duration(conf.CheckCountRequestTimeout)*time.Second)
	defer cancel()

	res, err := a.store.CountEntities(ctx, request, int64(conf.MaxMatchedItems))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}
