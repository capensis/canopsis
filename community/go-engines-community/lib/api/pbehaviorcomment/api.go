package pbehaviorcomment

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type API interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type api struct {
	store Store
}

func NewApi(
	store Store,
) API {
	return &api{
		store: store,
	}
}

// Create
// @Param body body Request true "body"
// @Success 201 {object} Response
func (a *api) Create(c *gin.Context) {
	var request Request

	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))

		return
	}

	response, err := a.store.Insert(c, request)
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (a *api) Delete(c *gin.Context) {
	ok, err := a.store.Delete(c, c.Param("id"))
	if err != nil {
		valErr := common.ValidationError{}
		if errors.As(err, &valErr) {
			c.AbortWithStatusJSON(http.StatusBadRequest, valErr.ValidationErrorResponse())
			return
		}
		panic(err)
	}

	if !ok {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
