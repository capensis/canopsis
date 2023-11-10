db.default_entities.createIndex({enabled: 1}, {name: "enabled_1"});
db.default_entities.createIndex({type: 1}, {name: "type_1"});

db.periodical_alarm.createIndex({t: 1}, {name: "t_1"});
db.periodical_alarm.createIndex({d: 1}, {name: "d_1"});
db.periodical_alarm.createIndex({"v.meta": 1, "v.creation_date": 1},{name: "v.meta_1_v.creation_date_1", partialFilterExpression: {"v.meta": {$exists: true}}});
db.periodical_alarm.createIndex({"v.resolved": 1}, {name: "v.resolved_1"});
db.periodical_alarm.createIndex({"v.creation_date": 1}, {name: "v.creation_date_1"});
db.periodical_alarm.createIndex({"v.last_event_date": 1}, {name: "v.last_event_date_1"});
db.periodical_alarm.createIndex({"v.last_update_date": 1}, {name: "v.last_update_date_1"});
db.periodical_alarm.createIndex({"v.parents": 1}, {name: "v.parents_1"});

db.resolved_alarms.createIndex({t: 1}, {name: "t_1"});
db.resolved_alarms.createIndex({d: 1}, {name: "d_1"});
db.resolved_alarms.createIndex({"v.meta": 1, "v.creation_date": 1}, {name: "v.meta_1_v.creation_date_1", partialFilterExpression: {"v.meta": {$exists: true}}});
db.resolved_alarms.createIndex({"v.creation_date": 1}, {name: "v.creation_date_1"});
db.resolved_alarms.createIndex({"v.last_event_date": 1}, {name: "v.last_event_date_1"});
db.resolved_alarms.createIndex({"v.last_update_date": 1}, {name: "v.last_update_date_1"});

db.pbehavior.createIndex({name: 1}, {name: "name_1", unique: true});

db.pbehavior_type.createIndex({priority: 1}, {name: "priority_1", unique: true});

db.junit_test_suite.createIndex({test_suite_id: 1}, {name: "test_suite_id_1"});
db.junit_test_suite.createIndex({entity_id: 1}, {name: "entity_id_1"});
db.junit_test_suite.createIndex({filename: 1}, {name: "filename_1"});

db.junit_test_case_media.createIndex({relative_filepath: 1}, {name: "relative_filepath_1"});

db.instruction_execution.createIndex({
    instruction: 1,
    status: 1,
}, {name: "instruction_1_status_1"});
db.instruction_execution.createIndex({
    status: 1,
    completed_at: 1,
}, {name: "status_1_completed_at_1"});
db.instruction_execution.createIndex({alarm: 1}, {name: "alarm_1"});

db.instruction_week_stats.createIndex({
    instruction: 1,
    date: 1,
}, {name: "instruction_1_date_1"});

db.instruction_mod_stats.createIndex({
    instruction: 1,
    date: 1,
}, {name: "instruction_1_date_1"});

db.instruction_rating.createIndex({instruction: 1}, {name: "instruction_1"});

db.job_history.createIndex({
    job: 1,
    status: 1,
}, {name: "job_1_status_1"});
db.job_history.createIndex({
    next_exec: 1,
    status: 1,
}, {name: "next_exec_1_status_1"});
db.job_history.createIndex({execution: 1}, {name: "execution_1"});

db.default_rights.createIndex({
    crecord_type: 1,
    role: 1,
}, {name: "crecord_type_1_role_1"});

db.action_scenario.createIndex({priority: 1}, {name: "priority_1"});

db.action_log.createIndex({"action": 1})
db.action_log.createIndex({"value_type": 1, "value_id": 1})

db.entity_category.createIndex({name: 1}, {name: "name_1"});

db.idle_rule.createIndex({priority: 1}, {name: "priority_1"});

db.job_config.createIndex({created: 1}, {name: "created_1"});

db.job.createIndex({created: 1}, {name: "created_1"});

db.instruction.createIndex({created: 1}, {name: "created_1"});

db.flapping_rule.createIndex({priority: 1}, {name: "priority_1"});

db.resolve_rule.createIndex({priority: 1}, {name: "priority_1"});

db.userpreferences.createIndex({user: 1, widget: 1}, {name: "user_1_widget_1"});

db.views.createIndex({group_id: 1}, {name: "group_id_1"});

db.viewtabs.createIndex({view: 1}, {name: "view_1"});

db.widgets.createIndex({tab: 1}, {name: "tab_1"});

// Can be removed if index is added
var collectionNames = db.getCollectionNames();
if (!collectionNames.includes('pbehavior_reason')) {
    db.createCollection("pbehavior_reason");
}
if (!collectionNames.includes('pbehavior_exception')) {
    db.createCollection("pbehavior_exception");
}
if (!collectionNames.includes('meta_alarm_rules')) {
    db.createCollection("meta_alarm_rules");
}
if (!collectionNames.includes('dynamic_infos')) {
    db.createCollection("dynamic_infos");
}
if (!collectionNames.includes('eventfilter')) {
    db.createCollection("eventfilter");
}
if (!collectionNames.includes('broadcast_message')) {
    db.createCollection("broadcast_message");
}
if (!collectionNames.includes('viewgroups')) {
    db.createCollection("viewgroups");
}
if (!collectionNames.includes('view_playlist')) {
    db.createCollection("view_playlist");
}
