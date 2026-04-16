db.meta_alarm_rules.updateMany({}, {$set: {enabled: true}});
db.flapping_rule.updateMany({}, {$set: {enabled: true}});
db.default_snmprules.updateMany({}, {$set: {enabled: true}});
db.resolve_rule.updateMany({}, {$set: {enabled: true}});
db.user.updateMany({}, {$rename: {enable: "enabled"}});
