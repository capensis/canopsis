package template

import (
	"net/http"
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"github.com/gin-gonic/gin"
)

type API interface {
	ValidateDeclareTicketRules(c *gin.Context)
	ValidateScenarios(c *gin.Context)
	GetEnvVars(c *gin.Context)
}

type api struct {
	validator              validator.Validator
	templateConfigProvider config.TemplateConfigProvider
}

func NewApi(validator validator.Validator, templateConfigProvider config.TemplateConfigProvider) API {
	return &api{
		validator:              validator,
		templateConfigProvider: templateConfigProvider,
	}
}

// ValidateDeclareTicketRules
// @Param body body Request true "body"
// @Success 200 {object} Response
func (a *api) ValidateDeclareTicketRules(c *gin.Context) {
	var request Request

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	isValid, errReport, wrnReports, err := a.validator.ValidateDeclareTicketRuleTemplate(request.Text)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, Response{
		IsValid:  isValid,
		Err:      errReport,
		Warnings: wrnReports,
	})
}

// ValidateScenarios
// @Param body body Request true "body"
// @Success 200 {object} Response
func (a *api) ValidateScenarios(c *gin.Context) {
	var request Request

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	isValid, errReport, wrnReports, err := a.validator.ValidateScenarioTemplate(request.Text)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, Response{
		IsValid:  isValid,
		Err:      errReport,
		Warnings: wrnReports,
	})
}

func (a *api) GetEnvVars(c *gin.Context) {
	envVars := a.templateConfigProvider.Get().Vars
	response := make([]string, len(envVars))
	i := 0
	for v := range envVars {
		response[i] = "." + template.EnvVar + "." + v
		i++
	}

	sort.Strings(response)
	c.JSON(http.StatusOK, response)
}
