db.periodical_alarm.updateMany({"v.comments": {$ne: null}}, {"$unset": {"v.comments": ""}})
db.resolved_alarms.updateMany({"v.comments": {$ne: null}}, {"$unset": {"v.comments": ""}})
