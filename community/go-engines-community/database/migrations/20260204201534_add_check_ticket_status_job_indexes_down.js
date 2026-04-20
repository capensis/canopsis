db.runCommand({collMod: "webhook_check_ticket_status", changeStreamPreAndPostImages: {enabled: false}})

db.webhook_check_ticket_status.dropIndex("ticket_id_ticket_system_name_1");
db.webhook_check_ticket_status.dropIndex("cmd_1");

db.permission.deleteMany({_id: {$in: ["api_ticket_status_job_management", "models_job_management"]}});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_ticket_status_job_management": "",
        "permissions.models_job_management": "",
    }
});
