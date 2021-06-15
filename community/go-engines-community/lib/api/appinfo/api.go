package appinfo

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	libsecurity "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security"
	"github.com/gin-gonic/gin"
)

var defaultInterval = IntervalUnit{
	Interval: 3,
	Unit:     "s",
}

type api struct {
	store  Store
	Config libsecurity.Config
}

func NewApi(
	store Store,
	config libsecurity.Config,
) *api {
	return &api{
		store:  store,
		Config: config,
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
	userInterface, err := a.store.RetrieveUserInterfaceConfig()
	if err != nil {
		panic(err)
	}
	tz, err := a.store.RetrieveTimezoneConf()
	if err != nil {
		panic(err)
	}
	version, err := a.store.RetrieveCanopsisVersionConfig()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, AppInfoResponse{
		UserInterfaceConf:   userInterface,
		TimezoneConf:        tz,
		CanopsisVersionConf: version,
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
	var config = LoginConfig{
		Providers:  make(map[string]int),
		Casconfig:  nil,
		Ldapconfig: nil,
	}
	for _, p := range a.Config.Security.AuthProviders {
		config.Providers[p] = 1
	}

	loginServices, err := a.store.RetrieveObjectConfig()
	if err != nil {
		panic(err)
	}

	for _, ser := range loginServices {
		switch ser.CrecordName {
		case Casconfig:
			ser.Fields[CrecordName] = ser.CrecordName
			config.Casconfig = ser.Fields
		case Ldapconfig:
			config.Ldapconfig = &struct {
				Enable bool `json:"enable"`
			}{Enable: ser.Enable}
		}
	}

	userInterface, err := a.store.RetrieveUserInterfaceConfig()
	if err != nil {
		panic(err)
	}

	version, err := a.store.RetrieveCanopsisVersionConfig()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, LoginConfigResponse{
		LoginConfig:         config,
		UserInterfaceConf:   userInterface,
		CanopsisVersionConf: version,
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
	var request UserInterfaceConf
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}
	if request.PopupTimeout == nil {
		request.PopupTimeout = &PopupTimeout{
			Info:  &defaultInterval,
			Error: &defaultInterval,
		}
	} else if request.PopupTimeout.Error == nil {
		request.PopupTimeout.Error = &defaultInterval
	} else if request.PopupTimeout.Info == nil {
		request.PopupTimeout.Info = &defaultInterval
	}

	err := a.store.UpdateUserInterfaceConfig(&request)
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
	err := a.store.DeleteUserInterfaceConfig()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusNoContent, nil)
}
