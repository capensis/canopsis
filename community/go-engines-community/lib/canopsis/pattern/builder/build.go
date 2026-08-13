package builder

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/api/patternfields"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/datetime"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

func Build(input Input, promptContext string) (Result, error) {
	collection, ok := CollectionByPromptContext(promptContext)
	if !ok {
		return Result{}, fmt.Errorf("unknown prompt context %q", promptContext)
	}

	forbiddenAlarmFields := patternfields.GetForbiddenFieldsInAlarmPattern(collection)
	onlyAbsTimeAlarmFields := patternfields.GetOnlyAbsoluteTimeCondFieldsInAlarmPattern(collection)
	forbiddenEntityFields := patternfields.GetForbiddenFieldsInEntityPattern(collection)
	if promptContext == promptContextStateSettingsServiceInherited {
		forbiddenEntityFields = patternfields.GetForbiddenFieldsInInheritedEntityPattern(
			mongo.StateSettingsMongoCollection,
			statesetting.RuleTypeService,
		)
	}

	var result Result

	for i, p := range input.Patterns {
		switch p.Type {
		case PatternTypeAlarm, PatternTypeEntity, PatternTypePbehavior, PatternTypeEvent, PatternTypeWeatherService:
		default:
			return Result{}, fmt.Errorf("pattern %d: unknown pattern type %q", i, p.Type)
		}

		groups := make([][]pattern.FieldCondition, 0, len(p.Groups))
		for j, g := range p.Groups {
			if len(g) == 0 {
				return Result{}, fmt.Errorf("pattern %d group %d: %w", i, j, pattern.ErrEmptyGroup)
			}

			conditions := make([]pattern.FieldCondition, 0, len(g))
			for k, c := range g {
				fc, err := buildFieldCondition(p.Type, c)
				if err != nil {
					return Result{}, fmt.Errorf("pattern %d group %d condition %d: %w", i, j, k, err)
				}

				switch p.Type {
				case PatternTypeAlarm:
					if pattern.IsForbiddenAlarmField(fc, forbiddenAlarmFields, onlyAbsTimeAlarmFields) {
						return Result{}, fmt.Errorf("pattern %d group %d condition %d: field %q is forbidden in %q for context %q: %w",
							i, j, k, fc.Field, p.Type, promptContext, pattern.ErrForbiddenField)
					}
				case PatternTypeEntity:
					if pattern.IsForbiddenEntityField(fc, forbiddenEntityFields) {
						field := fc.Field
						if field == "" && fc.Alias != "" {
							field = DynamicFieldAlias
						}

						return Result{}, fmt.Errorf("pattern %d group %d condition %d: field %q is forbidden in %q for context %q: %w",
							i, j, k, field, p.Type, promptContext, pattern.ErrForbiddenField)
					}
				}

				conditions = append(conditions, fc)
			}

			groups = append(groups, conditions)
		}

		switch p.Type {
		case PatternTypeAlarm:
			result.AlarmPattern = append(result.AlarmPattern, groups...)
		case PatternTypeEntity:
			result.EntityPattern = append(result.EntityPattern, groups...)
		case PatternTypePbehavior:
			result.PbehaviorPattern = append(result.PbehaviorPattern, groups...)
		case PatternTypeEvent:
			result.EventPattern = append(result.EventPattern, groups...)
		case PatternTypeWeatherService:
			result.WeatherServicePattern = append(result.WeatherServicePattern, groups...)
		}
	}

	return result, nil
}

func buildFieldCondition(patternType string, c ConditionSpec) (pattern.FieldCondition, error) {
	typ, field, alias, fieldType, err := resolveField(patternType, c)
	if err != nil {
		return pattern.FieldCondition{}, err
	}

	cond, err := buildCondition(typ, c.Operator, c.Value)
	if err != nil {
		return pattern.FieldCondition{}, fmt.Errorf("field %q: %w", c.Field, err)
	}

	if err := validateStringEnumValues(c.Field, c.Operator, c.Value); err != nil {
		return pattern.FieldCondition{}, fmt.Errorf("field %q: %w", c.Field, err)
	}

	if err := validateIntEnumValues(c.Field, c.Operator, c.Value); err != nil {
		return pattern.FieldCondition{}, fmt.Errorf("field %q: %w", c.Field, err)
	}

	return pattern.FieldCondition{
		Field:     field,
		FieldType: fieldType,
		Condition: cond,
		Alias:     alias,
	}, nil
}

