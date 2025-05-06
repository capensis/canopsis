db.periodical_alarm.updateMany({}, {$unset: {"v.max_state": ""}})
db.resolved_alarms.updateMany({}, {$unset: {"v.max_state": ""}})
