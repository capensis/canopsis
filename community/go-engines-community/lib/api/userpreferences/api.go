package userpreferences

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/logger"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/widget"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type API interface {
	Get(c *gin.Context)
	Update(c *gin.Context)
}

type api struct {
	store        Store
	widgetStore  widget.Store
	enforcer     security.Enforcer
	actionLogger logger.ActionLogger
}

func NewApi(
	store Store,
	widgetStore widget.Store,
	enforcer security.Enforcer,
	actionLogger logger.ActionLogger,
) API {
	return &api{
		store:        store,
		widgetStore:  widgetStore,
		enforcer:     enforcer,
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
	widgetId := c.Param("id")

	ok, err := a.checkAccess(c.Request.Context(), widgetId, userId)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	response, err := a.store.Find(c.Request.Context(), userId, widgetId)
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

	ok, err := a.checkAccess(c.Request.Context(), request.Widget, userId)
	if err != nil {
		panic(err)
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
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

	err = a.actionLogger.Action(c, logger.LogEntry{
		Action:    action,
		ValueType: logger.ValueTypeUserPreferences,
		ValueID:   response.Widget,
	})
	if err != nil {
		a.actionLogger.Err(err, "failed to log action")
	}

	c.JSON(http.StatusOK, response)
}

func (a *api) checkAccess(ctx context.Context, widgetId, userId string) (bool, error) {
	viewIds, err := a.widgetStore.FindViewIds(ctx, []string{widgetId})
	if err != nil || len(viewIds) == 0 {
		return false, err
	}

	for _, viewId := range viewIds {
		ok, err := a.enforcer.Enforce(userId, viewId, model.PermissionRead)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}
