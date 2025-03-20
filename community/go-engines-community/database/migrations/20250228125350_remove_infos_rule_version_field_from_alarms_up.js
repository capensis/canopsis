db.periodical_alarm.updateMany({"v.infos_rule_version": {$exists: true}}, {$unset: {"v.infos_rule_version": ""}})
