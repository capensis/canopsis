db.runCommand({collMod: "webhook_history", changeStreamPreAndPostImages: {enabled: false}})
db.runCommand({collMod: "webhook_token_history", changeStreamPreAndPostImages: {enabled: false}})

db.permission.deleteMany({_id: {$in: ["api_webhook_job_management", "api_webhook_token_job_management", "api_remediation_job_management"]}});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_webhook_job_management": "",
        "permissions.api_webhook_token_job_management": "",
        "permissions.api_remediation_job_management": "",
    }
});

db.permission.updateOne({_id: "models_job_management"}, {
    $unset: {
        "api_permissions.api_webhook_job_management": "",
        "api_permissions.api_webhook_token_job_management": "",
        "api_permissions.api_remediation_job_management": "",
    }
});
