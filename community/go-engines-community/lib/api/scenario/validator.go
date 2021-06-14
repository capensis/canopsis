package scenario

import (
	"net/http"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	// Validate triggers
	if len(r.Triggers) > 0 {
		validTriggers := []string{
			string(types.AlarmChangeTypeCreate),
			string(types.AlarmChangeTypeStateIncrease),
			string(types.AlarmChangeTypeStateDecrease),
			string(types.AlarmChangeTypeChangeState),
			string(types.AlarmChangeTypeUpdateStatus),
			string(types.AlarmChangeTypeAck),
			string(types.AlarmChangeTypeAckremove),
			string(types.AlarmChangeTypeCancel),
			string(types.AlarmChangeTypeUncancel),
			string(types.AlarmChangeTypeComment),
			string(types.AlarmChangeTypeDone),
			string(types.AlarmChangeTypeDeclareTicket),
			string(types.AlarmChangeTypeDeclareTicketWebhook),
			string(types.AlarmChangeTypeAssocTicket),
			string(types.AlarmChangeTypeSnooze),
			string(types.AlarmChangeTypeUnsnooze),
			string(types.AlarmChangeTypeResolve),
			string(types.AlarmChangeTypeActivate),
			string(types.AlarmChangeTypePbhEnter),
			string(types.AlarmChangeTypePbhLeave),
		}
		param := strings.Join(validTriggers, " ")
		for _, trigger := range r.Triggers {
			found := false
			for _, v := range validTriggers {
				if v == trigger {
					found = true
					break
				}
			}

			if !found {
				sl.ReportError(r.Triggers, "Triggers", "triggers", "oneof", param)
			}
		}
	}

	// Validate disableDuringPeriods
	if len(r.DisableDuringPeriods) > 0 {
		validPeriods := []string{
			pbehavior.TypeMaintenance,
			pbehavior.TypePause,
			pbehavior.TypeInactive,
		}
		param := strings.Join(validPeriods, " ")
		for _, period := range r.DisableDuringPeriods {
			found := false
			for _, v := range validPeriods {
				if v == period {
					found = true
					break
				}
			}

			if !found {
				sl.ReportError(r.DisableDuringPeriods, "DisableDuringPeriods", "disableDuringPeriods", "oneof", param)
			}
		}
	}
}

func ValidateActionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ActionRequest)

	// Validate type
	if r.Type != "" {
		validTypes := []string{
			types.ActionTypeAck,
			types.ActionTypeAckRemove,
			types.ActionTypeAssocTicket,
			types.ActionTypeCancel,
			types.ActionTypeChangeState,
			types.ActionTypePbehavior,
			types.ActionTypeSnooze,
			types.ActionTypeWebhook,
		}
		param := strings.Join(validTypes, " ")
		found := false
		for _, v := range validTypes {
			if v == r.Type {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.Type, "Type", "Type", "oneof", param)
		}
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

func ValidateChangeStateParametersRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ChangeStateParametersRequest)

	// Validate state
	if r.State != nil {
		validTypes := []types.CpsNumber{
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
			if v == *r.State {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.State, "State", "State", "oneof", param)
		}
	}
}

func ValidatePbehaviorParametersRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(PbehaviorParametersRequest)

	// Validate rrule
	if r.RRule != "" {
		_, err := rrule.StrToROption(r.RRule)
		if err != nil {
			sl.ReportError(r.RRule, "RRule", "RRule", "rrule", "")
		}
	}

	// Validate time
	if r.Tstart == nil && r.Tstop == nil && r.StartOnTrigger == nil && r.Duration == nil {
		sl.ReportError(r.Tstart, "Tstart", "Tstart", "required_or", "StartOnTrigger")
		sl.ReportError(r.StartOnTrigger, "StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
	}
	if (r.Tstart != nil || r.Tstop != nil) && (r.StartOnTrigger != nil || r.Duration != nil) {
		sl.ReportError(r.Tstart, "Tstart", "Tstart", "required_or", "StartOnTrigger")
		sl.ReportError(r.StartOnTrigger, "StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
	}

	if r.Tstart != nil && r.Tstop != nil && *r.Tstop < *r.Tstart {
		sl.ReportError(r.Tstop, "Tstop", "Tstop", "gtfield", "Tstart")
	}
}

func ValidateWebhookRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(WebhookRequest)

	if r.Method != "" {
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
			if v == strings.ToUpper(r.Method) {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.Method, "Method", "Method", "oneof", param)
		}
	}
}
