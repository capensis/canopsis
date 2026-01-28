package viewgroup

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/middleware"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type api struct {
	store          Store
	errorResponder response.ErrorResponder
}

func NewApi(store Store, errorResponder response.ErrorResponder) crud.API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]ViewGroup}
func (a *api) List(c *gin.Context) {
	var query ListRequest
	query.Query = pagination.GetDefaultQuery()
	var err error
	query.UserID, err = authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if err = binding.Bind(c, &query); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	var authorizedIds []string
	ids, ok := c.Get(middleware.AuthorizedIds)
	if ok {
		if s, ok := ids.([]string); ok {
			authorizedIds = s
		}
	}

	var ownedIds []string
	ids, ok = c.Get(middleware.OwnedIds)
	if ok {
		if s, ok := ids.([]string); ok {
			ownedIds = s
		}
	}

	viewgroups, err := a.store.Find(c, query, authorizedIds, ownedIds)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(query.Query, viewgroups)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} ViewGroup
func (a *api) Get(c *gin.Context) {
	viewgroup, err := a.store.GetOneBy(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if viewgroup == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, viewgroup)
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} ViewGroup
func (a *api) Create(c *gin.Context) {
	request := EditRequest{}
	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	group, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, group)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} ViewGroup
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	group, err := a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if group == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, group)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.store.Delete(c, c.Param("id"), userID)
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
