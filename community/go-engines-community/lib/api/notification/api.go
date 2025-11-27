package notification

import (
	"errors"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/pagination"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type API interface {
	List(c *gin.Context)
	UpdateSettings(c *gin.Context)
	GetSettings(c *gin.Context)
}

type api struct {
	store          Store
	errorResponder httperror.Responder
}

func NewApi(store Store, errorResponder httperror.Responder) API {
	return &api{
		store:          store,
		errorResponder: errorResponder,
	}
}

// List
// @Success 200 {object} pagination.ListResponse{data=[]Response}
func (a *api) List(c *gin.Context) {
	r := pagination.GetDefaultQuery()
	if err := validation.Bind(c, &r); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	userID := c.MustGet(authctx.UserKey).(string)
	roleIDs := c.MustGet(authctx.Roles).([]string)
	aggregationResult, err := a.store.Find(c, r, userID, roleIDs)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	res := pagination.NewResponse(r, &aggregationResult)
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
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, notification)
}

// UpdateSettings
// @Param body body UpdateSettingsRequest true "body"
// @Success 200 {object} SettingsResponse
func (a *api) UpdateSettings(c *gin.Context) {
	request := UpdateSettingsRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	notification, err := a.store.UpdateSettings(c, request)
	if errors.Is(err, mongo.ErrNoDocuments) {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	} else if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, notification)
}
