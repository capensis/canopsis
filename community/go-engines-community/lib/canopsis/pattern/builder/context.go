package builder

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/view"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

const promptContextStateSettingsServiceInherited = "state_settings_service_inherited"

var collectionByPromptContext = map[string]string{
	"idle_rule":           mongo.IdleRuleMongoCollection,
	"scenario":            mongo.ScenarioCollection,
	"flapping_rule":       mongo.FlappingRuleMongoCollection,
	"resolve_rule":        mongo.ResolveRuleMongoCollection,
	"alarm_tag":           mongo.AlarmTagCollection,
	"link_rule":           mongo.LinkRuleMongoCollection,
	"instruction":         mongo.InstructionMongoCollection,
	"dynamic_infos":       mongo.DynamicInfosRulesMongoCollection,
	"meta_alarm_rule":     mongo.MetaAlarmRulesMongoCollection,
	"declare_ticket_rule": mongo.DeclareTicketRuleCollection,
	"pbehavior":           mongo.PbehaviorMongoCollection,
	"entity_service":      mongo.EntityMongoCollection,
	"state_settings":      mongo.StateSettingsMongoCollection,
	promptContextStateSettingsServiceInherited: mongo.StateSettingsMongoCollection,
	"kpi_filter":   mongo.KpiFilterMongoCollection,
	"eventfilter":  mongo.EventFilterRuleCollection,
	"event_record": mongo.EventRecordsMongoCollection,

	view.LLMContextPrefixWidgetFilter + view.WidgetTypeAlarmsList:          mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeContextExplorer:     mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeServiceWeather:      mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeMap:                 mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeAlarmsCounter:       mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeAlarmsStatsCalendar: mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeAvailability:        mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeUserStatistics:      mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeAlarmStatistics:     mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeBarChart:            mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeLineChart:           mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypePieChart:            mongo.WidgetFiltersMongoCollection,
	view.LLMContextPrefixWidgetFilter + view.WidgetTypeNumbers:             mongo.WidgetFiltersMongoCollection,

	"corporate_alarm_pattern":           mongo.PatternMongoCollection,
	"corporate_entity_pattern":          mongo.PatternMongoCollection,
	"corporate_pbehavior_pattern":       mongo.PatternMongoCollection,
	"corporate_weather_service_pattern": mongo.PatternMongoCollection,
}

func CollectionByPromptContext(promptContext string) (string, bool) {
	c, ok := collectionByPromptContext[promptContext]

	return c, ok
}
