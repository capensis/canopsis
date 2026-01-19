db.webhook_history.createIndex({ status: 1, created_at: 1 }, { name: "status_1_created_at_1" });
db.webhook_history.createIndex({ status: 1, launched_at: 1 }, { name: "status_1_launched_at_1" });
db.webhook_history.createIndex({ next_exec: 1 }, { name: "next_exec_1" });
