package appinfo

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"github.com/gin-gonic/gin"
)

type api struct {
	store Store
}

func NewApi(
	store Store,
) *api {
	return &api{
		store: store,
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
// @Router /internal/app_info [get]
func (a *api) GetAppInfo(c *gin.Context) {
	userInterface, err := a.store.RetrieveUserInterfaceConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}
	tz, err := a.store.RetrieveTimezoneConf(c.Request.Context())
	if err != nil {
		panic(err)
	}
	version, err := a.store.RetrieveVersionConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, AppInfoResponse{
		UserInterfaceConf: userInterface,
		TimezoneConf:      tz,
		VersionConf:       version,
	})
}

// Get login information
// @Summary Get login information
// @Description Get login information
// @Tags internal
// @ID internal-get-login-info
// @Produce json
// @Success 200 {object} LoginConfigResponse
// @Router /internal/login_info [get]
func (a *api) LoginInfo(c *gin.Context) {
	login, err := a.store.RetrieveLoginConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}

	userInterface, err := a.store.RetrieveUserInterfaceConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}

	version, err := a.store.RetrieveVersionConfig(c.Request.Context())
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, LoginConfigResponse{
		LoginConfig:       login,
		UserInterfaceConf: userInterface,
		VersionConf:       version,
	})
}

// update user interface
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
		MaxMatchedItems:           config.DefaultMaxMatchedItems,
		CheckCountRequestTimeout:  config.DefaultCheckCountRequestTimeout,
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

// delete user interface
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
