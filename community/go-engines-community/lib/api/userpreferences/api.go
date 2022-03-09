package userpreferences

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin"
)

type API interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type api struct {
	store        Store
	actionLogger logger.ActionLogger
}

func NewApi(store Store, actionLogger logger.ActionLogger) API {
	return &api{
		store:        store,
		actionLogger: actionLogger,
	}
}

// Get user preferences by widget id
// @Summary Get user preferences by widget id
// @Description Get user preferences by widget id
// @Tags userpreference
// @ID userpreference-get-by-id
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param id path string true "widget id"
// @Success 200 {object} Response
// @Failure 404 {object} common.ErrorResponse
// @Router /user-preferences/{id} [get]
func (a *api) Get(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	response, err := a.store.Find(c.Request.Context(), userId, c.Param("id"))
	if err != nil {
		panic(err)
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update user preferences by widget id
// @Summary Update user preferences by widget id
// @Description Update user preferences by widget id
// @Tags userpreference
// @ID userpreference-update-by-widget-id
// @Accept json
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /user-preferences [put]
func (a api) Update(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response, isNew, err := a.store.Update(c.Request.Context(), userId, request)
	if err != nil {
		panic(err)
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	action := logger.ActionUpdate
	if isNew {
		action = logger.ActionCreate
	}

	err = a.actionLogger.Action(context.Background(), userId, logger.LogEntry{
		Action:    action,
		ValueType: logger.ValueTypeUserPreferences,
		ValueID:   response.Widget,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, response)
}
