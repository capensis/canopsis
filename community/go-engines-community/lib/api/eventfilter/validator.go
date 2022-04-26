package eventfilter

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"strings"
)

type eventfilterValidator struct {
	dbClient mongo.DbClient
}

func (v *eventfilterValidator) ValidateEventFilter(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.Type == eventfilter.RuleTypeChangeEntity &&
		r.Config.Component == "" &&
		r.Config.Resource == "" &&
		r.Config.Connector == "" &&
		r.Config.ConnectorName == "" {
		sl.ReportError(r.Config, "config", "Config", "required", "")
	}

	if r.Type == eventfilter.RuleTypeEnrichment {
		if len(r.Config.Actions) == 0 {
			sl.ReportError(r.Config.Actions, "actions", "Actions", "required", "")
		}

		validOutcome := []string{
			eventfilter.OutcomePass,
			eventfilter.OutcomeDrop,
			eventfilter.OutcomeBreak,
		}
		if r.Config.OnSuccess == "" {
			sl.ReportError(r.Config.OnSuccess, "on_success", "OnSuccess", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnSuccess {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnSuccess, "on_success", "OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}

		if r.Config.OnFailure == "" {
			sl.ReportError(r.Config.OnFailure, "on_failure", "OnFailure", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnFailure {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnFailure, "on_failure", "OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}
}

func NewValidator(dbClient mongo.DbClient) *eventfilterValidator {
	return &eventfilterValidator{dbClient: dbClient}
}
