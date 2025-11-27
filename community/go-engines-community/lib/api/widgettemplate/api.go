package widgettemplate

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/crud"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

func NewApi(store Store, errorResponder httperror.Responder) crud.API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

type api struct {
	store          Store
	errorResponder httperror.Responder
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

	tpls, err := a.store.Find(c, r)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r.Query, tpls)
	c.JSON(http.StatusOK, res)
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	tpl, err := a.store.GetOneById(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if tpl == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, tpl)
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

	tpl, err := a.store.Insert(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusCreated, tpl)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	request := EditRequest{
		ID: c.Param("id"),
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	tpl, err := a.store.GetOneById(c, request.ID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if tpl == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if tpl.Type != request.Type {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{"type": "Type cannot be changed"}})
		return
	}

	tpl, err = a.store.Update(c, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if tpl == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, tpl)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"), c.MustGet(authctx.UserKey).(string))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
