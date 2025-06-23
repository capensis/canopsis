db.entity_infos_property.dropIndex("key_1");
db.entity_infos_property.dropIndex("alias_1");
db.entity_infos_property.dropIndex("type_1");

db.alarm_tag.dropIndex("aliases_1");
db.default_entities.dropIndex("aliases_1");
db.eventfilter.dropIndex("aliases_1");
db.flapping_rule.dropIndex("aliases_1");
db.idle_rule.dropIndex("aliases_1");
db.link_rule.dropIndex("aliases_1");
db.pattern.dropIndex("aliases_1");
db.pbehavior.dropIndex("aliases_1");
db.resolve_rule.dropIndex("aliases_1");
db.action_scenario.dropIndex("aliases_1");
db.widget_filters.dropIndex("aliases_1");
db.meta_alarm_rules.dropIndex("aliases_1");
db.declare_ticket_rule.dropIndex("aliases_1");
db.dynamic_infos.dropIndex("aliases_1");
db.instruction.dropIndex("aliases_1");
db.kpi_filter.dropIndex("aliases_1");

db.permission.deleteOne({_id: "api_entity_infos_property"});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_entity_infos_property": "",
    }
});
