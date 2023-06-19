db.default_rights.deleteMany({
    _id: "api_alarm_update",
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_alarm_update": "",
    }
});
