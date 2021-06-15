package notification

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
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

// Get notification
// @Summary Get notification settings
// @Description Get notification settings
// @Tags notification
// @ID notification-get
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param page query integer true "current page"
// @Param limit query integer true "items per page"
// @Success 200 {object} common.PaginatedListResponse{data=[]Notification}
// @Router /notification [get]
func (a *api) Get(c *gin.Context) {
	notification, err := a.store.Get()
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}

// Update notification
// @Summary Update notification settings
// @Description Update notification settings
// @Tags notification
// @ID notification-update
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param id path string true "type id"
// @Param body body Notification true "body"
// @Success 200 {object} Notification
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /notification [put]
func (a *api) Update(c *gin.Context) {
	request := Notification{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	notification, err := a.store.Update(request)
	if err == mongo.ErrNoDocuments {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}
