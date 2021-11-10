package userpreferences

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"github.com/gin-gonic/gin"
)

type API interface {
	List(c *gin.Context)
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

// Find all user preferences
// @Summary Find all user preferences
// @Description Find all user preferences
// @Tags userpreference
// @ID userpreference-find-all
// @Produce json
// @Security JWTAuth
// @Security BasicAuth
// @Success 200 {array} Response
// @Router /user-preferences [get]
func (a api) List(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	response, err := a.store.Find(c.Request.Context(), userId)
	if err != nil {
		panic(err)
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
// @Param id path string true "widget id"
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
// @Failure 400 {object} common.ValidationErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /user-preferences/{id} [put]
func (a api) Update(c *gin.Context) {
	userId := c.MustGet(auth.UserKey).(string)
	request := EditRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	widgetId := c.Param("id")
	response, isNew, err := a.store.Update(c.Request.Context(), userId, widgetId, request)
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    action,
		ValueType: logger.ValueTypeUserPreferences,
		ValueID:   widgetId,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, response)
}
