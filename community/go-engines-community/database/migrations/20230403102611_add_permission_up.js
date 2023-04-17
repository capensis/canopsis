if (!db.default_rights.findOne({_id: "listalarm_fastRemoveAlarm"})) {
    db.default_rights.insertOne({
        _id: "listalarm_fastRemoveAlarm",
        crecord_name: "listalarm_fastRemoveAlarm",
        crecord_type: "action",
        description: "Rights on listalarm: fast remove alarm"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.listalarm_fastRemoveAlarm": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
