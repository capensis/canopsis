package userpreferences

import (
	"context"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/authctx"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
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
	store          Store
	widgetStore    widget.Store
	enforcer       security.Enforcer
	errorResponder httperror.Responder
}

func NewApi(
	store Store,
	widgetStore widget.Store,
	enforcer security.Enforcer,
	errorResponder httperror.Responder,
) API {
	return &api{
		store:          store,
		widgetStore:    widgetStore,
		enforcer:       enforcer,
		errorResponder: errorResponder,
	}
}

// Get
// @Success 200 {object} Response
func (a *api) Get(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	widgetId := c.Param("id")

	ok, err := a.checkAccess(c, widgetId, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	response, err := a.store.Find(c, userID, widgetId)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update
// @Param body body EditRequest true "body"
// @Success 200 {object} Response
func (a *api) Update(c *gin.Context) {
	userID, err := authctx.GetUserKey(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	request := EditRequest{}
	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	ok, err := a.checkAccess(c, request.Widget, userID)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	if !ok {
		c.AbortWithStatusJSON(http.StatusForbidden, common.ForbiddenResponse)
		return
	}

	response, err := a.store.Update(c, userID, request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, common.NotFoundResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (a *api) checkAccess(ctx context.Context, widgetId, userID string) (bool, error) {
	tabInfos, err := a.widgetStore.FindTabPrivacySettings(ctx, []string{widgetId})
	if err != nil || len(tabInfos) == 0 {
		return false, err
	}

	for _, tabInfo := range tabInfos {
		if tabInfo.IsPrivate && tabInfo.Author == userID {
			continue
		}

		ok, err := a.enforcer.Enforce(userID, tabInfo.View, model.PermissionRead)
		if err != nil || !ok {
			return false, err
		}
	}

	return true, nil
}
