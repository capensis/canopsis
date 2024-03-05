package scenario

import (
	"strconv"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/common"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern/match"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
)

func ValidateActionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ActionRequest)

	if r.Type != "" {
		validateActionParametersRequest(sl, r.Type, r.Parameters)
	}

	if r.CorporateEntityPattern == "" && len(r.EntityPattern) > 0 &&
		!match.ValidateEntityPattern(r.EntityPattern, common.GetForbiddenFieldsInEntityPattern(mongo.ScenarioMongoCollection)) {
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "entity_pattern", "")
	}

	if r.CorporateAlarmPattern == "" && len(r.AlarmPattern) > 0 &&
		!match.ValidateAlarmPattern(r.AlarmPattern,
			common.GetForbiddenFieldsInAlarmPattern(mongo.ScenarioMongoCollection),
			common.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(mongo.ScenarioMongoCollection),
		) {
		sl.ReportError(r.EntityPattern, "AlarmPattern", "AlarmPattern", "alarm_pattern", "")
	}

	if len(r.EntityPattern) == 0 && r.CorporateEntityPattern == "" &&
		len(r.AlarmPattern) == 0 && r.CorporateAlarmPattern == "" {
		sl.ReportError(r.AlarmPattern, "AlarmPattern", "AlarmPattern", "required_or", "EntityPattern")
		sl.ReportError(r.EntityPattern, "EntityPattern", "EntityPattern", "required_or", "AlarmPattern")
	}
}

func validateActionParametersRequest(sl validator.StructLevel, t string, params action.Parameters) {
	switch t {
	case types.ActionTypeAssocTicket:
		if params.Ticket == "" {
			sl.ReportError(params.Ticket, "Parameters.Ticket", "Ticket", "required", "")
		}
	case types.ActionTypeChangeState:
		if params.State == nil {
			sl.ReportError(params.State, "Parameters.State", "State", "required", "")
		} else {
			validTypes := []types.CpsNumber{
				types.AlarmStateOK,
				types.AlarmStateMinor,
				types.AlarmStateMajor,
				types.AlarmStateCritical,
			}
			param := ""
			for i := range validTypes {
				param += strconv.Itoa(int(validTypes[i]))
				if i < len(validTypes)-1 {
					param += " "
				}
			}

			found := false
			for _, v := range validTypes {
				if v == *params.State {
					found = true
					break
				}
			}

			if !found {
				sl.ReportError(params.State, "Parameters.State", "State", "oneof", param)
			}
		}
	case types.ActionTypeSnooze:
		if params.Duration == nil {
			sl.ReportError(params.Duration, "Parameters.Duration", "Duration", "required", "")
		}
	case types.ActionTypePbehavior:
		if params.Name == "" {
			sl.ReportError(params.Name, "Parameters.Name", "Name", "required", "")
		}
		if params.Reason == "" {
			sl.ReportError(params.Reason, "Parameters.Reason", "Reason", "required", "")
		}
		if params.Type == "" {
			sl.ReportError(params.Type, "Parameters.Type", "Type", "required", "")
		}
		// Validate rrule
		if params.RRule != "" {
			_, err := rrule.StrToROption(params.RRule)
			if err != nil {
				sl.ReportError(params.RRule, "Parameters.RRule", "RRule", "rrule", "")
			}
		}

		// Validate time
		if params.Tstart == nil && params.Tstop != nil {
			sl.ReportError(params.Tstart, "Parameters.Tstart", "Tstart", "required_with", "Tstop")
		}
		if params.Tstart != nil && params.Tstop == nil {
			sl.ReportError(params.Tstop, "Parameters.Tstop", "Tstop", "required_with", "Tstart")
		}
		if params.Duration == nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Duration, "Parameters.Duration", "Duration", "required_with", "StartOnTrigger")
		}
		if params.Duration != nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.StartOnTrigger, "Parameters.StartOnTrigger", "StartOnTrigger", "required_with", "Duration")
		}
		if params.Tstart == nil && params.Tstop == nil && params.Duration == nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.Tstart, "Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Tstart, "Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.Tstop != nil && params.Tstop.Before(*params.Tstart) {
			sl.ReportError(params.Tstop, "Parameters.Tstop", "Tstop", "gtfield", "Tstart")
		}
	case types.ActionTypeWebhook:
		if params.Request == nil {
			sl.ReportError(params.Request, "Parameters.Request", "Request", "required", "")
		}
	}
}
