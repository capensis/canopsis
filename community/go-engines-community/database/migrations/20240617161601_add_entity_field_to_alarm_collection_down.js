db.periodical_alarm.updateMany({}, {$unset: {entity: ""}});

db.periodical_alarm.dropIndex("entity_id_1");
db.periodical_alarm.dropIndex("entity_enabled_1");
db.periodical_alarm.dropIndex("entity_type_1");
db.periodical_alarm.dropIndex("entity_connector_1");
db.periodical_alarm.dropIndex("entity_component_1");
db.periodical_alarm.dropIndex("entity_services_1");
db.periodical_alarm.dropIndex("entity_type_service_1")
