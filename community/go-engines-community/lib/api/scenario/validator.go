package scenario

import (
	"net/http"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateActionRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(ActionRequest)

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
