package notification

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
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

// Get
// @Success 200 {object} common.PaginatedListResponse{data=[]Notification}
func (a *api) Get(c *gin.Context) {
	notification, err := a.store.Get(c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}

// Update
// @Param body body Notification true "body"
// @Success 200 {object} Notification
func (a *api) Update(c *gin.Context) {
	request := Notification{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	notification, err := a.store.Update(c, request)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}
