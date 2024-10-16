db.action_scenario.createIndex({enabled: 1, priority: 1}, {name: "enabled_1_priority_1"});
db.alarm_tag.createIndex({type: 1, value: 1}, {name: "type_1_value_1"});
db.default_entities.createIndex({"pbehavior_info.canonical_type": 1}, {name: "pbh_canonical_type_1", partialFilterExpression: {"pbehavior_info.canonical_type": {$exists: true}}});
db.default_importgraph.createIndex({status: 1}, {name: "status_1"});
db.eventfilter.createIndex({enabled: 1, type: 1}, {name: "enabled_1_type_1"});
db.export_task.createIndex({status: 1}, {name: "status_1"});
db.instruction.createIndex({enabled: 1, type: 1, status: 1}, {name: "enabled_1_type_1_status_1"});
db.periodical_alarm.createIndex({not_acked_since: 1, not_acked_metric_type: 1}, {name: "not_acked_since_1_not_acked_metric_type_1"});
db.user.createIndex({name: 1, source: 1}, {name: "name_1_source_1"});
db.viewgroups.createIndex({is_private: 1, author: 1}, {name: "is_private_1_author_1"});
db.webhook_history.createIndex({execution: 1}, {name: "execution_1"});