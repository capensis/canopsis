db.default_entities.dropIndex("enabled_1");

db.periodical_alarm.dropIndex("t_1");
db.periodical_alarm.dropIndex("d_1");
db.periodical_alarm.dropIndex("v.meta_1_v.creation_date_1");
db.periodical_alarm.dropIndex("v.resolved_1");
db.periodical_alarm.dropIndex("v.creation_date_1");
db.periodical_alarm.dropIndex("v.last_event_date_1");
db.periodical_alarm.dropIndex("v.last_update_date_1");

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
