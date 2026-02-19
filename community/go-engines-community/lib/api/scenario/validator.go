package scenario

import (
	"slices"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/template"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
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

func (v *Validator) ValidateActionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ActionRequest)

	if r.Type != "" {
		v.validateActionParametersRequest(sl, r.Type, r.Parameters)
	}

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, patternfields.GetForbiddenFieldsInEntityPattern(mongo.ScenarioCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			patternfields.GetForbiddenFieldsInAlarmPattern(mongo.ScenarioCollection),
			patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ScenarioCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}

func (v *Validator) validateActionParametersRequest(sl validator.StructLevel, t string, params action.Parameters) {
	switch t {
	case types.ActionTypeAssocTicket:
		if params.Ticket == "" {
			sl.ReportError(params.Ticket, "Ticket", "Parameters.Ticket", "required", "")
		}
	case types.ActionTypeChangeState:
		if params.State == nil {
			sl.ReportError(params.State, "State", "Parameters.State", "required", "")
		} else {
			validTypes := []types.CpsNumber{
				types.AlarmStateOK,
				types.AlarmStateMinor,
				types.AlarmStateMajor,
				types.AlarmStateCritical,
			}
			var param strings.Builder
			for i := range validTypes {
				param.WriteString(strconv.Itoa(int(validTypes[i])))
				if i < len(validTypes)-1 {
					param.WriteString(" ")
				}
			}

			if !slices.Contains(validTypes, *params.State) {
				sl.ReportError(params.State, "State", "Parameters.State", "oneof", param.String())
			}
		}
	case types.ActionTypeSnooze:
		if params.Duration == nil {
			sl.ReportError(params.Duration, "Duration", "Parameters.Duration", "required", "")
		}
	case types.ActionTypePbehavior:
		if params.Name == "" {
			sl.ReportError(params.Name, "Name", "Parameters.Name", "required", "")
		}
		if params.Reason == "" {
			sl.ReportError(params.Reason, "Reason", "Parameters.Reason", "required", "")
		}
		if params.Type == "" {
			sl.ReportError(params.Type, "Type", "Parameters.Type", "required", "")
		}
		// Validate rrule
		if params.RRule != "" {
			_, err := rrule.StrToROption(params.RRule)
			if err != nil {
				sl.ReportError(params.RRule, "RRule", "Parameters.RRule", "rrule", "")
			}
		}

		// Validate time
		if params.Tstart == nil && params.Tstop != nil {
			sl.ReportError(params.Tstart, "Tstart", "Parameters.Tstart", "required_with", "Tstop")
		}
		if params.Tstart != nil && params.Tstop == nil {
			sl.ReportError(params.Tstop, "Tstop", "Parameters.Tstop", "required_with", "Tstart")
		}
		if params.Duration == nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Duration, "Duration", "Parameters.Duration", "required_with", "StartOnTrigger")
		}
		if params.Duration != nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.StartOnTrigger, "StartOnTrigger", "Parameters.StartOnTrigger", "required_with", "Duration")
		}
		if params.Tstart == nil && params.Tstop == nil && params.Duration == nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.Tstart, "Tstart", "Parameters.Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "StartOnTrigger", "Parameters.StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Tstart, "Tstart", "Parameters.Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "StartOnTrigger", "Parameters.StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.Tstop != nil && params.Tstop.Before(*params.Tstart) {
			sl.ReportError(params.Tstop, "Tstop", "Parameters.Tstop", "gtfield", "Tstart")
		}
	case types.ActionTypeWebhook:
		if params.Request == nil {
			sl.ReportError(params.Request, "Request", "Parameters.Request", "required", "")
		} else {
			for k, header := range params.Request.Headers {
				if header != "" {
					parsedValue := v.templateExecutor.Parse(header)
					if parsedValue.Err != nil {
						sl.ReportError(header, k, "Parameters.Request.Headers."+k, "template", "")
					}
				}
			}
		}

		if params.StopOnFail == nil {
			sl.ReportError(params.StopOnFail, "StopOnFail", "Parameters.StopOnFail", "required", "")
		}

		if params.StopOnSuccess == nil {
			sl.ReportError(params.StopOnSuccess, "StopOnSuccess", "Parameters.StopOnSuccess", "required", "")
		}
	}
}
