if (!db.permission.findOne({_id: "models_remediationInstructionApprove"})) {
    db.permission.insertOne({
        _id: "models_remediationInstructionApprove",
        name: "models_remediationInstructionApprove",
        description: "Instructions - instruction approve",
        groups: ["technical", "technical_admin", "technical_admin_general"],
        api_permissions: {
            api_instruction_approve: 1
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.models_remediationInstructionApprove": 1
        }
    });
}
