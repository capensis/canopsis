db.eventfilter_failure.updateMany({rule_updated: {$ne: null}}, {$unset: {rule_updated: ""}});
db.eventfilter_failure.dropIndex("t_1");
db.eventfilter_failure.dropIndex("rule_1_rule_updated_1");
db.eventfilter_failure.createIndex({rule: 1}, {name: "rule_1"});

db.user_notification.updateMany({"rule.updated": {$ne: null}}, {$unset: {"rule.updated": ""}});
