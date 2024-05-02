db.default_entities.updateMany(
    {"pbehavior_info.id": {$nin: [null, ""]}},
    {
        $set: {
            "pbehavior_info.author": "system"
        }
    }
);

db.periodical_alarm.updateMany(
    {"v.pbehavior_info.id": {$nin: [null, ""]}},
    {
        $set: {
            "v.pbehavior_info.author": "system"
        }
    }
);

db.resolved_alarms.updateMany(
    {"v.pbehavior_info.id": {$nin: [null, ""]}},
    {
        $set: {
            "v.pbehavior_info.author": "system"
        }
    }
);
