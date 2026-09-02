package patternfields

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"

var disabledFieldsInAlarmPattern = map[string][]string{
	mongo.IdleRuleMongoCollection:          {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.MetaAlarmRulesMongoCollection:    {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.FlappingRuleMongoCollection:      {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.ResolveRuleMongoCollection:       {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.ScenarioCollection:               {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.InstructionMongoCollection:       {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.DeclareTicketRuleCollection:      {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.LinkRuleMongoCollection:          {"v.last_event_date", "v.last_update_date", "v.resolved"},
	mongo.DynamicInfosRulesMongoCollection: {"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "v.infos"},
	mongo.AlarmTagCollection:               {"v.last_event_date", "v.last_update_date", "v.resolved", "v.duration", "tags"},
	mongo.WidgetFiltersMongoCollection:     {},
}

var onlyAbsoluteTimeCondFieldsInAlarmPattern = map[string][]string{
	mongo.IdleRuleMongoCollection:          {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.DynamicInfosRulesMongoCollection: {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.MetaAlarmRulesMongoCollection:    {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.FlappingRuleMongoCollection:      {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.ResolveRuleMongoCollection:       {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.ScenarioCollection:               {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.InstructionMongoCollection:       {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.DeclareTicketRuleCollection:      {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.LinkRuleMongoCollection:          {"v.creation_date", "v.ack.t", "v.activation_date"},
	mongo.AlarmTagCollection:               {"v.creation_date", "v.ack.t", "v.activation_date"},
}

var disabledFieldsInEntityPattern = map[string][]string{
	mongo.StateSettingsMongoCollection:     {"last_event_date", "component", "component_infos"},
	mongo.EntityMongoCollection:            {"last_event_date", "connector", "component_infos"},
	mongo.PbehaviorMongoCollection:         {"last_event_date"},
	mongo.IdleRuleMongoCollection:          {"last_event_date"},
	mongo.DynamicInfosRulesMongoCollection: {"last_event_date"},
	mongo.MetaAlarmRulesMongoCollection:    {"last_event_date"},
	mongo.FlappingRuleMongoCollection:      {"last_event_date"},
	mongo.ResolveRuleMongoCollection:       {"last_event_date"},
	mongo.ScenarioCollection:               {"last_event_date"},
	mongo.InstructionMongoCollection:       {"last_event_date"},
	mongo.KpiFilterMongoCollection:         {"last_event_date"},
	mongo.DeclareTicketRuleCollection:      {"last_event_date"},
	mongo.LinkRuleMongoCollection:          {"last_event_date"},
	mongo.AlarmTagCollection:               {"last_event_date"},
	mongo.EventFilterRuleCollection:        {},
	mongo.WidgetFiltersMongoCollection:     {},
}

var disabledFieldsInEventPattern = map[string][]string{
	mongo.EventFilterRuleCollection:   {},
	mongo.EventRecordsMongoCollection: {},
}

var disabledFieldsInPbehaviorPattern = map[string][]string{
	mongo.WidgetFiltersMongoCollection: {},
}

var disabledFieldsInWeatherServicePattern = map[string][]string{
	mongo.WidgetFiltersMongoCollection: {},
}

func GetDisabledFieldsInAlarmPattern() map[string][]string {
	return disabledFieldsInAlarmPattern
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern() map[string][]string {
	return onlyAbsoluteTimeCondFieldsInAlarmPattern
}

func GetDisabledFieldsInEntityPattern() map[string][]string {
	return disabledFieldsInEntityPattern
}

func GetDisabledFieldsInEventPattern() map[string][]string {
	return disabledFieldsInEventPattern
}

func GetDisabledFieldsInPbehaviorPattern() map[string][]string {
	return disabledFieldsInPbehaviorPattern
}

func GetDisabledFieldsInWeatherServicePattern() map[string][]string {
	return disabledFieldsInWeatherServicePattern
}
