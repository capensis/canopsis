// Remove indexes no longer needed
db.webhook_history.getIndexes().forEach(function (index) {
    if (index.name === "status_1_created_at_1" || index.name === "status_1_launched_at_1") {
        db.webhook_history.dropIndex(index.name);
    }
});

db.webhook_history.createIndex({
    status: 1,
    last_ping: 1,
    created_at: 1
}, {name: "status_1_last_ping_1_created_at_1"});

db.webhook_token_history.createIndex({
    rule: 1,
    created_at: -1
}, {name: "rule_1_created_at_-1"});