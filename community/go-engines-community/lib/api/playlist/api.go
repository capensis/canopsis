package playlist

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/viewtab"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type api struct {
	store          Store
	tabStore       viewtab.Store
	enforcer       security.Enforcer
	errorResponder httperror.Responder
}

func NewApi(
	store Store,
	tabStore viewtab.Store,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) crud.API {
	return &api{
		store:          store,
		tabStore:       tabStore,
		enforcer:       enforcer,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Playlist}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ids, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		if s, ok := ids.([]string); ok {
			r.Ids = s
		}
	}

	playlists, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, playlists)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Playlist
func (a *api) Get(c *gin.Context) {
	playlist, err := a.store.GetByID(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if playlist == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, playlist)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Playlist
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.checkAccess(c, request.TabsList, request.Author)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	playlist, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, playlist)
}

// Update
// @Param body body Playlist true "body"
// @Success 200 {object} Playlist
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
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
	ok, err := a.checkAccess(c, request.TabsList, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

		return
	}

	playlist, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if playlist == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, playlist)
}

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
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

func (a *api) checkAccess(ctx context.Context, tabIds []string, userID string) (bool, error) {
	tabs, err := a.tabStore.Find(ctx, tabIds)
	if err != nil || len(tabs) != len(tabIds) {
		return false, err
	}

	for _, tab := range tabs {
		ok, err := a.enforcer.Enforce(userID, tab.View, model.PermissionRead)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}
