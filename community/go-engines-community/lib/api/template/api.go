package template

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"github.com/gin-gonic/gin"
)

type API interface {
	ValidateDeclareTicketRules(c *gin.Context)
	ValidateScenarios(c *gin.Context)
}

type api struct {
	validator validator.Validator
}

func NewApi(validator validator.Validator) API {
	return &api{validator: validator}
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

	isValid, errReport, wrnReports := a.validator.ValidateDeclareTicketRuleTemplate(request.Text)
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

	isValid, errReport, wrnReports := a.validator.ValidateScenarioTemplate(request.Text)
	c.JSON(http.StatusOK, Response{
		IsValid:  isValid,
		Err:      errReport,
		Warnings: wrnReports,
	})
}
