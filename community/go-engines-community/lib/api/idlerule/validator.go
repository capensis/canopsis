package idlerule

import (
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/eventfilter/pattern"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/idlerule"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/pbehavior"
	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
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
