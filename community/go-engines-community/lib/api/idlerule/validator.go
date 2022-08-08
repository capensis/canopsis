package idlerule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
	"github.com/teambition/rrule-go"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
)

func ValidateEditRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(EditRequest)

	if r.Priority != nil && *r.Priority < 0 {
		sl.ReportError(r.Priority, "Priority", "Priority", "min", "0")
	}

	validateType(sl, r.Type)
	validateAlarmRule(sl, r)
	validateEntityRule(sl, r)
	validateDisableDuringPeriods(sl, r.DisableDuringPeriods)
}

func ValidateCountPatternRequest(sl validator.StructLevel) {
	r := sl.Current().Interface().(CountByPatternRequest)

	validateAlarmPatterns(sl, r.AlarmPatterns)
	validateEntityPatterns(sl, r.EntityPatterns)
}

func validateType(sl validator.StructLevel, t string) {
	if t == "" {
		return
	}
	validTypes := []string{
		idlerule.RuleTypeAlarm,
		idlerule.RuleTypeEntity,
	}
	param := strings.Join(validTypes, " ")
	found := false
	for _, v := range validTypes {
		if v == t {
			found = true
			break
		}
	}

	if !found {
		sl.ReportError(t, "Type", "Type", "oneof", param)
	}
}

func validateAlarmRule(sl validator.StructLevel, r EditRequest) {
	if r.Type != idlerule.RuleTypeAlarm {
		return
	}

	entityPatternsIsSet := validateEntityPatterns(sl, r.EntityPatterns)
	alarmPatternsIsSet := validateAlarmPatterns(sl, r.AlarmPatterns)
	if !entityPatternsIsSet && !alarmPatternsIsSet {
		sl.ReportError(r.Type, "AlarmPatterns", "AlarmPatterns", "required_or", "EntityPatterns")
		sl.ReportError(r.Type, "EntityPatterns", "EntityPatterns", "required_or", "AlarmPatterns")
	}

	if r.AlarmCondition == "" {
		sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "required", "")
	} else {
		validValues := []string{
			idlerule.RuleAlarmConditionLastEvent,
			idlerule.RuleAlarmConditionLastUpdate,
		}
		param := strings.Join(validValues, " ")
		found := false
		for _, v := range validValues {
			if v == r.AlarmCondition {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "oneof", param)
		}
	}

	if r.Operation == nil {
		sl.ReportError(r.Operation, "Operation", "Operation", "required", "")
	} else {
		validOperationTypes := []string{
			types.ActionTypeAck,
			types.ActionTypeAckRemove,
			types.ActionTypeCancel,
			types.ActionTypeAssocTicket,
			types.ActionTypeChangeState,
			types.ActionTypeSnooze,
			types.ActionTypePbehavior,
		}
		param := strings.Join(validOperationTypes, " ")
		found := false
		for _, v := range validOperationTypes {
			if v == r.Operation.Type {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(r.Operation.Type, "Operation.Type", "Type", "oneof", param)
		}

		validateOperationParametersRequest(sl, r.Operation.Type, r.Operation.Parameters)
	}
}

func validateEntityRule(sl validator.StructLevel, r EditRequest) {
	if r.Type != idlerule.RuleTypeEntity {
		return
	}

	entityPatternsIsSet := validateEntityPatterns(sl, r.EntityPatterns)
	if !entityPatternsIsSet {
		sl.ReportError(r.Type, "EntityPatterns", "EntityPatterns", "required", "")
	}

	alarmPatternsIsSet := validateAlarmPatterns(sl, r.AlarmPatterns)
	if alarmPatternsIsSet {
		sl.ReportError(r.Type, "AlarmPatterns", "AlarmPatterns", "must_be_empty", "")
	}

	if r.Operation != nil {
		sl.ReportError(r.Operation, "Operation", "Operation", "must_be_empty", "")
	}

	if r.AlarmCondition != "" {
		sl.ReportError(r.AlarmCondition, "AlarmCondition", "AlarmCondition", "must_be_empty", "")
	}
}

func validateEntityPatterns(sl validator.StructLevel, patterns pattern.EntityPatternList) bool {
	patternsIsSet := false
	if patterns.IsSet() {
		if !patterns.IsValid() {
			sl.ReportError(patterns, "EntityPatterns", "EntityPatterns", "entitypattern_invalid", "")
		} else {
			query := patterns.AsMongoDriverQuery()["$or"].([]bson.M)
			if len(query) > 0 {
				patternsIsSet = true
				for _, q := range query {
					if len(q) == 0 {
						sl.ReportError(patterns, "EntityPatterns", "EntityPatterns", "entitypattern_contains_empty", "")
						break
					}
				}
			}
		}
	}

	return patternsIsSet
}

