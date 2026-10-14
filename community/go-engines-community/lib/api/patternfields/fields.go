package patternfields

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/statesetting"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
)

var disabledFieldsInAlarmPattern = map[string]map[string]bool{
	mongo.IdleRuleMongoCollection:          {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.MetaAlarmRulesMongoCollection:    {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.FlappingRuleMongoCollection:      {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.ResolveRuleMongoCollection:       {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.ScenarioCollection:               {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.InstructionMongoCollection:       {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.DeclareTicketRuleCollection:      {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.LinkRuleMongoCollection:          {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true},
	mongo.DynamicInfosRulesMongoCollection: {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true, "v.duration": true, "v.infos": true},
	mongo.AlarmTagCollection:               {"v.last_event_date": true, "v.last_update_date": true, "v.resolved": true, "v.duration": true, "tags": true},
	mongo.WidgetFiltersMongoCollection:     {},
}

var onlyAbsoluteTimeCondFieldsInAlarmPattern = map[string]map[string]bool{
	mongo.IdleRuleMongoCollection:          {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.DynamicInfosRulesMongoCollection: {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.MetaAlarmRulesMongoCollection:    {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.FlappingRuleMongoCollection:      {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.ResolveRuleMongoCollection:       {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.ScenarioCollection:               {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.InstructionMongoCollection:       {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.DeclareTicketRuleCollection:      {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.LinkRuleMongoCollection:          {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
	mongo.AlarmTagCollection:               {"v.creation_date": true, "v.ack.t": true, "v.activation_date": true},
}

var disabledFieldsInEntityPattern = map[string]map[string]bool{
	mongo.StateSettingsMongoCollection:     {"last_event_date": true, "component": true, "type": true, "component_infos": true},
	mongo.EntityMongoCollection:            {"last_event_date": true, "connector": true},
	mongo.PbehaviorMongoCollection:         {"last_event_date": true},
	mongo.IdleRuleMongoCollection:          {"last_event_date": true},
	mongo.DynamicInfosRulesMongoCollection: {"last_event_date": true},
	mongo.MetaAlarmRulesMongoCollection:    {"last_event_date": true},
	mongo.FlappingRuleMongoCollection:      {"last_event_date": true},
	mongo.ResolveRuleMongoCollection:       {"last_event_date": true},
	mongo.ScenarioCollection:               {"last_event_date": true},
	mongo.InstructionMongoCollection:       {"last_event_date": true},
	mongo.KpiFilterMongoCollection:         {"last_event_date": true},
	mongo.DeclareTicketRuleCollection:      {"last_event_date": true},
	mongo.LinkRuleMongoCollection:          {"last_event_date": true},
	mongo.AlarmTagCollection:               {"last_event_date": true},
	mongo.EventFilterRuleCollection:        {},
	mongo.WidgetFiltersMongoCollection:     {},
}

var disabledFieldsInInheritedEntityPattern = map[string]map[string]map[string]bool{
	mongo.StateSettingsMongoCollection: {
		statesetting.RuleTypeService:   {"last_event_date": true},
		statesetting.RuleTypeComponent: {"last_event_date": true, "component": true, "type": true, "component_infos": true},
	},
}

var disabledFieldsInEventPattern = map[string]map[string]bool{
	mongo.EventFilterRuleCollection:   {},
	mongo.EventRecordsMongoCollection: {},
}

var disabledFieldsInPbehaviorPattern = map[string]map[string]bool{
	mongo.WidgetFiltersMongoCollection: {},
}

var disabledFieldsInWeatherServicePattern = map[string]map[string]bool{
	mongo.WidgetFiltersMongoCollection: {},
}

func GetDisabledFieldsInAlarmPattern() map[string]map[string]bool {
	return disabledFieldsInAlarmPattern
}

func GetOnlyAbsoluteTimeCondFieldsInAlarmPattern() map[string]map[string]bool {
	return onlyAbsoluteTimeCondFieldsInAlarmPattern
}

func GetDisabledFieldsInEntityPattern() map[string]map[string]bool {
	return disabledFieldsInEntityPattern
}

func GetDisabledFieldsInInheritedEntityPattern() map[string]map[string]map[string]bool {
	return disabledFieldsInInheritedEntityPattern
}

func GetDisabledFieldsInEventPattern() map[string]map[string]bool {
	return disabledFieldsInEventPattern
}

func GetDisabledFieldsInPbehaviorPattern() map[string]map[string]bool {
	return disabledFieldsInPbehaviorPattern
}

func GetDisabledFieldsInWeatherServicePattern() map[string]map[string]bool {
	return disabledFieldsInWeatherServicePattern
}
