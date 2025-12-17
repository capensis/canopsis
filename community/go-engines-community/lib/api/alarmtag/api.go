package alarmtag

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/bulk"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

type API interface {
	crud.API
	BulkDelete(c *gin.Context)
	ListLabels(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

func NewApi(store Store, errorResponder httperror.Responder) API {
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

	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	tags, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response := pagination.NewResponse(r.Query, tags)
	c.JSON(http.StatusOK, response)
}

// ListLabels
// @Success 200 {object} pagination.ListResponse{data=[]LabelResponse}
func (a *api) ListLabels(c *gin.Context) {
	var r ListLabelsRequest
	r.Query = pagination.GetDefaultQuery()
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	tags, err := a.store.FindLabels(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response := pagination.NewResponse(r.Query, tags)
	c.JSON(http.StatusOK, response)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	response, err := a.store.GetByID(c, c.Param("id"))
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

// Create
// @Param body body CreateRequest true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	request := CreateRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Create(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, response)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Update(c, request)
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

func (a *api) Delete(c *gin.Context) {
	id := c.Param("id")

	tag, err := a.store.GetByID(c, id)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if tag == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

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

	bulk.Handler(c, func(request BulkDeleteRequestItem) (string, error) {
		ok, err := a.store.Delete(c, request.ID, userID)
		if err != nil {
			return "", err
		}

		if !ok {
			return "", httperror.ErrNotFound
		}

		return request.ID, nil
	}, a.errorResponder)
}
