package file

import (
	"errors"
	"fmt"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

func NewApi(enforcer security.Enforcer, store Store) Crud {
	return &api{
		store:    store,
		enforcer: enforcer,
	}
}

type Crud interface {
	List(*gin.Context)
	Create(*gin.Context)
	Get(*gin.Context)
	Delete(*gin.Context)
}

type api struct {
	enforcer security.Enforcer
	store    Store
}

// Create
// @Success 200 {array} File
func (a *api) Create(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		panic(err)
	}

	request := CreateRequest{}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	res, err := a.store.Create(c.Request.Context(), request.Public, form)
	if err != nil {
		validationError := ValidationError{}
		if errors.As(err, &validationError) {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.ValidationErrorResponse{
				Errors: map[string]string{
					validationError.field: validationError.Error(),
				},
			})
			return
		}

		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

func (a *api) Get(c *gin.Context) {
	m, err := a.store.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if m == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	if !m.IsPublic {
		user := c.MustGet(auth.UserKey)
		ok, err := a.enforcer.Enforce(user, apisecurity.ObjFile, model.PermissionRead)
		if err != nil {
			panic(err)
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
	res, err := a.store.List(c.Request.Context(), c.QueryArray("id"))
	if err != nil || res == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	for _, f := range res {
		if !f.IsPublic {
			user := c.MustGet(auth.UserKey)
			ok, err := a.enforcer.Enforce(user, apisecurity.ObjFile, model.PermissionRead)
			if err != nil {
				panic(err)
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
	ok, err := a.store.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.JSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.Status(http.StatusNoContent)
}
