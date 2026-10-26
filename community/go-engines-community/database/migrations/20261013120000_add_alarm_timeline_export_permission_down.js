db.permission.deleteOne({_id: "listalarm_exportAlarmTimeline"});

db.role.updateMany({}, {
    $unset: {
        "permissions.listalarm_exportAlarmTimeline": ""
    }
});
