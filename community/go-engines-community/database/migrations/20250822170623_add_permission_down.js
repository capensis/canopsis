db.permission.deleteOne({_id: "models_remediationInstructionApprove"});

db.role.updateMany({}, {
    $unset: {
        "permissions.models_remediationInstructionApprove": "",
    }
});