func resolveField(patternType string, c ConditionSpec) (typ, field, alias, fieldType string, err error) {
	if patternType == PatternTypeEntity && c.Field == DynamicFieldAlias {
		if c.Key == "" {
			return "", "", "", "", errors.New(`field "alias" requires key`)
		}

		typ, fieldType, err = resolveDynamicFieldType(c)
		if err != nil {
			return "", "", "", "", err
		}

		return typ, "", c.Key, fieldType, nil
	}

	if prefix, ok := dynamicFieldPrefixes[patternType][c.Field]; ok {
		if c.Key == "" {
			return "", "", "", "", fmt.Errorf("field %q requires key", c.Field)
		}

		field = prefix + c.Key

		if c.Field == DynamicFieldTicketData {
			return pattern.FieldTypeString, field, "", "", nil
		}

		typ, fieldType, err = resolveDynamicFieldType(c)
		if err != nil {
			return "", "", "", "", err
		}

		return typ, field, "", fieldType, nil
	}

	if ops, ok := mixedFields[patternType][c.Field]; ok {
		fieldTyp, ok := ops[c.Operator]
		if !ok {
			return "", "", "", "", fmt.Errorf("field %q: %w: %s", c.Field, pattern.ErrUnsupportedConditionType, c.Operator)
		}

		return fieldTyp, c.Field, "", "", nil
	}

	typ, ok := staticFields[patternType][c.Field]
	if !ok {
		return "", "", "", "", fmt.Errorf("%w: %s", pattern.ErrUnsupportedField, c.Field)
	}

	return typ, c.Field, "", "", nil
}

func resolveDynamicFieldType(c ConditionSpec) (string, string, error) {
	if c.FieldType == "" {
		if c.Operator != pattern.ConditionExist {
			return "", "", fmt.Errorf(`field %q: field_type is required for operator %q (only %q can omit it)`, c.Field, c.Operator, pattern.ConditionExist)
		}

		return pattern.FieldTypeReference, "", nil
	}

	typ, err := parseFieldType(c.FieldType)
	if err != nil {
		return "", "", err
	}

	return typ, c.FieldType, nil
}

func parseFieldType(fieldType string) (string, error) {
	switch fieldType {
	case pattern.FieldTypeString,
		pattern.FieldTypeInt,
		pattern.FieldTypeBool,
		pattern.FieldTypeStringArray,
		pattern.FieldTypeTimestamp:
		return fieldType, nil
	default:
		return "", fmt.Errorf("%w: %s", pattern.ErrUnsupportedFieldType, fieldType)
	}
}

func buildCondition(typ string, operator string, v ValueSpec) (pattern.Condition, error) {
	switch typ {
	case pattern.FieldTypeString:
		return buildStringCondition(operator, v)
	case pattern.FieldTypeInt:
		return buildIntCondition(operator, v)
	case pattern.FieldTypeBool:
		return buildBoolCondition(operator, v)
	case pattern.FieldTypeReference:
		return buildRefCondition(operator, v)
	case pattern.FieldTypeTimestamp:
		return buildTimeCondition(operator, v)
	case pattern.FieldTypeDuration:
		return buildDurationCondition(operator, v)
	case pattern.FieldTypeTags:
		return buildTagsCondition(operator, v)
	case pattern.FieldTypeStringArray:
		return buildStringArrayCondition(operator, v)
	default:
		return pattern.Condition{}, fmt.Errorf("unsupported field type: %s", typ)
	}
}

func buildStringCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	switch operator {
	case pattern.ConditionEqual, pattern.ConditionNotEqual,
		pattern.ConditionContain, pattern.ConditionNotContain,
		pattern.ConditionBeginWith, pattern.ConditionNotBeginWith,
		pattern.ConditionEndWith, pattern.ConditionNotEndWith:
		if v.String == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a string value", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewStringCondition(operator, *v.String), nil
	case pattern.ConditionRegexp:
		if v.String == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a string value", pattern.ErrWrongConditionValue, operator)
		}

		if *v.String == "" {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a non-empty regexp string", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewRegexpCondition(operator, *v.String)
	case pattern.ConditionIsOneOf, pattern.ConditionIsNotOneOf:
		if err := validateNonEmptyStrings(v.Strings, operator); err != nil {
			return pattern.Condition{}, err
		}

		return pattern.NewStringArrayCondition(operator, v.Strings), nil
	case pattern.ConditionExist:
		if v.Boolean == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a boolean value", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewBoolCondition(operator, *v.Boolean), nil
	default:
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}
}

func buildIntCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	switch operator {
	case pattern.ConditionEqual, pattern.ConditionNotEqual, pattern.ConditionGT, pattern.ConditionLT:
		if v.Integer == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires an integer value", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewIntCondition(operator, *v.Integer), nil
	default:
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}
}

func buildBoolCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	if operator != pattern.ConditionEqual {
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}

	if v.Boolean == nil {
		return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a boolean value", pattern.ErrWrongConditionValue, operator)
	}

	return pattern.NewBoolCondition(operator, *v.Boolean), nil
}

func buildRefCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	if operator != pattern.ConditionExist {
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}

	if v.Boolean == nil {
		return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a boolean value", pattern.ErrWrongConditionValue, operator)
	}

	return pattern.NewBoolCondition(operator, *v.Boolean), nil
}

func buildTagsCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	switch operator {
	case pattern.ConditionHasEvery, pattern.ConditionHasOneOf, pattern.ConditionHasNot,
		pattern.ConditionHasLabels, pattern.ConditionHasNotLabels:
		if err := validateNonEmptyStrings(v.Strings, operator); err != nil {
			return pattern.Condition{}, err
		}

		return pattern.NewStringArrayCondition(operator, v.Strings), nil
	case pattern.ConditionIsEmpty:
		if v.Boolean == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a boolean value", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewBoolCondition(operator, *v.Boolean), nil
	default:
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}
}

func buildStringArrayCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	switch operator {
	case pattern.ConditionHasEvery, pattern.ConditionHasOneOf, pattern.ConditionHasNot:
		if err := validateNonEmptyStrings(v.Strings, operator); err != nil {
			return pattern.Condition{}, err
		}

		return pattern.NewStringArrayCondition(operator, v.Strings), nil
	case pattern.ConditionIsEmpty:
		if v.Boolean == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a boolean value", pattern.ErrWrongConditionValue, operator)
		}

		return pattern.NewBoolCondition(operator, *v.Boolean), nil
	default:
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}
}

func buildDurationCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	if operator != pattern.ConditionGT && operator != pattern.ConditionLT {
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}

	if v.Duration == nil {
		return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a duration value", pattern.ErrWrongConditionValue, operator)
	}

	if err := validateDurationSpec(*v.Duration, "duration"); err != nil {
		return pattern.Condition{}, err
	}

	cond, err := pattern.NewDurationCondition(operator, v.Duration.toDurationWithUnit())
	if err != nil {
		return pattern.Condition{}, fmt.Errorf("%w: invalid duration for operator %q: %w", pattern.ErrWrongConditionValue, operator, err)
	}

	return cond, nil
}

func buildTimeCondition(operator string, v ValueSpec) (pattern.Condition, error) {
	switch operator {
	case pattern.ConditionTimeRelative:
		if v.RelativeTime == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires a relative_time value", pattern.ErrWrongConditionValue, operator)
		}

		return buildRelativeTimeCondition(*v.RelativeTime)
	case pattern.ConditionTimeAbsolute:
		if v.AbsoluteTime == nil {
			return pattern.Condition{}, fmt.Errorf("%w: operator %q requires an absolute_time value", pattern.ErrWrongConditionValue, operator)
		}

		if err := validateAbsoluteTime(*v.AbsoluteTime); err != nil {
			return pattern.Condition{}, err
		}

		return pattern.NewTimeIntervalCondition(operator, v.AbsoluteTime.From, v.AbsoluteTime.To), nil
	default:
		return pattern.Condition{}, fmt.Errorf("%w: %s", pattern.ErrUnsupportedConditionType, operator)
	}
}

