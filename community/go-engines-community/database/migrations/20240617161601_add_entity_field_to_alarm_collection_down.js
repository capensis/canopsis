db.periodical_alarm.updateMany({}, {$unset: {e: ""}});

db.periodical_alarm.dropIndex("e_id_1");
db.periodical_alarm.dropIndex("e_enabled_1");
db.periodical_alarm.dropIndex("e_type_1");
db.periodical_alarm.dropIndex("e_connector_1");
db.periodical_alarm.dropIndex("e_component_1");
db.periodical_alarm.dropIndex("e_services_1");
db.periodical_alarm.dropIndex("e_type_service_1")
