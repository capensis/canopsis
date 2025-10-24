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
		"v.ticket.a",
		"v.ticket.m",
		"v.ticket.ticket",
		"v.state.val",
		"v.status.val",
		"v.total_state_changes",
		"t",
		"v.creation_date",
		"v.last_event_date",
		"v.last_update_date",
		"v.ack.t",
		"v.ticket.t",
		"v.state.t",
		"v.status.t",
		"v.resolved",
		"v.activation_date",
		"v.duration",
		"v.current_state_duration",
		"v.snooze_duration",
		"v.pbh_inactive_duration",
		"v.active_duration",
		"v.events_count",
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
		"imported",
		"import_source",
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

var widgetExportColumns = map[string]map[string][]string{
	WidgetTypeAlarmsList: {
		"widgetExportColumns": []string{
			"assigned_instructions",
		},
	},
}

func init() {
	for _, columns := range columnsByType {
		sort.Strings(columns)
	}
	for wt, v := range widgetExportColumns {
		for param := range v {
			sort.Strings(widgetExportColumns[wt][param])
		}
	}
}

func GetWidgetTemplateParameters() map[string]map[string][]string {
	return map[string]map[string][]string{
		WidgetTypeAlarmsList: {
			WidgetTemplateTypeAlarmColumns: {
				"widgetColumns",
				"widgetGroupColumns",
				"widgetExportColumns",
			},
			WidgetTemplateTypeEntityColumns: {
				"serviceDependenciesColumns",
			},
			WidgetTemplateTypeAlarmMoreInfos: {
				"moreInfoTemplate",
			},
			WidgetTemplateTypeAlarmExportToPDF: {
				"exportPdfTemplate",
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
				"widgetExportColumns",
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
			WidgetTemplateTypeAlarmMoreInfos: {
				"alarmsList.moreInfoTemplate",
			},
		},
		WidgetTypeAlarmsCounter: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsList.widgetColumns",
			},
			WidgetTemplateTypeAlarmMoreInfos: {
				"alarmsList.moreInfoTemplate",
			},
		},
		WidgetTypeAlarmsStatsCalendar: {
			WidgetTemplateTypeAlarmColumns: {
				"alarmsList.widgetColumns",
			},
			WidgetTemplateTypeAlarmMoreInfos: {
				"alarmsList.moreInfoTemplate",
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

func IsValidWidgetExportColumn(widgetType, param, column string) bool {
	if columns, ok := widgetExportColumns[widgetType][param]; ok {
		idx := sort.SearchStrings(columns, column)
		if idx < len(columns) && columns[idx] == column {
			return true
		}
	}

	return false
}
