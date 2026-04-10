db.periodical_alarm.updateMany(
    {},
    [
        {
            $set: {
                "v.initial_status": "$v.status.val"
            }
        }
    ]
);

db.periodical_alarm.updateMany(
    {"v.status.val": 5},
    [
        {
            $set: {
                "v.no_events_date": "$v.status.t"
            }
        }
    ]
);

db.resolved_alarms.updateMany(
    {},
    [
        {
            $set: {
                "v.initial_status": "$v.status.val"
            }
        }
    ]
);
