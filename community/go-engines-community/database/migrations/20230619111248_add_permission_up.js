if (!db.default_rights.findOne({_id: "listalarm_unCancel"})) {
    db.default_rights.insertOne({
        _id: "listalarm_unCancel",
        crecord_name: "listalarm_unCancel",
        crecord_type: "action",
        description: "Rights on listalarm: uncancel alarm"
    });
    db.default_rights.updateMany(
        {
            crecord_type: "role",
            "rights.listalarm_removeAlarm.checksum": 1,
        },
        {
            $set: {
                "rights.listalarm_unCancel": {
                    checksum: 1,
                    crecord_type: "right"
                }
            }
        });
}
