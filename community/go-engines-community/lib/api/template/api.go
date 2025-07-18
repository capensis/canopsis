package template

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"github.com/gin-gonic/gin"
)

type API interface {
	GetEnvVars(c *gin.Context)
}

type api struct {
	templateConfigProvider config.TemplateConfigProvider
}

func NewApi(templateConfigProvider config.TemplateConfigProvider) API {
	return &api{
		templateConfigProvider: templateConfigProvider,
	}
}

func (a *api) GetEnvVars(c *gin.Context) {
	c.JSON(http.StatusOK, a.templateConfigProvider.Get().Vars)
}
