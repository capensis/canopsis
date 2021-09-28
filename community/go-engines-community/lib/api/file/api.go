package file

import (
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

// Upload files
// @Summary Upload files
// @Description Upload files
// @Tags files
// @ID file-upload
// @Accept multipart/form-data
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param files formData string true "request"
// @Param public query bool false "file visibility"
// @Success 200 {array} File
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /file [post]
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
		panic(err)
	}

	c.JSON(http.StatusOK, res)
}

// Get file
// @Summary Get file by ID
// @Description Get file content by ID or download with file name
// @Tags files
// @ID files-get
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "file id"
// @Success 200 {object} http.Response
// @Failure 404 {object} common.ErrorResponse
// @Router /file/{id} [get]
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

// List all files
// @Summary List files by ID
// @Description Get list of file objects by ID
// @Tags files
// @ID files-list
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Success 200 {object} []File
// @Failure 400 {object} common.ValidationErrorResponse
// @Success 404 {object} http.Response{content=string} "Not Found"
// @Router /file [get]
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

// Delete file
// @Summary Delete file
// @Description Delete file by ID
// @Tags files
// @ID files-delete
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "file id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /file/{id} [delete]
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
