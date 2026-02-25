package webhook

import (
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/webhook"
	"github.com/go-playground/validator/v10"
)

type CheckTicketStatusValidator struct {
	templateExecutor template.Executor
	validStatuses    string
}

func NewCheckTicketStatusValidator(templateExecutor template.Executor) *CheckTicketStatusValidator {
	validStatuses := []string{
		strconv.Itoa(webhook.TicketStatusOpen),
		strconv.Itoa(webhook.TicketStatusAssigned),
		strconv.Itoa(webhook.TicketStatusInProgress),
		strconv.Itoa(webhook.TicketStatusClosed),
	}
	return &CheckTicketStatusValidator{
		templateExecutor: templateExecutor,
		validStatuses:    strings.Join(validStatuses, " "),
	}
}

func (v *CheckTicketStatusValidator) ValidateCheckTicketStatus(sl validator.StructLevel) {
	r := sl.Current().Interface().(webhook.CheckTicketStatus)

	if r.TicketStatus == "" && r.TicketStatusTpl == "" {
		sl.ReportError(r.TicketStatus, "TicketStatus", "TicketStatus", "required_or", "TicketStatus")
		sl.ReportError(r.TicketStatusTpl, "TicketStatusTpl", "TicketStatusTpl", "required_or", "TicketStatusTpl")
	}

	for k, header := range r.Request.Headers {
		if header != "" {
			parsedValue := v.templateExecutor.Parse(header)
			if parsedValue.Err != nil {
				sl.ReportError(header, k, "Headers."+k, "template", "")
			}
		}
	}

	hasClosedMapping := false

	if len(r.StatusMapping) == 0 {
		sl.ReportError(r.StatusMapping, "StatusMapping", "StatusMapping", "required", "")
	} else {
		for _, val := range r.StatusMapping {
			switch val {
			case webhook.TicketStatusOpen, webhook.TicketStatusAssigned, webhook.TicketStatusInProgress:
			case webhook.TicketStatusClosed:
				hasClosedMapping = true
			default:
				sl.ReportError(r.StatusMapping, "StatusMapping", "StatusMapping", "oneof", v.validStatuses)
				return
			}
		}

		if !hasClosedMapping {
			sl.ReportError(r.StatusMapping, "StatusMapping", "StatusMapping", "required_closed_mapping", "")
		}
	}
}
