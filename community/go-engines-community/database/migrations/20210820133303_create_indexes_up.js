db.default_entities.createIndex({enabled: 1}, {name: "enabled_1"});

db.periodical_alarm.createIndex({d: 1}, {name: "d_1"});
db.periodical_alarm.createIndex(
    {
        "v.meta": 1,
        "v.creation_date": 1,
    },
    {
        name: "v.meta_1_v.creation_date_1",
        partialFilterExpression: {
            "v.meta": {$exists: true}
        }
    }
);
db.periodical_alarm.createIndex({"v.resolved": 1}, {name: "v.resolved_1"});

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
