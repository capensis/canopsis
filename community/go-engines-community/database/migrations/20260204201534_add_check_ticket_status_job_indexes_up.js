if (db.getCollectionNames().includes("webhook_check_ticket_status")) {
    db.runCommand({collMod: "webhook_check_ticket_status", changeStreamPreAndPostImages: {enabled: true}})
} else {
    db.createCollection("webhook_check_ticket_status", {changeStreamPreAndPostImages: {enabled: true}})
}

db.webhook_check_ticket_status.createIndex(
    {ticket_id: 1, ticket_system_name: 1},
    {name: "ticket_id_ticket_system_name_1", unique: true}
)
db.webhook_check_ticket_status.createIndex({cmd: 1}, {name: "cmd_1", sparse: true})

if (!db.permission.findOne({_id: "api_ticket_status_job_management"})) {
    db.permission.insertOne({
        _id: "api_ticket_status_job_management",
        name: "api_ticket_status_job_management",
        description: "Ticket status job management",
        groups: ["api", "api_general"]
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_ticket_status_job_management": 1
        }
    });
}

if (!db.permission.findOne({_id: "models_job_management"})) {
    db.permission.insertOne({
        _id: "models_job_management",
        name: "models_job_management",
        description: "Jobs management",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {api_ticket_status_job_management: 1}
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_job_management": 1
        }
    });
}
