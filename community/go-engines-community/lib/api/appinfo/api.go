package appinfo

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/httperror"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/validation"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
)

type API interface {
	GetAppInfo(c *gin.Context)
	UpdateUserInterface(c *gin.Context)
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

// GetAppInfo
// @Success 200 {object} AppInfoResponse
func (a *api) GetAppInfo(c *gin.Context) {
	response := AppInfoResponse{}
	var err error

	response.UserInterfaceConf, err = a.store.RetrieveUserInterfaceConfig(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	response.VersionConf, err = a.store.RetrieveVersionConfig(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	response.Login = a.store.RetrieveLoginConfig()
	response.GlobalConf, err = a.store.RetrieveGlobalConfig(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	remediation, err := a.store.RetrieveRemediationConfig(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	response.Remediation = &remediation

	response.Maintenance, err = a.store.RetrieveMaintenanceState(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response.DefaultColorTheme, err = a.store.RetrieveDefaultColorTheme(c, response.UserInterfaceConf.DefaultColorTheme)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	response.SerialName, err = a.store.RetrieveSerialName(c)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserInterface
// @Success 200 {object} UserInterfaceConf
func (a *api) UpdateUserInterface(c *gin.Context) {
	request := UserInterfaceConf{
		MaxMatchedItems:          config.UserInterfaceMaxMatchedItems,
		CheckCountRequestTimeout: config.UserInterfaceCheckCountRequestTimeout,
	}

	if err := validation.Bind(c, &request); err != nil {
		a.errorResponder.Respond(c, err)

		return
	}

	err := a.store.UpdateUserInterfaceConfig(c, &request)
	if err != nil {
		a.errorResponder.Respond(c, err)

		return
	}
	c.JSON(http.StatusOK, request)
}
