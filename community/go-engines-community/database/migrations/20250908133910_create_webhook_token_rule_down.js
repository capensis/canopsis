db.webhook_token_rule.drop();

db.permission.deleteMany({
    _id: {
        $in: [
            "api_webhook_token_rule",
            "models_externalAuthTokens",
        ]
    }
});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_webhook_token_rule": "",
        "permissions.models_externalAuthTokens": "",
    }
});
