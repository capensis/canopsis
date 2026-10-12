package apitoken

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/binding"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/ui/response"
	"github.com/gin-gonic/gin"
)

type API interface {
	Create(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder response.ErrorResponder
}

func NewAPI(store Store, errorResponder response.ErrorResponder) API {
	return &api{
		store:          store,
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
	if err := binding.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Find(c, r, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, pagination.NewResponse(r.Query, &res))
}

// Create
// @Param body body EditRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var r EditRequest
	if err := binding.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Insert(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, res)
}

func (a *api) Delete(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	err = a.store.Delete(c, c.Param("id"), userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.Status(http.StatusNoContent)
}