func validateAlarmPatterns(sl validator.StructLevel, patterns pattern.AlarmPatternList) bool {
	patternsIsSet := false
	if patterns.IsSet() {
		if !patterns.IsValid() {
			sl.ReportError(patterns, "AlarmPatterns", "AlarmPatterns", "alarmpattern_invalid", "")
		} else {
			query := patterns.AsMongoDriverQuery()["$or"].([]bson.M)
			if len(query) > 0 {
				patternsIsSet = true
				for _, q := range query {
					if len(q) == 0 {
						sl.ReportError(patterns, "AlarmPatterns", "AlarmPatterns", "alarmpattern_contains_empty", "")
						break
					}
				}
			}
		}
	}

	return patternsIsSet
}

func validateDisableDuringPeriods(sl validator.StructLevel, disableDuringPeriods []string) {
	if len(disableDuringPeriods) == 0 {
		return
	}
	validPeriods := []string{
		pbehavior.TypeMaintenance,
		pbehavior.TypePause,
		pbehavior.TypeInactive,
	}
	param := strings.Join(validPeriods, " ")
	for _, period := range disableDuringPeriods {
		found := false
		for _, v := range validPeriods {
			if v == period {
				found = true
				break
			}
		}

		if !found {
			sl.ReportError(disableDuringPeriods, "DisableDuringPeriods", "disableDuringPeriods", "oneof", param)
		}
	}
}

func validateOperationParametersRequest(sl validator.StructLevel, t string, params idlerule.Parameters) {
	switch t {
	case types.ActionTypeAssocTicket:
		if params.Ticket == "" {
			sl.ReportError(params.Ticket, "Operation.Parameters.Ticket", "Ticket", "required", "")
		}
	case types.ActionTypeChangeState:
		if params.State == nil {
			sl.ReportError(params.State, "Operation.Parameters.State", "State", "required", "")
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
				sl.ReportError(params.State, "Operation.Parameters.State", "State", "oneof", param)
			}
		}
	case types.ActionTypeSnooze:
		if params.Duration == nil {
			sl.ReportError(params.Duration, "Operation.Parameters.Duration", "Duration", "required", "")
		}
	case types.ActionTypePbehavior:
		if params.Name == "" {
			sl.ReportError(params.Name, "Operation.Parameters.Name", "Name", "required", "")
		}
		if params.Reason == "" {
			sl.ReportError(params.Reason, "Operation.Parameters.Reason", "Reason", "required", "")
		}
		if params.Type == "" {
			sl.ReportError(params.Type, "Operation.Parameters.Type", "Type", "required", "")
		}
		// Validate rrule
		if params.RRule != "" {
			_, err := rrule.StrToROption(params.RRule)
			if err != nil {
				sl.ReportError(params.RRule, "Operation.Parameters.RRule", "RRule", "rrule", "")
			}
		}

		// Validate time
		if params.Tstart == nil && params.Tstop != nil {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_with", "Tstop")
		}
		if params.Tstart != nil && params.Tstop == nil {
			sl.ReportError(params.Tstop, "Operation.Parameters.Tstop", "Tstop", "required_with", "Tstart")
		}
		if params.Duration == nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Duration, "Operation.Parameters.Duration", "Duration", "required_with", "StartOnTrigger")
		}
		if params.Duration != nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_with", "Duration")
		}
		if params.Tstart == nil && params.Tstop == nil && params.Duration == nil && (params.StartOnTrigger == nil || !*params.StartOnTrigger) {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.StartOnTrigger != nil && *params.StartOnTrigger {
			sl.ReportError(params.Tstart, "Operation.Parameters.Tstart", "Tstart", "required_or", "StartOnTrigger")
			sl.ReportError(params.StartOnTrigger, "Operation.Parameters.StartOnTrigger", "StartOnTrigger", "required_or", "Tstart")
		}
		if params.Tstart != nil && params.Tstop != nil && params.Tstop.Before(*params.Tstart) {
			sl.ReportError(params.Tstop, "Operation.Parameters.Tstop", "Tstop", "gtfield", "Tstart")
		}
	}
}
