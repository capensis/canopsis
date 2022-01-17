package appinfo

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/auth"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	apisecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/model"
	"github.com/gin-gonic/gin"
)

type api struct {
	enforcer security.Enforcer
	store    Store
}

func NewApi(
	enforcer security.Enforcer,
	store Store,
) *api {
	return &api{
		enforcer: enforcer,
		store:    store,
	}
}

// Get application information
// @Summary Get application information
// @Description Get application information
// @Tags internal
// @ID internal-get-app-info
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 200 {object} AppInfoResponse
// @Router /app-info [get]
func (a *api) GetAppInfo(c *gin.Context) {
	response := AppInfoResponse{}
	var err error

	response.UserInterfaceConf, err = a.store.RetrieveUserInterfaceConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}
	response.VersionConf, err = a.store.RetrieveVersionConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}
	response.Login, err = a.store.RetrieveLoginConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}

	user, ok := c.Get(auth.UserKey)
	if ok {
		ok, err := a.enforcer.Enforce(user.(string), apisecurity.PermAppInfoRead, model.PermissionCan)
		if err != nil {
			panic(err)
		}

		if ok {
			response.GlobalConf, err = a.store.RetrieveGlobalConfig(c.Request.Context())
			if err != nil {
				panic(err)
			}

			remediation, err := a.store.RetrieveRemediationConfig(c.Request.Context())
			if err != nil {
				panic(err)
			}
			response.Remediation = &remediation
		}
	}

	c.JSON(http.StatusOK, response)
}

// Update user interface
// @Summary update user interface
// @Description update user interface
// @Tags internal
// @ID internal-update-user-interface
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Security BasicAuth
// @Param body body UserInterfaceConf true "body"
// @Success 200 {object} UserInterfaceConf
// @Failure 400 {object} common.ValidationErrorResponse
// @Router /internal/user_interface [post]
// @Router /internal/user_interface [put]
func (a *api) UpdateUserInterface(c *gin.Context) {
	request := UserInterfaceConf{
		MaxMatchedItems:          config.DefaultMaxMatchedItems,
		CheckCountRequestTimeout: config.DefaultCheckCountRequestTimeout,
	}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.UpdateUserInterfaceConfig(c.Request.Context(), &request)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, request)
}

// Delete user interface
// @Summary delete user interface
// @Description delete user interface
// @Tags internal
// @ID internal-delete-user-interface
// @Security ApiKeyAuth
// @Security BasicAuth
// @Success 204
// @Router /internal/user_interface [delete]
func (a *api) DeleteUserInterface(c *gin.Context) {
	err := a.store.DeleteUserInterfaceConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusNoContent, nil)
}
