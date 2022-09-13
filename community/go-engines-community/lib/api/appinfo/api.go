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
	response.Login, err = a.store.RetrieveLoginConfig(c)
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
			response.GlobalConf, err = a.store.RetrieveGlobalConfig(c)
			if err != nil {
				panic(err)
			}

			remediation, err := a.store.RetrieveRemediationConfig(c)
			if err != nil {
				panic(err)
			}
			response.Remediation = &remediation
		}
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

func (a *api) GetApiSecurity(c *gin.Context) {
	response, err := a.store.RetrieveApiSecurityConfig(c)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, response)
}

func (a *api) UpdateApiSecurity(c *gin.Context) {
	var request map[string]apisecurity.AuthMethodConf

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response, err := a.store.UpdateApiSecurityConfig(c, request)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, response)
}
