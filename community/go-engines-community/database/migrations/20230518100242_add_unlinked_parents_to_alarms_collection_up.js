db.periodical_alarm.updateMany(
    {"v.unlinked_parents": {$exists: false}},
    {$set: {"v.unlinked_parents": []}},
);

if (!db.default_rights.findOne({_id: "listalarm_metaAlarmGroup"})) {
    db.default_rights.insertOne({
        _id: "listalarm_metaAlarmGroup",
        crecord_name: "listalarm_metaAlarmGroup",
        crecord_type: "action",
        description: "Rights on listalarm: Meta alarm actions"
    });
    db.default_rights.updateMany(
        {
            crecord_type: "role",
            "rights.listalarm_manualMetaAlarmGroup.checksum": 1,
        },
        {
            $set: {
                "rights.listalarm_metaAlarmGroup": {
                    checksum: 1,
                    crecord_type: "right"
                }
            }
        },
    );
}
