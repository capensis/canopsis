db.periodical_alarm.updateMany({}, {$unset: {"v.max_state": "", "v.initial_state": ""}})
db.resolved_alarms.updateMany({}, {$unset: {"v.max_state": "", "v.initial_state": ""}})
