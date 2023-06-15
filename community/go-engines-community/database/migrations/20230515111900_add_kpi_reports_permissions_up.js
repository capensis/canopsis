if (!db.default_rights.findOne({_id: "userStatistics_interval"})) {
    db.default_rights.insertOne({
        _id: "userStatistics_interval",
        crecord_name: "userStatistics_interval",
        crecord_type: "action",
        desc: "Rights on user statistics: Interval"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.userStatistics_interval": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "userStatistics_listFilters"})) {
    db.default_rights.insertOne({
        _id: "userStatistics_listFilters",
        crecord_name: "userStatistics_listFilters",
        crecord_type: "action",
        desc: "Rights on user statistics: List filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.userStatistics_listFilters": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "userStatistics_editFilter"})) {
    db.default_rights.insertOne({
        _id: "userStatistics_editFilter",
        crecord_name: "userStatistics_editFilter",
        crecord_type: "action",
        desc: "Rights on user statistics: Edit filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.userStatistics_editFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "userStatistics_addFilter"})) {
    db.default_rights.insertOne({
        _id: "userStatistics_addFilter",
        crecord_name: "userStatistics_addFilter",
        crecord_type: "action",
        desc: "Rights on user statistics: Add filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.userStatistics_addFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "userStatistics_userFilter"})) {
    db.default_rights.insertOne({
        _id: "userStatistics_userFilter",
        crecord_name: "userStatistics_userFilter",
        crecord_type: "action",
        desc: "Rights on user statistics: User filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.userStatistics_userFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}

if (!db.default_rights.findOne({_id: "alarmStatistics_interval"})) {
    db.default_rights.insertOne({
        _id: "alarmStatistics_interval",
        crecord_name: "alarmStatistics_interval",
        crecord_type: "action",
        desc: "Rights on alarm statistics: Interval"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.alarmStatistics_interval": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "alarmStatistics_listFilters"})) {
    db.default_rights.insertOne({
        _id: "alarmStatistics_listFilters",
        crecord_name: "alarmStatistics_listFilters",
        crecord_type: "action",
        desc: "Rights on alarm statistics: List filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.alarmStatistics_listFilters": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "alarmStatistics_editFilter"})) {
    db.default_rights.insertOne({
        _id: "alarmStatistics_editFilter",
        crecord_name: "alarmStatistics_editFilter",
        crecord_type: "action",
        desc: "Rights on alarm statistics: Edit filters"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.alarmStatistics_editFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "alarmStatistics_addFilter"})) {
    db.default_rights.insertOne({
        _id: "alarmStatistics_addFilter",
        crecord_name: "alarmStatistics_addFilter",
        crecord_type: "action",
        desc: "Rights on alarm statistics: Add filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.alarmStatistics_addFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
if (!db.default_rights.findOne({_id: "alarmStatistics_userFilter"})) {
    db.default_rights.insertOne({
        _id: "alarmStatistics_userFilter",
        crecord_name: "alarmStatistics_userFilter",
        crecord_type: "action",
        desc: "Rights on alarm statistics: User filter"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.alarmStatistics_userFilter": {
                checksum: 1,
                crecord_type: "right"
            }
        }
    });
}
