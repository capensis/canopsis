db.default_entities.updateMany(
    {"pbehavior_info.id": {$nin: [null, ""]}},
    {
        $unset: {
            "pbehavior_info.author": ""
        }
    }
);

db.periodical_alarm.updateMany(
    {"v.pbehavior_info.id": {$nin: [null, ""]}},
    {
        $unset: {
            "v.pbehavior_info.author": ""
        }
    }
);

db.resolved_alarms.updateMany(
    {"v.pbehavior_info.id": {$nin: [null, ""]}},
    {
        $unset: {
            "v.pbehavior_info.author": ""
        }
    }
);
