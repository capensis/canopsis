if (!db.permission.findOne({_id: "listalarm_remediationInstructionsFilter"})) {
    db.permission.insertOne({
        "_id": "listalarm_remediationInstructionsFilter",
        "groups": [
            "widgets",
            "widgets_alarmslist",
            "widgets_alarmslist_widgetsettings"
        ],
        "api_permissions": {
            "api_instruction": 4
        },
        "name": "listalarm_remediationInstructionsFilter",
        "description": "Set filters by remediation instructions"
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.listalarm_remediationInstructionsFilter": 1
        }
    });
}
