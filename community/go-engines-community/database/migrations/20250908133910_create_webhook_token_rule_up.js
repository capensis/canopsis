if (db.getCollectionNames().includes("webhook_token_rule")) {
    db.runCommand({collMod: "webhook_token_rule", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("webhook_token_rule", {changeStreamPreAndPostImages: {enabled: true}})
}

db.webhook_token_rule.createIndex({name: 1}, {name: "name_1", unique: true});

if (!db.permission.findOne({_id: "api_webhook_token_rule"})) {
    db.permission.insertOne({
        _id: "api_webhook_token_rule",
        name: "api_webhook_token_rule",
        type: "CRUD",
        description: "Webhook token rule",
        groups: ["api", "api_rules"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_webhook_token_rule": 15
        }
    });
}

if (!db.permission.findOne({_id: "modals_externalAuthTokens"})) {
    db.permission.insertOne({
        _id: "modals_externalAuthTokens",
        name: "modals_externalAuthTokens",
        type: "CRUD",
        description: "Webhook token rule",
        groups: ["technical", "technical_admin", "technical_admin_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.modals_externalAuthTokens": 15
        }
    });
}
