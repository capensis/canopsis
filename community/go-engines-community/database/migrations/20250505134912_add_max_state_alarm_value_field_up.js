db.periodical_alarm.updateMany(
    {},
    [
        {
            $set: {
                "v.max_state": "$v.state.val"
            }
        }
    ]
);

db.resolved_alarms.updateMany(
    {},
    [
        {
            $set: {
                "v.max_state": "$v.state.val"
            }
        }
    ]
);
