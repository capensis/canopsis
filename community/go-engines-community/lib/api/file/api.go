package file

import (
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ErrorResponse{Error: "Files are missing."})
		return
	}

	request := CreateRequest{}
	if err := validation.BindQuery(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res, err := a.store.Create(c, request.Public, form)
	if err != nil {
		validationError := common.ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, validationError.ValidationErrorResponse())
			return
		}

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
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
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
			c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
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
		c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{Errors: map[string]string{
			"id": "ID is missing.",
		}})

		return
	}

	res, err := a.store.List(c, ids)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if len(res) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
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
				c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
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
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
