package view

import (
	"sort"
	"strings"
)

var columnsByType = map[string][]string{
	WidgetTemplateTypeAlarmColumns: {
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
		"entity.last_pbehavior_date",
	},
	WidgetTemplateTypeEntityColumns: {
		"_id",
		"name",
		"category.name",
		"type",
		"component",
		"connector",
		"impact_level",
		"last_event_date",
		"last_pbehavior_date",
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

var columnsPrefixByType = map[string][]string{
	WidgetTemplateTypeAlarmColumns: {
		"v.infos",
		"links",
		// Entity
		"entity.infos",
		"entity.component_infos",
	},
	WidgetTemplateTypeEntityColumns: {
		"infos",
		"component_infos",
		"links",
	},
}

func init() {
	for _, columns := range columnsByType {
		sort.Strings(columns)
	}
}

func GetWidgetTemplateParameters() map[string]map[string][]string {
	return map[string]map[string][]string{
		WidgetTypeAlarmsList: {
			WidgetTemplateTypeAlarmColumns: {
				"widgetColumns",
				"widgetGroupColumns",
			},
			WidgetTemplateTypeEntityColumns: {
				"serviceDependenciesColumns",
			},
			WidgetTemplateTypeAlarmMoreInfos: {
				"moreInfoTemplate",
			},
		},
		WidgetTypeContextExplorer: {
			WidgetTemplateTypeAlarmColumns: {
				"activeAlarmsColumns",
				"resolvedAlarmsColumns",
			},
			WidgetTemplateTypeEntityColumns: {
				"widgetColumns",
				"serviceDependenciesColumns",
			},
		},
		WidgetTypeServiceWeather: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsList.widgetColumns",
			},
			WidgetTemplateTypeEntityColumns: {
				"serviceDependenciesColumns",
			},
			WidgetTemplateTypeServiceWeatherItem: {
				"blockTemplate",
			},
			WidgetTemplateTypeServiceWeatherModal: {
				"modalTemplate",
			},
			WidgetTemplateTypeServiceWeatherEntity: {
				"entityTemplate",
			},
		},
		WidgetTypeAlarmsCounter: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsList.widgetColumns",
			},
		},
		WidgetTypeAlarmsStatsCalendar: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsList.widgetColumns",
			},
		},
		WidgetTypeMap: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsColumns",
			},
			WidgetTemplateTypeEntityColumns: {
				"entitiesColumns",
			},
		},
	}
}

func IsValidWidgetColumn(t, column string) bool {
	columns := columnsByType[t]
	if len(columns) == 0 {
		return false
	}

	idx := sort.SearchStrings(columns, column)
	if idx < len(columns) && columns[idx] == column {
		return true
	}

	prefixes := columnsPrefixByType[t]
	for _, prefix := range prefixes {
		if column == prefix || strings.HasPrefix(column, prefix+".") {
			return true
		}
	}

	return false
}
