package template

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/config"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"github.com/gin-gonic/gin"
)

type API interface {
	ValidateDeclareTicketRules(c *gin.Context)
	ValidateScenarios(c *gin.Context)
	ValidateEventFilterRules(c *gin.Context)
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
// @Param body body []Request true "body"
// @Success 200 {array} Response
func (a *api) ValidateDeclareTicketRules(c *gin.Context) {
	var request []Request
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response := make([]Response, len(request))
	for i, r := range request {
		isValid, errReport, wrnReports, err := a.validator.ValidateDeclareTicketRuleTemplate(r.Text)
		if err != nil {
			panic(err)
		}

		response[i] = Response{
			IsValid:  isValid,
			Err:      errReport,
			Warnings: wrnReports,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ValidateScenarios
// @Param body body []Request true "body"
// @Success 200 {array} Response
func (a *api) ValidateScenarios(c *gin.Context) {
	var request []Request
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response := make([]Response, len(request))
	for i, r := range request {
		isValid, errReport, wrnReports, err := a.validator.ValidateScenarioTemplate(r.Text)
		if err != nil {
			panic(err)
		}

		response[i] = Response{
			IsValid:  isValid,
			Err:      errReport,
			Warnings: wrnReports,
		}
	}

	c.JSON(http.StatusOK, response)
}

// ValidateEventFilterRules
// @Param body body []Request true "body"
// @Success 200 {array} Response
func (a *api) ValidateEventFilterRules(c *gin.Context) {
	var request []Request
	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	response := make([]Response, len(request))
	for i, r := range request {
		isValid, errReport, wrnReports, err := a.validator.ValidateEventFilterRuleTemplate(r.Text)
		if err != nil {
			panic(err)
		}

		response[i] = Response{
			IsValid:  isValid,
			Err:      errReport,
			Warnings: wrnReports,
		}
	}

	c.JSON(http.StatusOK, response)
}

func (a *api) GetEnvVars(c *gin.Context) {
	c.JSON(http.StatusOK, a.templateConfigProvider.Get().Vars)
}
