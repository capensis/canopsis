package pbehaviorcomment

import (
	"git.canopsis.net/canopsis/go-engines/lib/api/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type API interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type api struct {
	transformer ModelTransformer
	store       Store
}

func NewApi(
	transformer ModelTransformer,
	store Store,
) API {
	return &api{
		transformer: transformer,
		store:       store,
	}
}

// Create pbehavior comment
// @Summary Create pbehavior comment
// @Description Create pbehavior comment
// @Tags pbehavior-comments
// @ID pbehavior-comments-create
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body Request true "body"
// @Success 201 {object} pbehavior.Comment
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /pbehavior-comments [post]
func (a *api) Create(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	model := a.transformer.TransformRequestToModel(&request)
	ok, err := a.store.Insert(request.Pbehavior, model)
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, model)
}

// Delete pbehavior comment by id
// @Summary Delete pbehavior comment by id
// @Description Delete pbehavior comment by id
// @Tags pbehavior-comments
// @ID pbehavior-comment-delete-by-id
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "pbehavior comment id"
// @Success 204
// @Failure 404 {object} common.ErrorResponse
// @Router /pbehavior-comments/:id [delete]
func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c.Param("id"))
	if err != nil {
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
