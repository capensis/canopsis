package scenario

import (
	"net/http"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/action"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateActionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ActionRequest)

	if r.Type != "" {
		validateActionParametersRequest(sl, r.Type, r.Parameters)
	}

	// Validate patterns
	alarmPatternsIsSet := false
	if r.AlarmPatterns.IsSet() {
		if !r.AlarmPatterns.IsValid() {
			alarmPatternsIsSet = true
			sl.ReportError(r.AlarmPatterns, "AlarmPatterns", "AlarmPatterns", "alarmpattern_invalid", "")
		} else {
			query := r.AlarmPatterns.AsMongoDriverQuery()["$or"].([]bson.M)
			if len(query) > 0 {
				alarmPatternsIsSet = true
				for _, q := range query {
					if len(q) == 0 {
						sl.ReportError(r.AlarmPatterns, "AlarmPatterns", "AlarmPatterns", "alarmpattern_contains_empty", "")
						break
					}
				}
			}
		}
	}
	entityPatternsIsSet := false
	if r.EntityPatterns.IsSet() {
		if !r.EntityPatterns.IsValid() {
			entityPatternsIsSet = true
			sl.ReportError(r.EntityPatterns, "EntityPatterns", "EntityPatterns", "entitypattern_invalid", "")
		} else {
			query := r.EntityPatterns.AsMongoDriverQuery()["$or"].([]bson.M)
			if len(query) > 0 {
				entityPatternsIsSet = true
				for _, q := range query {
					if len(q) == 0 {
						sl.ReportError(r.EntityPatterns, "EntityPatterns", "EntityPatterns", "entitypattern_contains_empty", "")
						break
					}
				}
			}
		}
	}

	if !entityPatternsIsSet && !alarmPatternsIsSet {
		sl.ReportError(r.Type, "AlarmPatterns", "AlarmPatterns", "required", "")
		sl.ReportError(r.Type, "EntityPatterns", "EntityPatterns", "required", "")
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
		} else if params.Request.Method != "" {
			validMethods := []string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
				http.MethodConnect,
				http.MethodOptions,
				http.MethodTrace,
			}
			param := strings.Join(validMethods, " ")

			found := false
			for _, v := range validMethods {
				if v == strings.ToUpper(params.Request.Method) {
					found = true
					break
				}
			}

			if !found {
				sl.ReportError(params.Request.Method, "Parameters.Request.Method", "Method", "oneof", param)
			}
		}
	}
}
