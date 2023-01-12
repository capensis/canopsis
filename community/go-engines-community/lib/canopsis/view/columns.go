package view

import (
	"sort"
	"strings"
)

var columnsByWidget = map[string][]string{
	WidgetTemplateTypeAlarm: {
		"_id",
		"v.display_name",
		"v.output",
		"v.long_output",
		"v.initial_output",
		"v.initial_long_output",
		"v.connector",
		"v.connector_name",
		"v.component",
		"v.resource",
		"v.last_comment.m",
		"v.ack.a",
		"v.ack.m",
		"v.ack.initiator",
		"v.state.m",
		"v.status.m",
		"v.state.val",
		"v.status.val",
		"v.total_state_changes",
		"t",
		"v.creation_date",
		"v.last_event_date",
		"v.last_update_date",
		"v.ack.t",
		"v.state.t",
		"v.status.t",
		"v.resolved",
		"v.activation_date",
		"v.duration",
		"v.current_state_duration",
		"v.snooze_duration",
		"v.pbh_inactive_duration",
		"v.active_duration",
		"tags",
		"extra_details",
		"impact_state",
		// Entity
		"entity._id",
		"entity.name",
		"entity.category.name",
		"entity.type",
		"entity.component",
		"entity.connector",
		"entity.impact_level",
		"entity.ko_events",
		"entity.ok_events",
	},
	WidgetTemplateTypeEntity: {
		"_id",
		"name",
		"category.name",
		"type",
		"component",
		"connector",
		"impact_level",
		"last_event_date",
		"ko_events",
		"ok_events",
		"pbehavior_info",
		"state",
		"impact_state",
		"status",
		"idle_since",
		"enabled",
	},
}

var columnsPrefixByWidget = map[string][]string{
	WidgetTemplateTypeAlarm: {
		"v.infos",
		"links",
		// Entity
		"entity.infos",
		"entity.component_infos",
	},
	WidgetTemplateTypeEntity: {
		"infos",
		"component_infos",
		"links",
	},
}

func init() {
	for _, columns := range columnsByWidget {
		sort.Strings(columns)
	}
}

func GetWidgetColumnParameters() map[string]map[string][]string {
	return map[string]map[string][]string{
		WidgetTypeAlarmsList: {
			WidgetTemplateTypeAlarm: {
				"widgetColumns",
				"widgetGroupColumns",
			},
			WidgetTemplateTypeEntity: {
				"serviceDependenciesColumns",
			},
		},
		WidgetTypeContextExplorer: {
			WidgetTemplateTypeAlarm: {
				"activeAlarmsColumns",
				"resolvedAlarmsColumns",
			},
			WidgetTemplateTypeEntity: {
				"widgetColumns",
				"serviceDependenciesColumns",
			},
		},
		WidgetTypeServiceWeather: {
			WidgetTemplateTypeAlarm: {
				"alarmsList.widgetColumns",
			},
			WidgetTemplateTypeEntity: {
				"serviceDependenciesColumns",
			},
		},
		WidgetTypeAlarmsCounter: {
			WidgetTemplateTypeAlarm: {
				"alarmsList.widgetColumns",
			},
		},
		WidgetTypeAlarmsStatsCalendar: {
			WidgetTemplateTypeAlarm: {
				"alarmsList.widgetColumns",
			},
		},
		WidgetTypeMap: {
			WidgetTemplateTypeAlarm: {
				"alarmsColumns",
			},
			WidgetTemplateTypeEntity: {
				"entitiesColumns",
			},
		},
	}
}

func IsValidWidgetColumn(t, column string) bool {
	columns := columnsByWidget[t]
	if len(columns) == 0 {
		return false
	}

	idx := sort.SearchStrings(columns, column)
	if idx < len(columns) && columns[idx] == column {
		return true
	}

	prefixes := columnsPrefixByWidget[t]
	for _, prefix := range prefixes {
		if column == prefix || strings.HasPrefix(column, prefix+".") {
			return true
		}
	}

	return false
}
