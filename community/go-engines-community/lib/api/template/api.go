package template

import (
	"net/http"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template/validator"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/gin-gonic/gin"
)

type API interface {
	ValidateDeclareTicket(c *gin.Context)
}

type api struct {
	validator validator.Validator
}

func NewApi(validator validator.Validator) API {
	return &api{validator: validator}
}

type AlarmWithEntity struct {
	types.Alarm `bson:",inline"`
	Entity      types.Entity `bson:"entity" json:"entity"`
}

// ValidateDeclareTicket
// @Param body body Request true "body"
// @Success 200 {object} Response
func (a *api) ValidateDeclareTicket(c *gin.Context) {
	var request Request

	if err := c.ShouldBind(&request); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationErrorResponse(err, request))
		return
	}

	isValid, report := a.validator.ValidateDeclareTicketTemplate(request.Text)
	c.JSON(http.StatusOK, Response{
		IsValid: isValid,
		Report:  report,
	})
}
