// revert deleted indexes
db.webhook_history.createIndex({ status: 1, created_at: 1 }, { name: "status_1_created_at_1" });
db.webhook_history.createIndex({ status: 1, launched_at: 1 }, { name: "status_1_launched_at_1" });

db.webhook_history.dropIndex("status_1_last_ping_1_created_at_1");
db.webhook_token_history.dropIndex("rule_1_created_at_-1");