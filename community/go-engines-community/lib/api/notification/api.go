package notification

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type API interface {
	List(c *gin.Context)
	UpdateSettings(c *gin.Context)
	GetSettings(c *gin.Context)
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

// List
// @Success 200 {object} common.PaginatedListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	r := pagination.GetDefaultQuery()
	if err := c.ShouldBind(&r); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, r))
		return
	}

	userID := c.MustGet(auth.UserKey).(string)
	roleIDs := c.MustGet(auth.RolesKey).([]string)
	aggregationResult, err := a.store.Find(c, r, userID, roleIDs)
	if err != nil {
		panic(err)
	}

	res, err := common.NewPaginatedResponse(r, &aggregationResult)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetSettings
// @Success 200 {object} SettingsResponse
func (a *api) GetSettings(c *gin.Context) {
	notification, err := a.store.GetSettings(c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}

// UpdateSettings
// @Param body body UpdateSettingsRequest true "body"
// @Success 200 {object} SettingsResponse
func (a *api) UpdateSettings(c *gin.Context) {
	request := UpdateSettingsRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	notification, err := a.store.UpdateSettings(c, request)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, notification)
}
