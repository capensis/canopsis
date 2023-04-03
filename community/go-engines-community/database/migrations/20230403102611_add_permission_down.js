db.default_rights.deleteMany({
    _id: "listalarm_fastRemoveAlarm",
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.listalarm_fastRemoveAlarm": "",
    }
});
