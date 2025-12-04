package user

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/metrics"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.BulkAPI
	Patch(c *gin.Context)
	BulkPatch(c *gin.Context)
}

type api struct {
	store             Store
	metricMetaUpdater metrics.MetaUpdater
	errorResponder    httperror.Responder
}

func NewApi(
	store Store,
	metricMetaUpdater metrics.MetaUpdater,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:             store,
		metricMetaUpdater: metricMetaUpdater,
		errorResponder:    errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]User}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	roleIDs, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	users, err := a.store.Find(c, query, userID, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, users)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} User
func (a *api) Get(c *gin.Context) {
	user, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} User
func (a *api) Create(c *gin.Context) {
	var request CreateRequest
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	roleIDs, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	user, err := a.store.Insert(c, request, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.metricMetaUpdater.UpdateById(c, user.ID)

	c.JSON(http.StatusCreated, user)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} User
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	roleIDs, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	user, err := a.store.Update(c, BulkUpdateRequestItem(request), userID, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.metricMetaUpdater.UpdateById(c, user.ID)

	c.JSON(http.StatusOK, user)
}

// Patch
// @Param body body PatchRequest true "body"
// @Success 200 {object} User
func (a *api) Patch(c *gin.Context) {
	request := PatchRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	roleIDs, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	user, err := a.store.Patch(c, BulkPatchRequestItem(request), userID, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.metricMetaUpdater.UpdateById(c, user.ID)

	c.JSON(http.StatusOK, user)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	roleIDs, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, BulkDeleteRequestItem{ID: id}, userID, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	a.metricMetaUpdater.DeleteById(c, id)

	c.Status(http.StatusNoContent)
}

// BulkCreate
// @Param body body []CreateRequest true "body"
func (a *api) BulkCreate(c *gin.Context) {
	userIDs := make([]string, 0)
	requestRoles, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request CreateRequest) (string, error) {
		user, err := a.store.Insert(c, request, requestRoles)
		if err != nil {
			return "", err
		}

		userIDs = append(userIDs, user.ID)
		return user.ID, nil
	}, a.errorResponder)
	a.metricMetaUpdater.UpdateById(c, userIDs...)
}

// BulkUpdate
// @Param body body []BulkUpdateRequestItem true "body"
func (a *api) BulkUpdate(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	userIDs := make([]string, 0)
	requestRoles, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request BulkUpdateRequestItem) (string, error) {
		user, err := a.store.Update(c, request, userID, requestRoles)
		if err != nil {
			return "", err
		}

		if user == nil {
			return "", httperror.ErrNotFound
		}

		userIDs = append(userIDs, user.ID)

		return user.ID, nil
	}, a.errorResponder)
	a.metricMetaUpdater.UpdateById(c, userIDs...)
}

// BulkDelete
// @Param body body []BulkDeleteRequestItem true "body"
func (a *api) BulkDelete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	roles, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	userIDs := make([]string, 0)
	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request, userID, roles)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		userIDs = append(userIDs, request.ID)
		return request.ID, nil
	}, a.errorResponder)

	a.metricMetaUpdater.DeleteById(c, userIDs...)
}

// BulkPatch
// @Param body body []BulkPatchRequestItem true "body"
func (a *api) BulkPatch(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	userIDs := make([]string, 0)
	requestRoles, err := authctx.GetRoles(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	bulk.Handler(c, func(request BulkPatchRequestItem) (string, error) {
		user, err := a.store.Patch(c, request, userID, requestRoles)
		if err != nil {
			return "", err
		}

		if user == nil {
			return "", httperror.ErrNotFound
		}

		userIDs = append(userIDs, user.ID)

		return user.ID, nil
	}, a.errorResponder)
	a.metricMetaUpdater.UpdateById(c, userIDs...)
}
