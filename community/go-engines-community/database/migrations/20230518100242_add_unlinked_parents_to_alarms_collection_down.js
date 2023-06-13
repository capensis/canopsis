db.periodical_alarm.updateMany(
    {"v.unlinked_parents": {$exists: true}},
    {$unset: {"v.unlinked_parents": ""}},
);

db.default_rights.deleteMany({
    _id: "listalarm_metaAlarmGroup",
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.listalarm_metaAlarmGroup": "",
    }
});
