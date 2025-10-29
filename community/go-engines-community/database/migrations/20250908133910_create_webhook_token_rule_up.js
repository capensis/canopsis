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

    db.permission.updateMany(
        {
            _id: {
                $in: [
                    "models_exploitation_scenario",
                    "models_exploitation_declareTicketRule"
                ]
            }
        },
        {
            $set: {
                "api_permissions.api_webhook_token_rule": 4
            }
        }
    );
}

if (!db.permission.findOne({_id: "models_externalAuthTokens"})) {
    db.permission.insertOne({
        _id: "models_externalAuthTokens",
        name: "models_externalAuthTokens",
        type: "CRUD",
        description: "Webhook token rule",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {api_webhook_token_rule: 0}
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_externalAuthTokens": 15
        }
    });
}
