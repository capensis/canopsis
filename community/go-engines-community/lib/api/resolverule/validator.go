package resolverule

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/eventfilter/oldpattern"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

func ValidateEditRequest(sl validator.StructLevel) {
	var r = sl.Current().Interface().(EditRequest)
	validateEntityPatterns(sl, r.EntityPatterns)
	validateAlarmPatterns(sl, r.AlarmPatterns)
}

func validateEntityPatterns(sl validator.StructLevel, patterns oldpattern.EntityPatternList) bool {
	patternsIsSet := false
	if patterns.IsSet() {
		if !patterns.IsValid() {
			patternsIsSet = true
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

func validateAlarmPatterns(sl validator.StructLevel, patterns oldpattern.AlarmPatternList) bool {
	patternsIsSet := false
	if patterns.IsSet() {
		if !patterns.IsValid() {
			patternsIsSet = true
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
