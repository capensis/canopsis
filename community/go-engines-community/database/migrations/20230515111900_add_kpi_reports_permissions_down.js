db.map.drop();

db.default_rights.deleteMany({
    _id: {
        $in: [
            "userStatistics_interval",
            "userStatistics_listFilters",
            "userStatistics_editFilter",
            "userStatistics_addFilter",
            "userStatistics_userFilter",
            "alarmStatistics_interval",
            "alarmStatistics_listFilters",
            "alarmStatistics_editFilter",
            "alarmStatistics_addFilter",
            "alarmStatistics_userFilter",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.userStatistics_interval": "",
        "rights.userStatistics_listFilters": "",
        "rights.userStatistics_editFilter": "",
        "rights.userStatistics_addFilter": "",
        "rights.userStatistics_userFilter": "",
        "rights.alarmStatistics_interval": "",
        "rights.alarmStatistics_listFilters": "",
        "rights.alarmStatistics_editFilter": "",
        "rights.alarmStatistics_addFilter": "",
        "rights.alarmStatistics_userFilter": "",
    }
});
