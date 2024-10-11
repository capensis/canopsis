package appinfo

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
)

type API interface {
	GetAppInfo(c *gin.Context)
	UpdateUserInterface(c *gin.Context)
	DeleteUserInterface(c *gin.Context)
}

type api struct {
	store Store
}

func NewApi(store Store) API {
	return &api{
		store: store,
	}
}

// GetAppInfo
// @Success 200 {object} AppInfoResponse
func (a *api) GetAppInfo(c *gin.Context) {
	response := AppInfoResponse{}
	var err error

	response.UserInterfaceConf, err = a.store.RetrieveUserInterfaceConfig(c)
	if err != nil {
		panic(err)
	}
	response.VersionConf, err = a.store.RetrieveVersionConfig(c)
	if err != nil {
		panic(err)
	}
	response.Login = a.store.RetrieveLoginConfig()
	response.GlobalConf, err = a.store.RetrieveGlobalConfig(c)
	if err != nil {
		panic(err)
	}

	remediation, err := a.store.RetrieveRemediationConfig(c)
	if err != nil {
		panic(err)
	}
	response.Remediation = &remediation

	response.Maintenance, err = a.store.RetrieveMaintenanceState(c)
	if err != nil {
		panic(err)
	}

	response.DefaultColorTheme, err = a.store.RetrieveDefaultColorTheme(c)
	if err != nil {
		panic(err)
	}

	response.SerialName, err = a.store.RetrieveSerialName(c)
	if err != nil {
		panic(err)
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

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	err := a.store.UpdateUserInterfaceConfig(c, &request)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, request)
}

func (a *api) DeleteUserInterface(c *gin.Context) {
	err := a.store.DeleteUserInterfaceConfig(c)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusNoContent, nil)
}
