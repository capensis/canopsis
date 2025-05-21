db.widgets.updateMany({"type": "AlarmsList"}, {"$unset": {"parameters.remediationInstructionsFilters": ""}})

db.userpreferences.updateMany({"content.remediationInstructionsFilters": {"$ne": null}}, {"$unset": {"content.remediationInstructionsFilters": ""}})

db.permission.deleteOne({"_id": "listalarm_remediationInstructionsFilter"})

db.role.updateMany({"permissions.listalarm_remediationInstructionsFilter": {"$ne": null}}, {
    $unset: {
        "permissions.listalarm_remediationInstructionsFilter": "",
    }
});
