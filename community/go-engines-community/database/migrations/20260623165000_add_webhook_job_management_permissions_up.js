if (db.getCollectionNames().includes("webhook_history")) {
    db.runCommand({collMod: "webhook_history", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("webhook_history", {changeStreamPreAndPostImages: {enabled: true}})
}

if (db.getCollectionNames().includes("webhook_token_history")) {
    db.runCommand({collMod: "webhook_token_history", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("webhook_token_history", {changeStreamPreAndPostImages: {enabled: true}})
}

if (!db.permission.findOne({_id: "api_webhook_job_management"})) {
    db.permission.insertOne({
        _id: "api_webhook_job_management",
        name: "api_webhook_job_management",
        description: "Webhook job management",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_webhook_job_management": 1
        }
    });
}

if (!db.permission.findOne({_id: "api_webhook_token_job_management"})) {
    db.permission.insertOne({
        _id: "api_webhook_token_job_management",
        name: "api_webhook_token_job_management",
        description: "Webhook token job management",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_webhook_token_job_management": 1
        }
    });
}

if (!db.permission.findOne({_id: "api_remediation_job_management"})) {
    db.permission.insertOne({
        _id: "api_remediation_job_management",
        name: "api_remediation_job_management",
        description: "Remediation job management",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_remediation_job_management": 1
        }
    });
}

if (!db.permission.findOne({_id: "models_job_management"})) {
    db.permission.insertOne({
        _id: "models_job_management",
        name: "models_job_management",
        description: "Jobs management",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {
            api_webhook_job_management: 1,
            api_webhook_token_job_management: 1,
            api_remediation_job_management: 1
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_job_management": 1
        }
    });
} else {
    db.permission.updateOne({_id: "models_job_management"}, {
        $set: {
            "api_permissions.api_webhook_job_management": 1,
            "api_permissions.api_webhook_token_job_management": 1,
            "api_permissions.api_remediation_job_management": 1
        }
    });
}
