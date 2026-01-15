db.periodical_alarm.updateMany({}, {$unset: {"v.initial_status": "", "v.no_events_date": ""}})
db.resolved_alarms.updateMany({}, {$unset: {"v.initial_status": "", "v.no_events_date": ""}})
