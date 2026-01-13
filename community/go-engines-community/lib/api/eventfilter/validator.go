package eventfilter

import (
	"context"
	"strconv"
	"strings"

	apiexternaldata "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/externaldatatable"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/externaldata"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
)

type Validator struct {
	templateExecutor template.Executor
}

func NewValidator(templateExecutor template.Executor) *Validator {
	return &Validator{
		templateExecutor: templateExecutor,
	}
}

func (v *Validator) ValidateEditRequest(ctx context.Context, sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)
	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if r.Config.Component == "" &&
			r.Config.Resource == "" &&
			r.Config.Connector == "" &&
			r.Config.ConnectorName == "" &&
			r.Config.Upstream == "" {
			sl.ReportError(r.Config, "Config", "Config", "required", "")
		}
	case eventfilter.RuleTypeEnrichment:
		if len(r.Config.Actions) == 0 {
			sl.ReportError(r.Config.Actions, "Actions", "Config.Actions", "required", "")
		}

		for i, action := range r.Config.Actions {
			switch action.Type {
			case eventfilter.ActionSetFieldFromTemplate,
				eventfilter.ActionSetEntityInfoFromTemplate,
				eventfilter.ActionSetTagsFromTemplate:
				structNs := "Config.Actions." + strconv.Itoa(i) + ".Value"
				strVal, ok := action.Value.(string)
				if !ok {
					sl.ReportError(action.Value, "Value", structNs, "value_string", "")
					continue
				}

				if strVal != "" {
					parsedValue := v.templateExecutor.Parse(strVal)
					if parsedValue.Err != nil {
						sl.ReportError(strVal, "Value", structNs, "template", "")
					}
				}
			}
		}

		validOutcome := []string{
			eventfilter.OutcomePass,
			eventfilter.OutcomeDrop,
			eventfilter.OutcomeBreak,
		}
		if r.Config.OnSuccess == "" {
			sl.ReportError(r.Config.OnSuccess, "OnSuccess", "Config.OnSuccess", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnSuccess {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnSuccess, "OnSuccess", "Config.OnSuccess", "oneof", strings.Join(validOutcome, " "))
			}
		}

		if r.Config.OnFailure == "" {
			sl.ReportError(r.Config.OnFailure, "OnFailure", "Config.OnFailure", "required_if", "Type enrichment")
		} else {
			switch r.Config.OnFailure {
			case eventfilter.OutcomePass, eventfilter.OutcomeDrop, eventfilter.OutcomeBreak:
			default:
				sl.ReportError(r.Config.OnFailure, "OnFailure", "Config.OnFailure", "oneof", strings.Join(validOutcome, " "))
			}
		}
	}

	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required", "")
		}

		if len(r.EntityPattern) > 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "must_be_empty", "")
		}
	default:
		if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" && len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required_or", "EntityPattern")
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "EventPattern")
		}
	}

	if r.Start == nil && r.Stop != nil {
		sl.ReportError(r.Start, "Start", "Start", "required_with", "Stop")
	}

	if r.Stop == nil && r.Start != nil {
		sl.ReportError(r.Stop, "Stop", "Stop", "required_with", "Start")
	}

	if r.Stop != nil && r.Start != nil && r.Start.Unix() >= r.Stop.Unix() {
		sl.ReportError(r.Stop, "Stop", "Stop", "gtfield", "Start")
	}

	if r.RRule != "" && r.Stop == nil && r.Start == nil {
		sl.ReportError(r.Start, "Start", "Start", "required_with", "RRule")
		sl.ReportError(r.Stop, "Stop", "Stop", "required_with", "RRule")
	}

	if r.RRule != "" && !v.checkRrule(r.RRule) {
		sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
	}

	apiexternaldata.ValidateRefParameters(sl, v.templateExecutor, r.ExternalData, []string{externaldata.RefTypeAPI, externaldata.RefTypeTable})
}

func (v *Validator) ValidateTemplateRuleRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(TemplateRuleRequest)
	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if r.Config.Component == "" &&
			r.Config.Resource == "" &&
			r.Config.Connector == "" &&
			r.Config.ConnectorName == "" {
			sl.ReportError(r.Config, "Config", "Config", "required", "")
		}
	case eventfilter.RuleTypeEnrichment:
		if len(r.Config.Actions) == 0 {
			sl.ReportError(r.Config.Actions, "Actions", "Config.Actions", "required", "")
		}
	}

	switch r.Type {
	case eventfilter.RuleTypeChangeEntity:
		if len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required", "")
		}

		if len(r.EntityPattern) > 0 {
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "must_be_empty", "")
		}
	default:
		if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" && len(r.EventPattern) == 0 {
			sl.ReportError(r.EventPattern, "EventPattern", "EventPattern", "required_or", "EntityPattern")
			sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "EventPattern")
		}
	}
}

func (v *Validator) checkRrule(r string) bool {
	_, err := rrule.StrToROption(r)
	return err == nil
}
