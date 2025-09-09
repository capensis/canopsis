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
            "permissions.api_external_api_webhook_token_ruledata_table": 15
        }
    });
}
