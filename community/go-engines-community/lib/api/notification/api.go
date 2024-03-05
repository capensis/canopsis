package notification

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	Update(c *gin.Context)
	Get(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// Get
// @Success 200 {object} common.PaginatedListResponse{data=[]Notification}
func (a *api) Get(c *gin.Context) {
	notification, err := a.store.Get(c.Request.Context())
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

	notification, err := a.store.Update(c.Request.Context(), request)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}
