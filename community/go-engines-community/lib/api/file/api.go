package file

import (
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(*gin.Context)
	Create(*gin.Context)
	Get(*gin.Context)
	Delete(*gin.Context)
}

func NewApi(enforcer security.Enforcer, store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		enforcer:       enforcer,
		errorResponder: errorResponder,
	}
}

type api struct {
	enforcer       security.Enforcer
	store          Store
	errorResponder httperror.Responder
}

// Create
// @Success 200 {array} File
func (a *api) Create(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		a.errorResponder.Respond(c, validation.NewInvalidRequestBodyError(err))

		return
	}

	request := CreateRequest{}
	if err = validation.BindQuery(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Create(c, request.Public, form)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Get(c *gin.Context) {
	m, err := a.store.Get(c, c.Param("id"))
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if m == nil {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	if !m.IsPublic {
		userID, err := authctx.GetUserKey(c)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		ok, err := a.enforcer.Enforce(userID, apisecurity.ObjFile, model.PermissionRead)
		if err != nil {
			a.errorResponder.Respond(c, err)

			return
		}

		if !ok {
			a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

			return
		}
	}

	c.Header("Etag", fmt.Sprintf("%q", m.Etag))
	path := a.store.GetFilepath(*m)
	c.FileAttachment(path, m.FileName)
}

// List
// @Success 200 {object} []File
func (a *api) List(c *gin.Context) {
	ids := c.QueryArray("id")
	if len(ids) == 0 {
		err := validation.NewSingleError("required", "id", "id", nil)
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.List(c, ids)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if len(res) == 0 {
		a.errorResponder.Respond(c, httperror.ErrNotFound)

		return
	}

	for _, f := range res {
		if !f.IsPublic {
			userID, err := authctx.GetUserKey(c)
			if err != nil {
				a.errorResponder.Respond(c, err)

				return
			}

			ok, err := a.enforcer.Enforce(userID, apisecurity.ObjFile, model.PermissionRead)
			if err != nil {
				a.errorResponder.Respond(c, err)

				return
			}

			if !ok {
				a.errorResponder.Respond(c, httperror.NewForbiddenError(""))

				return
			}

			break
		}
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
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
