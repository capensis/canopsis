db.periodical_alarm.updateMany({}, {
    $unset: {
        "v.inactive_duration": "",
        "v.inactive_start": "",
    }
});

db.resolved_alarms.updateMany({}, {
    $unset: {
        "v.duration": "",
        "v.current_state_duration": "",
        "v.inactive_duration": "",
        "v.active_duration": "",
    }
});
