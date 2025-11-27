package entitycomment

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
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
	r.Entity = c.Param("entity")

	if err := validation.Bind(c, &r); err != nil {
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
	userID := c.MustGet(authctx.UserKey).(string)
	username := c.MustGet(authctx.Username).(string)

	var request Request

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Insert(c, request, userID, username)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Update
// @Param body body UpdateRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID := c.MustGet(authctx.UserKey).(string)
	username := c.MustGet(authctx.Username).(string)

	request := UpdateRequest{
		ID: c.Param("id"),
	}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response, err := a.store.Update(c, request, userID, username)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}