func buildRelativeTimeCondition(r RelativeTimeSpec) (pattern.Condition, error) {
	switch r.Kind {
	case relativeTimeKindNewerThan:
		if r.From != nil || r.To != nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "newer_than" must use duration only (omit from/to); use kind "between" for a from/to window`, pattern.ErrWrongConditionValue)
		}
		if r.Duration == nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "newer_than" requires duration`, pattern.ErrWrongConditionValue)
		}

		if err := validateDurationSpec(*r.Duration, `relative_time "newer_than" duration`); err != nil {
			return pattern.Condition{}, err
		}

		cond, err := pattern.NewDurationCondition(pattern.ConditionTimeRelative, r.Duration.toDurationWithUnit())
		if err != nil {
			return pattern.Condition{}, fmt.Errorf("%w: invalid relative_time newer_than: %w", pattern.ErrWrongConditionValue, err)
		}

		return cond, nil
	case relativeTimeKindOlderThan:
		if r.From != nil || r.To != nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "older_than" must use duration only (omit from/to); use kind "between" for a from/to window`, pattern.ErrWrongConditionValue)
		}
		if r.Duration == nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "older_than" requires duration`, pattern.ErrWrongConditionValue)
		}

		if err := validateDurationSpec(*r.Duration, `relative_time "older_than" duration`); err != nil {
			return pattern.Condition{}, err
		}

		cond, err := pattern.NewDurationCondition(
			pattern.ConditionTimeRelative,
			datetime.DurationWithUnit{},
			r.Duration.toDurationWithUnit(),
		)
		if err != nil {
			return pattern.Condition{}, fmt.Errorf("%w: invalid relative_time older_than: %w", pattern.ErrWrongConditionValue, err)
		}

		return cond, nil
	case relativeTimeKindBetween:
		if r.Duration != nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "between" must use from and to only (omit duration); use newer_than/older_than for a single threshold`, pattern.ErrWrongConditionValue)
		}
		if r.From == nil || r.To == nil {
			return pattern.Condition{}, fmt.Errorf(`%w: relative_time kind "between" requires from and to`, pattern.ErrWrongConditionValue)
		}

		if err := validateDurationSpec(*r.From, `relative_time "between" from`); err != nil {
			return pattern.Condition{}, err
		}

		if err := validateDurationSpec(*r.To, `relative_time "between" to`); err != nil {
			return pattern.Condition{}, err
		}

		if err := validateRelativeBetween(*r.From, *r.To); err != nil {
			return pattern.Condition{}, err
		}

		return pattern.NewDurationCondition(
			pattern.ConditionTimeRelative,
			r.From.toDurationWithUnit(),
			r.To.toDurationWithUnit(),
		)
	default:
		return pattern.Condition{}, fmt.Errorf("%w: unsupported relative_time kind %q; use newer_than, older_than, or between", pattern.ErrWrongConditionValue, r.Kind)
	}
}

func validateAbsoluteTime(a AbsoluteTimeSpec) error {
	if a.From <= 0 || a.To <= 0 {
		return fmt.Errorf("%w: absolute_time from and to must be positive Unix timestamps in seconds (not milliseconds); got from=%d to=%d", pattern.ErrWrongConditionValue, a.From, a.To)
	}
	const maxUnixSeconds = 100000000000
	if a.From >= maxUnixSeconds || a.To >= maxUnixSeconds {
		return fmt.Errorf("%w: absolute_time from and to must be Unix timestamps in seconds, not milliseconds; got from=%d to=%d (example seconds: 1710000000)", pattern.ErrWrongConditionValue, a.From, a.To)
	}
	if a.To <= a.From {
		return fmt.Errorf("%w: absolute_time to (%d) must be greater than from (%d); they cannot be equal or reversed", pattern.ErrWrongConditionValue, a.To, a.From)
	}

	return nil
}

func validateDurationSpec(d DurationSpec, what string) error {
	if d.Value <= 0 {
		return fmt.Errorf("%w: %s value must be a positive integer, got %d", pattern.ErrWrongConditionValue, what, d.Value)
	}

	switch d.Unit {
	case datetime.DurationUnitSecond,
		datetime.DurationUnitMinute,
		datetime.DurationUnitHour,
		datetime.DurationUnitDay,
		datetime.DurationUnitWeek:
		return nil
	default:
		return fmt.Errorf("%w: %s has unknown unit %q; use one of: s, m, h, d, w", pattern.ErrWrongConditionValue, what, d.Unit)
	}
}

func validateRelativeBetween(from, to DurationSpec) error {
	fromSec, err := from.toDurationWithUnit().To(datetime.DurationUnitSecond)
	if err != nil {
		return fmt.Errorf("%w: relative_time \"between\" from: %w", pattern.ErrWrongConditionValue, err)
	}

	toSec, err := to.toDurationWithUnit().To(datetime.DurationUnitSecond)
	if err != nil {
		return fmt.Errorf("%w: relative_time \"between\" to: %w", pattern.ErrWrongConditionValue, err)
	}

	if fromSec.Value == toSec.Value {
		return fmt.Errorf("%w: relative_time \"between\" from and to must not be equal", pattern.ErrWrongConditionValue)
	}

	if fromSec.Value < toSec.Value {
		return fmt.Errorf("%w: relative_time \"between\" requires from > to, got from=%d%s to=%d%s", pattern.ErrWrongConditionValue, from.Value, from.Unit, to.Value, to.Unit)
	}

	return nil
}

func validateNonEmptyStrings(values []string, operator string) error {
	if len(values) == 0 {
		return fmt.Errorf("%w: operator %q requires a non-empty string list value", pattern.ErrWrongConditionValue, operator)
	}

	for i, s := range values {
		if s == "" {
			return fmt.Errorf("%w: operator %q string list must not contain empty strings (index %d)", pattern.ErrWrongConditionValue, operator, i)
		}
	}

	return nil
}

func validateStringEnumValues(field, operator string, v ValueSpec) error {
	allowed, ok := stringEnumValues[field]
	if !ok {
		return nil
	}

	switch operator {
	case pattern.ConditionEqual, pattern.ConditionNotEqual:
		if v.String == nil {
			return nil
		}

		if !slices.Contains(allowed, *v.String) {
			return fmt.Errorf("%w: %q is not a valid value for %q; allowed: %s",
				pattern.ErrWrongConditionValue, *v.String, field, strings.Join(allowed, ", "))
		}
	case pattern.ConditionIsOneOf, pattern.ConditionIsNotOneOf:
		for _, s := range v.Strings {
			if !slices.Contains(allowed, s) {
				return fmt.Errorf("%w: %q is not a valid value for %q; allowed: %s",
					pattern.ErrWrongConditionValue, s, field, strings.Join(allowed, ", "))
			}
		}
	case pattern.ConditionExist:
		return nil
	default:
		return fmt.Errorf("%w: field %q only supports eq, neq, is_one_of, is_not_one_of, exist (got %q)",
			pattern.ErrUnsupportedConditionType, field, operator)
	}

	return nil
}

func validateIntEnumValues(field, operator string, v ValueSpec) error {
	if v.Integer == nil {
		return nil
	}

	allowed, ok := intEnumValues[field]
	if !ok {
		return nil
	}

	switch operator {
	case pattern.ConditionEqual, pattern.ConditionNotEqual:
		if !slices.Contains(allowed, *v.Integer) {
			return fmt.Errorf("%w: %d is not a valid value for %q; allowed: %v",
				pattern.ErrWrongConditionValue, *v.Integer, field, allowed)
		}
	case pattern.ConditionGT:
		maxVal := slices.Max(allowed)
		if *v.Integer >= maxVal {
			return fmt.Errorf("%w: gt %d matches no values of %q (allowed: %v); threshold must be < %d",
				pattern.ErrWrongConditionValue, *v.Integer, field, allowed, maxVal)
		}
	case pattern.ConditionLT:
		minVal := slices.Min(allowed)
		if *v.Integer <= minVal {
			return fmt.Errorf("%w: lt %d matches no values of %q (allowed: %v); threshold must be > %d",
				pattern.ErrWrongConditionValue, *v.Integer, field, allowed, minVal)
		}
	}

	return nil
}
