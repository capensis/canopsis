if (!db.permission.findOne({_id: "listalarm_exportAlarmTimeline"})) {
    db.permission.insertOne({
        _id: "listalarm_exportAlarmTimeline",
        name: "listalarm_exportAlarmTimeline",
        description: "Export alarm timeline as CSV file",
        groups: ["widgets", "widgets_alarmslist", "widgets_alarmslist_alarmactions"],
        api_permissions: {
            api_alarm_read: 1
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.listalarm_exportAlarmTimeline": 1
        }
    });
}
