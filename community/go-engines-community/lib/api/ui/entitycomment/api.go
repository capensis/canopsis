package entitycomment

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
}

type api struct {
	store          Store
	errorResponder response.ErrorResponder
}

func NewApi(store Store, errorResponder response.ErrorResponder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	var r ListRequest
	r.Query = pagination.GetDefaultQuery()
	r.Entity = c.Param("entity")

	if err := binding.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	comments, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response := pagination.NewResponse(r.Query, comments)
	c.JSON(http.StatusOK, response)
}

// Create
// @Param body body Request true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	var request Request

	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Insert(c, request, userID, username)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusCreated, response)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	username, err := authctx.GetUsername(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := binding.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Update(c, request, userID, username)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	c.JSON(http.StatusOK, response)
}
