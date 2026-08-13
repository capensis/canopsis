package builder

import (
	"maps"
	"slices"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	PatternTypeAlarm          = "alarm_pattern"
	PatternTypeEntity         = "entity_pattern"
	PatternTypePbehavior      = "pbehavior_pattern"
	PatternTypeEvent          = "event_pattern"
	PatternTypeWeatherService = "weather_service_pattern"

	DynamicFieldAlias      = "alias"
	DynamicFieldTicketData = "v.ticket.ticket_data.*"
)

// PatternTypes returns the supported pattern type names in stable order.
func PatternTypes() []string {
	return []string{
		PatternTypeAlarm,
		PatternTypeEntity,
		PatternTypePbehavior,
		PatternTypeEvent,
		PatternTypeWeatherService,
	}
}

func StaticFieldsOfType(fieldType string) []string {
	set := make(map[string]struct{})
	for _, byField := range staticFields {
		for field, typ := range byField {
			if typ == fieldType {
				set[field] = struct{}{}
			}
		}
	}

	return slices.Sorted(maps.Keys(set))
}

func ExistReferenceFields() []string {
	set := make(map[string]struct{})
	for _, byField := range staticFields {
		for field, typ := range byField {
			if typ == pattern.FieldTypeReference {
				set[field] = struct{}{}
			}
		}
	}
	for _, byField := range mixedFields {
		for field, byOp := range byField {
			if _, ok := byOp[pattern.ConditionExist]; ok {
				set[field] = struct{}{}
			}
		}
	}

	return slices.Sorted(maps.Keys(set))
}

func TimestampFields() []string {
	set := make(map[string]struct{})
	for _, byField := range staticFields {
		for field, typ := range byField {
			if typ == pattern.FieldTypeTimestamp {
				set[field] = struct{}{}
			}
		}
	}
	for _, byField := range mixedFields {
		for field, byOp := range byField {
			if _, ok := byOp[pattern.ConditionTimeRelative]; ok {
				set[field] = struct{}{}
			}
			if _, ok := byOp[pattern.ConditionTimeAbsolute]; ok {
				set[field] = struct{}{}
			}
		}
	}

	return slices.Sorted(maps.Keys(set))
}

func TypedDynamicFields() []string {
	set := map[string]struct{}{
		DynamicFieldAlias: {},
	}
	for _, byField := range dynamicFieldPrefixes {
		for field := range byField {
			if field == DynamicFieldTicketData {
				continue
			}
			set[field] = struct{}{}
		}
	}

	return slices.Sorted(maps.Keys(set))
}
