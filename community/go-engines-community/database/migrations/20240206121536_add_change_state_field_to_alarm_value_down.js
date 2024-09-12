db.periodical_alarm.updateMany({"v.change_state": {$ne: null}}, {$unset: {"v.change_state": ""}})
