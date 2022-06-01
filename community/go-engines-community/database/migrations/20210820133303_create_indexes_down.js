db.default_entities.dropIndex("enabled_1");

db.periodical_alarm.dropIndex("t_1");
db.periodical_alarm.dropIndex("d_1");
db.periodical_alarm.dropIndex("v.meta_1_v.creation_date_1");
db.periodical_alarm.dropIndex("v.resolved_1");
db.periodical_alarm.dropIndex("v.creation_date_1");
db.periodical_alarm.dropIndex("v.last_event_date_1");
db.periodical_alarm.dropIndex("v.last_update_date_1");

db.resolved_alarms.dropIndex("t_1");
db.resolved_alarms.dropIndex("d_1");
db.resolved_alarms.dropIndex("v.meta_1_v.creation_date_1");
db.resolved_alarms.dropIndex("v.creation_date_1");
db.resolved_alarms.dropIndex("v.last_event_date_1");
db.resolved_alarms.dropIndex("v.last_update_date_1");

db.pbehavior_type.dropIndex("priority_1");

db.junit_test_suite.dropIndex("test_suite_id_1");
db.junit_test_suite.dropIndex("entity_id_1");
db.junit_test_suite.dropIndex("filename_1");

db.junit_test_case_media.dropIndex("relative_filepath_1");

db.instruction_execution.dropIndex("instruction_1_status_1");
db.instruction_execution.dropIndex("status_1_completed_at_1");
db.instruction_execution.dropIndex("alarm_1");

db.instruction_week_stats.dropIndex("instruction_1_date_1");

db.instruction_mod_stats.dropIndex("instruction_1_date_1");

db.job_history.dropIndex("job_1_status_1");
db.job_history.dropIndex("next_exec_1_status_1");
db.job_history.dropIndex("execution_1");

db.default_rights.dropIndex("crecord_type_1_role_1");

db.action_scenario.dropIndex("priority_1");

db.entity_category.dropIndex("name_1");

db.idle_rule.dropIndex("priority_1");

db.job_config.dropIndex("created_1");

db.job.dropIndex("created_1");

db.instruction.dropIndex("created_1");

db.flapping_rule.dropIndex("priority_1");

db.resolve_rule.dropIndex("priority_1");
