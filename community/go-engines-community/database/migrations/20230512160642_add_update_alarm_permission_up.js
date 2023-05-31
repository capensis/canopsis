if (!db.default_rights.findOne({_id: "api_alarm_update"})) {
    db.default_rights.insertOne({
        _id: "api_alarm_update",
        crecord_name: "api_alarm_update",
        crecord_type: "action",
        description: "Update alarms"
    });
    db.default_rights.updateMany(
        {
            crecord_type: "role",
            "rights.api_event.checksum": 1,
        },
        {
            $set: {
                "rights.api_alarm_update": {
                    checksum: 1,
                    crecord_type: "right"
                }
            }
        },
    );
}
